package sql

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/tunebook/tools/query"
)

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

type CompileResult struct {
	Where exp.Expression
	Order []exp.OrderedExpression
}

func (c *Compiler) Compile(plan *query.Plan) (*CompileResult, error) {
	res := &CompileResult{}

	if plan.Filter != nil {
		where, err := c.compileFilter(plan.Filter)
		if err != nil {
			return nil, err
		}

		res.Where = where
	}

	if len(plan.OrderBy) > 0 {
		order, err := c.compileOrdering(plan.OrderBy)
		if err != nil {
			return nil, err
		}
		res.Order = order
	}

	return res, nil
}

func (c *Compiler) compileFilter(node query.FilterNode) (goqu.Expression, error) {
	switch n := node.(type) {
	case *query.AndNode:
		return c.compileAnd(n)
	case *query.OrNode:
		return c.compileOr(n)
	case *query.NotNode:
		return c.compileNot(n)
	case *query.ComparisonNode:
		return c.compileComparison(n)
	case *query.IsNullNode:
		return c.compileIsNull(n)
	case *query.ContainsNode:
		return c.compileContains(n)
	case *query.InNode:
		return c.compileIn(n)
	case *query.HasNode:
		return c.compileHas(n)
	}
	return nil, fmt.Errorf("unknown filter node: %T", node)
}

func (c *Compiler) compileAnd(n *query.AndNode) (goqu.Expression, error) {
	left, err := c.compileFilter(n.Left)
	if err != nil {
		return nil, err
	}

	right, err := c.compileFilter(n.Right)
	if err != nil {
		return nil, err
	}

	return goqu.And(left, right), nil
}

func (c *Compiler) compileOr(n *query.OrNode) (goqu.Expression, error) {
	left, err := c.compileFilter(n.Left)
	if err != nil {
		return nil, err
	}

	right, err := c.compileFilter(n.Right)
	if err != nil {
		return nil, err
	}

	return goqu.Or(left, right), nil
}

func (c *Compiler) compileNot(n *query.NotNode) (goqu.Expression, error) {
	expr, err := c.compileFilter(n.Expr)
	if err != nil {
		return nil, err
	}

	return goqu.L("NOT (?)", expr), nil
}

func (c *Compiler) compileComparison(n *query.ComparisonNode) (goqu.Expression, error) {
	col := c.getFieldIdentifier(n.Field)

	switch n.Operator {
	case query.OpEqual:
		return col.Eq(n.Value.Value), nil
	case query.OpNotEqual:
		return col.Neq(n.Value.Value), nil
	case query.OpGreater:
		return col.Gt(n.Value.Value), nil
	case query.OpGreaterEqual:
		return col.Gte(n.Value.Value), nil
	case query.OpLess:
		return col.Lt(n.Value.Value), nil
	case query.OpLessEqual:
		return col.Lte(n.Value.Value), nil
	case query.OpLike:
		return col.Like(n.Value.Value), nil
	}

	return nil, fmt.Errorf("unsupported operator: %v", n.Operator)
}

func (c *Compiler) compileIsNull(n *query.IsNullNode) (goqu.Expression, error) {
	col := c.getFieldIdentifier(n.Field)

	if n.Not {
		return col.IsNotNull(), nil
	}
	return col.IsNull(), nil
}

func (c *Compiler) compileContains(n *query.ContainsNode) (goqu.Expression, error) {
	col := c.getFieldIdentifier(n.Field)
	pattern := "%" + fmt.Sprint(n.Value.Value) + "%"
	return col.Like(pattern), nil
}

func (c *Compiler) compileIn(n *query.InNode) (goqu.Expression, error) {
	col := c.getFieldIdentifier(n.Field)

	if len(n.Values) == 0 {
		if n.Not {
			return goqu.L("1=1"), nil
		}
		return goqu.L("1=0"), nil
	}

	values := make([]interface{}, len(n.Values))
	for i, v := range n.Values {
		values[i] = v.Value
	}

	if n.Not {
		return col.NotIn(values...), nil
	}
	return col.In(values...), nil
}

func (c *Compiler) compileHas(n *query.HasNode) (goqu.Expression, error) {
	if n.Field.Relation == nil {
		return nil, fmt.Errorf("field '%s' is not a relation field", n.Field.Name)
	}

	rel := n.Field.Relation

	// Build the EXISTS subquery
	// SELECT 1 FROM join_table WHERE join_table.join_foreign_key = main_table.id AND join_table.join_reference IN (value)
	joinTable := goqu.From(rel.JoinTable)

	// Create the join table condition - use IN with the value directly
	joinCol := goqu.I(rel.JoinTable + "." + rel.JoinReference)
	existsCondition := joinCol.In(n.Value.Value)

	// Create the EXISTS subquery
	existsSubquery := joinTable.Select(goqu.L("1")).Where(
		goqu.And(
			goqu.I(rel.JoinTable+"."+rel.JoinForeignKey).Eq(goqu.I("tracks.id")),
			existsCondition,
		),
	)

	// Use goqu's EXISTS support
	existsExpr := goqu.L("EXISTS ?", existsSubquery)

	if n.Not {
		return goqu.L("NOT ?", existsExpr), nil
	}

	return existsExpr, nil
}

func (c *Compiler) compileOrdering(orderings []query.Ordering) ([]exp.OrderedExpression, error) {
	result := make([]exp.OrderedExpression, 0, len(orderings))

	for _, o := range orderings {
		switch o := o.(type) {
		case *query.FieldOrdering:
			col := c.getFieldIdentifier(o.Field)
			if o.Dir == query.DirDesc {
				result = append(result, col.Desc())
			} else {
				result = append(result, col.Asc())
			}
		case *query.RandomOrdering:
			result = append(result, goqu.Func("RANDOM").Asc())
		case *query.ShuffleOrdering:
			return nil, fmt.Errorf("ShuffleOrdering not yet supported")
		case *query.ScoreOrdering:
			return nil, fmt.Errorf("ScoreOrdering not yet supported")
		default:
			return nil, fmt.Errorf("unknown ordering type: %T", o)
		}
	}

	return result, nil
}

func (c *Compiler) getFieldIdentifier(f *query.Field) exp.IdentifierExpression {
	if col, ok := f.Meta["column"]; ok {
		return goqu.I(fmt.Sprint(col))
	}
	return goqu.I(f.Name)
}
