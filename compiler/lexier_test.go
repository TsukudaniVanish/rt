package compiler

import (
	"testing"
	"github.com/stretchr/testify/assert"

	"domain"
)

func TestConsumeWord(t *testing.T) {
	type TestItem struct {
		Arg string 
		ExpectWord string 
		ExpectRes string 
	}
	table := map[string]TestItem {
		"parenthesis": {"(gandr 1 2)", "(", "gandr 1 2)"},
		"word": {"1 2)", "1", " 2)"},
		"one word": {"12", "12", ""},
		"punctuator": {"2)", "2", ")"},
	}

	for key, item := range table{
		w, r := consumeWord(item.Arg)
		if w != item.ExpectWord {
			t.Fatalf("%s: expect %s but got %s", key, item.ExpectWord, w)
		}
		if r != item.ExpectRes {
			t.Fatalf("%s: expect %s but got %s", key, item.ExpectRes, r)
		}
	}
}

func TestSkipWhiteSpace(t *testing.T) {
	type TestItem struct {
		Arg string
		Expect string 
	}
	table := map[string]TestItem {
		"white head one" : {" 12 )", "12 )"},
		"many white space" : {"  	12)", "12)"},
		"no white head" : {"1)", "1)"},
	}
	for key, item := range table {
		result := skipWhiteSpace(item.Arg)
		if result != item.Expect {
			t.Fatalf("%s: expect %s but got %s", key, item.Expect, result)
		}
	}
}

func TestLexicalAnalyze(t *testing.T) {
	type TestItem struct {
		Arg string
		IsSuccess bool  
		ExpectLen int 
		ExpectObjects []domain.LexicalObject
	}
	table := map[string]TestItem {
		"success": {
			Arg: "(gandr 1 2)",
			IsSuccess: true,
			ExpectLen: 5,
			ExpectObjects: []domain.LexicalObject{
				{Type: domain.Symbol, Label: domain.SymbolParenthesisOpen, Word: "("},
				{Type: domain.KeyWord, Label: domain.KeyWordGandr, Word: "gandr"},
				{Type: domain.NumberLiteral, Label: domain.LabelNull, Word: "1"},
				{Type:domain.NumberLiteral, Label: domain.LabelNull, Word: "2"},
				{Type: domain.Symbol, Label: domain.SymbolParenthesisClose, Word: ")"},
			},
		},
		"error": {
			Arg: "(gundr 1 2)",
			IsSuccess: false,
			ExpectLen: 0,
			ExpectObjects: nil,
		},
	}
	for key, item := range table {
		tokens, err := lexicalAnalyze(item.Arg)
		if item.IsSuccess && err != nil {
			t.Fatal(err)
		}

		if !item.IsSuccess{
			assert.NotNil(t, err, key)
		}
		assert.Len(t, tokens, item.ExpectLen, key)
		for i := range tokens {
			assert.Equal(t, item.ExpectObjects[i], tokens[i], key)
		}
	}

}