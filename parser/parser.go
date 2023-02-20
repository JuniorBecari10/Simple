package parser

import (
  "simple/token"
  "simple/ast"
)

type Parser struct {
  tokens []token.Token
  cursor int
}

func New(tokens []token.Token) *Parser {
  return &Parser { tokens, 0 }
}

func (this *Parser) advance() {
  this.cursor++
}

func (this *Parser) token() {
  return this.tokens[this.cursor]
}

func (this *Parser) nextStatement() ast.Statement {
  if this.token().Type == token.End {
    return ast.EndStatement {}
  }
  
  return ast.ErrorStatement { "Unknown statement. tokens: " + this.tokens }
}

func Parse(tokens []token.Token) []ast.Statement {
  p := New(tokens)
  stats := []ast.Statement {}
  
  st := p.nextStatement()
  _, ok := st.(ast.EndStatement)
  for !ok {
    stats = append(stats, st)
    
    st := p.nextStatement()
    _, ok := st.(ast.EndStatement)
  }
  
  return stats
}