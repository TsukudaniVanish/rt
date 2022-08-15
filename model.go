package main

import (
	"sort"
	"unicode/utf8"
    "strings"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/ui"
)

const (
    playerName = "You"
    UIWidth = 80
    UIHight = 24
)

const (
	colorFOV gruid.Color = iota + 1
	colorPlayer 
	colorEnemy
    colorConsumable
    colorLogPlayerAttack
    colorLogEnemyAttack
    colorLogSpecial
    colorStatusHealthy
    colorStatusWounded
)
type ActionType string 
const (
	NoAction ActionType = "no action"
	ActionBump ActionType = "action bump"
    ActionDrop ActionType = "action drop"
    ActionInventory ActionType = "action inventory"
    ActionPickup ActionType = "action pickup"
    ActionWait ActionType = "action wait"
	ActionQuit ActionType = "action quit"
    ActionViewMessage = "action view message"
)

type UIMode int 
const (
    modeNormal UIMode = iota
    modeEnd 
    modeMessageViewer
    modeInventoryActivate
    modeInventoryDrop
)

type Model struct {
	Grid gruid.Grid
	Game *Game
	Action UIAction
    Mode UIMode
    Inventory *ui.Menu
    LogLabel *ui.Label
    StatusLabel *ui.Label
    DescLabel *ui.Label // label for description
    Viewer *ui.Pager
    MousePos gruid.Point
}


type UIAction struct {
	Type ActionType
	Delta gruid.Point
}

func (m *Model)Update(msg gruid.Msg) (eff gruid.Effect){
	m.Action = UIAction{}
    switch m.Mode{
    case modeEnd:
        switch msg := msg.(type) {
            case gruid.MsgKeyDown:
                switch msg.Key{
                case gruid.KeyEscape:
                    eff = gruid.End()
                    return 
                }
        }
        return nil
    case modeMessageViewer:
        m.Viewer.Update(msg)
        if m.Viewer.Action() == ui.PagerQuit {
            m.Mode = modeNormal
        }
        return nil

    case modeInventoryDrop, modeInventoryActivate:
        m.updateInventory(msg)
        return nil
    }
	switch msg := msg.(type) {
	case gruid.MsgInit:
        m.LogLabel = &ui.Label{}
        m.StatusLabel = &ui.Label{}
        m.DescLabel = &ui.Label{Box: &ui.Box{}}
        m.InitializeMessageViewer()

		m.Game = &Game{}

		// init map
		size := m.Grid.Size()
        size.Y -= 3 // for log and status
		m.Game.Map = NewMap(size)
        m.Game.PR = paths.NewPathRange(gruid.NewRange(0, 0, size.X, size.Y))
		m.Game.ECS = NewEcs()

		// init player
		m.Game.ECS.PlayerID = m.Game.ECS.AddEntity(NewPlayer(), m.Game.Map.RandFloor())
        m.Game.ECS.Statuses[m.Game.ECS.PlayerID] = &Status{
            HP: 30, MaxHP: 30, Power: 5, Defence: 2,
        }
        m.Game.ECS.Styles[m.Game.ECS.PlayerID] = Style{Rune: '@', Color: colorPlayer}
        m.Game.ECS.Name[m.Game.ECS.PlayerID] = playerName
        m.Game.ECS.Inventories[m.Game.ECS.PlayerID] = &Inventory{}

		m.Game.UpdateFOV()

		// add enemies 
		m.Game.SpawnEnemies()

        // add Items 
        m.Game.PlaceItems()
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
    case gruid.MsgMouse:
        if msg.Action == gruid.MouseMove {
            m.MousePos = msg.P
        }
	}
	eff =  m.handleAction()
	return 
}

func (m *Model)Draw() (grid gruid.Grid) {
    mapGrid := m.Grid.Slice(m.Grid.Range().Shift(0, 2, 0, -1))
    switch m.Mode {
    case modeMessageViewer:
        m.Grid.Copy(m.Viewer.Draw())
        grid = m.Grid
        return 
    case modeInventoryDrop, modeInventoryActivate:
        mapGrid.Copy(m.Inventory.Draw())
        grid = m.Grid
        return 
    }
    

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
        mapGrid.Set(it.P(), c)
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
		c := mapGrid.At(p)
		c.Rune, c.Style.Fg = g.ECS.GetStyle(i)
		mapGrid.Set(p, c)
	}
    m.DrawNames(mapGrid)
    m.DrawLog(m.Grid.Slice(m.Grid.Range().Lines(0, 2)))
    m.DrawStatus(m.Grid.Slice(m.Grid.Range().Line(m.Grid.Size().Y - 1)))
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
    case "m":
        m.Action = UIAction{Type: ActionViewMessage}
    case "i":
        m.Action = UIAction{Type: ActionInventory}
    case "D":
        m.Action = UIAction{Type: ActionDrop}
    case "g":
        m.Action = UIAction{Type: ActionPickup}
	}

}

