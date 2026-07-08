package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []TokenType
		values []string
	}{
		{
			name:   "simple comparison",
			input:  "year >= 1970",
			tokens: []TokenType{TokenIdent, TokenGte, TokenInt, TokenEOF},
			values: []string{"year", ">=", "1970", ""},
		},
		{
			name:   "string literal",
			input:  `name = "rock"`,
			tokens: []TokenType{TokenIdent, TokenEq, TokenString, TokenEOF},
			values: []string{"name", "=", `"rock"`, ""},
		},
		{
			name:   "and expression",
			input:  `genre = "rock" and year >= 1970`,
			tokens: []TokenType{TokenIdent, TokenEq, TokenString, TokenAnd, TokenIdent, TokenGte, TokenInt, TokenEOF},
			values: []string{"genre", "=", `"rock"`, "and", "year", ">=", "1970", ""},
		},
		{
			name:   "contains keyword",
			input:  `tags contains "rock"`,
			tokens: []TokenType{TokenIdent, TokenContains, TokenString, TokenEOF},
			values: []string{"tags", "contains", `"rock"`, ""},
		},
		{
			name:   "is null",
			input:  "deleted_at is null",
			tokens: []TokenType{TokenIdent, TokenIs, TokenNull, TokenEOF},
			values: []string{"deleted_at", "is", "null", ""},
		},
		{
			name:   "is not null",
			input:  "deleted_at is not null",
			tokens: []TokenType{TokenIdent, TokenIs, TokenNot, TokenNull, TokenEOF},
			values: []string{"deleted_at", "is", "not", "null", ""},
		},
		{
			name:   "in operator",
			input:  `genre in ("rock", "pop")`,
			tokens: []TokenType{TokenIdent, TokenIn, TokenLParen, TokenString, TokenComma, TokenString, TokenRParen, TokenEOF},
			values: []string{"genre", "in", "(", `"rock"`, ",", `"pop"`, ")", ""},
		},
		{
			name:   "not in operator",
			input:  `genre not in ("rock", "pop")`,
			tokens: []TokenType{TokenIdent, TokenNot, TokenIn, TokenLParen, TokenString, TokenComma, TokenString, TokenRParen, TokenEOF},
			values: []string{"genre", "not", "in", "(", `"rock"`, ",", `"pop"`, ")", ""},
		},
		{
			name:   "all comparison operators",
			input:  "a = 1 != 2 > 3 >= 4 < 5 <= 6",
			tokens: []TokenType{TokenIdent, TokenEq, TokenInt, TokenNeq, TokenInt, TokenGt, TokenInt, TokenGte, TokenInt, TokenLt, TokenInt, TokenLte, TokenInt, TokenEOF},
			values: []string{"a", "=", "1", "!=", "2", ">", "3", ">=", "4", "<", "5", "<=", "6", ""},
		},
		{
			name:   "parentheses",
			input:  "(a = 1)",
			tokens: []TokenType{TokenLParen, TokenIdent, TokenEq, TokenInt, TokenRParen, TokenEOF},
			values: []string{"(", "a", "=", "1", ")", ""},
		},
		{
			name:   "boolean literals",
			input:  "active = true and deleted = false",
			tokens: []TokenType{TokenIdent, TokenEq, TokenTrue, TokenAnd, TokenIdent, TokenEq, TokenFalse, TokenEOF},
			values: []string{"active", "=", "true", "and", "deleted", "=", "false", ""},
		},
		{
			name:   "float literal",
			input:  "rating >= 3.5",
			tokens: []TokenType{TokenIdent, TokenGte, TokenFloat, TokenEOF},
			values: []string{"rating", ">=", "3.5", ""},
		},
		{
			name:   "or expression",
			input:  `genre = "rock" or genre = "pop"`,
			tokens: []TokenType{TokenIdent, TokenEq, TokenString, TokenOr, TokenIdent, TokenEq, TokenString, TokenEOF},
			values: []string{"genre", "=", `"rock"`, "or", "genre", "=", `"pop"`, ""},
		},
		{
			name:   "not prefix",
			input:  `not genre = "rock"`,
			tokens: []TokenType{TokenNot, TokenIdent, TokenEq, TokenString, TokenEOF},
			values: []string{"not", "genre", "=", `"rock"`, ""},
		},
		{
			name:   "like operator",
			input:  `title like "%love%"`,
			tokens: []TokenType{TokenIdent, TokenLike, TokenString, TokenEOF},
			values: []string{"title", "like", `"%love%"`, ""},
		},
		{
			name:   "case insensitive keywords",
			input:  `genre = "rock" AND year >= 1970`,
			tokens: []TokenType{TokenIdent, TokenEq, TokenString, TokenAnd, TokenIdent, TokenGte, TokenInt, TokenEOF},
			values: []string{"genre", "=", `"rock"`, "AND", "year", ">=", "1970", ""},
		},
		{
			name:   "string with escape",
			input:  `name = "rock\"band"`,
			tokens: []TokenType{TokenIdent, TokenEq, TokenString, TokenEOF},
			values: []string{"name", "=", `"rock"band"`, ""},
		},
		{
			name:   "empty input",
			input:  "",
			tokens: []TokenType{TokenEOF},
			values: []string{""},
		},
		{
			name:   "whitespace only",
			input:  "   \t\n  ",
			tokens: []TokenType{TokenEOF},
			values: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tokens, err := l.Scan()
			if err != nil {
				t.Fatal(err)
			}

			if len(tokens) != len(tt.tokens) {
				t.Fatalf("token count: got %d, want %d\ngot: %v", len(tokens), len(tt.tokens), tokens)
			}

			for i, tok := range tokens {
				if tok.Type != tt.tokens[i] {
					t.Errorf("token[%d] type: got %v, want %v", i, tok.Type, tt.tokens[i])
				}
				if tok.Value != tt.values[i] {
					t.Errorf("token[%d] value: got %q, want %q", i, tok.Value, tt.values[i])
				}
			}
		})
	}
}

func TestLexerErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unterminated string", `name = "rock`},
		{"invalid escape", `name = "rock\x"`},
		{"unexpected character", "year @ 1970"},
		{"bang without equals", "year ! 1970"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			_, err := l.Scan()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestLexerPositions(t *testing.T) {
	input := "year >= 1970"
	l := New(input)
	tokens, err := l.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].Pos.Line != 1 || tokens[0].Pos.Column != 1 {
		t.Errorf("first token position: got %d:%d, want 1:1", tokens[0].Pos.Line, tokens[0].Pos.Column)
	}

	if tokens[1].Pos.Line != 1 || tokens[1].Pos.Column != 6 {
		t.Errorf("second token position: got %d:%d, want 1:6", tokens[1].Pos.Line, tokens[1].Pos.Column)
	}
}

func TestLexerMultilinePositions(t *testing.T) {
	input := "year >= 1970\nand name = \"rock\""
	l := New(input)
	tokens, err := l.Scan()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, tok := range tokens {
		if tok.Type == TokenAnd {
			found = true
			if tok.Pos.Line != 2 {
				t.Errorf("'and' token line: got %d, want 2", tok.Pos.Line)
			}
		}
	}
	if !found {
		t.Error("'and' token not found")
	}
}
