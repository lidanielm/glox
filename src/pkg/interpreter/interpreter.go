package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/internal/stmt"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Interpreter struct {
	env *Env
	globals *Env
	locals map[ast.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnv()
	globals.Define("clock", &ClockFn{})
	env := globals
	locals := make(map[ast.Expr]int)
	return &Interpreter{env: env, globals: globals, locals: locals}
}

func (ip *Interpreter) Interpret(stmts []stmt.Stmt) error {
    for _, stmt := range stmts {
		err := ip.execute(stmt)
		if err != nil {
			runtimeError(err)
			return err
		}
	}

	return nil
}

/** VISIT METHODS */
func (ip *Interpreter) VisitLiteralExpr(literal ast.Literal) (any, error) {
	return literal.Value, nil
}

func (ip *Interpreter) VisitGroupingExpr(grouping ast.Grouping) (any, error) {
	return ip.evaluate(grouping.Expression)
}

func (ip *Interpreter) VisitUnaryExpr(unary ast.Unary) (any, error) {
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
		return !isTruthy(right), nil
    default:
        return nil, lox_error.NewRuntimeError(unary.Operator, "Invalid operator.")
	}
}

func (ip *Interpreter) VisitBinaryExpr(binary ast.Binary) (any, error) {
    left, lerr := ip.evaluate(binary.Left)
    right, rerr := ip.evaluate(binary.Right)

	// Evaluate both subexpressions first but report the first error
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

		if right.(float64) == 0 {
			return nil, lox_error.NewRuntimeError(binary.Operator, "Invalid divison by zero.")
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
        return !isEqual(left, right), nil
    case token.EQUAL_EQUAL:
        return !isEqual(left, right), nil
    default:
        return nil, lox_error.NewRuntimeError(binary.Operator, "Invalid operator.")
    }
}

func (ip *Interpreter) VisitTernaryExpr(ternary ast.Ternary) (any, error) {
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


func (ip *Interpreter) VisitVariableExpr(expr ast.Variable) (any, error) {
	return ip.lookUpVariable(expr.Name, expr)
}


func (ip *Interpreter) VisitAssignExpr(expr ast.Assign) (any, error) {
	value, err := ip.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := ip.locals[expr]
	if ok {
		err = ip.env.AssignAt(distance, expr.Name, value)
		if err != nil {
			return nil, err
		}
	} else {
		err = ip.globals.Assign(expr.Name, value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}


func (ip *Interpreter) VisitLogicalExpr(expr ast.Logical) (any, error) {
	leftVal, err := ip.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.OR:
		if isTruthy(leftVal) {
			return leftVal, nil
		}
	case token.AND:
		if !isTruthy(leftVal) {
			return leftVal, nil
		}
	}

	return ip.evaluate(expr.Right)
}


func (ip *Interpreter) VisitCallExpr(expr ast.Call) (any, error) {
	callee, err := ip.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}
	for _, argument := range expr.Arguments {
		value, err := ip.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, value)
	}

	callableFn, ok := callee.(Callable)
	if !ok {
		return nil, lox_error.NewRuntimeError(expr.Paren, "Can only call functions and classes.")
	}

	if len(arguments) != callableFn.Arity() {
		return nil, lox_error.NewRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", callableFn.Arity(), len(arguments)))
	}

	return callableFn.Call(ip, arguments)

}

func (ip *Interpreter) VisitGetExpr(expr ast.Get) (any, error) {
	object, err := ip.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if object, ok := object.(Instance); ok {
		return object.Get(expr.Name.Lexeme)
	}

	return nil, lox_error.NewRuntimeError(expr.Name, "Only instances have fields.")
}

func (ip *Interpreter) VisitSetExpr(expr ast.Set) (any, error) {
	object, err := ip.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if object, ok := object.(Instance); ok {
		value, err := ip.evaluate(expr.Value)
		if err != nil {
			return nil, err
		}

		object.Set(expr.Name.Lexeme, value)
	}

	return nil, lox_error.NewRuntimeError(expr.Name, "Only instances have properties.")
}

func (ip *Interpreter) VisitThisExpr(expr ast.This) (any, error) {
	return ip.lookUpVariable(expr.Keyword, expr)
}

func (ip *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(ip)
}


/** STATEMENT METHODS */
func (ip *Interpreter) VisitExpressionStmt(stmt stmt.Expression) error {
	_, err := ip.evaluate(stmt.Expr)
	if err != nil {
		return err
	}

	return nil
}


