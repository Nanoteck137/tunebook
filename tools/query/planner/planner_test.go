package planner

import (
	"strings"
	"testing"

	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/ast"
	"github.com/nanoteck137/tunebook/tools/query/parser"
	"github.com/nanoteck137/tunebook/tools/query/schema"
)

func testSchema() *schema.Schema {
	return schema.New().
		AddField("name", query.TypeString).
		AddField("genre", query.TypeString).
		AddField("title", query.TypeString).
		AddField("description", query.TypeString).
		AddField("year", query.TypeInt).
		AddField("rating", query.TypeFloat).
		AddField("duration", query.TypeInt).
		AddField("active", query.TypeBool).
		AddField("deleted_at", query.TypeString, schema.Nullable())
}

func parseExpr(t *testing.T, input string) ast.Expr {
	t.Helper()
	p := parser.New(input)
	expr, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	return expr
}

func TestPlanner(t *testing.T) {
	s := testSchema()
	pl := New(s)

	tests := []struct {
		name  string
		input string
		check func(t *testing.T, plan *query.Plan)
	}{
		{
			name:  "equal string",
			input: `name = "rock"`,
			check: func(t *testing.T, plan *query.Plan) {
				cmp, ok := plan.Filter.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
				}
				if cmp.Field.Name != "name" {
					t.Errorf("field: got %q, want %q", cmp.Field.Name, "name")
				}
				if cmp.Operator != query.OpEqual {
					t.Errorf("op: got %v, want OpEqual", cmp.Operator)
				}
				if cmp.Value.Value != "rock" {
					t.Errorf("value: got %v, want %q", cmp.Value.Value, "rock")
				}
			},
		},
		{
			name:  "greater equal int",
			input: `year >= 1970`,
			check: func(t *testing.T, plan *query.Plan) {
				cmp, ok := plan.Filter.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
				}
				if cmp.Field.Name != "year" {
					t.Errorf("field: got %q, want %q", cmp.Field.Name, "year")
				}
				if cmp.Operator != query.OpGreaterEqual {
					t.Errorf("op: got %v, want OpGreaterEqual", cmp.Operator)
				}
				if cmp.Value.Value != int64(1970) {
					t.Errorf("value: got %v, want 1970", cmp.Value.Value)
				}
			},
		},
		{
			name:  "float comparison",
			input: `rating >= 3.5`,
			check: func(t *testing.T, plan *query.Plan) {
				cmp, ok := plan.Filter.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
				}
				if cmp.Field.Name != "rating" {
					t.Errorf("field: got %q, want %q", cmp.Field.Name, "rating")
				}
				if cmp.Value.Value != 3.5 {
					t.Errorf("value: got %v, want 3.5", cmp.Value.Value)
				}
			},
		},
		{
			name:  "int to float promotion",
			input: `rating >= 3`,
			check: func(t *testing.T, plan *query.Plan) {
				cmp, ok := plan.Filter.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
				}
				if cmp.Value.Type != query.TypeFloat {
					t.Errorf("value type: got %v, want TypeFloat", cmp.Value.Type)
				}
				if cmp.Value.Value != float64(3) {
					t.Errorf("value: got %v, want 3.0", cmp.Value.Value)
				}
			},
		},
		{
			name:  "contains",
			input: `description contains "classic"`,
			check: func(t *testing.T, plan *query.Plan) {
				cnt, ok := plan.Filter.(*query.ContainsNode)
				if !ok {
					t.Fatalf("expected ContainsNode, got %T", plan.Filter)
				}
				if cnt.Field.Name != "description" {
					t.Errorf("field: got %q, want %q", cnt.Field.Name, "description")
				}
				if cnt.Value.Value != "classic" {
					t.Errorf("value: got %v, want %q", cnt.Value.Value, "classic")
				}
			},
		},
		{
			name:  "like",
			input: `title like "%love%"`,
			check: func(t *testing.T, plan *query.Plan) {
				cmp, ok := plan.Filter.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
				}
				if cmp.Operator != query.OpLike {
					t.Errorf("op: got %v, want OpLike", cmp.Operator)
				}
			},
		},
		{
			name:  "and",
			input: `genre = "rock" and year >= 1970`,
			check: func(t *testing.T, plan *query.Plan) {
				and, ok := plan.Filter.(*query.AndNode)
				if !ok {
					t.Fatalf("expected AndNode, got %T", plan.Filter)
				}
				left, ok := and.Left.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("left: expected ComparisonNode, got %T", and.Left)
				}
				if left.Field.Name != "genre" {
					t.Errorf("left field: got %q, want %q", left.Field.Name, "genre")
				}
				right, ok := and.Right.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("right: expected ComparisonNode, got %T", and.Right)
				}
				if right.Field.Name != "year" {
					t.Errorf("right field: got %q, want %q", right.Field.Name, "year")
				}
			},
		},
		{
			name:  "or",
			input: `genre = "rock" or genre = "pop"`,
			check: func(t *testing.T, plan *query.Plan) {
				or, ok := plan.Filter.(*query.OrNode)
				if !ok {
					t.Fatalf("expected OrNode, got %T", plan.Filter)
				}
				if or.Left == nil || or.Right == nil {
					t.Fatal("left or right is nil")
				}
			},
		},
		{
			name:  "not",
			input: `not genre = "rock"`,
			check: func(t *testing.T, plan *query.Plan) {
				not, ok := plan.Filter.(*query.NotNode)
				if !ok {
					t.Fatalf("expected NotNode, got %T", plan.Filter)
				}
				if not.Expr == nil {
					t.Fatal("expr is nil")
				}
			},
		},
		{
			name:  "is null",
			input: `deleted_at is null`,
			check: func(t *testing.T, plan *query.Plan) {
				isNull, ok := plan.Filter.(*query.IsNullNode)
				if !ok {
					t.Fatalf("expected IsNullNode, got %T", plan.Filter)
				}
				if isNull.Field.Name != "deleted_at" {
					t.Errorf("field: got %q, want %q", isNull.Field.Name, "deleted_at")
				}
				if isNull.Not {
					t.Error("expected Not=false")
				}
			},
		},
		{
			name:  "is not null",
			input: `deleted_at is not null`,
			check: func(t *testing.T, plan *query.Plan) {
				isNull, ok := plan.Filter.(*query.IsNullNode)
				if !ok {
					t.Fatalf("expected IsNullNode, got %T", plan.Filter)
				}
				if !isNull.Not {
					t.Error("expected Not=true")
				}
			},
		},
		{
			name:  "field = null converts to is null",
			input: `deleted_at = null`,
			check: func(t *testing.T, plan *query.Plan) {
				isNull, ok := plan.Filter.(*query.IsNullNode)
				if !ok {
					t.Fatalf("expected IsNullNode, got %T", plan.Filter)
				}
				if isNull.Not {
					t.Error("expected Not=false")
				}
			},
		},
		{
			name:  "field != null converts to is not null",
			input: `deleted_at != null`,
			check: func(t *testing.T, plan *query.Plan) {
				isNull, ok := plan.Filter.(*query.IsNullNode)
				if !ok {
					t.Fatalf("expected IsNullNode, got %T", plan.Filter)
				}
				if !isNull.Not {
					t.Error("expected Not=true")
				}
			},
		},
		{
			name:  "in",
			input: `genre in ("rock", "pop", "jazz")`,
			check: func(t *testing.T, plan *query.Plan) {
				in, ok := plan.Filter.(*query.InNode)
				if !ok {
					t.Fatalf("expected InNode, got %T", plan.Filter)
				}
				if in.Field.Name != "genre" {
					t.Errorf("field: got %q, want %q", in.Field.Name, "genre")
				}
				if len(in.Values) != 3 {
					t.Fatalf("values len: got %d, want 3", len(in.Values))
				}
				if in.Not {
					t.Error("expected Not=false")
				}
			},
		},
		{
			name:  "not in",
			input: `genre not in ("rock", "pop")`,
			check: func(t *testing.T, plan *query.Plan) {
				in, ok := plan.Filter.(*query.InNode)
				if !ok {
					t.Fatalf("expected InNode, got %T", plan.Filter)
				}
				if !in.Not {
					t.Error("expected Not=true")
				}
			},
		},
		{
			name:  "in with ints",
			input: `year in (1970, 1980, 1990)`,
			check: func(t *testing.T, plan *query.Plan) {
				in, ok := plan.Filter.(*query.InNode)
				if !ok {
					t.Fatalf("expected InNode, got %T", plan.Filter)
				}
				if len(in.Values) != 3 {
					t.Fatalf("values len: got %d, want 3", len(in.Values))
				}
				if in.Values[0].Value != int64(1970) {
					t.Errorf("values[0]: got %v, want 1970", in.Values[0].Value)
				}
			},
		},
		{
			name:  "bool comparison",
			input: `active = true`,
			check: func(t *testing.T, plan *query.Plan) {
				cmp, ok := plan.Filter.(*query.ComparisonNode)
				if !ok {
					t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
				}
				if cmp.Value.Value != true {
					t.Errorf("value: got %v, want true", cmp.Value.Value)
				}
			},
		},
		{
			name:  "complex nested",
			input: `(genre = "rock" or genre = "pop") and year >= 1970 and deleted_at is not null`,
			check: func(t *testing.T, plan *query.Plan) {
				if plan.Filter == nil {
					t.Fatal("filter is nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			plan, err := pl.Plan(expr)
			if err != nil {
				t.Fatal(err)
			}
			tt.check(t, plan)
		})
	}
}

