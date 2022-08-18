package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

    "game"
    "domain"
    "save"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type ActionType string 
const (
	NoAction ActionType = "no action"
	ActionBump ActionType = "action bump"
    ActionDrop ActionType = "action drop"
    ActionInventory ActionType = "action inventory"
    ActionPickup ActionType = "action pickup"
    ActionWait ActionType = "action wait"
    ActionSave ActionType = "action save"
	ActionQuit ActionType = "action quit"
    ActionViewMessage ActionType = "action view message"
    ActionExamine ActionType = "action examine a map"
)

type UIMode int 
const (
    modeNormal UIMode = iota
    modeEnd 
    modeMenu
    modeMessageViewer
    modeInventoryActivate
    modeInventoryDrop
    modeTargetting 
    modeExamination // map examination mode 
)

type MenuEntry int 
const (
    MenuNewGame MenuEntry = iota
    MenuContinue
    MenuQuit
    
)

type Model struct {
	Grid gruid.Grid
	Game *game.Game
	Action UIAction
    Mode UIMode
    Inventory *ui.Menu
    GameMenu *ui.Menu
    MenuInfoLabel *ui.Label // for menu info (errors) 
    LogLabel *ui.Label
    StatusLabel *ui.Label
    DescLabel *ui.Label // label for description
    Viewer *ui.Pager
    Target Targetting
}

type Targetting struct {
    Position gruid.Point // target position in ui (* != map position)
    ItemID int // item to use after select a target
    Radius int 
}

type UIAction struct {
	Type ActionType
	Delta gruid.Point
}

// init ... initialize Model
// only called at App Initialize
func (m *Model) init() (eff gruid.Effect) {
    m.MenuInfoLabel = &ui.Label{}
    m.LogLabel = &ui.Label{}
    m.StatusLabel = &ui.Label{}
    m.DescLabel = &ui.Label{Box: &ui.Box{}}
    m.InitializeMessageViewer()
    m.Mode = modeMenu

    menuEntries := []ui.MenuEntry{
        MenuNewGame: {Text: ui.Text("(N)ew game"), Keys: []gruid.Key{"N", "n"}},
        MenuContinue: {Text: ui.Text("(C)ontinue last game"), Keys: []gruid.Key{"C", "c"}},
        MenuQuit: {Text: ui.Text("(Q)uit game"), Keys: []gruid.Key{"Q", "q", gruid.KeyEscape}},
    }

    m.GameMenu = ui.NewMenu(ui.MenuConfig{
        Grid: gruid.NewGrid(domain.UIWidth / 2, len(menuEntries) + 2),
        Entries: menuEntries,
        Box: &ui.Box{Title: ui.Text("Game Menu")},
    })
    return
}

func (m *Model)Update(msg gruid.Msg) (eff gruid.Effect){
    switch msg.(type) {
    case gruid.MsgInit:
        return m.init()
    }
    // reset last action
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
    case modeMenu:
        return m.updateMenu(msg)
    case modeMessageViewer:
        m.Viewer.Update(msg)
        if m.Viewer.Action() == ui.PagerQuit {
            m.Mode = modeNormal
        }
        return nil

    case modeInventoryDrop, modeInventoryActivate:
        m.updateInventory(msg)
        return nil
    case modeTargetting, modeExamination:
        println("update targetting!")
        m.updateTargetting(msg)
        println(fmt.Sprintf("%v", m.Target.Position))
        return nil
    default: // modeNormal
        switch msg := msg.(type) {
        case gruid.MsgKeyDown:
            m.updateMsgKeyDown(msg)
        case gruid.MsgMouse:
            if msg.Action == gruid.MouseMove {
                m.Target.Position = msg.P
            }
        }
        eff =  m.handleAction()   
        return
    }
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
    case "x":
        m.Action = UIAction{Type: ActionExamine}
    case "S":
        m.Action = UIAction{Type: ActionSave}
	}

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
            radius, err := m.Game.TargetingRadius(m.Game.ECS.PlayerID, n)
            if err != nil {
                if err.Error() == domain.ErrNoTargeting { // no targetting
                    err = m.Game.InventoryUseItem(m.Game.ECS.PlayerID, n)
                } else { // error 
                    m.Game.Logf("%v", domain.ColorLogSpecial, err)
                }
            } else { // change mode to targetting 
                m.Target = Targetting{
                    ItemID: n,
                    Position: m.Game.ECS.PlayerPosition().Shift(0, domain.LogLines),
                    Radius: radius,
                }
                m.Mode = modeTargetting
                return
            }
        }
        if err != nil {
            m.Game.Logf("%v", domain.ColorLogSpecial, err)
        } else {
            m.Game.EndTurn()
        }
        m.Mode = modeNormal
    }
}