func (ip *Interpreter) VisitPrintStmt(stmt stmt.Print) error {
	val, err := ip.evaluate(stmt.Expr)
	if err != nil {
		return err
	}

	fmt.Println(stringify(val))
	return nil
}


func (ip *Interpreter) VisitVarStmt(stmt stmt.Var) error {
	if stmt.Initializer != nil {
		value, err := ip.evaluate(stmt.Initializer)
		if err != nil {
			return err
		}
		ip.env.Define(stmt.Name.Lexeme, value)
	} else {
		ip.env.Define(stmt.Name.Lexeme, nil)
	}

	return nil
}


func (ip *Interpreter) VisitBlockStmt(stmt stmt.Block) error {
	newEnv := NewEnv().WithParent(ip.env)
	return ip.executeBlock(stmt.Statements, newEnv)
}

// Wrapper for Go conditional control flow
func (ip *Interpreter) VisitIfStmt(stmt stmt.If) error {
	truthy, err := ip.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	if isTruthy(truthy) {
		return ip.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return ip.execute(stmt.ElseBranch)
	}

	return nil
}

func (ip *Interpreter) VisitWhileStmt(stmt stmt.While) error {
	for {
		eval, err := ip.evaluate(stmt.Condition)
		if err != nil {
			return err
		}
		if !isTruthy(eval) {
			break
		}

		var continueLoop bool
		err = ip.execute(stmt.Body)
		if err != nil {
			if _, ok := err.(lox_error.BreakError); ok {
				return nil
			} else if _, ok := err.(lox_error.ContinueError); ok {
				continueLoop = true
			} else {
				return err
			}
		}

		if stmt.Increment != nil {
			_, incErr := ip.evaluate(stmt.Increment)
			if incErr != nil {
				return incErr
			}
		}

		if continueLoop {
			continue
		}
	}
	return nil
}

func (ip *Interpreter) VisitFunctionStmt(stmt stmt.Function) error {
	function := NewFunction(stmt, ip.env, false)
	ip.env.Define(stmt.Name.Lexeme, function)
	return nil
}

func (ip *Interpreter) VisitBreakStmt(stmt stmt.Break) error {
	return lox_error.BreakError{} // Custom error to signal breaking
}

func (ip *Interpreter) VisitContinueStmt(stmt stmt.Continue) error {
	return lox_error.ContinueError{}
}

func (ip *Interpreter) VisitReturnStmt(stmt stmt.Return) error {
	var value any
	if stmt.Value != nil {
		var err error
		value, err = ip.evaluate(stmt.Value)
		if err != nil {
			return err
		}
	}

	return lox_error.ReturnError{Value: value}
}

func (ip *Interpreter) VisitClassStmt(stmt stmt.Class) error {
	ip.env.Define(stmt.Name.Lexeme, nil)

	// Bind methods to class
	methods := make(map[string]*Function)
	for _, method := range stmt.Methods {
		isInitializer := method.Name.Lexeme == "init"
		fn := NewFunction(method, ip.env, isInitializer)
		methods[method.Name.Lexeme] = fn
	}

	class := NewClass(stmt.Name.Lexeme, methods)
	ip.env.Assign(stmt.Name, class)
	return nil
}


func (ip *Interpreter) execute(stmt stmt.Stmt) error {
	return stmt.Accept(ip)
}

func (ip *Interpreter) resolve(expr ast.Expr, depth int) {
	ip.locals[expr] = depth
}

func (ip *Interpreter) executeBlock(statements []stmt.Stmt, env *Env) error {
	previous := ip.env
	ip.env = env
	for _, stmt := range statements {
		err := ip.execute(stmt)
		if err != nil {
			ip.env = previous
			return err
		}
	}

	ip.env = previous
	return nil
}

func (ip *Interpreter) lookUpVariable(name token.Token, expr ast.Expr) (any, error) {
	distance, ok := ip.locals[expr]
	if ok {
		return ip.env.GetAt(distance, name)
	} else {
		return ip.globals.Get(name)
	}
}

/** HELPER METHODS */
func isTruthy(expr any) bool {
	if expr == nil {
		return false
	}

	if val, ok := expr.(bool); ok {
		return val
	}
	return true
}

func isEqual(left any, right any) bool {
    if left == nil && right == nil {
        return true
    }

    if left == nil {
        return false
    }

    return left == right
}

func runtimeError(err error) {
	fmt.Println(err.Error())
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