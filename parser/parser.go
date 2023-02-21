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
  
  if len(this.tokens) >= 2 && this.token().Type == token.Identifier && this.tokens[this.cursor + 1].Type == token.Assign {
    return this.parseVarDeclStatement()
  }
  
  return ast.ErrorStatement { "Unknown statement. tokens: " + this.tokens }
}

func (this *Parser) parseVarDeclStatement() ast.Statement {
  stat := ast.VarDeclStatement {}
  id := ast.Identifier { Token: this.token(), Name: this.Token.Content }
  
  stat.id = id
  this.advance()
  
  if this.token().Type != token.Assign {
    return ErrorStatement { "Syntax error when declaring a variable. Examples: a = 10; message = 'Hello'." }
  }
  
  this.advance()
  stat.Value = parseExpression()
  
  return stat
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