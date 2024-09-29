package parser

import (
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
)

type parser struct {
	tokens    []lexer.Token
	currIndex int
}

func (p *parser) consume(predicate func(lexer.Token) bool) (*lexer.Token, *file.Error) {
	if p.currIndex >= len(p.tokens) {
		return nil, &file.Error{
			Location: file.SourceLocation{Line: -1, Col: -1},
			Message:  "incomplete token stream",
		}
	}

	currToken := p.tokens[p.currIndex]
	p.currIndex++

	if !predicate(currToken) {
		return nil, &file.Error{
			Location: currToken.Location,
			Message:  "unexpected token",
		}
	}
	return &currToken, nil
}

func (p *parser) parseNumber() (*NumberNode, *file.Error) {
	currToken, err := p.consume(func(token lexer.Token) bool {
		if len(token.Source) > 0 &&
			(unicode.IsDigit(token.Source[0]) || strings.ContainsRune("-+", token.Source[0])) {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	source := currToken.Source

	if strings.Contains(source.String(), "/") {
		items := strings.Split(source.String(), "/")
		n, _ := strconv.Atoi(items[0])
		d, _ := strconv.Atoi(items[1])
		return NewNumberNode(currToken.Location, n, d), nil
	} else if strings.Contains(source.String(), ".") {
		items := strings.Split(source.String(), ".")
		n, _ := strconv.Atoi(items[0] + items[1])
		d := math.Pow10(len(items[1]))
		return NewNumberNode(currToken.Location, n, int(d)), nil
	} else {
		n, _ := strconv.Atoi(source.String())
		return NewNumberNode(currToken.Location, n, 1), nil
	}
}

func (p *parser) parseExpr() (ExprNode, *file.Error) {
	if p.currIndex >= len(p.tokens) {
		return nil, &file.Error{
			Location: file.SourceLocation{Line: -1, Col: -1},
			Message:  "incomplete token stream",
		}
	}

	currToken := p.tokens[p.currIndex]
	if len(currToken.Source) > 0 &&
		(unicode.IsDigit(currToken.Source[0]) || strings.ContainsRune("-+", currToken.Source[0])) {
		return p.parseNumber()
	} else {
		return nil, &file.Error{Location: currToken.Location, Message: "unrecognized token"}
	}
}

func Parse(tokens []lexer.Token) (ExprNode, *file.Error) {
	parser := &parser{
		tokens:    tokens,
		currIndex: 0,
	}

	expr, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}
	if parser.currIndex < len(parser.tokens) {
		return nil, &file.Error{
			Location: parser.tokens[parser.currIndex].Location,
			Message:  "redundant token(s)",
		}
	}
	return expr, nil
}
