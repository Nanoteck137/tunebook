package parser

import (
	"fmt"
	"strconv"

	"github.com/nanoteck137/tunebook/tools/query/ast"
	"github.com/nanoteck137/tunebook/tools/query/lexer"
)

type ParseError struct {
	Pos     lexer.Position
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Column, e.Message)
}

type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
	peeked  bool
	peekTok lexer.Token
}

func New(input string) *Parser {
	return &Parser{
		lexer: lexer.New(input),
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	// Get first token
	if err := p.advance(); err != nil {
		return nil, err
	}

	if p.current.Type == lexer.TokenEOF {
		return nil, p.errorf("empty expression")
	}

	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	if p.current.Type != lexer.TokenEOF {
		return nil, p.errorf("unexpected token: %s", p.current.Value)
	}

	return expr, nil
}

func (p *Parser) parseOr() (ast.Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TokenOr {
		if err := p.advance(); err != nil {
			return nil, err
		}
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: ast.OpOr, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAnd() (ast.Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	for p.current.Type == lexer.TokenAnd {
		if err := p.advance(); err != nil {
			return nil, err
		}
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: ast.OpAnd, Right: right}
	}

	return left, nil
}

func (p *Parser) parseNot() (ast.Expr, error) {
	if p.current.Type == lexer.TokenNot {
		if err := p.advance(); err != nil {
			return nil, err
		}
		expr, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: ast.OpNot, Expr: expr}, nil
	}

	return p.parseComparison()
}

func (p *Parser) parseComparison() (ast.Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	tok := p.current

	if tok.Type == lexer.TokenIs {
		if err := p.advance(); err != nil {
			return nil, err
		}
		not := false
		if p.current.Type == lexer.TokenNot {
			if err := p.advance(); err != nil {
				return nil, err
			}
			not = true
		}
		if err := p.expect(lexer.TokenNull); err != nil {
			return nil, err
		}
		return &ast.IsNullExpr{Field: left, Not: not}, nil
	}

	if tok.Type == lexer.TokenNot {
		// Peek ahead to see if next token is "in"
		nextTok, err := p.peek()
		if err != nil {
			return nil, err
		}
		if nextTok.Type == lexer.TokenIn {
			// Consume the "not" token
			if err := p.advance(); err != nil {
				return nil, err
			}
			// Consume the "in" token
			if err := p.advance(); err != nil {
				return nil, err
			}
			values, err := p.parseInList()
			if err != nil {
				return nil, err
			}
			return &ast.InExpr{Field: left, Values: values, Not: true}, nil
		}
		return nil, p.errorf("expected 'in' after 'not'")
	}

	if tok.Type == lexer.TokenIn {
		if err := p.advance(); err != nil {
			return nil, err
		}
		values, err := p.parseInList()
		if err != nil {
			return nil, err
		}
		return &ast.InExpr{Field: left, Values: values, Not: false}, nil
	}

	var op ast.BinaryOp
	switch tok.Type {
	case lexer.TokenEq:
		op = ast.OpEq
	case lexer.TokenNeq:
		op = ast.OpNeq
	case lexer.TokenGt:
		op = ast.OpGt
	case lexer.TokenGte:
		op = ast.OpGte
	case lexer.TokenLt:
		op = ast.OpLt
	case lexer.TokenLte:
		op = ast.OpLte
	case lexer.TokenContains:
		op = ast.OpContains
	case lexer.TokenLike:
		op = ast.OpLike
	case lexer.TokenHas:
		op = ast.OpHas
	default:
		return left, nil
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	right, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	return &ast.BinaryExpr{Left: left, Op: op, Right: right}, nil
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	tok := p.current

	switch tok.Type {
	case lexer.TokenString:
		if err := p.advance(); err != nil {
			return nil, err
		}
		s, err := strconv.Unquote(tok.Value)
		if err != nil {
			return nil, p.errorf("invalid string: %s", tok.Value)
		}
		return &ast.StringLit{Value: s}, nil

	case lexer.TokenInt:
		if err := p.advance(); err != nil {
			return nil, err
		}
		i, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			return nil, p.errorf("invalid integer: %s", tok.Value)
		}
		return &ast.IntLit{Value: i}, nil

	case lexer.TokenFloat:
		if err := p.advance(); err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, p.errorf("invalid float: %s", tok.Value)
		}
		return &ast.FloatLit{Value: f}, nil

	case lexer.TokenTrue:
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &ast.BoolLit{Value: true}, nil

	case lexer.TokenFalse:
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &ast.BoolLit{Value: false}, nil

	case lexer.TokenNull:
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &ast.NullLit{}, nil

	case lexer.TokenIdent:
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &ast.FieldRef{Name: tok.Value}, nil

	case lexer.TokenLParen:
		if err := p.advance(); err != nil {
			return nil, err
		}
		expr, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.TokenRParen); err != nil {
			return nil, err
		}
		return expr, nil

	default:
		return nil, p.errorf("unexpected token: %s", tok.Value)
	}
}

func (p *Parser) parseInList() ([]ast.Expr, error) {
	if err := p.expect(lexer.TokenLParen); err != nil {
		return nil, err
	}

	var items []ast.Expr

	if p.current.Type != lexer.TokenRParen {
		item, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		items = append(items, item)

		for p.current.Type == lexer.TokenComma {
			if err := p.advance(); err != nil {
				return nil, err
			}
			item, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	}

	if err := p.expect(lexer.TokenRParen); err != nil {
		return nil, err
	}

	return items, nil
}

func (p *Parser) advance() error {
	if p.peeked {
		p.current = p.peekTok
		p.peeked = false
		return nil
	}
	tok, err := p.lexer.Next()
	if err != nil {
		return err
	}
	p.current = tok
	return nil
}

func (p *Parser) peek() (lexer.Token, error) {
	if p.peeked {
		return p.peekTok, nil
	}
	tok, err := p.lexer.Next()
	if err != nil {
		return lexer.Token{}, err
	}
	p.peeked = true
	p.peekTok = tok
	return tok, nil
}

func (p *Parser) expect(tt lexer.TokenType) error {
	if p.current.Type != tt {
		return p.errorf("expected %s, got %s", tt, p.current.Value)
	}
	return p.advance()
}

func (p *Parser) errorf(format string, args ...interface{}) error {
	pos := p.current.Pos
	return &ParseError{
		Pos:     pos,
		Message: fmt.Sprintf(format, args...),
	}
}