func (m *Model)updateTargetting(msg gruid.Msg) {
    mapRange := m.getMapRange()
    if !m.Target.Position.In(mapRange) {
        m.Target.Position = m.Game.ECS.PlayerPosition().Add(mapRange.Min)
    }
    p := m.convertUiPositionToMapPosition(m.Target.Position)
    switch msg := msg.(type) {
    case gruid.MsgKeyDown:
        switch msg.Key {
        case gruid.KeyArrowLeft, "a":
           p = p.Shift(-1, 0)
        case gruid.KeyArrowRight, "d":
           p = p.Shift(1, 0)
        case gruid.KeyArrowUp, "w":
           p = p.Shift(0, -1)
        case gruid.KeyArrowDown, "s":
           p = p.Shift(0, 1)
        case gruid.KeyEnter, ".":
            if m.Mode == modeExamination {
                break
            }
            m.activateTarget(p)
            return 
        case gruid.KeyEscape, "q":
            m.Target = Targetting{}
            m.Mode = modeNormal
            return
        }
        m.Target.Position = m.convertMapPositionToUiPosition(p)
    case gruid.MsgMouse:
        switch msg.Action{
        case gruid.MouseMove:
            m.Target.Position = msg.P
        case gruid.MouseMain:
            if m.Mode == modeExamination {
                break
            }
            m.activateTarget(p)
            return 
        }
    }
}

func (m *Model) updateMenu(msg gruid.Msg) (eff gruid.Effect){
    rg := m.getUIRange().Intersect(m.getUIRange().Add(m.menuAnchor()))
    m.GameMenu.Update(rg.RelMsg(msg))
    switch m.GameMenu.Action() {
    case ui.MenuMove:
        m.MenuInfoLabel.SetText("")
    case ui.MenuInvoke:
        m.MenuInfoLabel.SetText("")
        switch m.GameMenu.Active() {
        case int(MenuNewGame):
            m.Game = game.NewGame()
            m.Mode = modeNormal
        case int(MenuContinue):
            m.loadGame()
        case int(MenuQuit):
            eff = gruid.End()
            return 
        }
    case ui.MenuQuit:
        eff = gruid.End()
        return
    }
    return nil 
}

func (m *Model) activateTarget(p gruid.Point) {
    err := m.Game.InventoryUseItemWithTarget(m.Game.ECS.PlayerID, m.Target.ItemID, &p)
    if err != nil {
        m.Game.Logf("%v", domain.ColorLogSpecial, err)
    } else {
        m.Game.EndTurn()
    }
    m.Target = Targetting{}
    m.Mode = modeNormal
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
    case ActionExamine:
        m.Mode = modeExamination
        m.Target.Position = m.Game.ECS.PlayerPosition().Shift(0, domain.LogLines)
    case ActionQuit:
		eff = gruid.End()
    case ActionSave:
        m.saveGame()
	}
    if m.Game.ECS.PlayerDead() {
        m.Game.Logf("You Died -- press Escape to quit", domain.ColorLogSpecial)
        m.Mode = modeEnd
        return nil 
    }
	return
}

func (m *Model) InitializeMessageViewer() {
    m.Viewer = ui.NewPager(ui.PagerConfig{
        Grid: gruid.NewGrid(domain.UIWidth, domain.UIHight),
        Box: &ui.Box{},
    })
}

func (m *Model)Draw() (grid gruid.Grid) {
    mapGrid := m.Grid.Slice(m.getMapRange())
    switch m.Mode {
    case modeMenu:
        grid = m.DrawMenu()
        return 
    case modeMessageViewer:
        m.Grid.Copy(m.Viewer.Draw())
        grid = m.Grid
        return 
    case modeInventoryDrop, modeInventoryActivate:
        mapGrid.Copy(m.Inventory.Draw())
        grid = m.Grid
        return 
    }
    
    // init grid
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
			c.Style.Bg = domain.ColorFOV
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
    // for examine mode 
    if m.Mode == modeExamination || m.Mode == modeTargetting {
        p := m.convertUiPositionToMapPosition(m.Target.Position)
        c := mapGrid.At(p)
        c.Rune = '+'
        mapGrid.Set(p, c)
    }


    // draw ui's
    m.DrawLog(m.Grid.Slice(m.Grid.Range().Lines(0, domain.LogLines)))
    m.DrawStatus(m.Grid.Slice(m.Grid.Range().Lines(m.Grid.Size().Y -domain.StatusLines, m.Grid.Size().Y)))
    grid = m.Grid
	return
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
    st.Fg = domain.ColorStatusHealthy
    g := m.Game
    statusPlayer := g.ECS.Statuses[g.ECS.PlayerID]
    if statusPlayer.HP < statusPlayer.MaxHP / 2 {
        st.Fg = domain.ColorStatusWounded
    }
    m.StatusLabel.Content = ui.Textf("HP: %d/%d", statusPlayer.HP, statusPlayer.MaxHP)
    m.StatusLabel.Box = &ui.Box{Title: ui.Text("Status")}
    m.StatusLabel.Draw(gd)
}

