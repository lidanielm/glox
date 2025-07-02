package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Interpreter struct {}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (ip Interpreter) Interpret(expr ast.Expr) error {
    value, err := ip.evaluate(expr)
    if err != nil {
        return err
    }
    fmt.Println(stringify(value))
    return nil
}

/** VISIT METHODS */
func (ip Interpreter) VisitLiteralExpr(literal ast.Literal) (any, error) {
	return literal.Value, nil
}

func (ip Interpreter) VisitGroupingExpr(grouping ast.Grouping) (any, error) {
	return ip.evaluate(grouping.Expression)
}

func (ip Interpreter) VisitUnaryExpr(unary ast.Unary) (any, error) {
	right, err := ip.evaluate(unary.Right)
    if err != nil {
        return nil, err
    }

	switch unary.Operator.Type {
	case token.MINUS:
        if !isNumber(right) {
            return nil, lox_error.NewRuntimeError(unary.Operator, "Operand must be a number.")
        }
		return -right.(float64), nil
	case token.BANG:
        if !isBool(right) {
            return nil, lox_error.NewRuntimeError(unary.Operator, "Operand must be a boolean.")
        }
		return !ip.isTruthy(right), nil
    default:
        return nil, lox_error.NewRuntimeError(unary.Operator, "Invalid operator.")
	}
}

func (ip Interpreter) VisitBinaryExpr(binary ast.Binary) (any, error) {
    left, lerr := ip.evaluate(binary.Left)
    right, rerr := ip.evaluate(binary.Right)

    if lerr != nil || rerr != nil {
        return nil, lerr
    }

    switch binary.Operator.Type {
    case token.MINUS:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) - right.(float64), nil
    case token.STAR:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) * right.(float64), nil
    case token.SLASH:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) / right.(float64), nil
    case token.PLUS:
        if isNumber(left, right) {
            return left.(float64) + right.(float64), nil
        }

        if isString(left, right) {
            return left.(string) + right.(string), nil
        }

        return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be two numbers or two strings.")
    case token.GREATER:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) > right.(float64), nil
    case token.GREATER_EQUAL:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) >= right.(float64), nil
    case token.LESS:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) < right.(float64), nil
    case token.LESS_EQUAL:
        if !isNumber(left, right) {
            return nil, lox_error.NewRuntimeError(binary.Operator, "Operands must be numbers.")
        }
        return left.(float64) <= right.(float64), nil
    case token.BANG_EQUAL:
        return !ip.isEqual(left, right), nil
    case token.EQUAL_EQUAL:
        return !ip.isEqual(left, right), nil
    default:
        return nil, lox_error.NewRuntimeError(binary.Operator, "Invalid operator.")
    }
}

func (ip Interpreter) VisitTernaryExpr(ternary ast.Ternary) (any, error) {
    condition, cerr := ip.evaluate(ternary.Condition)
    left, lerr := ip.evaluate(ternary.Left)
    right, rerr := ip.evaluate(ternary.Right)

    if cerr != nil || lerr != nil || rerr != nil {
        return nil, cerr
    }

    switch (ternary.Operator1.Type) {
    case token.INTERRO:
        switch (ternary.Operator2.Type) {
        case token.COLON:
            if condition.(bool) {
                return left, nil
            } else {
                return right, nil
            }
        default:
            return nil, lox_error.NewRuntimeError(ternary.Operator2, "Invalid operator.")
        }
    default:
        return nil, lox_error.NewRuntimeError(ternary.Operator1, "Invalid operator.")
    }
}


func (ip Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(ip)
}

/** HELPER METHODS */
func (ip Interpreter) isTruthy(expr any) bool {
	if expr == nil {
		return false
	}

	if val, ok := expr.(bool); ok {
		return val
	}
	return true
}

func (ip Interpreter) isEqual(left any, right any) bool {
    if left == nil && right == nil {
        return true
    }

    if left == nil {
        return false
    }

    return left == right
}

func isNumber(vs ...any) bool {
    for _, v := range vs {
        switch reflect.TypeOf(v).Kind() {
            case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
                 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
                 reflect.Float32, reflect.Float64:
                continue
            default:
                return false
            }
    }
    
    return true
}

func isBool(vs ...any) bool {
    for _, v := range vs {
        switch reflect.TypeOf(v).Kind() {
        case reflect.Bool:
            continue
        default:
            return false
        }
    }

    return true
}

func isString(vs ...any) bool {
    for _, v := range vs {
        switch reflect.TypeOf(v).Kind() {
        case reflect.String:
            continue
        default:
            return false
        }
    }

    return true
}

func stringify(value any) string {
    if value == nil {
        return "nil"
    }

    return fmt.Sprintf("%v", value)
}