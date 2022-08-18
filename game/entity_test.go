package game 

import "testing"

func TestItemStyle(t *testing.T) {
	g := NewGame()
	for i := range g.ECS.Positions{
		e := g.ECS.Entities[i]
		switch e.(type) {
		case Consumable:
			r, _ := g.ECS.GetStyle(i)
			if r == '!' || r == '?' {
				t.Fatalf("rune: %c", r)
			}
		}
	}
}