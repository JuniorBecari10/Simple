package ast

type Node interface {
  node()
}

type Statement interface {
  Node
  stat()
}

type Program struct {
  Statements []Statement
}

func (p Program) node() {}

type EndStatement struct {}
func (es EndStatement) stat() {}

type ErrorStatement struct {
  msg string
}

func (es ErrorStatement) stat() {}