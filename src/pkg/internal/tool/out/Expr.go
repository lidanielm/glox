package out

type Expr interface {
}

type Binary struct {
	Expr left
	Token operator
	Expr right
}

func NewExpr(Expr left, Token operator, Expr right) *Binary {
	return &Binary{Expr: left, Token: operator, Expr: right}
}
