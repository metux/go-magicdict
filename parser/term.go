package parser

import (
    "fmt"
)

type TermType byte

const (
    TermLiteral = TermType(1)
    TermRef     = TermType(2)
    TermFunc    = TermType(3)
)

func (tt TermType) String() string {
    switch tt {
        case TermLiteral: return "literal"
        case TermRef:     return "ref"
        case TermFunc:    return "func"
    }
    return fmt.Sprintf("%d", tt)
}

type Term struct {
    Type        TermType
    Literal     string
    Expr        Expression
}

type Expression [] Term

func NewExpr() Expression {
    return Expression(make(Expression, 0))
}

func (e * Expression) AddLiteral(s string) {
    *e = append(*e, Term { Type: TermLiteral, Literal: s })
}

func (e * Expression) AddRef(sub Expression) {
    *e = append(*e, Term { Type: TermRef, Expr: sub })
}

func (e * Expression) AddFunc(sub Expression) {
    *e = append(*e, Term { Type: TermFunc, Expr: sub })
}