func (m *Model) DrawNames(gd gruid.Grid) {
    maprg := m.getMapRange()
    if !m.Target.Position.In(maprg) {
        return 
    }
    p := m.Target.Position.Sub(maprg.Min)
    names := []string{}
    for i, q := range m.Game.ECS.Positions {
        if q != p || !m.Game.InFOV(q) {
            continue 
        }
        name := m.Game.ECS.GetName(i)
        if name != "" {
            names = append(names, name)
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
    if p.X + 1 + width >= domain.UIWidth {
        rg = rg.Shift(-1 -width, 0, -1 -width, 0)
    }
    if p.Y + 2> domain.MapHight {
        rg = rg.Shift(0, -1, 0, -1)
    }
    if p.Y -1 <0 {
        rg = rg.Shift(0, 1, 0, -1)
    }
    slice := gd.Slice(rg)
    m.DescLabel.Content = ui.Text(text)
    m.DescLabel.Draw(slice)
}

func (m *Model)DrawMenu() (grid gruid.Grid) {
    m.Grid.Fill(gruid.Cell{Rune: ' '})
    m.Grid.Slice(m.GameMenu.Bounds().Add(m.menuAnchor())).Copy(m.GameMenu.Draw())
    m.MenuInfoLabel.Draw(m.Grid.Slice(m.Grid.Range().Line(12).Shift(10, 0, 0, 0)))
    grid = m.Grid
    return 
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
            if err.Error() == domain.ErrNoShow {
                continue 
            }
            g.Logf("Could not pickup: %v", domain.ColorStatusWounded, err)
            return 
        }
        g.Logf("You pickup: %v", domain.ColorStatusHealthy, g.ECS.Name[i])
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
        Grid: gruid.NewGrid(40, domain.MapHight),
        Box: &ui.Box{Title: ui.Text(title),},
        Entries: entries,
    })
}

func (m *Model) getMapRange() gruid.Range {
    return gruid.NewRange(0, domain.LogLines, domain.UIWidth, domain.UIHight - domain.StatusLines)
}

func (m *Model) getUIRange() gruid.Range {
    return gruid.NewRange(0, 0, domain.UIWidth, domain.UIHight)
}

func (m *Model) convertUiPositionToMapPosition(uipos gruid.Point) gruid.Point {
    mrg := m.getMapRange()
    return uipos.Sub(mrg.Min)
}

func (m *Model) convertMapPositionToUiPosition (mapPos gruid.Point) gruid.Point {
    mrg := m.getMapRange()
    return mapPos.Add(mrg.Min)
}

func (m *Model) menuAnchor() (p gruid.Point) {
    p = gruid.Point{X: 10, Y:6}
    return
}

func (m *Model) saveGame() {
    data, err := save.EncodeNoGzip(m.Game)
    if err != nil {
        m.Game.Logf("could not save game", domain.ColorStatusWounded)
        log.Fatal(err)
        return
    }
    
    err = save.SaveFile("save", data)
    if err != nil {
        m.Game.Logf("could not save game", domain.ColorStatusWounded)
        log.Fatal(err)
        return 
    }
    m.Game.Logf("game saved successfully!", domain.ColorLogSpecial)
}

func (m *Model)loadGame() {
    data, err := save.LoadFile("save")
    if err != nil {
        m.MenuInfoLabel.SetText(err.Error())
        return
    }

    g, err := save.DecodeNoGzip(data)
    if err != nil {
        m.MenuInfoLabel.SetText(err.Error())
        return
    }
    m.Game = g 
    m.Game.Map.SetRand(rand.New(rand.NewSource(time.Now().UnixNano())))
    m.Mode = modeNormal

    m.Game.Logf("load game successfully!", domain.ColorLogSpecial)
}
