package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)
const (
	Wall rl.Cell = iota
	Floor 
)

type GameMap struct {
	Grid rl.Grid
}

func NewMap(size gruid.Point) (gmap *GameMap) {
	gmap = &GameMap{}
	gmap.Grid = rl.NewGrid(size.X, size.Y)
	gmap.Grid.Fill(Floor)
	for i := 0; i < 3; i++ {
		// We add a few walls. We'll deal with map generation
		// in the next part of the tutorial.
		gmap.Grid.Set(gruid.Point{X: 30 + i, Y: 12}, Wall)
	}
	return
}

func (gmap *GameMap)IsWalkable(p gruid.Point) (isWalkable bool) {
	isWalkable = gmap.Grid.At(p) == Floor
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