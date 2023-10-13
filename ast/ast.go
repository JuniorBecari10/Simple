package ast

import (
	"csimple/token"
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
  Token       token.Token // print keyword
  BreakLine   bool
  Expression  ExpressionNode
}

// Syntax: goto :<label>
type GotoStatement struct {
  Token token.Token // goto keyword
  Label string
}

// Syntax: if <expression> goto :<label>
type IfStatement struct {
  Token      token.Token // goto keyword
  Expression ExpressionNode
  Label      string
}

type ExpressionStatement struct {
  Expression ExpressionNode
}

type LabelStatement struct {
  Name string
}

type ExitStatement struct {
  Code ExpressionNode
}

type RetStatement struct {}

type EndStatement struct {}

type ErrorStatement struct {
  Msg string
}

func (vs VarDeclStatement)    stat() {}
func (os OperationStatement)  stat() {}
func (ps PrintStatement)      stat() {}
func (es EndStatement)        stat() {}
func (ls LabelStatement)      stat() {}
func (es ExitStatement)       stat() {}
func (rs RetStatement)        stat() {}
func (gs GotoStatement)       stat() {}
func (is IfStatement)         stat() {}
func (es ExpressionStatement) stat() {}
func (es ErrorStatement)      stat() {}

func (vs VarDeclStatement)    node() {}
func (os OperationStatement)  node() {}
func (ps PrintStatement)      node() {}
func (es EndStatement)        node() {}
func (ls LabelStatement)      node() {}
func (es ExitStatement)       node() {}
func (rs RetStatement)        node() {}
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

type ExecNode struct {
  Command ExpressionNode
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
func (e ExecNode)      exNode() {}