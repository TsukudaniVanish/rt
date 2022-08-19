// lexier ... do lexical analyze of magic spell
package compiler

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"domain"
)
func lexicalAnalyze(arg string) (tokens []domain.LexicalObject, err error) {	
	input := arg
	for len(input) > 0 {
		input = skipWhiteSpace(input)
		w, rest := consumeWord(input)
		lo, err := getLexicalObjectFromWord(w)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, lo)
		input = rest
	}
	return
}


// isWhiteSpace ... if head of rune is whitespace then returns true
func isWhiteSpace(r rune) bool {
	switch r {
	case '\t', ' ', '\v', '\n':
		return true 
	default:
		return false
	}
}

func isPunctuator(r rune) bool {
	switch r {
	case '(', ')':
		return true
	default:
		return false
	}
}

// skipWhiteSpace .. skip whiteSpace 
func skipWhiteSpace(s string) string{
	pos := 0
	for len(s) > pos {
		r, size := utf8.DecodeRuneInString(s[pos:])
		if !isWhiteSpace(r) {
			return s[pos:]
		}
		pos += size
	}
	return s 
}

// consumeWord ... split a word from string. this function expects string not begin with white space
func consumeWord (s string) (word, rest string) {
	// single character
	r, size := utf8.DecodeRuneInString(s)
	switch  r {
	case '(':
		word = "("
		rest = s[size:]
		return 
	case ')':
		word = ")"
		rest = s[size:]
		return 
	}
  
	// word 
	pos := 0
	for len(s) > pos{
		r, size := utf8.DecodeRuneInString(s[pos:])
		if isWhiteSpace(r) || isPunctuator(r) {
			break
		}
		pos += size
	}
	word = s[0:pos]
	rest = s[pos:]
	return 
}

func getLexicalObjectFromWord(w string) (lo domain.LexicalObject, err error) {
	lo.Label = domain.LabelNull
	lo.Word = w 
	// check keyword 
	switch w {
	case "gandr":
		lo.Type = domain.KeyWord
		lo.Label = domain.KeyWordGandr
		return 
	case "seiethr":
		lo.Type = domain.KeyWord
		lo.Label = domain.KeyWordSeiethr
		return 
	}

	// check symbol 
	switch w {
	case "(":
		lo.Type = domain.Symbol
		lo.Label = domain.SymbolParenthesisOpen
		return 
	case ")":
		lo.Type = domain.Symbol
		lo.Label = domain.SymbolParenthesisClose
		return 
	}

	// check number literal
	_, err = strconv.Atoi(w)
	if err != nil {
		err = fmt.Errorf("%s this is not word", w)
	}
	lo.Type = domain.NumberLiteral
	return 
}