package loot

import (
	"math/rand"
	"sync"
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

func (t *Table) Add(item interface{}, weight int) {
	t.Lock()
	defer t.Unlock()
	d := Drop{Item: item, Weight: weight}
	t.drops = append(t.drops, d)
}

func (t Table) TotalWeight() int {
	t.RLock()
	defer t.RUnlock()
	var totalWeight int
	for _, d := range t.drops {
		totalWeight += d.Weight
	}
	return totalWeight
}

func (t Table) Roll(seed int64) (int, interface{}) {
	t.RLock()
	defer t.RUnlock()
	totalWeight := t.TotalWeight()

	r := rand.New(rand.NewSource(seed))
	roll := r.Intn(totalWeight + 1)

	var weight int
	var drop int
	for i, d := range t.drops {
		weight += d.Weight
		if roll <= weight {
			drop = i
			break
		}
	}
	return drop, t.drops[drop].Item
}
