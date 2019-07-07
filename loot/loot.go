package loot

import (
	"math/rand"
	"sync"
	"time"
)

// Drop represents an item drop.
type Drop struct {
	Item   interface{}
	Weight int
}

// Table represents a loot table.
type Table struct {
	sync.RWMutex
	drops []Drop
}

// DefaultTable is the default loot table that is going to be used when the user
// specifies no drops for the new table.
var DefaultTable = &Table{drops: []Drop{}}

// NewTable creates a new loot table based on the user specified drops. If no
// drops are specified an empty default table is used.
func NewTable(d []Drop) *Table {
	if d == nil {
		return DefaultTable
	}
	return &Table{drops: d}
}

// Drops returns the item drops of a loot table.
func (t *Table) Drops() []Drop {
	t.RLock()
	defer t.RUnlock()
	return t.drops
}

// Add adds an item to the loot table with a specific weight chance.
func (t *Table) Add(item interface{}, weight int) {
	t.Lock()
	defer t.Unlock()
	d := Drop{Item: item, Weight: weight}
	t.drops = append(t.drops, d)
}

// TotalWeight reports the total weight of the items in the loot table.
func (t *Table) TotalWeight() int {
	t.RLock()
	defer t.RUnlock()
	var totalWeight int
	for _, d := range t.drops {
		totalWeight += d.Weight
	}
	return totalWeight
}

// Len returns the length of the table i.e. how many items it holds.
func (t *Table) Len() int {
	t.RLock()
	defer t.RUnlock()
	if t.drops == nil {
		return 0
	}
	return len(t.drops)
}

// Roll uses a random seed to randomly select an item from the loot table.
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

// RollDecreaseWeight returns a random item and decreases its weight.
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

// Sim simulates a number of rolls.
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
