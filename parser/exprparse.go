package parser

// TODO: implement functions
// TODO: list operations

import (
	"fmt"
	"strings"
	"log"

	"github.com/metux/go-magicdict/utils"
)

const (
	tokenEscDollar = "\\$"
	tokenRefStart  = "${"
	tokenRefEnd    = "}"
	tokenFuncStart = "$("
	tokenFuncEnd   = ")"
)

func parseRef(s []string, expr Expression, strict bool) ([]string, Expression, error) {
	sub := NewExpr()
	err := error(nil)
	end := tokenRefEnd

	for len(s) > 0 {
		head, tail := utils.ListHead(s)
		switch head {
		case tokenRefStart:
			if tail, sub, err = parseRef(tail, sub, strict); err != nil {
				return tail, sub, err
			}
		case tokenFuncStart:
			if tail, sub, err = parseFunc(tail, sub, strict); err != nil {
				return tail, sub, err
			}
		case end:
			expr.AddRef(sub)
			return tail, expr, nil
		case tokenFuncEnd:
			e := fmt.Errorf("unexpected \"%s\" in ref", head)
			if strict {
				return tail, sub, e
			}
		case tokenEscDollar:
			sub.AddLiteral("$")
		default:
			sub.AddLiteral(head)
		}
		s = tail
	}

	e := fmt.Errorf("missing token \"%s\"", end)
	if strict {
		return s, expr, e
	}
	return s, expr, nil
}

func parseFunc(s []string, expr Expression, strict bool) ([]string, Expression, error) {
	sub := NewExpr()
	err := error(nil)
	end := tokenFuncEnd

	for len(s) > 0 {
		head, tail := utils.ListHead(s)
		switch head {
		case tokenRefStart:
			if tail, sub, err = parseRef(tail, sub, strict); err != nil {
				return tail, sub, err
			}
		case tokenRefEnd:
			e := fmt.Errorf("unexpected \"%s\" in func", head)
			if strict {
				return tail, sub, e
			}
		case tokenFuncStart:
			if tail, sub, err = parseFunc(tail, sub, strict); err != nil {
				return tail, sub, err
			}
		case end:
			expr.AddFunc(sub)
			return tail, expr, nil
		case tokenEscDollar:
			sub.AddLiteral("$")
		default:
			sub.AddLiteral(head)
		}
		s = tail
	}

	e := fmt.Errorf("parse fail - missing token \"%s\"", end)
	if strict {
		return s, expr, e
	}
	return s, expr, nil
}

func ParseExpression(text string, strict bool) (Expression, error) {

	tokens := []string{tokenEscDollar, tokenRefStart, tokenRefEnd, tokenFuncStart, tokenFuncEnd}
	s := utils.SplitTokens(text, tokens)

	elems := NewExpr()
	err := error(nil)

	for len(s) > 0 {
		head, tail := utils.ListHead(s)
		switch head {
		case tokenRefStart:
			if tail, elems, err = parseRef(tail, elems, strict); err != nil {
				return elems, err
			}
		case tokenFuncStart:
			if tail, elems, err = parseFunc(tail, elems, strict); err != nil {
				return elems, err
			}
		case tokenEscDollar:
			elems.AddLiteral("$")
		/* tokenRefEnd and tokenFuncEnd will fall through here, correctly */
		default:
			elems.AddLiteral(head)
		}
		s = tail
	}
	return elems, nil
}

// FIXME: doesn't support recursive references, nor mult-element refs
// FIXME: should use the real parser and check for exactly one ref term
func ParseSimpleRefExpr(text string) (string, bool) {
	if strings.HasPrefix(text, tokenRefStart) && strings.HasSuffix(text, tokenRefEnd) {
		log.Println("expr text: ", text)
		parsed := text[len(tokenRefStart):len(text)-len(tokenRefEnd)]
		log.Println("parsed: ", parsed)
		return parsed, true
	}
	return text, false
}
