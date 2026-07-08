package parser

import (
	"testing"
)

func parse(t *testing.T, input string) string {
	t.Helper()
	p := New(input)
	expr, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	return expr.String()
}

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple comparison",
			input: "year >= 1970",
			want:  "(year >= 1970)",
		},
		{
			name:  "equal string",
			input: `name = "rock"`,
			want:  `(name = "rock")`,
		},
		{
			name:  "not equal",
			input: `genre != "pop"`,
			want:  `(genre != "pop")`,
		},
		{
			name:  "greater than",
			input: "year > 2000",
			want:  "(year > 2000)",
		},
		{
			name:  "less than",
			input: "rating < 3.5",
			want:  "(rating < 3.5)",
		},
		{
			name:  "less equal",
			input: "duration <= 300",
			want:  "(duration <= 300)",
		},
		{
			name:  "contains",
			input: `tags contains "rock"`,
			want:  `(tags contains "rock")`,
		},
		{
			name:  "like",
			input: `title like "%love%"`,
			want:  `(title like "%love%")`,
		},
		{
			name:  "and",
			input: `genre = "rock" and year >= 1970`,
			want:  `((genre = "rock") and (year >= 1970))`,
		},
		{
			name:  "or",
			input: `genre = "rock" or genre = "pop"`,
			want:  `((genre = "rock") or (genre = "pop"))`,
		},
		{
			name:  "not prefix",
			input: `not genre = "rock"`,
			want:  `(not (genre = "rock"))`,
		},
		{
			name:  "not with and",
			input: `not genre = "rock" and year > 1970`,
			want:  `((not (genre = "rock")) and (year > 1970))`,
		},
		{
			name:  "parentheses",
			input: `(genre = "rock" or genre = "pop") and year >= 1970`,
			want:  `(((genre = "rock") or (genre = "pop")) and (year >= 1970))`,
		},
		{
			name:  "nested parentheses",
			input: `(genre = "rock") and (year >= 1970)`,
			want:  `((genre = "rock") and (year >= 1970))`,
		},
		{
			name:  "is null",
			input: "deleted_at is null",
			want:  "(deleted_at is null)",
		},
		{
			name:  "is not null",
			input: "deleted_at is not null",
			want:  "(deleted_at is not null)",
		},
		{
			name:  "in operator",
			input: `genre in ("rock", "pop", "jazz")`,
			want:  `(genre in ("rock", "pop", "jazz"))`,
		},
		{
			name:  "not in operator",
			input: `genre not in ("rock", "pop")`,
			want:  `(genre not in ("rock", "pop"))`,
		},
		{
			name:  "in with ints",
			input: "year in (1970, 1980, 1990)",
			want:  "(year in (1970, 1980, 1990))",
		},
		{
			name:  "empty in",
			input: "genre in ()",
			want:  "(genre in ())",
		},
		{
			name:  "boolean true",
			input: "active = true",
			want:  "(active = true)",
		},
		{
			name:  "boolean false",
			input: "deleted = false",
			want:  "(deleted = false)",
		},
		{
			name:  "float literal",
			input: "rating >= 3.5",
			want:  "(rating >= 3.5)",
		},
		{
			name:  "complex expression",
			input: `(genre = "rock" or genre = "pop") and year >= 1970 and deleted_at is null`,
			want:  `((((genre = "rock") or (genre = "pop")) and (year >= 1970)) and (deleted_at is null))`,
		},
		{
			name:  "double not",
			input: `not not genre = "rock"`,
			want:  `(not (not (genre = "rock")))`,
		},
		{
			name:  "precedence: and before or",
			input: `a = 1 or b = 2 and c = 3`,
			want:  `((a = 1) or ((b = 2) and (c = 3)))`,
		},
		{
			name:  "precedence: not before and",
			input: `not a = 1 and b = 2`,
			want:  `((not (a = 1)) and (b = 2))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"whitespace only", "   "},
		{"unexpected token", ")"},
		{"missing right side", "year >="},
		{"missing left side", "= 1970"},
		{"unclosed paren", "(year >= 1970"},
		{"extra close paren", "year >= 1970)"},
		{"not without expression", "not"},
		{"not at comparison level without in", "genre not"},
		{"not at comparison level wrong keyword", `genre not "rock"`},
		{"in without list", "genre in"},
		{"in without close paren", `genre in ("rock"`},
		{"missing value after comma", `genre in ("rock",)`},
		{"double operator", "year >= >= 1970"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			_, err := p.Parse()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestParseErrorFormat(t *testing.T) {
	input := `genre = "rock" and`
	p := New(input)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}

	if parseErr.Pos.Line != 1 {
		t.Errorf("error line: got %d, want 1", parseErr.Pos.Line)
	}

	errStr := parseErr.Error()
	if errStr == "" {
		t.Error("error string is empty")
	}
}
