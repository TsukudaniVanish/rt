package main 

import (
    "errors"
    "fmt"

	"github.com/anaseto/gruid"
)

type Consumable interface {
    // Activate makes use of item 
    Activate (g *Game, a ItemAction) error
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
    si := g.ECS.Statuses[a.Actor]
    if si == nil {
        err = fmt.Errorf("%s cannot use %s", g.ECS.Name[a.Actor], p.Name)
        return 
    }
    hp := si.Heal(p.Amount)
    if hp <= 0 {
        err = errors.New("Your health is already full")
        return
    }
   return
}
