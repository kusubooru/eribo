package eribo

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kusubooru/eribo/flist"
	"github.com/kusubooru/eribo/loot"
)

type Player struct {
	Name   string
	Role   flist.Role
	Status flist.Status
	Fave   bool
}

func (p Player) String() string {
	return fmt.Sprintf("Name: %q, Status: %q, Role: %q, Fave: %v", p.Name, p.Status, p.Role, p.Fave)
}

type PlayerMap struct {
	sync.RWMutex
	m map[string]*Player
}

func NewPlayerMap() *PlayerMap {
	m := make(map[string]*Player)
	return &PlayerMap{m: m}
}

func (c *PlayerMap) SetPlayer(p *Player) {
	c.Lock()
	defer c.Unlock()
	c.m[p.Name] = p
}

func (c *PlayerMap) SetPlayerRole(playerName string, role flist.Role) {
	c.Lock()
	defer c.Unlock()
	if p, ok := c.m[playerName]; ok {
		p.Role = role
	}
}

func (c *PlayerMap) SetPlayerStatus(playerName string, status flist.Status) {
	c.Lock()
	defer c.Unlock()
	if p, ok := c.m[playerName]; ok {
		p.Status = status
	}
}

func (c *PlayerMap) SetPlayerFave(playerName string, fave bool) {
	c.Lock()
	defer c.Unlock()
	if p, ok := c.m[playerName]; ok {
		p.Fave = fave
	}
}

func (c *PlayerMap) DelPlayer(playerName string) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, playerName)
}

func (c *PlayerMap) GetPlayer(playerName string) (*Player, bool) {
	c.RLock()
	defer c.RUnlock()
	p, ok := c.m[playerName]
	return p, ok
}

func (c *PlayerMap) ForEach(fn func(name string, p *Player)) {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.m {
		fn(k, v)
	}
}

type ChannelMap struct {
	sync.RWMutex
	m         map[string]*PlayerMap
	lothm     map[string]*Loth
	lastLothm map[string]*Loth
}

func NewChannelMap() *ChannelMap {
	m := make(map[string]*PlayerMap)
	lothm := make(map[string]*Loth)
	lastLothm := make(map[string]*Loth)
	return &ChannelMap{m: m, lothm: lothm, lastLothm: lastLothm}
}

func (c *ChannelMap) ForEach(fn func(channel string, pm *PlayerMap)) {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.m {
		fn(k, v)
	}
}

func (c *ChannelMap) DelPlayer(channel, playerName string) {
	c.Lock()
	defer c.Unlock()
	if pm, ok := c.m[channel]; ok {
		pm.DelPlayer(playerName)
	}
	if loth, ok := c.lothm[channel]; ok {
		if loth.Name == playerName {
			delete(c.lothm, channel)
		}
	}
}

func (c *ChannelMap) DelPlayerAllChannels(playerName string) {
	c.Lock()
	defer c.Unlock()
	for channel := range c.m {
		c.m[channel].DelPlayer(playerName)
		if loth, ok := c.lothm[channel]; ok {
			if loth.Name == playerName {
				delete(c.lothm, channel)
			}
		}
	}
}

func (c *ChannelMap) SetPlayer(channel string, p *Player) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[channel]; !ok {
		pm := NewPlayerMap()
		pm.SetPlayer(p)
		c.m[channel] = pm
		return
	}
	c.m[channel].SetPlayer(p)
}

func (c *ChannelMap) SetPlayerStatus(channel, playerName string, status flist.Status) {
	c.Lock()
	defer c.Unlock()
	if pm, ok := c.m[channel]; ok {
		if p, ok := pm.GetPlayer(playerName); ok {
			p.Status = status
			c.m[channel].SetPlayer(p)
		}
	}
}

