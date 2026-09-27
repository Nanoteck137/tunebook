package query

import (
	"testing"

	"github.com/nanoteck137/tunebook/tools/query/parser"
)

func TestAndFilters(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "both empty",
			left:  "",
			right: "",
			want:  "",
		},
		{
			name:  "left empty passes right through",
			left:  "",
			right: `name contains "x"`,
			want:  `name contains "x"`,
		},
		{
			name:  "right empty passes left through",
			left:  `name contains "x"`,
			right: "",
			want:  `name contains "x"`,
		},
		{
			name:  "whitespace is trimmed",
			left:  "  ",
			right: `name contains "x"`,
			want:  `name contains "x"`,
		},
		{
			name:  "simple operands are grouped",
			left:  `name contains "a"`,
			right: `name contains "b"`,
			want:  `((name contains "a") and (name contains "b"))`,
		},
		{
			name:  "or in the left operand cannot absorb the right",
			left:  `name contains "a" or name contains "b"`,
			right: `name contains "c"`,
			want:  `(((name contains "a") or (name contains "b")) and (name contains "c"))`,
		},
		{
			name:  "or in the right operand",
			left:  `name contains "a"`,
			right: `name contains "b" or name contains "c"`,
			want:  `((name contains "a") and ((name contains "b") or (name contains "c")))`,
		},
		{
			name:  "nested parens in the left operand are preserved",
			left:  `(name contains "a" and name contains "b")`,
			right: `name contains "c"`,
			want:  `(((name contains "a") and (name contains "b")) and (name contains "c"))`,
		},
		{
			name:  "not operand",
			left:  `not name contains "a"`,
			right: `name contains "b"`,
			want:  `((not (name contains "a")) and (name contains "b"))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AndFilters(tt.left, tt.right)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("AndFilters(%q, %q)\n got: %s\nwant: %s",
					tt.left, tt.right, got, tt.want)
			}

			// The whole point of going through the AST is that the result
			// must survive a round trip back through the parser.
			if got != "" {
				if _, err := parser.New(got).Parse(); err != nil {
					t.Errorf("result does not re-parse: %v", err)
				}
			}
		})
	}
}

func TestAndFiltersInvalid(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
	}{
		{
			name:  "invalid left",
			left:  `name contains`,
			right: `name contains "b"`,
		},
		{
			name:  "invalid right",
			left:  `name contains "a"`,
			right: `name contains`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := AndFilters(tt.left, tt.right); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}
