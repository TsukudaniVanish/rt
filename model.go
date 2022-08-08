package main

import (
	"github.com/anaseto/gruid"
)

const (
	ActionMovement = "action movement"
	ActionQuit = "action quit"
)

type Model struct {
	Grid gruid.Grid
	Game Game
	Action UIAction
}

type Game struct {
	PlayerPosition gruid.Point
}

type UIAction struct {
	Type string
	Delta gruid.Point
}

func (m *Model)Update(msg gruid.Msg) (eff gruid.Effect){
	m.Action = UIAction{}
	switch msg := msg.(type) {
	case gruid.MsgInit:
		m.Game.PlayerPosition = m.Grid.Size().Div(2)
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	}
	eff =  m.handleAction()
	return 
}

func (m *Model)Draw() (grid gruid.Grid) {
	it := m.Grid.Iterator()
	for it.Next() {
		switch {
		case it.P() == m.Game.PlayerPosition:
			it.SetCell(gruid.Cell{Rune: '@'})
		default:
			it.SetCell(gruid.Cell{Rune: ' '}) 
		}
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
	return
}