func (c *ChannelMap) GetPlayer(playerName string) (*Player, []string) {
	c.RLock()
	defer c.RUnlock()
	channels := make([]string, 0)
	var player *Player
	for channel := range c.m {
		if pm, ok := c.m[channel]; ok {
			if p, ok := pm.GetPlayer(playerName); ok {
				channels = append(channels, channel)
				if player == nil {
					player = p
				}
			}
		}
	}
	return player, channels
}

func (c *ChannelMap) PlayerMap(channel string) (*PlayerMap, bool) {
	c.RLock()
	defer c.RUnlock()
	pm, ok := c.m[channel]
	return pm, ok
}

func (c *ChannelMap) GetActivePlayers() *PlayerMap {
	c.RLock()
	defer c.RUnlock()
	actives := NewPlayerMap()
	for _, pm := range c.m {
		pm.ForEach(func(k string, v *Player) {
			if v.Status.IsActive() {
				actives.SetPlayer(v)
			}
		})
	}
	return actives
}

func (c *ChannelMap) Loth(channel string) *Loth {
	c.RLock()
	defer c.RUnlock()
	return c.lothm[channel]
}

func (c *ChannelMap) ChooseLoth(playerName, channel, botName string, d time.Duration, lowNames []string) (*Loth, bool, []*Player) {
	c.RLock()
	loth := c.lothm[channel]
	lastLoth := c.lastLothm[channel]
	c.RUnlock()
	targets := make([]*Player, 0)
	if loth != nil && !loth.Expired() {
		return loth, false, targets
	}
	pm := c.GetActivePlayers()
	pm.ForEach(func(name string, p *Player) {
		c.RLock()
		defer c.RUnlock()
		if p.Role == flist.RoleFullDom || p.Role == "" {
			return
		}
		if !p.Fave {
			return
		}
		if p.Name == botName {
			return
		}
		// Avoid choosing the same loth two times in a row.
		if lastLoth != nil && p.Name == lastLoth.Name {
			return
		}
		targets = append(targets, p)
	})
	if len(targets) == 0 {
		return nil, false, targets
	}
	target := randTarget(playerName, targets, lowNames)
	if target == nil {
		return nil, false, targets
	}
	c.Lock()
	defer c.Unlock()
	newLoth := NewLoth(target, d)
	c.lothm[channel] = newLoth
	c.lastLothm[channel] = newLoth
	return c.lothm[channel], true, targets
}

func randTarget(playerName string, targets []*Player, lowNames []string) *Player {
	t := &loot.Table{}
	for _, p := range targets {
		var weight int
		switch p.Role {
		case flist.RoleSomeDom:
			weight = 10
		case flist.RoleSwitch:
			weight = 40
		case flist.RoleSomeSub:
			weight = 45
		case flist.RoleFullSub:
			weight = 50
		}
		// Give sub players a slightly higher chance for malfunction.
		if p.Name == playerName {
			switch p.Role {
			case flist.RoleSomeSub:
				weight = 8
			case flist.RoleFullSub:
				weight = 10
			default:
				weight = 5
			}
		}
		// Give list of players who do not enjoy loth, a much lower chance to
		// be chosen.
		for _, name := range lowNames {
			if p.Name == name {
				weight = 3
			}
		}
		t.Add(p, weight)
	}
	seed := time.Now().UnixNano()
	_, loth := t.Roll(seed)
	p, ok := loth.(*Player)
	if !ok {
		return nil
	}
	return p
}

func (c *ChannelMap) GetChannel(channel string) (*PlayerMap, bool) {
	c.RLock()
	defer c.RUnlock()
	pm, ok := c.m[channel]
	return pm, ok
}

func (c *ChannelMap) Find(playerName, channelName string) []*Player {
	targets := make([]*Player, 0)
	pm, ok := c.PlayerMap(channelName)
	if !ok {
		return targets
	}
	pm.ForEach(func(name string, p *Player) {
		c.RLock()
		defer c.RUnlock()
		if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(playerName)) {
			targets = append(targets, p)
		}
	})
	return targets
}
