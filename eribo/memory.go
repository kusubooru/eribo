package eribo

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kusubooru/eribo/flist"
)

type Player struct {
	Name   string
	Role   flist.Role
	Status flist.Status
}

func (p Player) String() string {
	return fmt.Sprintf("Name: %q, Status: %q, Role: %q", p.Name, p.Status, p.Role)
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
	m     map[string]*PlayerMap
	lothm map[string]*Loth
}

func NewChannelMap() *ChannelMap {
	m := make(map[string]*PlayerMap)
	lothm := make(map[string]*Loth)
	return &ChannelMap{m: m, lothm: lothm}
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

func (c ChannelMap) GetActivePlayers() *PlayerMap {
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

func (c *ChannelMap) ChooseLoth(channel, botName string, d time.Duration) (*Loth, bool) {
	loth := c.lothm[channel]
	if loth != nil && !loth.Expired() {
		return loth, false
	}
	victims := make([]*Player, 0)
	pm := c.GetActivePlayers()
	pm.ForEach(func(name string, p *Player) {
		c.RLock()
		defer c.RUnlock()
		if p.Role == flist.RoleFullDom || p.Role == flist.RoleSomeDom || p.Role == "" {
			return
		}
		if p.Name == botName {
			return
		}
		victims = append(victims, p)
	})
	if len(victims) == 0 {
		return nil, false
	}
	victim := randVictim(victims)
	c.Lock()
	defer c.Unlock()
	c.lothm[channel] = NewLoth(victim, d)
	return c.lothm[channel], true
}

func newRand(n int) int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Intn(n)
}

func randVictim(victims []*Player) *Player {
	return victims[newRand(len(victims))]
}

func (c ChannelMap) GetChannel(channel string) (*PlayerMap, bool) {
	c.RLock()
	defer c.RUnlock()
	pm, ok := c.m[channel]
	return pm, ok
}
