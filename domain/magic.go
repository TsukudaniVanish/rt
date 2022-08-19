package domain

import "github.com/anaseto/gruid"


type LexicalObjectType int 
const (
	NotLexicalObject LexicalObjectType = iota
	Symbol
	KeyWord
	NumberLiteral
)

type LexicalObjectLabel int 
const (
	LabelNull LexicalObjectLabel = iota 
	SymbolParenthesisOpen
	SymbolParenthesisClose
	KeyWordGandr
	KeyWordSeiethr	
)

type LexicalObject struct {
	Type LexicalObjectType
	// Label ... this field has meaning when type is symbol or keyword
	Label LexicalObjectLabel
	//Unit of this object
	Word string 
}

type NodeType int 
const (
	Atom NodeType = iota 
	Pair
)

type DataLabel int 
const (
	Null DataLabel = iota 
	Number 
	String 
	List 
	OperatorGandr
	OperatorSeiethr
)

type NodeData struct {
	Label DataLabel 
	Number int 
	Text string
	// this field have meaning when Label is List  
	Length int 
}

type Node struct {
	Type NodeType
	Data NodeData 
	// Left and Right have meaning when type is Pair
	Left *Node 
	Right *Node
}

type Magic struct {
	// actor ... caster of this magic 
	Actor int 
	// amount of mana
	Amount int	
	// Target ... target of magic
	Target gruid.Point
	// Radius ... Radius of magic 
	Radius int
	// Name ... Name of magic <-- concatenating atoms
	Name string
}

var MagicArrow = Magic {
	Amount: 4,
	Radius: 0,
	Name: "gandr",
}