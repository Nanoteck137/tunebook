package planner

import (
	"fmt"

	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/ast"
	"github.com/nanoteck137/tunebook/tools/query/schema"
)

type Planner struct {
	schema *schema.Schema
}

func New(s *schema.Schema) *Planner {
	return &Planner{schema: s}
}

func (p *Planner) Plan(expr ast.Expr) (*query.Plan, error) {
	filter, err := p.planFilter(expr)
	if err != nil {
		return nil, err
	}

	return &query.Plan{
		Filter: filter,
	}, nil
}

func (p *Planner) planFilter(expr ast.Expr) (query.FilterNode, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		return p.planBinary(e)
	case *ast.UnaryExpr:
		return p.planUnary(e)
	case *ast.IsNullExpr:
		return p.planIsNull(e)
	case *ast.InExpr:
		return p.planIn(e)
	default:
		return nil, fmt.Errorf("unexpected expression: %T", expr)
	}
}

func (p *Planner) planBinary(e *ast.BinaryExpr) (query.FilterNode, error) {
	switch e.Op {
	case ast.OpAnd:
		left, err := p.planFilter(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := p.planFilter(e.Right)
		if err != nil {
			return nil, err
		}
		return &query.AndNode{Left: left, Right: right}, nil

	case ast.OpOr:
		left, err := p.planFilter(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := p.planFilter(e.Right)
		if err != nil {
			return nil, err
		}
		return &query.OrNode{Left: left, Right: right}, nil

	default:
		return p.planComparison(e)
	}
}

func (p *Planner) planComparison(e *ast.BinaryExpr) (query.FilterNode, error) {
	fieldRef, ok := e.Left.(*ast.FieldRef)
	if !ok {
		return nil, fmt.Errorf("left side of comparison must be a field")
	}

	field, ok := p.schema.Field(fieldRef.Name)
	if !ok {
		return nil, fmt.Errorf("unknown field: %s", fieldRef.Name)
	}

	if _, isNull := e.Right.(*ast.NullLit); isNull {
		if !field.Nullable {
			return nil, fmt.Errorf("field '%s' is not nullable", field.Name)
		}
		not := e.Op == ast.OpNeq
		return &query.IsNullNode{Field: field, Not: not}, nil
	}

	op, err := p.resolveOperator(e.Op, field.Type)
	if err != nil {
		return nil, fmt.Errorf("field '%s': %w", field.Name, err)
	}

	value, err := p.resolveValue(e.Right, field.Type)
	if err != nil {
		return nil, fmt.Errorf("field '%s': %w", field.Name, err)
	}

	switch op {
	case opContains:
		return &query.ContainsNode{Field: field, Value: value}, nil
	case opLike:
		return &query.ComparisonNode{
			Field:    field,
			Operator: query.OpLike,
			Value:    value,
		}, nil
	case opHas:
		return &query.HasNode{
			Field:    field,
			Operator: query.OpEqual,
			Value:    value,
		}, nil
	default:
		return &query.ComparisonNode{
			Field:    field,
			Operator: opToQuery(op),
			Value:    value,
		}, nil
	}
}

func (p *Planner) planUnary(e *ast.UnaryExpr) (query.FilterNode, error) {
	switch e.Op {
	case ast.OpNot:
		expr, err := p.planFilter(e.Expr)
		if err != nil {
			return nil, err
		}
		return &query.NotNode{Expr: expr}, nil
	default:
		return nil, fmt.Errorf("unexpected unary operator: %v", e.Op)
	}
}

func (p *Planner) planIsNull(e *ast.IsNullExpr) (query.FilterNode, error) {
	fieldRef, ok := e.Field.(*ast.FieldRef)
	if !ok {
		return nil, fmt.Errorf("is null requires a field")
	}

	field, ok := p.schema.Field(fieldRef.Name)
	if !ok {
		return nil, fmt.Errorf("unknown field: %s", fieldRef.Name)
	}

	if !field.Nullable {
		return nil, fmt.Errorf("field '%s' is not nullable", field.Name)
	}

	return &query.IsNullNode{Field: field, Not: e.Not}, nil
}

func (p *Planner) planIn(e *ast.InExpr) (query.FilterNode, error) {
	fieldRef, ok := e.Field.(*ast.FieldRef)
	if !ok {
		return nil, fmt.Errorf("in requires a field")
	}

	field, ok := p.schema.Field(fieldRef.Name)
	if !ok {
		return nil, fmt.Errorf("unknown field: %s", fieldRef.Name)
	}

	values := make([]query.Value, 0, len(e.Values))
	for _, v := range e.Values {
		val, err := p.resolveValue(v, field.Type)
		if err != nil {
			return nil, fmt.Errorf("field '%s': %w", field.Name, err)
		}
		values = append(values, val)
	}

	return &query.InNode{Field: field, Values: values, Not: e.Not}, nil
}

func (p *Planner) resolveValue(expr ast.Expr, fieldType query.Type) (query.Value, error) {
	switch lit := expr.(type) {
	case *ast.StringLit:
		if fieldType != query.TypeString && fieldType != query.TypeRelation {
			return query.Value{}, fmt.Errorf("expected %s, got string", typeName(fieldType))
		}
		return query.Value{Type: query.TypeString, Value: lit.Value}, nil

	case *ast.IntLit:
		if fieldType != query.TypeInt && fieldType != query.TypeFloat {
			return query.Value{}, fmt.Errorf("expected %s, got integer", typeName(fieldType))
		}
		if fieldType == query.TypeFloat {
			return query.Value{Type: query.TypeFloat, Value: float64(lit.Value)}, nil
		}
		return query.Value{Type: query.TypeInt, Value: lit.Value}, nil

	case *ast.FloatLit:
		if fieldType != query.TypeFloat {
			return query.Value{}, fmt.Errorf("expected %s, got float", typeName(fieldType))
		}
		return query.Value{Type: query.TypeFloat, Value: lit.Value}, nil

	case *ast.BoolLit:
		if fieldType != query.TypeBool {
			return query.Value{}, fmt.Errorf("expected %s, got bool", typeName(fieldType))
		}
		return query.Value{Type: query.TypeBool, Value: lit.Value}, nil

	case *ast.FieldRef:
		return query.Value{}, fmt.Errorf("field references not supported in values")

	default:
		return query.Value{}, fmt.Errorf("unexpected literal type: %T", expr)
	}
}

type internalOp int

const (
	opEq internalOp = iota
	opNeq
	opGt
	opGte
	opLt
	opLte
	opContains
	opLike
	opHas
)

func (p *Planner) resolveOperator(op ast.BinaryOp, fieldType query.Type) (internalOp, error) {
	switch op {
	case ast.OpEq:
		if fieldType == query.TypeRelation {
			return 0, fmt.Errorf("operator '=' cannot be used on relation field, use 'has' instead")
		}
		return opEq, nil
	case ast.OpNeq:
		if fieldType == query.TypeRelation {
			return 0, fmt.Errorf("operator '!=' cannot be used on relation field, use 'not has' instead")
		}
		return opNeq, nil
	case ast.OpGt:
		if !comparable(fieldType) {
			return 0, fmt.Errorf("operator '>' cannot be used on %s field", typeName(fieldType))
		}
		return opGt, nil
	case ast.OpGte:
		if !comparable(fieldType) {
			return 0, fmt.Errorf("operator '>=' cannot be used on %s field", typeName(fieldType))
		}
		return opGte, nil
	case ast.OpLt:
		if !comparable(fieldType) {
			return 0, fmt.Errorf("operator '<' cannot be used on %s field", typeName(fieldType))
		}
		return opLt, nil
	case ast.OpLte:
		if !comparable(fieldType) {
			return 0, fmt.Errorf("operator '<=' cannot be used on %s field", typeName(fieldType))
		}
		return opLte, nil
	case ast.OpContains:
		if fieldType != query.TypeString {
			return 0, fmt.Errorf("operator 'contains' cannot be used on %s field", typeName(fieldType))
		}
		return opContains, nil
	case ast.OpLike:
		if fieldType != query.TypeString {
			return 0, fmt.Errorf("operator 'like' cannot be used on %s field", typeName(fieldType))
		}
		return opLike, nil
	case ast.OpHas:
		if fieldType != query.TypeRelation {
			return 0, fmt.Errorf("operator 'has' can only be used on relation fields")
		}
		return opHas, nil
	default:
		return 0, fmt.Errorf("unexpected operator: %v", op)
	}
}

func opToQuery(op internalOp) query.Operator {
	switch op {
	case opEq:
		return query.OpEqual
	case opNeq:
		return query.OpNotEqual
	case opGt:
		return query.OpGreater
	case opGte:
		return query.OpGreaterEqual
	case opLt:
		return query.OpLess
	case opLte:
		return query.OpLessEqual
	default:
		return query.OpEqual
	}
}

func comparable(t query.Type) bool {
	switch t {
	case query.TypeInt, query.TypeFloat, query.TypeTime:
		return true
	}
	return false
}

func typeName(t query.Type) string {
	switch t {
	case query.TypeString:
		return "string"
	case query.TypeInt:
		return "integer"
	case query.TypeFloat:
		return "float"
	case query.TypeBool:
		return "boolean"
	case query.TypeTime:
		return "time"
	case query.TypeRelation:
		return "relation"
	default:
		return "unknown"
	}
}
