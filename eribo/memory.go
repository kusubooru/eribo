package eribo

import (
	"sync"

	"github.com/kusubooru/eribo/flist"
)

type Player struct {
	Name   string
	Role   flist.Role
	Status flist.Status
	//Status string
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

func (c *PlayerMap) ForEach(fn func(k string, v *Player)) {
	c.Lock()
	defer c.Unlock()
	for name, player := range c.m {
		fn(name, player)
	}
}

type ChannelMap struct {
	sync.RWMutex
	m map[string]*PlayerMap
}

func NewChannelMap() *ChannelMap {
	m := make(map[string]*PlayerMap)
	return &ChannelMap{m: m}
}

func (c *ChannelMap) DelPlayer(channel, playerName string) {
	c.Lock()
	defer c.Unlock()
	if pm, ok := c.m[channel]; ok {
		pm.DelPlayer(playerName)
	}
}

func (c *ChannelMap) DelPlayerAllChannels(playerName string) {
	c.Lock()
	defer c.Unlock()
	for channel := range c.m {
		c.m[channel].DelPlayer(playerName)
	}
}

func (c *ChannelMap) SetPlayer(channel string, p *Player) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[channel]; !ok {
		pm := NewPlayerMap()
		pm.SetPlayer(p)
		//c.m[channel] = map[string]*Player{p.Name: p}
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

func (c ChannelMap) GetChannel(channel string) (*PlayerMap, bool) {
	c.RLock()
	defer c.RUnlock()
	pm, ok := c.m[channel]
	return pm, ok
}
