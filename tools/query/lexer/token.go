package lexer

import "fmt"

type TokenType int

const (
	TokenIdent TokenType = iota
	TokenString
	TokenInt
	TokenFloat

	TokenAnd
	TokenOr
	TokenNot
	TokenContains
	TokenLike
	TokenIs
	TokenNull
	TokenIn
	TokenTrue
	TokenFalse
	TokenHas

	TokenEq
	TokenNeq
	TokenGt
	TokenGte
	TokenLt
	TokenLte

	TokenLParen
	TokenRParen
	TokenComma

	TokenEOF
)

func (t TokenType) String() string {
	switch t {
	case TokenIdent:
		return "IDENT"
	case TokenString:
		return "STRING"
	case TokenInt:
		return "INT"
	case TokenFloat:
		return "FLOAT"
	case TokenAnd:
		return "AND"
	case TokenOr:
		return "OR"
	case TokenNot:
		return "NOT"
	case TokenContains:
		return "CONTAINS"
	case TokenLike:
		return "LIKE"
	case TokenIs:
		return "IS"
	case TokenNull:
		return "NULL"
	case TokenIn:
		return "IN"
	case TokenTrue:
		return "TRUE"
	case TokenFalse:
		return "FALSE"
	case TokenHas:
		return "HAS"
	case TokenEq:
		return "="
	case TokenNeq:
		return "!="
	case TokenGt:
		return ">"
	case TokenGte:
		return ">="
	case TokenLt:
		return "<"
	case TokenLte:
		return "<="
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenComma:
		return ","
	case TokenEOF:
		return "EOF"
	}
	return "?"
}

type Position struct {
	Line   int
	Column int
	Offset int
}

type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}

func (t Token) String() string {
	if t.Value != "" {
		return fmt.Sprintf("%s(%s)", t.Type, t.Value)
	}
	return t.Type.String()
}
