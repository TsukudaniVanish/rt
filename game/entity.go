package game

import (
	"domain"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

// Entity Component System
type ECS struct {
	Entities  map[int]Entity
	Positions map[int]gruid.Point // key: index of entity value: position of entity in map
	PlayerID  int                 // index of player
    NextID int 
	Bodies int 

	Statuses map[int]*Status
	AI       map[int]*EnemyAI
	Name     map[int]string
    Styles map[int]Style
    Inventories map[int]*Inventory
}

func NewEcs() *ECS {
	return &ECS{
        Entities: map[int]Entity{},
		Positions: map[int]gruid.Point{},
		Statuses:  map[int]*Status{},
		AI:        map[int]*EnemyAI{},
		Name:      map[int]string{},
        Styles: map[int]Style{},
        Inventories: map[int]*Inventory{},
        NextID: 0,
	}
}

func (ecs *ECS) AddEntity(e Entity, p gruid.Point) (id int) {
	id = ecs.NextID
    ecs.Entities[id] = e 
	ecs.Positions[id] = p
    ecs.NextID++
	return 
}

func (ecs *ECS) RemoveEntity (id int) {
    delete(ecs.Entities, id)
    delete(ecs.Positions, id)
    delete(ecs.Statuses, id)
    delete(ecs.AI, id)
    delete(ecs.Name, id)
    delete(ecs.Styles, id)
    delete(ecs.Inventories, id)
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

func (ecs *ECS) GetStyle(i int) (r rune, c gruid.Color) {
	r = ecs.Styles[i].Rune
	c = ecs.Styles[i].Color
	if ecs.Dead(i) {
		r = '%'
		c = gruid.ColorDefault
	}
	return
}

func (ecs *ECS) GetName(i int) (name string) {
    name = ecs.Name[i]
    if ecs.Dead(i) {
        name = "corpse"
    }
    return 
}

func (ecs *ECS) GameClear() bool {
	return ecs.Bodies >= domain.EnemyNumber
}

// RenderOrder ... Priority of rendering
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
    case Consumable:
        ro = roItem
	}
	return
}

type Entity interface {}

type Player struct {
	FOV *rl.FOV // player'S field of view
}

func NewPlayer() (player *Player) {
	maxLOS := domain.MaxLOS
	player = &Player{}
	player.FOV = rl.NewFOV(gruid.NewRange(-maxLOS, -maxLOS, maxLOS+1, maxLOS+1))
	return
}

type Enemy struct {}
