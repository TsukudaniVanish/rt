package game

import (
	"math/rand"
	"time"

	"domain"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
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
	rand *rand.Rand
	Explored map[gruid.Point]bool // explored cells
}

func NewMap(size gruid.Point) (gmap *GameMap) {
	gmap = &GameMap{
		Grid: rl.NewGrid(size.X, size.Y),
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		Explored: make(map[gruid.Point]bool),
	}
	gmap.Generate()
	return
}

func (gmap *GameMap)SetRand(rand *rand.Rand) {
	gmap.rand = rand 
}

func (gmap *GameMap)IsWalkable(p gruid.Point) (isWalkable bool) {
	isWalkable = gmap.Grid.At(p) == domain.Floor && gmap.Grid.Contains(p)
	return 
}

func (gmap *GameMap)Rune(c rl.Cell) (r rune){
	switch c {
	case domain.Wall:
		r = '#'
	case domain.Floor:
		r = '.'
	}
	return
}

// Generate ... fills Grid attribute of gmap with a procedurally generated map
func (gmap *GameMap)Generate() {
	mapGen := rl.MapGen{Rand: gmap.rand, Grid: gmap.Grid}
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, WallsOutOfRange: true}, 
	}

	for {
		mapGen.CellularAutomataCave(domain.Wall, domain.Floor, 0.42, rules)

		freep := gmap.RandFloor() // random floor cell

		pr := paths.NewPathRange(gmap.Grid.Range())
		pr.CCMap(&Path{Map: gmap}, freep)
		ntiles := mapGen.KeepCC(pr, freep, domain.Wall)

		if ntiles > domain.MinCaveSize { // ensure map size for enemy generation
			break
		}
		
	}
}

func (gmap *GameMap)RandFloor() gruid.Point{
	size := gmap.Grid.Size()

	for {
		freep := gruid.Point{X: gmap.rand.Intn(size.X), Y: gmap.rand.Intn(size.Y)}
		if gmap.Grid.At(freep) == domain.Floor {
			return freep
		}
	}
}
