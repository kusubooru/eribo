package loot

import (
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
	rolls := 10000
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
