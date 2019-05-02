package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kusubooru/eribo/loot"
	"github.com/kusubooru/eribo/rp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		rolls   = flag.Int("rolls", 1000000, "sim rolls")
		tkSim   = flag.Bool("tk", false, "")
		tieSim  = flag.Bool("tie", false, "")
		tiedSim = flag.Bool("tied", false, "")
		tieHard = flag.Bool("tiehard", false, "")
	)
	flag.Parse()

	if *tkSim {
		fmt.Println(simTktools(*rolls))
		return nil
	}

	if *tieSim {
		toolType := ""
		if *tieHard {
			toolType = "hard"
		}
		fmt.Println(simTietools(*rolls, toolType))
		return nil
	}

	if *tiedSim {
		toolType := ""
		if *tieHard {
			toolType = "hard"
		}
		simTietoolsDecreasedWeight(*rolls, toolType)
		return nil
	}

	return nil
}

func simTktools(rolls int) string {
	table := &loot.Table{}
	for _, t := range rp.Tktools() {
		table.Add(t, t.Weight())
	}
	drops, pr := table.Sim(rolls)
	var buf bytes.Buffer
	buf.WriteString("\n")
	for i, t := range rp.Tktools() {
		buf.WriteString(fmt.Sprintf("%9s %35s = %5d, %.3f%%\n", t.Quality, t.Name(), drops[i], pr[i]*100.0))
	}
	return buf.String()
}

func simTietools(rolls int, toolType string) string {
	table := &loot.Table{}

	for _, t := range rp.Tietools(toolType) {
		table.Add(t, t.Quality.Weight())
	}
	drops, pr := table.Sim(rolls)
	var buf bytes.Buffer
	buf.WriteString("\n")
	for i, t := range rp.Tietools(toolType) {
		buf.WriteString(fmt.Sprintf("%9s %35s = %5d, %.3f%%\n", t.Quality, t.Name(), drops[i], pr[i]*100.0))
	}
	return buf.String()
}

func simTietoolsDecreasedWeight(rolls int, toolType string) {
	table := &loot.Table{}

	for _, t := range rp.Tietools(toolType) {
		table.Add(t, t.Quality.Weight())
	}

	for i := 0; i < rolls; i++ {
		seed := time.Now().UnixNano()
		_, item := table.RollDecreaseWeight(seed)
		if item == nil {
			break
		}
		tool, ok := item.(rp.Tietool)
		if !ok {
			fmt.Println("item was not tietool")
		}
		fmt.Printf("Rolled %9s, %35s, total weight = %d\n", tool.Quality, tool.Name(), table.TotalWeight())
	}
}
