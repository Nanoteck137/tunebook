package sql

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/tunebook/tools/query"
)

func newTestCompiler() *Compiler {
	return NewCompiler()
}

func strField(name string) *query.Field {
	return &query.Field{Name: name, Type: query.TypeString}
}

func intField(name string) *query.Field {
	return &query.Field{Name: name, Type: query.TypeInt}
}

func strVal(s string) query.Value {
	return query.Value{Type: query.TypeString, Value: s}
}

func intVal(i int) query.Value {
	return query.Value{Type: query.TypeInt, Value: i}
}

func toSQL(expr exp.Expression) (string, []interface{}, error) {
	if expr == nil {
		return "", nil, nil
	}
	// Use a dummy select to convert expression to SQL
	sql, args, err := goqu.From("dummy").Where(expr).ToSQL()
	if err != nil {
		return "", nil, err
	}
	// Extract just the WHERE clause
	// The SQL will be like: SELECT * FROM "dummy" WHERE (...)
	// We need to extract the part after WHERE
	whereIdx := len("SELECT * FROM \"dummy\" WHERE ")
	if len(sql) > whereIdx {
		return sql[whereIdx:], args, nil
	}
	return sql, args, nil
}

func orderToSQL(ord exp.OrderedExpression) (string, error) {
	// Use a dummy select to convert ordered expression to SQL
	sql, _, err := goqu.From("dummy").Order(ord).ToSQL()
	if err != nil {
		return "", err
	}
	// Extract just the ORDER BY clause
	// The SQL will be like: SELECT * FROM "dummy" ORDER BY ...
	orderIdx := len("SELECT * FROM \"dummy\" ORDER BY ")
	if len(sql) > orderIdx {
		return sql[orderIdx:], nil
	}
	return sql, nil
}

func TestCompileComparison(t *testing.T) {
	c := newTestCompiler()

	tests := []struct {
		name     string
		plan     *query.Plan
		wantSQL  string
		wantArgs []any
	}{
		{
			name: "equal string",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    strField("name"),
					Operator: query.OpEqual,
					Value:    strVal("rock"),
				},
			},
			wantSQL:  `("name" = 'rock')`,
			wantArgs: nil,
		},
		{
			name: "not equal",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    strField("name"),
					Operator: query.OpNotEqual,
					Value:    strVal("rock"),
				},
			},
			wantSQL:  `("name" != 'rock')`,
			wantArgs: nil,
		},
		{
			name: "greater equal int",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    intField("year"),
					Operator: query.OpGreaterEqual,
					Value:    intVal(1970),
				},
			},
			wantSQL:  `("year" >= 1970)`,
			wantArgs: nil,
		},
		{
			name: "less than",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    intField("year"),
					Operator: query.OpLess,
					Value:    intVal(2000),
				},
			},
			wantSQL:  `("year" < 2000)`,
			wantArgs: nil,
		},
		{
			name: "like",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    strField("name"),
					Operator: query.OpLike,
					Value:    strVal("rock%"),
				},
			},
			wantSQL:  `("name" LIKE 'rock%')`,
			wantArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := c.Compile(tt.plan)
			if err != nil {
				t.Fatal(err)
			}

			sql, args, err := toSQL(res.Where)
			if err != nil {
				t.Fatal(err)
			}

			if sql != tt.wantSQL {
				t.Errorf("SQL: got %q, want %q", sql, tt.wantSQL)
			}
			if len(args) != len(tt.wantArgs) {
				t.Errorf("Args len: got %d, want %d", len(args), len(tt.wantArgs))
			}
		})
	}
}

func TestCompileIsNull(t *testing.T) {
	c := newTestCompiler()

	t.Run("is null", func(t *testing.T) {
		plan := &query.Plan{
			Filter: &query.IsNullNode{
				Field: strField("deleted_at"),
				Not:   false,
			},
		}
		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := toSQL(res.Where)
		if err != nil {
			t.Fatal(err)
		}

		want := `("deleted_at" IS NULL)`
		if sql != want {
			t.Errorf("got %q, want %q", sql, want)
		}
	})

	t.Run("is not null", func(t *testing.T) {
		plan := &query.Plan{
			Filter: &query.IsNullNode{
				Field: strField("deleted_at"),
				Not:   true,
			},
		}
		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := toSQL(res.Where)
		if err != nil {
			t.Fatal(err)
		}

		want := `("deleted_at" IS NOT NULL)`
		if sql != want {
			t.Errorf("got %q, want %q", sql, want)
		}
	})
}

