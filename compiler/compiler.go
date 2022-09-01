package compiler

import (
	"domain"
	"fmt"

	"github.com/anaseto/gruid"
)

const (
	amountGandr = 5
	amountSeiethr = 10 
	radiusGandr = 0
	radiusSeiethr = 3
)


func Compile(arg string) (magic domain.Magic, err error) {
	tokens, err := lexicalAnalyze(arg)
	if err != nil {
		return 
	}

	ast, err := parse(tokens)
	if err != nil {
		return
	}

	magic, err = genMagic(ast)
	return 
}

func genMagic(ast *domain.Node) (magic domain.Magic, err error) {
	if ast == nil {
		return 
	}

	switch ast.Type {
	case domain.Root:
		magic, err = genRoot(ast)
		return 
	}
	return 
}

func genRoot(ast *domain.Node) (magic domain.Magic, err error) {
	left := ast.Left
	if left == nil {
		err = fmt.Errorf("runtime error: unexpected nil left at root")
		return 
	} 

	switch left.Type{
	case domain.Operator:
		magic, err = genOperator(left.Data.Label, ast.Right)
		return 
	default:
		err = fmt.Errorf("runtime error: expected operator")
		return 
	}
}

func genOperator(label domain.DataLabel, operands *domain.Node) (magic domain.Magic, err error) {
	switch label{
	case domain.OperatorGandr:
		// check operands
		var parameters []int64
		parameters, err = getParameters(operands, 2)
		if err != nil {
			return 
		}
		magic = domain.Magic{
			Amount: amountGandr,
			Target: gruid.Point{X: int(parameters[0]), Y: int(parameters[1])},
			Radius: radiusGandr,
			Name: "gandr",
		}
		return 
	case domain.OperatorSeiethr:
		// check operands
		var parameters []int64
		parameters, err = getParameters(operands, 2)
		if err != nil {
			return 
		}
		magic = domain.Magic{
			Amount: amountSeiethr,
			Target: gruid.Point{X: int(parameters[0]), Y: int(parameters[1])},
			Radius: radiusSeiethr,
			Name: "seiethr",
		}
		return 
	default:
		err = fmt.Errorf("runtime error: expected operator")
		return 
	}
}

func getParameters(operands *domain.Node, length int) (parameters []int64, err error) {
	if length > 0 {
		if operands == nil {
			err = fmt.Errorf("runtime error: expected list of length %d but got null list", length)
			return 
		}
		if operands.Type != domain.ArrayHead {
			err = fmt.Errorf("runtime error: expected list")
			return 
		}
		operands = operands.Right

		i := 0
		for {
			if operands == nil {
				err = fmt.Errorf("internal error: unexpected nil at operands")
				return 
			}

			if operands.Type == domain.ArrayTail {
				if i != length {
					err = fmt.Errorf("runtime error: expect %d operands but got %d", length, i)
					return 
				}
				return 
			}

			if operands.Type == domain.ArrayItem {
				// error check
				if operands.Left == nil {
					err = fmt.Errorf("internal error: unexpected nil at ArrayItem.Left")
					return 
				}
				if operands.Left.Type != domain.Literal {
					err = fmt.Errorf("runtime error: unexpected type of operands")
					return 
				}

				parameters = append(parameters, operands.Left.Data.Number)
				operands = operands.Right
				i++
				continue 
			}
			err = fmt.Errorf("internal error: unexpected type %d at getOperands", operands.Type)
		}
	}
	return 
}