func TestPlannerErrors(t *testing.T) {
	s := testSchema()
	pl := New(s)

	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "unknown field",
			input:   `unknown = "value"`,
			wantErr: "unknown field 'unknown'",
		},
		{
			name:    "contains on int",
			input:   `year contains "rock"`,
			wantErr: "'contains' can only be used on string fields",
		},
		{
			name:    "like on int",
			input:   `year like "%1970%"`,
			wantErr: "'like' can only be used on string fields",
		},
		{
			name:    "greater on string",
			input:   `name > "rock"`,
			wantErr: "cannot use '>' on string field",
		},
		{
			name:    "greater equal on string",
			input:   `name >= "rock"`,
			wantErr: "cannot use '>=' on string field",
		},
		{
			name:    "less on string",
			input:   `name < "rock"`,
			wantErr: "cannot use '<' on string field",
		},
		{
			name:    "less equal on string",
			input:   `name <= "rock"`,
			wantErr: "cannot use '<=' on string field",
		},
		{
			name:    "string on int field",
			input:   `year = "rock"`,
			wantErr: "expected integer value, got string",
		},
		{
			name:    "int on string field",
			input:   `name = 1970`,
			wantErr: "expected string value, got integer",
		},
		{
			name:    "float on int field",
			input:   `year = 3.5`,
			wantErr: "expected integer value, got float",
		},
		{
			name:    "bool on string field",
			input:   `name = true`,
			wantErr: "expected string value, got boolean",
		},
		{
			name:    "is null on non-nullable",
			input:   `name is null`,
			wantErr: "field 'name' cannot be null",
		},
		{
			name:    "is not null on non-nullable",
			input:   `name is not null`,
			wantErr: "field 'name' cannot be null",
		},
		{
			name:    "field = null on non-nullable",
			input:   `name = null`,
			wantErr: "field 'name' cannot be null",
		},
		{
			name:    "in with wrong type",
			input:   `year in ("rock", "pop")`,
			wantErr: "expected integer value, got string",
		},
		{
			name:    "contains on float",
			input:   `rating contains "3"`,
			wantErr: "'contains' can only be used on string fields",
		},
		{
			name:    "greater on bool",
			input:   `active > true`,
			wantErr: "cannot use '>' on boolean field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			_, err := pl.Plan(expr)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error: got %q, want to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestPlannerFieldResolution(t *testing.T) {
	s := schema.New().
		AddField("display_name", query.TypeString, schema.Column("actual_column"))

	pl := New(s)

	expr := parseExpr(t, `display_name = "test"`)
	plan, err := pl.Plan(expr)
	if err != nil {
		t.Fatal(err)
	}

	cmp, ok := plan.Filter.(*query.ComparisonNode)
	if !ok {
		t.Fatalf("expected ComparisonNode, got %T", plan.Filter)
	}

	if cmp.Field.Name != "display_name" {
		t.Errorf("field name: got %q, want %q", cmp.Field.Name, "display_name")
	}

	col, ok := cmp.Field.Meta["column"]
	if !ok {
		t.Fatal("field meta missing 'column'")
	}
	if col != "actual_column" {
		t.Errorf("column meta: got %v, want %q", col, "actual_column")
	}
}