func TestCompileAnd(t *testing.T) {
	c := newTestCompiler()

	plan := &query.Plan{
		Filter: &query.AndNode{
			Left: &query.ComparisonNode{
				Field:    strField("name"),
				Operator: query.OpEqual,
				Value:    strVal("rock"),
			},
			Right: &query.ComparisonNode{
				Field:    intField("year"),
				Operator: query.OpGreaterEqual,
				Value:    intVal(1970),
			},
		},
	}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	sql, _, err := toSQL(res.Where)
	if err != nil {
		t.Fatal(err)
	}

	wantSQL := `(("name" = 'rock') AND ("year" >= 1970))`
	if sql != wantSQL {
		t.Errorf("SQL: got %q, want %q", sql, wantSQL)
	}
}

func TestCompileOr(t *testing.T) {
	c := newTestCompiler()

	plan := &query.Plan{
		Filter: &query.OrNode{
			Left: &query.ComparisonNode{
				Field:    strField("genre"),
				Operator: query.OpEqual,
				Value:    strVal("rock"),
			},
			Right: &query.ComparisonNode{
				Field:    strField("genre"),
				Operator: query.OpEqual,
				Value:    strVal("pop"),
			},
		},
	}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	sql, _, err := toSQL(res.Where)
	if err != nil {
		t.Fatal(err)
	}

	wantSQL := `(("genre" = 'rock') OR ("genre" = 'pop'))`
	if sql != wantSQL {
		t.Errorf("SQL: got %q, want %q", sql, wantSQL)
	}
}

func TestCompileNot(t *testing.T) {
	c := newTestCompiler()

	plan := &query.Plan{
		Filter: &query.NotNode{
			Expr: &query.ComparisonNode{
				Field:    strField("name"),
				Operator: query.OpEqual,
				Value:    strVal("rock"),
			},
		},
	}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	sql, _, err := toSQL(res.Where)
	if err != nil {
		t.Fatal(err)
	}

	wantSQL := `NOT (("name" = 'rock'))`
	if sql != wantSQL {
		t.Errorf("SQL: got %q, want %q", sql, wantSQL)
	}
}

func TestCompileContains(t *testing.T) {
	c := newTestCompiler()

	plan := &query.Plan{
		Filter: &query.ContainsNode{
			Field: strField("name"),
			Value: strVal("rock"),
		},
	}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	sql, _, err := toSQL(res.Where)
	if err != nil {
		t.Fatal(err)
	}

	wantSQL := `("name" LIKE '%rock%')`
	if sql != wantSQL {
		t.Errorf("SQL: got %q, want %q", sql, wantSQL)
	}
}

