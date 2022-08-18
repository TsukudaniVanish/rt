package game

import (
	"errors"
	"fmt"
	"strings"

    "domain"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)


type Game struct {
	ECS *ECS
	Map *GameMap
    PR *paths.PathRange
    Logs []LogEntry
}

func NewGame() (g *Game) {
   g = &Game{}

    // init map
    size := gruid.Point{X:domain.MapWidth, Y:domain.MapHight}
   g.Map = NewMap(size)
   g.PR = paths.NewPathRange(gruid.NewRange(0, 0, size.X, size.Y))
   g.ECS = NewEcs()

    // init player
   g.ECS.PlayerID =g.ECS.AddEntity(NewPlayer(),g.Map.RandFloor())
   g.ECS.Statuses[g.ECS.PlayerID] = &Status{
        HP: 30, MaxHP: 30, Power: 5, Defence: 2,
    }
   g.ECS.Styles[g.ECS.PlayerID] = Style{Rune: '@', Color: domain.ColorPlayer}
   g.ECS.Name[g.ECS.PlayerID] = domain.PlayerName
   g.ECS.Inventories[g.ECS.PlayerID] = &Inventory{}

   g.UpdateFOV()

    // add enemies 
   g.SpawnEnemies()

    // add Items 
   g.PlaceItems()
   return
}

// Bump ... player move or attack
func (g *Game)Bump (to gruid.Point) {
	if !g.Map.IsWalkable(to) {
		return 
	}

	if i,enemy := g.ECS.EnemyAt(to); enemy != nil {
        g.BumpAttack(g.ECS.PlayerID, i)
        if g.ECS.Dead(i) {
            g.ECS.Bodies ++
        }
        g.EndTurn()
		return
	}

	g.ECS.MovePlayer(to)
    g.EndTurn()
}

func (g *Game) EndTurn() {
    g.UpdateFOV()
    for i, e := range g.ECS.Entities{
        if g.ECS.PlayerDead(){
            return
        }
        switch e.(type) {
        case *Enemy:
            g.HandleMonsterTurn(i)
        case *Player:
            isHeal := g.Map.rand.Intn(100) < domain.HealRate
            if isHeal {
                g.ECS.Statuses[i].Heal(2)
            }
        }
    }
}

func (g *Game)UpdateFOV() {
	player := g.ECS.Player()
	playerPosition := g.ECS.PlayerPosition()

	// new range for fov
    maxLOS := domain.MaxLOS
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
	return g.ECS.Player().FOV.Visible(p) && paths.DistanceManhattan(playerPosition, p) <= domain.MaxLOS
}

func (g *Game)SpawnEnemies() {
	const numberOfEnemies = 12
	for i := 0; i < numberOfEnemies; i++{
		m := &Enemy{}
        const (
            orc = iota
            troll 
        )
        kind := orc

		// orc or troll
		switch {
		case g.Map.rand.Intn(100) < 80:
		default:
            kind = troll
		}
		p := g.FreeFloorTile()
        i := g.ECS.AddEntity(m, p)
        switch kind {
            case orc:
                g.ECS.Statuses[i] = &Status{
                    HP: 10, MaxHP: 10,Power: 3, Defence: 0,
                }
                g.ECS.Name[i] = "orc"
                g.ECS.Styles[i] = Style{Rune: 'o', Color: domain.ColorEnemy}
            case troll:
                g.ECS.Statuses[i] = &Status{
                    HP: 16, MaxHP: 16,Power: 5, Defence: 1,
                }
                g.ECS.Name[i] = "troll"
                g.ECS.Styles[i] = Style{Rune: 'T', Color: domain.ColorEnemy}

        }
        g.ECS.AI[i] = &EnemyAI{}
	}
}

func (g *Game) BumpAttack(i, j int) {
    si := g.ECS.Statuses[i]
    sj := g.ECS.Statuses[j]
    damage := si.Power - sj.Defence
    attackDesc := fmt.Sprintf("%s attacks %s", strings.Title(g.ECS.Name[i]), strings.Title(g.ECS.Name[j]))
    color := domain.ColorLogEnemyAttack
    if i == g.ECS.PlayerID {
        color = domain.ColorLogPlayerAttack
    }
    if damage > 0 {
        g.Logf("%s for %d damage", color, attackDesc, damage)
        sj.HP -= damage
    } else {
        g.Logf("%s\nbut does no damage", color, attackDesc)
    }
}

