package main

import (

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

const (
	maxLOS = 10
)

// Entity Component System
type ECS struct {
	Entities  []Entity
	Positions map[int]gruid.Point // key: index of entity value: position of entity
	PlayerID  int                 // index of player

	Statuses map[int]*Status
	AI       map[int]*EnemyAI
	Name     map[int]string
}

func NewEcs() *ECS {
	return &ECS{
		Positions: map[int]gruid.Point{},
		Statuses:  map[int]*Status{-1: nil},
		AI:        map[int]*EnemyAI{-1: nil},
		Name:      map[int]string{},
	}
}

func (ecs *ECS) AddEntity(e Entity, p gruid.Point) (id int) {
	id = len(ecs.Entities)
	ecs.Entities = append(ecs.Entities, e)
	ecs.Positions[id] = p
	return
}

func (ecs *ECS) MoveEntity(id int, p gruid.Point) {
	ecs.Positions[id] = p
}

func (ecs *ECS) MovePlayer(p gruid.Point) {
	ecs.MoveEntity(ecs.PlayerID, p)
}

func (ecs *ECS) Player() (player *Player) {
	player = ecs.Entities[ecs.PlayerID].(*Player)
	return
}

func (ecs *ECS) PlayerPosition() (p gruid.Point) {
	p = ecs.Positions[ecs.PlayerID]
	return
}

func (ecs *ECS) EnemyAt(p gruid.Point) (id int, enemy *Enemy) {
	for i, q := range ecs.Positions {
		if p != q || !ecs.Alive(i) {
			continue
		}
		switch e := ecs.Entities[i].(type) {
		case *Enemy:
			id = i
			enemy = e
			return
		}
	}
	id = -1 
	enemy = nil 
	return
}

func (ecs *ECS) NoBlockingEnemyAt(p gruid.Point) (noBlockingEnemy bool) {
	i, _:= ecs.EnemyAt(p)
	noBlockingEnemy = ecs.PlayerPosition() != p && !ecs.Alive(i)
	return
}

func (ecs *ECS) PlayerDead() bool {
	return ecs.Dead(ecs.PlayerID)
}

func (ecs *ECS) Alive(i int) (isAlive bool) {
	st := ecs.Statuses[i]
	isAlive = st != nil && st.HP > 0
	return

}

func (ecs *ECS) Dead(i int) (isDead bool) {
	st := ecs.Statuses[i]
	isDead = st != nil && st.HP <= 0
	return
}

func (ecs *ECS) Style(i int) (r rune, c gruid.Color) {
	r = ecs.Entities[i].Rune()
	c = ecs.Entities[i].Color()
	if ecs.Dead(i) {
		r = '%'
		c = gruid.ColorDefault
	}
	return
}

// RenderOrder ... Priority of rendaring
type RenderOrder int

const (
	roNone RenderOrder = iota
	roCorpse
	roItem
	roActor
)

func (ecs *ECS) GetRenderOrder(i int) (ro RenderOrder) {
	switch ecs.Entities[i].(type) {
	case *Player:
		ro = roActor
	case *Enemy:
		if ecs.Dead(i) {
			ro = roCorpse
		} else {
			ro = roActor
		}
	}
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
	player.FOV = rl.NewFOV(gruid.NewRange(-maxLOS, -maxLOS, maxLOS+1, maxLOS+1))
	return
}

func (p *Player) Rune() (r rune) {
	r = '@'
	return
}

func (p *Player) Color() (color gruid.Color) {
	color = colorPlayer
	return
}

type Enemy struct {
	Char rune
}

func (e *Enemy) Rune() (r rune) {
	r = e.Char
	return
}

func (e *Enemy) Color() (color gruid.Color) {
	color = colorEnemy
	return
}
