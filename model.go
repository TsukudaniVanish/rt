package main

import (
	"log"
    "sort"

	"github.com/anaseto/gruid"
    "github.com/anaseto/gruid/paths"
)

const playerName = "You"

const (
	colorFOV gruid.Color = iota + 1
	colorPlayer 
	colorEnemy
)
type ActionType string 
const (
	NoAction ActionType = "no action"
	ActionBump ActionType = "action bump"
    ActionWait ActionType = "action wait"
	ActionQuit ActionType = "action quit"
)

type Model struct {
	Grid gruid.Grid
	Game *Game
	Action UIAction
}


type UIAction struct {
	Type ActionType
	Delta gruid.Point
}

func (m *Model)Update(msg gruid.Msg) (eff gruid.Effect){
	m.Action = UIAction{}
	switch msg := msg.(type) {
	case gruid.MsgInit:
		m.Game = &Game{}

		// init map
		size := m.Grid.Size()
		m.Game.Map = NewMap(size)
        m.Game.PR = paths.NewPathRange(gruid.NewRange(0, 0, size.X, size.Y))
		m.Game.ECS = NewEcs()

		// init player
		m.Game.ECS.PlayerID = m.Game.ECS.AddEntity(NewPlayer(), m.Game.Map.RandFloor())
        m.Game.ECS.Statuses[m.Game.ECS.PlayerID] = &Status{
            HP: 30, MaxHP: 30, Power: 5, Defence: 2,
        }
        m.Game.ECS.Name[m.Game.ECS.PlayerID] = playerName

		m.Game.UpdateFOV()

		// add enemies 
		m.Game.SpawnEnemies()
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	}
	eff =  m.handleAction()
	return 
}

func (m *Model)Draw() (grid gruid.Grid) {
	m.Grid.Fill(gruid.Cell{Rune:' '})
	g := m.Game
	// draw map 
	it := g.Map.Grid.Iterator()
	for it.Next() {
		if !g.Map.Explored[it.P()] {
			continue
		}

		c := gruid.Cell{Rune: g.Map.Rune(it.Cell()),}
		if g.InFOV(it.P()) {
			c.Style.Bg = colorFOV
		}
		m.Grid.Set(it.P(), c)
	}

    // sort entity by RenderOrder
    sortedEntities := make([]int, 0, len(g.ECS.Entities))
    for i := range g.ECS.Entities{
        sortedEntities = append(sortedEntities, i)
    }
    sort.Slice(sortedEntities, func(i, j int) bool{
        return g.ECS.GetRenderOrder(sortedEntities[i]) < g.ECS.GetRenderOrder(sortedEntities[j])
    })

	// draw entity 
	for _, i := range sortedEntities{
		p := g.ECS.Positions[i]
		if !g.Map.Explored[p] || !g.InFOV(p) {
			continue
		}
		c := m.Grid.At(p)
		c.Rune, c.Style.Fg = g.ECS.Style(i)
		m.Grid.Set(p, c)
	}
	return m.Grid
}

func (m *Model)updateMsgKeyDown(msg gruid.MsgKeyDown) {
	pdelta := gruid.Point{}
	switch msg.Key {
	case gruid.KeyArrowLeft, "a":
		m.Action = UIAction{Type: ActionBump, Delta: pdelta.Shift(-1, 0)}
	case gruid.KeyArrowRight, "d":
		m.Action = UIAction{Type: ActionBump, Delta: pdelta.Shift(1, 0)}
	case gruid.KeyArrowUp, "w":
		m.Action = UIAction{Type: ActionBump, Delta: pdelta.Shift(0, -1)}
	case gruid.KeyArrowDown, "s":
		m.Action = UIAction{Type: ActionBump, Delta: pdelta.Shift(0, 1)}
    case gruid.KeyEnter, ".":
        m.Action = UIAction{Type: ActionWait}
	case gruid.KeyEscape:
		m.Action = UIAction{Type: ActionQuit}
	}

}

func (m *Model)handleAction() (eff gruid.Effect) {
	switch m.Action.Type{
	case ActionBump:
		np := m.Game.ECS.PlayerPosition().Add(m.Action.Delta)
		m.Game.Bump(np)
	case ActionWait:
        m.Game.EndTurn()
    case ActionQuit:
		eff = gruid.End()
	}
    if m.Game.ECS.PlayerDead() {
        log.Printf("You Died")
        eff = gruid.End()
    }
	return
}
