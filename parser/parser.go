package parser

import (
	"math"
	"strconv"
	"strings"
	"unicode"

	. "github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
)

var intrinsics = [...]string{
	"void",
	"isVoid", "isNum", "isStr", "isCont",
	"getline", "put",
	"reg", "go",
	"callCC", "eval", "exit",
}

type parser struct {
	tokens    []*lexer.Token
	currIndex int
}

func (p *parser) consume(predicate func(*lexer.Token) bool) (*lexer.Token, *file.Error) {
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
	return currToken, nil
}

func (p *parser) parseNumber() (*NumberNode, *file.Error) {
	currToken, err := p.consume(func(token *lexer.Token) bool {
		return len(token.Source) > 0 && token.Kind == lexer.Number
	})
	if err != nil {
		return nil, err
	}

	source := currToken.Source

	if strings.Contains(source, "/") {
		items := strings.Split(source, "/")
		n, _ := strconv.Atoi(items[0])
		d, _ := strconv.Atoi(items[1])
		return NewNumberNode(currToken.Location, n, d), nil
	} else if strings.Contains(source, ".") {
		items := strings.Split(source, ".")
		n, _ := strconv.Atoi(items[0] + items[1])
		d := math.Pow10(len(items[1]))
		return NewNumberNode(currToken.Location, n, int(d)), nil
	} else {
		n, _ := strconv.Atoi(source)
		return NewNumberNode(currToken.Location, n, 1), nil
	}
}

func (p *parser) parseString() (*StringNode, *file.Error) {
	currToken, err := p.consume(func(token *lexer.Token) bool {
		return len(token.Source) > 0 && token.Source[0] == '"'
	})
	if err != nil {
		return nil, err
	}
	s := ""
	currIndex := 1
	for currIndex < len(currToken.Source)-1 {
		char := currToken.Source[currIndex]
		currIndex++
		if char == '\\' {
			if currIndex < len(currToken.Source)-1 {
				nextChar := currToken.Source[currIndex]
				currIndex++
				if nextChar == '\\' {
					s += "\\"
				} else if nextChar == '"' {
					s += "\""
				} else if nextChar == 't' {
					s += "\t"
				} else if nextChar == 'n' {
					s += "\n"
				} else {
					return nil, &file.Error{Location: currToken.Location, Message: "unsupported escape sequence"}
				}
			} else {
				return nil, &file.Error{Location: currToken.Location, Message: "incomplete escape sequence"}
			}
		} else {
			s += string(char)
		}
	}
	return NewStringNode(currToken.Location, s), nil
}

func (p *parser) parseLambda() (*LambdaNode, *file.Error) {
	start, err := p.consume(func(token *lexer.Token) bool { return token.Source == "lambda" })
	if err != nil {
		return nil, err
	}
	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "(" })
	if err != nil {
		return nil, err
	}

	varList := []*VariableNode{}
	for p.currIndex < len(p.tokens) {
		currToken := p.tokens[p.currIndex]
		if len(currToken.Source) > 0 && currToken.Kind == lexer.Identifier {
			node, err := p.parseVariable()
			if err != nil {
				return nil, err
			}
			varList = append(varList, node)
		} else {
			break
		}
	}
	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == ")" })
	if err != nil {
		return nil, err
	}
	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "{" })
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "}" })
	if err != nil {
		return nil, err
	}

	return NewLambdaNode(start.Location, varList, expr), nil
}

func (p *parser) parseLetrec() (*LetrecNode, *file.Error) {
	start, err := p.consume(func(token *lexer.Token) bool { return token.Source == "letrec" })
	if err != nil {
		return nil, err
	}
	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "(" })
	if err != nil {
		return nil, err
	}

	varExprList := []*LetrecVarExprItem{}
	for p.currIndex < len(p.tokens) {
		currToken := p.tokens[p.currIndex]
		if len(currToken.Source) > 0 && currToken.Kind == lexer.Identifier {
			v, err := p.parseVariable()
			if err != nil {
				return nil, err
			}
			_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "=" })
			if err != nil {
				return nil, err
			}
			e, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			varExprList = append(varExprList, &LetrecVarExprItem{
				Variable: v,
				Expr:     e,
			})
		} else {
			break
		}
	}

	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == ")" })
	if err != nil {
		return nil, err
	}

	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "{" })
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "}" })
	if err != nil {
		return nil, err
	}
	return NewLetrecNode(start.Location, varExprList, expr), nil
}

