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
	Root NodeType = iota
	Operator 
	Literal
	ArrayHead
	ArrayItem
	ArrayTail
)

type DataLabel int 
const (
	Null DataLabel = iota 
	Number 
	String 
	OperatorGandr
	OperatorSeiethr
)

type NodeData struct {
	Label DataLabel 
	// Text ... this field have meaning when label is string or Null
	Text string
	// Number ... this field have meaning when label is number 
	Number int64
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