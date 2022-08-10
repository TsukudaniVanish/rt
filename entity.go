package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

const maxLOS = 10

// Entity Component System
type ECS struct {
	Entities []Entity
	Positions map[int] gruid.Point // key: index of entity value: position of entity
	PlayerID int // index of player
}

func NewEcs() *ECS{
	return &ECS{
		Positions: map[int] gruid.Point{},
	}
}

func (ecs *ECS)AddEntity(e Entity, p gruid.Point) (id int) {
	id = len(ecs.Entities)
	ecs.Entities = append(ecs.Entities, e)
	ecs.Positions[id] = p
	return 
}

func (ecs *ECS)MoveEntity(id int, p gruid.Point) {
	ecs.Positions[id] = p
}

func (ecs *ECS)MovePlayer(p gruid.Point) {
	ecs.MoveEntity(ecs.PlayerID, p)
}

func (ecs *ECS)Player() (player Entity) {
	player = ecs.Entities[ecs.PlayerID]
	return
}

func (ecs *ECS)PlayerPosition() (p gruid.Point){
	p = ecs.Positions[ecs.PlayerID]
	return
}

type Entity interface {
	Rune() rune
	Color() gruid.Color
}

type Player struct {
	FOV *rl.FOV // player'S field of view
}

func NewPlayer() (player *Player) {
	player = &Player{}
	player.FOV = rl.NewFOV(gruid.NewRange(-maxLOS, -maxLOS, maxLOS + 1, maxLOS + 1))
	return 
}

func (p *Player)Rune() (r rune){
	r = '@'
	return
}

func (p *Player)Color() (color gruid.Color) {
	color = gruid.ColorDefault
	return 
}