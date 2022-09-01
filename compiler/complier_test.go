package compiler

import (
	"domain"
	"testing"

	"github.com/anaseto/gruid"
	"github.com/stretchr/testify/assert"
)


func TestCompiler(t *testing.T) {
	type TestItem struct {
		Arg string 
		IsSuccess bool 
		ExpectedMagic domain.Magic
	}
	table := map[string] TestItem{
		"success: gandr 1 2": {
			Arg: "(gandr 1 2)",
			IsSuccess: true,
			ExpectedMagic: domain.Magic{
				Amount: amountGandr,
				Target: gruid.Point{X: 1, Y: 2},
				Radius: radiusGandr,
				Name: "gandr",
			},
		},
	}

	for key, item := range table {
		magic ,err := Compile(item.Arg)
		if item.IsSuccess && err != nil {
			t.Fatal(err)
		}

		if item.IsSuccess {
			assert.Equal(t, item.ExpectedMagic, magic, key)
		}
	}
}