func TestCompileIn(t *testing.T) {
	c := newTestCompiler()

	t.Run("in", func(t *testing.T) {
		plan := &query.Plan{
			Filter: &query.InNode{
				Field:  strField("genre"),
				Values: []query.Value{strVal("rock"), strVal("pop")},
				Not:    false,
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := toSQL(res.Where)
		if err != nil {
			t.Fatal(err)
		}

		wantSQL := `("genre" IN ('rock', 'pop'))`
		if sql != wantSQL {
			t.Errorf("SQL: got %q, want %q", sql, wantSQL)
		}
	})

	t.Run("not in", func(t *testing.T) {
		plan := &query.Plan{
			Filter: &query.InNode{
				Field:  strField("genre"),
				Values: []query.Value{strVal("rock"), strVal("pop")},
				Not:    true,
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := toSQL(res.Where)
		if err != nil {
			t.Fatal(err)
		}

		wantSQL := `("genre" NOT IN ('rock', 'pop'))`
		if sql != wantSQL {
			t.Errorf("SQL: got %q, want %q", sql, wantSQL)
		}
	})

	t.Run("empty in", func(t *testing.T) {
		plan := &query.Plan{
			Filter: &query.InNode{
				Field:  strField("genre"),
				Values: nil,
				Not:    false,
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := toSQL(res.Where)
		if err != nil {
			t.Fatal(err)
		}

		if sql != "1=0" {
			t.Errorf("SQL: got %q, want 1=0", sql)
		}
	})

	t.Run("empty not in", func(t *testing.T) {
		plan := &query.Plan{
			Filter: &query.InNode{
				Field:  strField("genre"),
				Values: nil,
				Not:    true,
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := toSQL(res.Where)
		if err != nil {
			t.Fatal(err)
		}

		if sql != "1=1" {
			t.Errorf("SQL: got %q, want 1=1", sql)
		}
	})
}

func TestCompileOrdering(t *testing.T) {
	c := newTestCompiler()

	t.Run("field asc", func(t *testing.T) {
		plan := &query.Plan{
			OrderBy: []query.Ordering{
				&query.FieldOrdering{
					Field: strField("name"),
					Dir:   query.DirAsc,
				},
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Order) != 1 {
			t.Fatalf("expected 1 ordering, got %d", len(res.Order))
		}

		sql, err := orderToSQL(res.Order[0])
		if err != nil {
			t.Fatal(err)
		}

		want := `"name" ASC`
		if sql != want {
			t.Errorf("got %q, want %q", sql, want)
		}
	})

	t.Run("field desc", func(t *testing.T) {
		plan := &query.Plan{
			OrderBy: []query.Ordering{
				&query.FieldOrdering{
					Field: intField("year"),
					Dir:   query.DirDesc,
				},
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Order) != 1 {
			t.Fatalf("expected 1 ordering, got %d", len(res.Order))
		}

		sql, err := orderToSQL(res.Order[0])
		if err != nil {
			t.Fatal(err)
		}

		want := `"year" DESC`
		if sql != want {
			t.Errorf("got %q, want %q", sql, want)
		}
	})

	t.Run("random", func(t *testing.T) {
		plan := &query.Plan{
			OrderBy: []query.Ordering{
				&query.RandomOrdering{},
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Order) != 1 {
			t.Fatalf("expected 1 ordering, got %d", len(res.Order))
		}

		sql, err := orderToSQL(res.Order[0])
		if err != nil {
			t.Fatal(err)
		}

		want := `RANDOM() ASC`
		if sql != want {
			t.Errorf("got %q, want %q", sql, want)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		plan := &query.Plan{
			OrderBy: []query.Ordering{
				&query.FieldOrdering{
					Field: strField("artist"),
					Dir:   query.DirAsc,
				},
				&query.FieldOrdering{
					Field: intField("year"),
					Dir:   query.DirDesc,
				},
			},
		}

		res, err := c.Compile(plan)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Order) != 2 {
			t.Fatalf("expected 2 orderings, got %d", len(res.Order))
		}

		sql1, _ := orderToSQL(res.Order[0])
		sql2, _ := orderToSQL(res.Order[1])

		want1 := `"artist" ASC`
		want2 := `"year" DESC`

		if sql1 != want1 {
			t.Errorf("order[0]: got %q, want %q", sql1, want1)
		}
		if sql2 != want2 {
			t.Errorf("order[1]: got %q, want %q", sql2, want2)
		}
	})
}

func TestCompileFieldMeta(t *testing.T) {
	c := newTestCompiler()

	field := &query.Field{
		Name: "display_name",
		Type: query.TypeString,
		Meta: map[string]any{
			"column": "actual_column_name",
		},
	}

	plan := &query.Plan{
		Filter: &query.ComparisonNode{
			Field:    field,
			Operator: query.OpEqual,
			Value:    strVal("test"),
		},
	}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	sql, _, err := toSQL(res.Where)
	if err != nil {
		t.Fatal(err)
	}

	wantSQL := `("actual_column_name" = 'test')`
	if sql != wantSQL {
		t.Errorf("got %q, want %q", sql, wantSQL)
	}
}

func TestCompileComplex(t *testing.T) {
	c := newTestCompiler()

	plan := &query.Plan{
		Filter: &query.AndNode{
			Left: &query.OrNode{
				Left: &query.ComparisonNode{
					Field:    strField("genre"),
					Operator: query.OpEqual,
					Value:    strVal("rock"),
				},
				Right: &query.ComparisonNode{
					Field:    strField("genre"),
					Operator: query.OpEqual,
					Value:    strVal("pop"),
				},
			},
			Right: &query.AndNode{
				Left: &query.ComparisonNode{
					Field:    intField("year"),
					Operator: query.OpGreaterEqual,
					Value:    intVal(1970),
				},
				Right: &query.IsNullNode{
					Field: strField("deleted_at"),
					Not:   true,
				},
			},
		},
		OrderBy: []query.Ordering{
			&query.FieldOrdering{
				Field: intField("year"),
				Dir:   query.DirDesc,
			},
		},
	}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	sql, _, err := toSQL(res.Where)
	if err != nil {
		t.Fatal(err)
	}

	wantSQL := `((("genre" = 'rock') OR ("genre" = 'pop')) AND (("year" >= 1970) AND ("deleted_at" IS NOT NULL)))`
	if sql != wantSQL {
		t.Errorf("SQL: got %q, want %q", sql, wantSQL)
	}

	if len(res.Order) != 1 {
		t.Fatalf("expected 1 ordering, got %d", len(res.Order))
	}

	orderSQL, _ := orderToSQL(res.Order[0])
	wantOrder := `"year" DESC`
	if orderSQL != wantOrder {
		t.Errorf("Order: got %q, want %q", orderSQL, wantOrder)
	}
}

func TestCompileEmptyPlan(t *testing.T) {
	c := newTestCompiler()

	plan := &query.Plan{}

	res, err := c.Compile(plan)
	if err != nil {
		t.Fatal(err)
	}

	if res.Where != nil {
		sql, _, _ := toSQL(res.Where)
		t.Errorf("Where: got %q, want nil", sql)
	}
	if len(res.Order) != 0 {
		t.Errorf("Order: got %d items, want 0", len(res.Order))
	}
}