func (g *Game)PlaceItems() {
    const numberOfPortions = 5
    const amount = 100
    for i := 0; i< numberOfPortions; i++{
        r := g.Map.rand.Float64()
        p := g.FreeFloorTile()

        switch {
        case r < 0.7: // portion
            name := "portion"
            id := g.ECS.AddEntity(&HealthPotion{Amount: amount, Name: name}, p)
            g.ECS.Styles[id] = Style{Rune: '!', Color: domain.ColorConsumable}
            g.ECS.Name[id] = name
        case r < 0.9: // magicArrow
            name := "magic arrow scroll"
            id := g.ECS.AddEntity(&MagicArrowScroll{Damage: 3, Range: 5}, p)
            g.ECS.Styles[id] = Style{Rune: '?', Color: domain.ColorConsumable}
            g.ECS.Name[id] = name
        default:
            name := "explode scroll"
            id := g.ECS.AddEntity(&ExplodeScroll{Damage: 100, Radius: 10}, p)
            g.ECS.Styles[id] = Style{Rune: '?', Color: domain.ColorConsumable}
            g.ECS.Name[id] = name
        }
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

func (g *Game) HandleMonsterTurn(i int) {
    if !g.ECS.Alive(i) {
        return 
    }
    p := g.ECS.Positions[i]
    ai := g.ECS.AI[i]
    aip := &AIPath{ Game: g}
    playerPosition := g.ECS.PlayerPosition()
    if paths.DistanceManhattan(p,playerPosition) == 1 {
        g.BumpAttack(i, g.ECS.PlayerID)
        return 
    }
    if !g.InFOV(p) {
        if len(ai.Path) < 1 {
            ai.Path = g.PR.AstarPath(aip, p, g.Map.RandFloor())
        }
        g.AIMove(i)
        return 
    }
    ai.Path = g.PR.AstarPath(aip, p, playerPosition)
    g.AIMove(i)
}

func (g *Game) AIMove(i int) {
    ai := g.ECS.AI[i]
    if len(ai.Path) > 0 && ai.Path[0] == g.ECS.Positions[i] {
        ai.Path = ai.Path[1:]
    }
    if len(ai.Path) > 0 && g.ECS.NoBlockingEnemyAt(ai.Path[0]){
        g.ECS.MoveEntity(i, ai.Path[0])
        ai.Path = ai.Path[1:]
    }
}

func (g *Game)log(e LogEntry) {
    if len(g.Logs) > 0 {
        if g.Logs[len(g.Logs) -1].Text == e.Text {
            return 
        }
    }
    g.Logs = append(g.Logs, e)
}

func (g *Game)Logf(format string, color gruid.Color, a ...interface{}) {
    e := LogEntry{
        Text: fmt.Sprintf(format, a ...),
        Color: color,
    }
    g.log(e)
}

// InventoryAdd ... add an item to actors's inventry 
func (g *Game) InventoryAdd(actor, i int) (err error) {
   switch g.ECS.Entities[i].(type) {
    case Consumable:
        inv := g.ECS.Inventories[actor]
        inv.Items = append(inv.Items, i)
        delete(g.ECS.Positions, i)
        return
   }
   err = errors.New(domain.ErrNoShow)
   return 
}

// InventoryRemove ... remove an item at itemID from actor's inventry
func (g *Game) InventoryRemove(actor, itemID int) (err error) {
    inv := g.ECS.Inventories[actor]
    i := inv.Items[itemID]
    inv.Items = inv.Items[:len(inv.Items) -1]
    g.ECS.Positions[i] = g.ECS.PlayerPosition()
    return 
}


// InventoryUseItem ... Use an item at itemID from actor's inventory 
func (g *Game)InventoryUseItem(actor, itemID int) (err error) {
    err = g.InventoryUseItemWithTarget(actor, itemID, nil)
    return
}

func (g *Game)InventoryUseItemWithTarget(actor, itemID int, target *gruid.Point) (err error) {
    inv := g.ECS.Inventories[actor]
    i := inv.Items[itemID]
    switch e := g.ECS.Entities[i].(type) {
    case Consumable:
        itemAction := ItemAction{Actor: actor, Target: target}
        err = e.Activate(g, itemAction)
        if err != nil {
            return
        }

    }
    inv.Items = inv.Items[:len(inv.Items) -1]
    return 
}

// TargetingRadius ... returns target radius of item at itemID of actor's inventory if it needs target
func (g *Game) TargetingRadius (actor int, itemID int) (radius int, err error) {
    inv := g.ECS.Inventories[actor]
    if len(inv.Items) <= itemID{
        err = errors.New("index is out of range")
        return 
    }
    item := inv.Items[itemID]
    switch e := g.ECS.Entities[item].(type) {
    case Targetter:
        radius = e.TargetRadius()
        return
    default:
        err = errors.New(domain.ErrNoTargeting)
        return  
    }
}