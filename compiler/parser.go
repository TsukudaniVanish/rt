package compiler

import (
	"domain"
	"fmt"
	"strconv"
)

func parse(tokens []domain.LexicalObject) (root *domain.Node, err error) {
	root = &domain.Node{
		Type: domain.Root,
	}
	err = pair(tokens, root)
	return 
}

// pair ... parse pair
func pair(tokens []domain.LexicalObject, ast *domain.Node) (err error) {
	if len(tokens) <= 0 {
		return 
	}
	tokens, ok := expect(tokens, "(")
	if !ok {
		err = fmt.Errorf("expect ( but got %s", tokens[0].Word)
		return 
	}

	for {
		if len(tokens) <= 0 {
			return 
		}
		var head domain.LexicalObject
		tokens, head = consume(tokens)
		if isOperator(head) {
			tokens, err = operator(head, tokens, ast)
			if err != nil {
				return 
			}
			var ok bool 
			tokens, ok = expect(tokens, ")")
			if !ok {
				err = fmt.Errorf("expect ) but got %s", tokens[0].Word)
				return
			}
			continue 
		}

		if head.Word == ")" {
			err = pair(tokens, ast)
			return 
		}

		err = fmt.Errorf("%s is not expected", head.Word)
		return
	}
}

// operator ... parse operator witch is first element of array and operands of it. add then to ast
func operator(operatorToken domain.LexicalObject, tokens []domain.LexicalObject, ast *domain.Node) (res []domain.LexicalObject, err error) {
	leafOperator := &domain.Node{
		Type: domain.Operator,
		Data: domain.NodeData{},
	}
	switch operatorToken.Label {
	case domain.KeyWordGandr:
		leafOperator.Data.Label = domain.OperatorGandr
	case domain.KeyWordSeiethr:
		leafOperator.Data.Label = domain.OperatorSeiethr
	}
	
	// parse operands
	operands, res , err := list(tokens, true)
	if err != nil {
		return 
	}

	ast.Left = leafOperator
	ast.Right = operands 
	return 
}

// list ... pase operand of operator 
func list(tokens []domain.LexicalObject, isOperand bool) (ast *domain.Node, res []domain.LexicalObject, err error) {
	if !isOperand {
		tokens, ok := expect(tokens, "(")
		if !ok {
			err = fmt.Errorf("expect ( but got %s", tokens[0].Word)
			return 
		}
	}

	ast = &domain.Node{
		Type: domain.ArrayHead,
		Right: &domain.Node{},
	}
	res, err = listElem(tokens, ast.Right)
	return 
}

func listElem(tokens []domain.LexicalObject, ast *domain.Node) (res[]domain.LexicalObject, err error) {
	if len(tokens) <= 0 {
		err = fmt.Errorf("unexpected eof at list")
		return 
	}


	ast.Type = domain.ArrayItem
	if tokens[0].Word == ")" {
		ast.Type = domain.ArrayTail
		return 
	}

	ast.Right = &domain.Node{}	
	left, tokens, err := atom(tokens)
	if err != nil {
		return 
	}
	ast.Left = left
	
	res, err = listElem(tokens, ast.Right)
	return 
}

// atom ... parse atom, primary lexeme of this language. return leaf node which represents the atom 
func atom(tokens []domain.LexicalObject) (leaf *domain.Node, res []domain.LexicalObject, err error) {
	res, head := consume(tokens)
	leaf = &domain.Node{
		Type: domain.Literal,
		Data: domain.NodeData{},
	}
	switch head.Type {
	case domain.NumberLiteral:
		leaf.Data.Label = domain.Number

		var number int64
		number, err = strconv.ParseInt(head.Word, 10, 64)
		if err != nil {
			err = fmt.Errorf("INTERNAL ERROR: FAILED TO PARSE NUMBER")
			return 
		}
		leaf.Data.Number = number
	default:
		err = fmt.Errorf("%s is not atom of expression", head.Word)
		return 
	}
	return 
}

// expect ... check head token is an expect object if it is ok then consume tokens
func expect(tokens []domain.LexicalObject, expectedWord string ) (res []domain.LexicalObject, ok bool) {
	res = tokens 
	if len(tokens) <= 0 {
		ok = true 
		return 
	}
	head := tokens[0]
	if head.Word == expectedWord {
		res = tokens[1:]
		ok = true 
		return 
	}
	ok = false 
	return 
}

// consume ... pop a head of tokens
func consume(tokens []domain.LexicalObject) (res []domain.LexicalObject, head domain.LexicalObject) {
	if len(tokens) <= 0 {
		return 
	}

	res = tokens[1:]
	head = tokens[0]
	return 
}

// isOperator ... check token is operator
func isOperator(token domain.LexicalObject) bool {
	switch token.Label {
	case domain.KeyWordGandr:
		return true 
	case domain.KeyWordSeiethr:
		return true
	default:
		return false 
	}
}