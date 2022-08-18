package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type Consumable interface {
    // Activate makes use of item 
    Activate (g *Game, a ItemAction) error
}

type Targetter interface {
    // TargetRadius returns radius of affected area of the target  
    TargetRadius() int 
}

type ItemAction struct {
    Actor int // index of entity 
    Target *gruid.Point
}

// implement of health portion 

type HealthPotion struct {
    Amount int
    Name string 
}

func (p *HealthPotion) Activate(g *Game, a ItemAction) (err error) {
    si, ok := g.ECS.Statuses[a.Actor]
    if !ok {
        err = fmt.Errorf("%s cannot use %s", g.ECS.Name[a.Actor], p.Name)
        return 
    }
    hp := si.Heal(p.Amount)
    if hp <= 0 {
        err = errors.New("your health is already full")
        return
    }
    g.Logf("You used portion", colorStatusHealthy)
   return
}

type MagicArrowScroll struct {
    Damage int
    Range int 
}

//
func (ms *MagicArrowScroll) Activate(g *Game, a ItemAction) (err error) {
    targetID := -1
    minDist := ms.Range + 1 
    for i := range g.ECS.Statuses {
        pos := g.ECS.Positions[i]
        if a.Actor == i || g.ECS.Dead(i) || !g.InFOV(pos) {
            continue
        }
        dist := paths.DistanceManhattan(g.ECS.Positions[a.Actor], pos)
        if dist < minDist {
            targetID = i 
            minDist = dist
        }
    }
    if targetID < 0 {
        err = errors.New("there is no enemy")
        return 
    }
    st, ok := g.ECS.Statuses[targetID]
    if ok {
        g.Logf("a magic lightning strikes %v", colorStatusHealthy, g.ECS.Name[targetID])
        st.Damage(ms.Damage)
    } else {
        log.Fatalf("could not find status of %d", targetID)
    }
    return 
}

type ExplodeScroll struct {
    Damage int 
    Radius int 
}

func (es *ExplodeScroll) Activate(g *Game, a ItemAction) (err error) {
    if a.Target == nil {
        err = errors.New("you have to choose a target")
        return 
    }
    p := *a.Target
    if !g.InFOV(p){
        err = errors.New("you cannot target where you cannot see")
        return 
    }
    hit := 0
    for i, st := range g.ECS.Statuses {
        q := g.ECS.Positions[i]
        if q == g.ECS.PlayerPosition() || g.ECS.Dead(i) {
            continue
        }
        g.Logf("%v is engulfed in vortex of mana", colorStatusHealthy, g.ECS.GetName(i))
        st.Damage(es.Damage)
        hit++
    }
    if hit == 0{
        err = errors.New("there is no enemy in range of explosion")
        return 
    }
    return 
}

func (es *ExplodeScroll) TargetRadius() (radius int) {
    radius = es.Radius
    return 
}