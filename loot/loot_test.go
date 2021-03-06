package loot

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestTable_Roll(t *testing.T) {
	drops := []Drop{
		{Item: "sword", Weight: 100},
		{Item: "shield", Weight: 400},
		{Item: "null", Weight: 499},
		{Item: "legendary", Weight: 1},
	}
	table := NewTable(drops)
	m := make(map[int]int, len(drops))
	// TODO(kusubooruji): this used to pass with 10000 rolls. Rework the numbers
	// and reduce the rolls for faster test.
	rolls := 1000000
	for k := 0; k < rolls; k++ {
		seed := time.Now().UnixNano()
		i, _ := table.Roll(seed)
		m[i]++
	}

	delta := 0.015
	for i, d := range drops {
		pr := float64(m[i]) / float64(rolls)
		expectedPr := float64(d.Weight) / float64(table.TotalWeight())
		got, want := math.Abs(float64(expectedPr)-pr), delta
		if got > want {
			t.Errorf("%d: Rolls: %d, Drops: %d, Pr: %f, Expected Pr: %f, Current delta: %f, wanted delta: %f", i, rolls, m[i], pr, expectedPr, got, want)
		}
	}
}

func TestQualityDrops(t *testing.T) {
	t.Skip("not really a test")

	d := []Drop{
		{Item: "poor", Weight: 8},
		{Item: "common", Weight: 40},
		{Item: "uncommon", Weight: 4},
	}
	table := NewTable(d)
	m := make(map[int]int, len(d))
	rolls := 1000
	for k := 0; k < rolls; k++ {
		seed := time.Now().UnixNano()
		i, _ := table.Roll(seed)
		m[i]++
	}

	drops, pr := table.Sim(rolls)
	for i := range drops {
		fmt.Printf("%s = %d, %.1f%%\n", d[i].Item, drops[i], pr[i]*100.0)
	}
}
