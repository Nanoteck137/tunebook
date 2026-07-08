package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Expr interface {
	exprNode()
	String() string
}

type StringLit struct {
	Value string
}

type IntLit struct {
	Value int64
}

type FloatLit struct {
	Value float64
}

type BoolLit struct {
	Value bool
}

type NullLit struct{}

type FieldRef struct {
	Name string
}

type BinaryOp int

const (
	OpAnd BinaryOp = iota
	OpOr
	OpEq
	OpNeq
	OpGt
	OpGte
	OpLt
	OpLte
	OpContains
	OpLike
	OpHas
)

func (op BinaryOp) String() string {
	switch op {
	case OpAnd:
		return "and"
	case OpOr:
		return "or"
	case OpEq:
		return "="
	case OpNeq:
		return "!="
	case OpGt:
		return ">"
	case OpGte:
		return ">="
	case OpLt:
		return "<"
	case OpLte:
		return "<="
	case OpContains:
		return "contains"
	case OpLike:
		return "like"
	case OpHas:
		return "has"
	}
	return "?"
}

type BinaryExpr struct {
	Left  Expr
	Op    BinaryOp
	Right Expr
}

type UnaryOp int

const (
	OpNot UnaryOp = iota
)

func (op UnaryOp) String() string {
	switch op {
	case OpNot:
		return "not"
	}
	return "?"
}

type UnaryExpr struct {
	Op   UnaryOp
	Expr Expr
}

type InExpr struct {
	Field  Expr
	Values []Expr
	Not    bool
}

type IsNullExpr struct {
	Field Expr
	Not   bool
}

func (*StringLit) exprNode()  {}
func (*IntLit) exprNode()     {}
func (*FloatLit) exprNode()   {}
func (*BoolLit) exprNode()    {}
func (*NullLit) exprNode()    {}
func (*FieldRef) exprNode()   {}
func (*BinaryExpr) exprNode() {}
func (*UnaryExpr) exprNode()  {}
func (*InExpr) exprNode()     {}
func (*IsNullExpr) exprNode() {}

func (l *StringLit) String() string {
	return strconv.Quote(l.Value)
}

func (l *IntLit) String() string {
	return strconv.FormatInt(l.Value, 10)
}

func (l *FloatLit) String() string {
	return strconv.FormatFloat(l.Value, 'f', -1, 64)
}

func (l *BoolLit) String() string {
	if l.Value {
		return "true"
	}
	return "false"
}

func (l *NullLit) String() string {
	return "null"
}

func (f *FieldRef) String() string {
	return f.Name
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Op, e.Right)
}

func (e *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", e.Op, e.Expr)
}

func (e *InExpr) String() string {
	var values []string
	for _, v := range e.Values {
		values = append(values, v.String())
	}
	not := ""
	if e.Not {
		not = "not "
	}
	return fmt.Sprintf("(%s %sin (%s))", e.Field, not, strings.Join(values, ", "))
}

func (e *IsNullExpr) String() string {
	if e.Not {
		return fmt.Sprintf("(%s is not null)", e.Field)
	}
	return fmt.Sprintf("(%s is null)", e.Field)
}
