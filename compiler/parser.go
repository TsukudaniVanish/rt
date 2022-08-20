package compiler

import "domain"

func parse(tokens []domain.LexicalObject) (root domain.Node) {

}

// pair .., parse pair 
func pair(tokens []domain.LexicalObject, prev domain.Node) (ast domain.Node) {

}

// operator ... parse operator witch is first element of pair 
func operator(tokens []domain.LexicalObject, prev domain.Node) (ast domain.Node) {

}

// list ... pase operand of operator 
func list(tokens []domain.LexicalObject, prev domain.Node) (ast domain.Node) {

}

// atom ... parse atom, primary lexeme of this language
func atom(tokens []domain.LexicalObject, prev domain.Node) (ast domain.Node) {

}