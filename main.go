package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Karocyt/Npupu/internal/heuristics"
	"github.com/Karocyt/Npupu/internal/parser"
	"github.com/Karocyt/Npupu/internal/solver"
)

func printError(e error) {
	if e == nil {
		return
	}
	fmt.Print("ERROR: ")
	fmt.Fprintln(os.Stderr, e)
	flag.Usage()
	os.Exit(1)
}

func parseCmd() (string, map[string]solver.ScoreFn, bool, bool, int, int) {
	var filename string
	var aStar, compare, uniform, display, classic bool
	var randomSize, shuffleCount int
	heuristic := 2

	hUsage := "Available heuristics:\n"
	for i, h := range heuristics.Functions {
		if i != 0 {
			hUsage += fmt.Sprintf("\t%d: %s\n", i, h.Name)
		}
	}

	flag.StringVar(&filename, "f", "", "filename of your input file")
	flag.IntVar(&heuristic, "h", 1, hUsage)
	flag.BoolVar(&aStar, "shortest", false, "uses A* algorithm to find an optimal solution")
	flag.BoolVar(&compare, "vs", false, "compare greedy search and Astar performance")
	flag.BoolVar(&uniform, "ref", false, "adds uniform-cost search for reference")
	flag.BoolVar(&display, "display", false, "force print of full solution in any case")
	flag.BoolVar(&classic, "classic", false, "uses an ascendant order solution instead of a snail one")
	flag.IntVar(&randomSize, "size", 3, "size of the generated puzzle if filename omitted, (3 minimum)")
	flag.IntVar(&shuffleCount, "mix", 200, "numbers of scramblings of the generated puzzle if filename omitted")

	flag.Parse()

	if heuristic < 1 || heuristic >= len(heuristics.Functions) || flag.NArg() != 0 ||
		(filename != "" && randomSize < 3) || shuffleCount < 0 {
		flag.Usage()
		os.Exit(1)
	}

	heuristicsMap := map[string]solver.ScoreFn{}
	if compare || !aStar {
		heuristicsMap["Greedy "+heuristics.Functions[heuristic].Name] = heuristics.Functions[heuristic].Fn
	}
	if compare || aStar {
		heuristicsMap["A* "+heuristics.Functions[heuristic].Name] = heuristics.Functions[heuristic].Astar
	}
	if uniform {
		heuristicsMap[heuristics.Functions[0].Name] = heuristics.Functions[0].Fn
	}

	return filename, heuristicsMap, display, classic, randomSize, shuffleCount
}

func main() {
	filename, heuristicsMap, display, classic, randomSize, shuffleCount := parseCmd()
	input, size, e := parser.Parse(filename)
	printError(e)
	solvers := make([]*solver.Solver, 0, 2)
	finalPos, finalGrid, input := solver.Init(size, classic, input, randomSize, shuffleCount)
	heuristics.Init(finalGrid, finalPos)
	for name, fn := range heuristicsMap {
		s := solver.New(input, size, fn, name)
		solvers = append(solvers, s)
		go s.Solve()
	}

	if len(solvers) == 1 {
		display = true
	}
	for i := range solvers {
		res, ok := <-solvers[i].Solution
		stats := <-solvers[i].Stats
		solvers[i].PrintRes(res, ok, stats, display)
	}
	os.Exit(0)
}
