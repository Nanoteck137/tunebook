package query

import (
	"fmt"
	"strings"

	"github.com/nanoteck137/tunebook/tools/query/ast"
	"github.com/nanoteck137/tunebook/tools/query/parser"
)

// AndFilters combines two filter expressions with a logical AND.
//
// The operands are parsed and rejoined through the AST instead of being
// concatenated as text. A bare "left and right" is not equivalent: a saved
// filter like `a or b` would absorb the right operand, since `and` binds
// tighter than `or`, turning `a or b and c` into `a or (b and c)`. Parsing
// first restores the grouping, and re-serialising the tree emits the explicit
// parentheses that make the result safe to parse again downstream.
func AndFilters(left, right string) (string, error) {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	switch {
	case left == "":
		return right, nil
	case right == "":
		return left, nil
	}

	l, err := parser.New(left).Parse()
	if err != nil {
		return "", fmt.Errorf("parse left filter: %w", err)
	}

	r, err := parser.New(right).Parse()
	if err != nil {
		return "", fmt.Errorf("parse right filter: %w", err)
	}

	return (&ast.BinaryExpr{
		Left:  l,
		Op:    ast.OpAnd,
		Right: r,
	}).String(), nil
}
