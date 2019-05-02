package loot

import (
	"math/rand"
	"sync"
	"time"
)

type Drop struct {
	Item   interface{}
	Weight int
}

type Table struct {
	sync.RWMutex
	drops []Drop
}

var DefaultTable = &Table{drops: []Drop{}}

func NewTable(d []Drop) *Table {
	if d == nil {
		return DefaultTable
	}
	return &Table{drops: d}
}

func (t *Table) Drops() []Drop {
	t.RLock()
	defer t.RUnlock()
	return t.drops
}

func (t *Table) Add(item interface{}, weight int) {
	t.Lock()
	defer t.Unlock()
	d := Drop{Item: item, Weight: weight}
	t.drops = append(t.drops, d)
}

func (t *Table) TotalWeight() int {
	t.RLock()
	defer t.RUnlock()
	var totalWeight int
	for _, d := range t.drops {
		totalWeight += d.Weight
	}
	return totalWeight
}

func (t *Table) Len() int {
	t.RLock()
	defer t.RUnlock()
	if t.drops == nil {
		return 0
	}
	return len(t.drops)
}

func (t *Table) Roll(seed int64) (int, interface{}) {
	t.RLock()
	defer t.RUnlock()
	totalWeight := t.TotalWeight()
	if totalWeight == 0 {
		return 0, nil
	}

	r := rand.New(rand.NewSource(seed))
	roll := r.Intn(totalWeight + 1)

	var weight int
	var drop int
	for i, d := range t.drops {
		weight += d.Weight
		if weight >= roll && roll <= weight {
			drop = i
			break
		}
	}
	return drop, t.drops[drop].Item
}

type namer interface {
	Name() string
}

func (t *Table) RollDecreaseWeight(seed int64) (int, interface{}) {
	roll, item := t.Roll(seed)
	// TODO(kusuboorujin): item can be nil. Maybe add a nil check in the future
	// or just not let item be nil.
	rolledItem, ok := item.(namer)
	if !ok {
		return roll, item
	}
	t.Lock()
	defer t.Unlock()
	for i := 0; i < len(t.drops); i++ {
		if t.drops[i].Item == nil {
			continue
		}
		n, ok := t.drops[i].Item.(namer)
		if !ok {
			continue
		}
		if n.Name() == rolledItem.Name() {
			t.drops[i].Weight--
			if t.drops[i].Weight < 0 {
				t.drops[i].Weight = 0
			}
		}
	}
	return roll, item
}

func (t *Table) Sim(rolls int) (map[int]int, map[int]float64) {
	t.RLock()
	defer t.RUnlock()
	dropsMap := make(map[int]int, t.Len())
	for k := 0; k < rolls; k++ {
		seed := time.Now().UnixNano()
		i, _ := t.Roll(seed)
		dropsMap[i]++
	}

	prMap := make(map[int]float64, t.Len())
	for i := range t.drops {
		prMap[i] = float64(dropsMap[i]) / float64(rolls)
	}
	return dropsMap, prMap
}
