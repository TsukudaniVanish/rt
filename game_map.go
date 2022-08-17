package main

import (
	"math/rand"
	"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)
const (
	Wall rl.Cell = iota
	Floor 
	minCaveSize = 400
    MapWidth = UIWidth
    MapHight = UIHight - LogLines - StatusLines
)

type Path struct {
	Map *GameMap
	NBs paths.Neighbors
}

// implement Pather
func (p *Path)Neighbors(q gruid.Point) (nbs []gruid.Point){
	nbs = p.NBs.Cardinal(q, func (r gruid.Point) bool  {
		return p.Map.IsWalkable(r)
	}) 
	return 
}

type GameMap struct {
	Grid rl.Grid
	Rand *rand.Rand
	Explored map[gruid.Point]bool // explored cells
}

func NewMap(size gruid.Point) (gmap *GameMap) {
	gmap = &GameMap{
		Grid: rl.NewGrid(size.X, size.Y),
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		Explored: make(map[gruid.Point]bool),
	}
	gmap.Generate()
	return
}

func (gmap *GameMap)IsWalkable(p gruid.Point) (isWalkable bool) {
	isWalkable = gmap.Grid.At(p) == Floor && gmap.Grid.Contains(p)
	return 
}

func (gmap *GameMap)Rune(c rl.Cell) (r rune){
	switch c {
	case Wall:
		r = '#'
	case Floor:
		r = '.'
	}
	return
}

// Generate ... fills Grid attribute of gmap with a procedurally generated map
func (gmap *GameMap)Generate() {
	mapGen := rl.MapGen{Rand: gmap.Rand, Grid: gmap.Grid}
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, WallsOutOfRange: true}, 
	}

	for {
		mapGen.CellularAutomataCave(Wall, Floor, 0.42, rules)

		freep := gmap.RandFloor() // random floor cell

		pr := paths.NewPathRange(gmap.Grid.Range())
		pr.CCMap(&Path{Map: gmap}, freep)
		ntiles := mapGen.KeepCC(pr, freep, Wall)

		if ntiles > minCaveSize { // ensure map size for enemy generation
			break
		}
		
	}
}

func (gmap *GameMap)RandFloor() gruid.Point{
	size := gmap.Grid.Size()

	for {
		freep := gruid.Point{X: gmap.Rand.Intn(size.X), Y: gmap.Rand.Intn(size.Y)}
		if gmap.Grid.At(freep) == Floor {
			return freep
		}
	}
}
