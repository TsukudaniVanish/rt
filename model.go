package main

import (
	"github.com/anaseto/gruid"
)

type ActionType string 
const (
	ActionMovement ActionType = "action movement"
	ActionQuit ActionType = "action quit"
)

type Model struct {
	Grid gruid.Grid
	Game Game
	Action UIAction
}

type Game struct {
	ECS *ECS
	Map *GameMap
}

type UIAction struct {
	Type ActionType
	Delta gruid.Point
}

func (m *Model)Update(msg gruid.Msg) (eff gruid.Effect){
	m.Action = UIAction{}
	switch msg := msg.(type) {
	case gruid.MsgInit:
		// init map
		size := m.Grid.Size()
		m.Game.Map = NewMap(size)
		m.Game.ECS = NewEcs()

		// init player
		m.Game.ECS.PlayerID = m.Game.ECS.AddEntity(&Player{}, size.Div(2))
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	}
	eff =  m.handleAction()
	return 
}

func (m *Model)Draw() (grid gruid.Grid) {
	m.Grid.Fill(gruid.Cell{Rune:' '})

	// draw map 
	it := m.Game.Map.Grid.Iterator()
	for it.Next() {
		m.Grid.Set(it.P(), gruid.Cell{Rune: m.Game.Map.Rune(it.Cell()),})
	}

	// draw entity 
	for i, e := range m.Game.ECS.Entities{
		m.Grid.Set(
			m.Game.ECS.Positions[i], 
			gruid.Cell{
				Rune:e.Rune(),
				Style: gruid.Style{Fg: e.Color()},
			},
		)
	}
	return m.Grid
}

func (m *Model)updateMsgKeyDown(msg gruid.MsgKeyDown) {
	pdelta := gruid.Point{}
	switch msg.Key {
	case gruid.KeyArrowLeft, "a":
		m.Action = UIAction{Type: ActionMovement, Delta: pdelta.Shift(-1, 0)}
	case gruid.KeyArrowRight, "d":
		m.Action = UIAction{Type: ActionMovement, Delta: pdelta.Shift(1, 0)}
	case gruid.KeyArrowUp, "w":
		m.Action = UIAction{Type: ActionMovement, Delta: pdelta.Shift(0, -1)}
	case gruid.KeyArrowDown, "s":
		m.Action = UIAction{Type: ActionMovement, Delta: pdelta.Shift(0, 1)}
	case gruid.KeyEscape:
		m.Action = UIAction{Type: ActionQuit}
	}

}

func (m *Model)handleAction() (eff gruid.Effect) {
	switch m.Action.Type{
	case ActionMovement:
		np := m.Game.ECS.PlayerPosition().Add(m.Action.Delta)
		if m.Game.Map.IsWalkable(np){
			m.Game.ECS.MovePlayer(np)
		}

	case ActionQuit:
		eff = gruid.End()
	}
	return
}