func (p *parser) parseIf() (*IfNode, *file.Error) {
	start, err := p.consume(func(token *lexer.Token) bool { return token.Source == "if" })
	if err != nil {
		return nil, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "then" })
	if err != nil {
		return nil, err
	}
	branch1, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "else" })
	if err != nil {
		return nil, err
	}
	branch2, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return NewIfNode(start.Location, cond, branch1, branch2), nil
}

func (p *parser) parseIntrinsic() (*IntrinsicNode, *file.Error) {
	currToken, err := p.consume(func(token *lexer.Token) bool {
		return len(token.Source) > 0 && token.Kind == lexer.Identifier
	})
	if err != nil {
		return nil, err
	}
	isIntrinsic := false
	for _, kw := range intrinsics {
		if kw == currToken.Source {
			isIntrinsic = true
			break
		}
	}
	if !isIntrinsic {
		return nil, &file.Error{Location: currToken.Location, Message: "incorrect intrinsic"}
	}
	return NewIntrinsicNode(currToken.Location, currToken.Source), nil
}

func (p *parser) parseVariable() (*VariableNode, *file.Error) {
	currToken, err := p.consume(func(token *lexer.Token) bool {
		return len(token.Source) > 0 && token.Kind == lexer.Identifier
	})
	if err != nil {
		return nil, err
	}

	for _, kw := range intrinsics {
		if kw == currToken.Source {
			return nil, &file.Error{Location: currToken.Location, Message: "incorrect variable name"}
		}
	}
	kind := Lexical
	if unicode.IsUpper([]rune(currToken.Source)[0]) {
		kind = Dynamic
	}
	return NewVariableNode(currToken.Location, currToken.Source, kind), nil
}

func (p *parser) parseCall() (*CallNode, *file.Error) {
	start, err := p.consume(func(token *lexer.Token) bool { return token.Source == "(" })
	if err != nil {
		return nil, err
	}

	if p.currIndex >= len(p.tokens) {
		return nil, &file.Error{Location: start.Location, Message: "incomplete token stream"}
	}
	currToken := p.tokens[p.currIndex]

	isIntrinsic := false
	for _, kw := range intrinsics {
		if kw == currToken.Source {
			isIntrinsic = true
			break
		}
	}
	var callee ExprNode
	if isIntrinsic {
		callee, err = p.parseIntrinsic()
		if err != nil {
			return nil, err
		}
	} else {
		callee, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	argList := []ExprNode{}
	for p.currIndex < len(p.tokens) {
		currToken = p.tokens[p.currIndex]
		if currToken.Source != ")" {
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			argList = append(argList, arg)
		} else {
			break
		}
	}

	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == ")" })
	if err != nil {
		return nil, err
	}

	return NewCallNode(start.Location, callee, argList), nil
}

func (p *parser) parseSequence() (*SequenceNode, *file.Error) {
	start, err := p.consume(func(token *lexer.Token) bool { return token.Source == "[" })
	if err != nil {
		return nil, err
	}
	exprList := []ExprNode{}
	for p.currIndex < len(p.tokens) {
		currToken := p.tokens[p.currIndex]
		if currToken.Source != "]" {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			exprList = append(exprList, expr)
		} else {
			break
		}
	}

	_, err = p.consume(func(token *lexer.Token) bool { return token.Source == "]" })
	if err != nil {
		return nil, err
	}

	if len(exprList) == 0 {
		return nil, &file.Error{Location: start.Location, Message: "zero-length sequence"}
	}

	return NewSequenceNode(start.Location, exprList), nil
}

func (p *parser) parseExpr() (ExprNode, *file.Error) {
	if p.currIndex >= len(p.tokens) {
		return nil, &file.Error{
			Location: file.SourceLocation{Line: -1, Col: -1},
			Message:  "incomplete token stream",
		}
	}

	currToken := p.tokens[p.currIndex]
	if len(currToken.Source) > 0 && currToken.Kind == lexer.Number {
		return p.parseNumber()
	} else if len(currToken.Source) > 0 && currToken.Kind == lexer.String {
		return p.parseString()
	} else if currToken.Source == "lambda" {
		return p.parseLambda()
	} else if currToken.Source == "letrec" {
		return p.parseLetrec()
	} else if currToken.Source == "if" {
		return p.parseIf()
	} else if len(currToken.Source) > 0 && currToken.Kind == lexer.Identifier {
		return p.parseVariable()
	} else if currToken.Source == "(" {
		return p.parseCall()
	} else if currToken.Source == "[" {
		return p.parseSequence()
	} else {
		return nil, &file.Error{Location: currToken.Location, Message: "unrecognized token"}
	}
}

func Parse(tokens []*lexer.Token) (ExprNode, *file.Error) {
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
