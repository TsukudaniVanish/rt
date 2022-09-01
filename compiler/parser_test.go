package compiler

import (
	"domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func recursiveEqualNode(t *testing.T, expect, node *domain.Node, key string) {
	if expect == nil {
		assert.Nil(t, node, key)
		return 
	}

	assert.NotNil(t, node, key)
	if node != nil {
		assert.Equal(t, expect.Type, node.Type, key)
		assert.Equal(t, expect.Data, node.Data, key)
		recursiveEqualNode(t, expect.Left, node.Left, key)
		recursiveEqualNode(t, expect.Right, node.Right, key)
		return 
	}
}

func TestParser(t *testing.T) {
	type TestItem struct {
		Tokens []domain.LexicalObject
		IsSuccess bool 
		ExpectedAST *domain.Node
	}
	table := map[string] TestItem{
		"success: (gandr 1 2)": {
			Tokens: []domain.LexicalObject{
				{Type: domain.Symbol, Label: domain.SymbolParenthesisOpen, Word: "("},
				{Type: domain.KeyWord, Label: domain.KeyWordGandr, Word: "gandr"},
				{Type: domain.NumberLiteral, Label: domain.LabelNull, Word: "1"},
				{Type:domain.NumberLiteral, Label: domain.LabelNull, Word: "2"},
				{Type: domain.Symbol, Label: domain.SymbolParenthesisClose, Word: ")"},
			},
			IsSuccess: true,
			ExpectedAST: &domain.Node{
				Type: domain.Root,
				Left: &domain.Node{
					Type: domain.Operator,
					Data: domain.NodeData{
						Label: domain.OperatorGandr,
					},
				},
				Right: &domain.Node{
					Type: domain.ArrayHead,
					Right: &domain.Node{
						Type: domain.ArrayItem,
						Left: &domain.Node{
							Type: domain.Literal,
							Data: domain.NodeData{
								Label: domain.Number,
								Number: 1,
							},
						},
						Right: &domain.Node{
							Type: domain.ArrayItem,
							Left: &domain.Node{
								Type: domain.Literal,
								Data: domain.NodeData{
									Label: domain.Number,
									Number: 2,
								},
							},
							Right: &domain.Node{
								Type: domain.ArrayTail,
							},
						},
					},
				},
			},
		},
	}

	for key, item := range table {
		ast, err := parse(item.Tokens)
		if item.IsSuccess && err != nil {
			t.Fatal(err)
		}

		if item.IsSuccess {
			assert.NotNil(t, ast, key)
			recursiveEqualNode(t, item.ExpectedAST, ast, key)
		}
	}
}