package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
}

func New(input string) *Lexer {
	return &Lexer{input: input, line: 1, col: 1}
}

func (l *Lexer) Scan() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.scanToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens, nil
}

func (l *Lexer) Next() (Token, error) {
	return l.scanToken()
}

func (l *Lexer) scanToken() (Token, error) {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF, Pos: l.position()}, nil
	}

	ch := l.input[l.pos]
	pos := l.position()

	switch {
	case ch == '"':
		return l.scanString()
	case ch == '(':
		l.advance()
		return Token{Type: TokenLParen, Value: "(", Pos: pos}, nil
	case ch == ')':
		l.advance()
		return Token{Type: TokenRParen, Value: ")", Pos: pos}, nil
	case ch == ',':
		l.advance()
		return Token{Type: TokenComma, Value: ",", Pos: pos}, nil
	case ch == '=':
		l.advance()
		return Token{Type: TokenEq, Value: "=", Pos: pos}, nil
	case ch == '!':
		l.advance()
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			return Token{Type: TokenNeq, Value: "!=", Pos: pos}, nil
		}
		return Token{}, l.errorfAt(pos, "expected '=' after '!'")
	case ch == '>':
		l.advance()
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			return Token{Type: TokenGte, Value: ">=", Pos: pos}, nil
		}
		return Token{Type: TokenGt, Value: ">", Pos: pos}, nil
	case ch == '<':
		l.advance()
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			return Token{Type: TokenLte, Value: "<=", Pos: pos}, nil
		}
		return Token{Type: TokenLt, Value: "<", Pos: pos}, nil
	case isDigit(ch):
		return l.scanNumber()
	case isLetter(ch) || ch == '_':
		return l.scanIdentOrKeyword()
	default:
		return Token{}, l.errorf("unexpected character: %c", ch)
	}
}

func (l *Lexer) scanString() (Token, error) {
	pos := l.position()
	l.advance()

	var buf strings.Builder
	buf.WriteByte('"')

	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		if ch == '\\' {
			l.advance()
			if l.pos >= len(l.input) {
				return Token{}, l.errorfAt(pos, "unterminated string")
			}
			esc := l.input[l.pos]
			switch esc {
			case '"', '\\':
				buf.WriteByte(esc)
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			default:
				return Token{}, l.errorfAt(pos, "invalid escape sequence: \\%c", esc)
			}
			l.advance()
			continue
		}

		if ch == '"' {
			buf.WriteByte('"')
			l.advance()
			return Token{Type: TokenString, Value: buf.String(), Pos: pos}, nil
		}

		buf.WriteByte(ch)
		l.advance()
	}

	return Token{}, l.errorfAt(pos, "unterminated string")
}

func (l *Lexer) scanNumber() (Token, error) {
	pos := l.position()
	start := l.pos

	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.advance()
	}

	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		if l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]) {
			l.advance()
			for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
				l.advance()
			}
			return Token{Type: TokenFloat, Value: l.input[start:l.pos], Pos: pos}, nil
		}
	}

	return Token{Type: TokenInt, Value: l.input[start:l.pos], Pos: pos}, nil
}

func (l *Lexer) scanIdentOrKeyword() (Token, error) {
	pos := l.position()
	start := l.pos

	for l.pos < len(l.input) && (isLetter(l.input[l.pos]) || isDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
		l.advance()
	}

	word := l.input[start:l.pos]

	switch strings.ToLower(word) {
	case "and":
		return Token{Type: TokenAnd, Value: word, Pos: pos}, nil
	case "or":
		return Token{Type: TokenOr, Value: word, Pos: pos}, nil
	case "not":
		return Token{Type: TokenNot, Value: word, Pos: pos}, nil
	case "contains":
		return Token{Type: TokenContains, Value: word, Pos: pos}, nil
	case "like":
		return Token{Type: TokenLike, Value: word, Pos: pos}, nil
	case "is":
		return Token{Type: TokenIs, Value: word, Pos: pos}, nil
	case "null":
		return Token{Type: TokenNull, Value: word, Pos: pos}, nil
	case "in":
		return Token{Type: TokenIn, Value: word, Pos: pos}, nil
	case "true":
		return Token{Type: TokenTrue, Value: word, Pos: pos}, nil
	case "false":
		return Token{Type: TokenFalse, Value: word, Pos: pos}, nil
	case "has":
		return Token{Type: TokenHas, Value: word, Pos: pos}, nil
	default:
		return Token{Type: TokenIdent, Value: word, Pos: pos}, nil
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && isWhitespace(l.input[l.pos]) {
		l.advance()
	}
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.pos++
	}
}

func (l *Lexer) position() Position {
	return Position{Line: l.line, Column: l.col, Offset: l.pos}
}

func (l *Lexer) errorf(format string, args ...any) error {
	return l.errorfAt(l.position(), format, args...)
}

func (l *Lexer) errorfAt(pos Position, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%d:%d: %s", pos.Line, pos.Column, msg)
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
