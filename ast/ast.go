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
  Name *Identifier
  Value ExpressionNode
}

type EndStatement struct {}

type ErrorStatement struct {
  msg string
}

func (vds VarDeclStatement) stat() {}
func (es EndStatement)      stat() {}
func (es ErrorStatement)    stat() {}

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
func (b BinNode)    exNode() {}
func (p PlusNode)   exNode() {}
func (m MinusNode)  exNode() {}