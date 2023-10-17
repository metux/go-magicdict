package macro

import (
	"strings"

	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
	"github.com/metux/go-magicdict/parser"
)

func resolveRef(ent parser.Term, v api.Entry, root api.Entry) (api.Entry, error) {
	str := ""
	for _, v := range ent.Expr {
		switch v.Type {
		case parser.TermLiteral:
			str = str + v.Literal
		default:
			return nil, api.ErrUnknownEntryType
		}
	}

	if strings.HasPrefix(str, "@@") {
		root = v
	}

	return root.Get(api.Key(str))
}

func ProcessVars(v api.Entry, root api.Entry) (api.Entry, error) {

	if !v.IsScalar() || !v.IsConst() {
		return v, nil
	}

	s := v.String()

	parsed, err := parser.ParseExpression(s, true)

	if err != nil {
		return nil, err
	}

	if len(parsed) == 0 {
		return nil, nil
	}

	// only one element ?
	if len(parsed) == 1 {
		switch parsed[0].Type {
		case parser.TermLiteral:
			return v, nil
		case parser.TermRef:
			return resolveRef(parsed[0], v, root)
		default:
			return nil, api.ErrUnknownEntryType
		}
	}

	// FIXME: implement functions
	// FIXME: merging lists:
	// if we have lists and literals, then merge the lists and add the literals as inviduals
	// FIXME: implement merging dicts
	// this only concat's strings
	retstr := ""
	for _, y := range parsed {
		switch y.Type {
		case parser.TermRef:
			if val, err := resolveRef(y, v, root); err == nil {
				// nil is treated as "" here
				if val != nil {
					retstr = retstr + val.String()
				}
			} else {
				return nil, err
			}
		case parser.TermLiteral:
			retstr = retstr + y.Literal
		default:
			return nil, api.ErrUnknownEntryType
		}
	}

	return core.NewScalarStr(retstr), nil
}
