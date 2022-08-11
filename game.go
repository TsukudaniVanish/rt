package main

import (
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type Game struct {
	ECS *ECS
	Map *GameMap
}


func (g *Game)Bump (to gruid.Point) {
	if !g.Map.IsWalkable(to) {
		return 
	}

	if enemy := g.ECS.EnemyAt(to); enemy != nil {
		log.Printf("You kicked the %s, much to its annoyance!\n", enemy.Name)
		return
	}

	g.ECS.MovePlayer(to)
	g.UpdateFOV() // update FOV
}

func (g *Game)UpdateFOV() {
	player := g.ECS.Player()
	playerPosition := g.ECS.PlayerPosition()

	// new range for fov
	rangeFOV := gruid.NewRange(-maxLOS, -maxLOS, maxLOS + 1, maxLOS + 1)
	player.FOV.SetRange(rangeFOV.Add(playerPosition).Intersect(g.Map.Grid.Range()))

	passible := func (p gruid.Point) bool {
		return g.Map.IsWalkable(p)
	}

	for _, p := range player.FOV.SSCVisionMap(playerPosition, maxLOS, passible, false){
		if paths.DistanceManhattan(p, playerPosition) > maxLOS {
			continue
		}
		if !g.Map.Explored[p] {
			g.Map.Explored[p] = true
		}
	}
}

func (g *Game) InFOV(p gruid.Point) bool {
	playerPosition := g.ECS.PlayerPosition()
	return g.ECS.Player().FOV.Visible(p) && paths.DistanceManhattan(playerPosition, p) <= maxLOS
}

func (g *Game)SpawnEnemies() {
	const numberOfEnemies = 6

	for i := 0; i < numberOfEnemies; i++{
		m := &Enemy{}

		// orc or troll
		switch {
		case g.Map.Rand.Intn(100) < 80:
			m.Name = "orc"
			m.Char = 'o'
		default:
			m.Name = "troll" 
			m.Char = 'T'
		}
		p := g.FreeFloorTile()
		g.ECS.AddEntity(m, p)
	}
}

func (g *Game)FreeFloorTile() (point gruid.Point) {
	for {
		p := g.Map.RandFloor()
		if g.ECS.NoBlockingEnemyAt(p){
			return p
		}
	}
}