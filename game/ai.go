package game

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)


type AIPath struct {
    Game *Game 
    NB paths.Neighbors
}

func (aip *AIPath) Neighbors(q gruid.Point) (nbs []gruid.Point) {
    nbs = aip.NB.Cardinal(q, func(r gruid.Point) bool{
        return aip.Game.Map.IsWalkable(r)
    })
    return
}


func (aip *AIPath) Cost(p, q gruid.Point) (cost int) {
    if aip.Game.ECS.NoBlockingEnemyAt(q) {
        // extra cost for blocked positions
        return 8
    }
    return 1
}

func (aip *AIPath) Estimation(p,q gruid.Point) int {
    return paths.DistanceManhattan(p, q)
}