func (m *Model)handleAction() (eff gruid.Effect) {
	switch m.Action.Type{
	case ActionBump:
		np := m.Game.ECS.PlayerPosition().Add(m.Action.Delta)
		m.Game.Bump(np)
    case ActionDrop:
        m.OpenInventory("Drop Item")
        m.Mode = modeInventoryDrop
    case ActionInventory:
        m.OpenInventory("Use item")
        m.Mode = modeInventoryActivate
    case ActionPickup:
        m.PickUpItem()
	case ActionWait:
        m.Game.EndTurn()
    case ActionViewMessage:
        m.Mode = modeMessageViewer
        lines := []ui.StyledText{}
        for _, e := range m.Game.Logs {
            st := gruid.Style{}
            st.Fg = e.Color
            lines = append(lines, ui.NewStyledText(e.String(), st))
        }
        m.Viewer.SetLines(lines)
    case ActionQuit:
		eff = gruid.End()
	}
    if m.Game.ECS.PlayerDead() {
        m.Game.Logf("You Died -- press Escape to quit", colorLogSpecial)
        m.Mode = modeEnd
        return nil 
    }
	return
}

func (m *Model) InitializeMessageViewer() {
    m.Viewer = ui.NewPager(ui.PagerConfig{
        Grid: gruid.NewGrid(UIWidth, UIHight),
        Box: &ui.Box{},
    })
}

func (m *Model) DrawLog(gd gruid.Grid) {
    j := 1 

    for i := len(m.Game.Logs) - 1; i >= 0; i --{
        if j < 0 {
            break
        }
        e := m.Game.Logs[i]
        st := gruid.Style{}
        st.Fg = e.Color
        m.LogLabel.Content = ui.NewStyledText(e.String(), st)
        m.LogLabel.Draw(gd.Slice(gd.Range().Line(j)))
        j--
    }
}


func (m *Model) DrawStatus(gd gruid.Grid) {
    st := gruid.Style{}
    st.Fg = colorStatusHealthy
    g := m.Game
    statusPlayer := g.ECS.Statuses[g.ECS.PlayerID]
    if statusPlayer.HP < statusPlayer.MaxHP / 2 {
        st.Fg = colorStatusWounded
    }
    m.StatusLabel.Content = ui.Textf("HP: %d/%d", statusPlayer.HP, statusPlayer.MaxHP)
    m.StatusLabel.Draw(gd)
}

func (m *Model) DrawNames(gd gruid.Grid) {
    maprg := gruid.NewRange(0, 2, UIWidth, UIWidth)
    if !m.MousePos.In(maprg) {
        return 
    }
    p := m.MousePos.Sub(gruid.Point{X:0, Y: 2})
    names := []string{}
    for i, q := range m.Game.ECS.Positions {
        if q != p || !m.Game.InFOV(q) {
            continue 
        }
        name, ok := m.Game.ECS.GetName(i)
        if ok {
            if m.Game.ECS.Alive(i) {
                names = append(names, name)
            } else {
                names = append(names, "corpse")
            }
        }
    }

    if len(names) == 0 {
        return 
    }   
    sort.Strings(names)

    text := strings.Join(names, ", ")
    width := utf8.RuneCountInString(text) + 2
    rg := gruid.NewRange(p.X +1 , p.Y -1, p.X + 1 + width, p.Y + 2)
    // if box is on edge. adjust place of the box
    if p.X + 1 + width >= UIWidth {
        rg = rg.Shift(-1 -width, 0, -1 -width, 0)
    }
    if p.Y + 2> MapHight {
        rg = rg.Shift(0, -1, 0, -1)
    }
    if p.Y -1 <0 {
        rg = rg.Shift(0, 1, 0, -1)
    }
    slice := gd.Slice(rg)
    m.DescLabel.Content = ui.Text(text)
    m.DescLabel.Draw(slice)
}

func (m *Model) PickUpItem() {
    g := m.Game 
    pp := g.ECS.PlayerPosition()
    // search item at pp 
    for i, p := range g.ECS.Positions{
        if p != pp {
            continue
        }
        err := g.InventoryAdd(g.ECS.PlayerID, i)
        if err != nil {
            if err.Error() == ErrNoShow {
                continue 
            }
            g.Logf("Could not pickup: %v", colorStatusWounded, err)
            return 
        }
        g.Logf("You pickup: %v", colorStatusHealthy, g.ECS.Name[i])
        g.EndTurn()
        return 
    }

}

func (m *Model)OpenInventory(title string) {
    inv := m.Game.ECS.Inventories[m.Game.ECS.PlayerID]
    entries := []ui.MenuEntry{}
    r := 'a'
    for _, it := range inv.Items {
        name := m.Game.ECS.Name[it]
        entries = append(entries, ui.MenuEntry{
            Text: ui.Text(string(r) + " - " + name),
            Keys: []gruid.Key{gruid.Key(r)},
        })
        r++
    }
    m.Inventory = ui.NewMenu(ui.MenuConfig{
        Grid: gruid.NewGrid(40, MapHight),
        Box: &ui.Box{Title: ui.Text(title),},
        Entries: entries,
    })
}

func (m *Model) updateInventory(msg gruid.Msg) {
    m.Inventory.Update(msg)
    switch m.Inventory.Action(){
    case ui.MenuQuit:
        m.Mode = modeNormal
        return 
    case ui.MenuInvoke:
        n := m.Inventory.Active()
        var err error 
        switch m.Mode{
        case modeInventoryDrop:
            err = m.Game.InventoryRemove(m.Game.ECS.PlayerID, n)
        case modeInventoryActivate:
            err = m.Game.InventoryUseItem(m.Game.ECS.PlayerID, n)
        }
        if err != nil {
            m.Game.Logf("%v", colorLogSpecial, err)
        } else {
            m.Game.EndTurn()
        }
        m.Mode = modeNormal
    }
}
