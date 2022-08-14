package main

import (
    "strings"
    "fmt"
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type Game struct {
	ECS *ECS
	Map *GameMap
    PR *paths.PathRange
    Logs []*LogEntry
}

// Bump ... player move or attack
func (g *Game)Bump (to gruid.Point) {
	if !g.Map.IsWalkable(to) {
		return 
	}

	if i,enemy := g.ECS.EnemyAt(to); enemy != nil {
        g.BumpAttack(g.ECS.PlayerID, i)
        g.EndTurn()
		return
	}

	g.ECS.MovePlayer(to)
    g.EndTurn()
    return 
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
        }
    }
    return 
}

func (g *Game)UpdateFOV() {
	player := g.ECS.Player()
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
	return g.ECS.Player().FOV.Visible(p) && paths.DistanceManhattan(playerPosition, p) <= maxLOS
}

func (g *Game)SpawnEnemies() {
	const numberOfEnemies = 6
	for i := 0; i < numberOfEnemies; i++{
		m := &Enemy{}

		// orc or troll
		switch {
		case g.Map.Rand.Intn(100) < 80:
			m.Char = 'o'
		default:
			m.Char = 'T'
		}
		p := g.FreeFloorTile()
        i := g.ECS.AddEntity(m, p)
        switch m.Char {
            case 'o':
                g.ECS.Statuses[i] = &Status{
                    HP: 10, MaxHP: 10,Power: 3, Defence: 0,
                }
                g.ECS.Name[i] = "orc"
            case 'T':
                g.ECS.Statuses[i] = &Status{
                    HP: 16, MaxHP: 16,Power: 5, Defence: 1,
                }
                g.ECS.Name[i] = "troll"

        }
        g.ECS.AI[i] = &EnemyAI{}
	}
}

func (g *Game) BumpAttack(i, j int) {
    si := g.ECS.Statuses[i]
    sj := g.ECS.Statuses[j]
    damage := si.Power - sj.Defence
    attackDesc := fmt.Sprintf("%s attacks %s", strings.Title(g.ECS.Name[i]), strings.Title(g.ECS.Name[j]))
    color := colorLogEnemyAttack
    if i == g.ECS.PlayerID {
        color = colorLogPlayerAttack
    }
    if damage > 0 {
        g.Logf("%s for %d damage", color, attackDesc, damage)
        sj.HP -= damage
    } else {
        g.Logf("%s\nbut does no damage", color, attackDesc)
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

func (g *Game)log(e *LogEntry) {
    if len(g.Logs) > 0 {
        if g.Logs[len(g.Logs) -1].Text == e.Text {
            return 
        }
    }
    g.Logs = append(g.Logs, e)
}

func (g *Game)Logf(format string, color gruid.Color, a ...interface{}) {
    e := &LogEntry{
        Text: fmt.Sprintf(format, a ...),
        Color: color,
    }
    g.log(e)
}
