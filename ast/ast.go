package ast

import (
  "simple/token"
)

const (
  TrueType  = "True"
  FalseType = "False"
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

// Syntax: <ident> +|-|*|/= <expression>
type OperationStatement struct {
  Name  Identifier
  Value ExpressionNode
  Op    string
}

// Syntax: print <expression>
type PrintStatement struct {
  Code string
  Line int
  
  BreakLine   bool
  Expression  ExpressionNode
}

// Syntax: goto :<label>
type GotoStatement struct {
  Code string
  Line int
  
  Label string
}

// Syntax: if <expression> goto :<label>
type IfStatement struct {
  Code string
  Line int
  
  Expression ExpressionNode
  Label      string
}

type ExpressionStatement struct {
  Code string
  Line int
  
  Expression ExpressionNode
}

type LabelStatement struct {
  Code string
  Line int
  
  Name string
}

type EndStatement struct {}

type ErrorStatement struct {
  Code string
  Line int
  
  Msg string
}

func (vs VarDeclStatement)    stat() {}
func (os OperationStatement)  stat() {}
func (ps PrintStatement)      stat() {}
func (es EndStatement)        stat() {}
func (ls LabelStatement)      stat() {}
func (gs GotoStatement)       stat() {}
func (is IfStatement)         stat() {}
func (es ExpressionStatement) stat() {}
func (es ErrorStatement)      stat() {}

func (vs VarDeclStatement)    node() {}
func (os OperationStatement)  node() {}
func (ps PrintStatement)      node() {}
func (es EndStatement)        node() {}
func (ls LabelStatement)      node() {}
func (gs GotoStatement)       node() {}
func (is IfStatement)         node() {}
func (es ExpressionStatement) node() {}
func (es ErrorStatement)      node() {}

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

type InputNode struct {
  Type string
}

type BoolNode struct {
  Type string
}

type FactorialNode struct {
  Node ExpressionNode
}

func (i Identifier)    exNode() {}
func (n NumberNode)    exNode() {}
func (s StringNode)    exNode() {}
func (b BinNode)       exNode() {}
func (p PlusNode)      exNode() {}
func (m MinusNode)     exNode() {}
func (i InputNode)     exNode() {}
func (b BoolNode)      exNode() {}
func (f FactorialNode) exNode() {}