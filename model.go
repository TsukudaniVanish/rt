package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

const colorFOV gruid.Color = iota + 1

type ActionType string 
const (
	ActionMovement ActionType = "action movement"
	ActionQuit ActionType = "action quit"
)

type Model struct {
	Grid gruid.Grid
	Game *Game
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
		m.Game = &Game{}

		// init map
		size := m.Grid.Size()
		m.Game.Map = NewMap(size)
		m.Game.ECS = NewEcs()

		// init player
		m.Game.ECS.PlayerID = m.Game.ECS.AddEntity(NewPlayer(), m.Game.Map.RandFloor())
		m.Game.UpdateFOV()
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

	// draw entity 
	for i, e := range g.ECS.Entities{
		p := g.ECS.Positions[i]
		if !g.Map.Explored[p] || !g.InFOV(p) {
			continue
		}
		c := m.Grid.At(p)
		c.Rune = e.Rune()
		c.Style.Fg = e.Color()
		m.Grid.Set(p, c)
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
		m.Game.MovePlayer(np)
	case ActionQuit:
		eff = gruid.End()
	}
	return
}

func (g *Game)MovePlayer (to gruid.Point) {
	if !g.Map.IsWalkable(to) {
		return 
	}

	g.ECS.MovePlayer(to)
	g.UpdateFOV() // update FOV
}

func (g *Game)UpdateFOV() {
	player := g.ECS.Player().(*Player)
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
	return g.ECS.Player().(*Player).FOV.Visible(p) && paths.DistanceManhattan(playerPosition, p) <= maxLOS
}