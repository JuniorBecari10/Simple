package ast

import (
  "simple/token"
)

type Node interface {
  node()
}

type Statement interface {
  Node
  stat()
}

// Syntax: <ident> = <expression>
type VarDeclStatement struct {
  Name  Identifier
  Value ExpressionNode
}

// Syntax: print <expression>
type PrintStatement struct {
  Token       token.Token // print keyword
  BreakLine   bool
  Expression  ExpressionNode
}

type ExpressionStatement struct {
  Expression ExpressionNode
}

type EndStatement struct {}

type ErrorStatement struct {
  Msg string
}

func (vs VarDeclStatement)    stat() {}
func (ps PrintStatement)      stat() {}
func (es EndStatement)        stat() {}
func (es ErrorStatement)      stat() {}
func (es ExpressionStatement) stat() {}

func (vs VarDeclStatement)    node() {}
func (ps PrintStatement)      node() {}
func (es EndStatement)        node() {}
func (es ErrorStatement)      node() {}
func (es ExpressionStatement) node() {}

// Expressions

type ExpressionNode interface {
  exNode()
}

type Identifier struct {
  Token token.Token
  Value string
}

type NumberNode struct {
  Value float64
}

type StringNode struct {
  Value string
}

type BinNode struct {
  NodeA ExpressionNode
  NodeB ExpressionNode
  Op string
}

type PlusNode struct {
  Value ExpressionNode
}

type MinusNode struct {
  Value ExpressionNode
}

func (i Identifier) exNode() {}
func (n NumberNode) exNode() {}
func (s StringNode) exNode() {}
func (b BinNode)    exNode() {}
func (p PlusNode)   exNode() {}
func (m MinusNode)  exNode() {}