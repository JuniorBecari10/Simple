package parser

import (
  "strconv"
  
  "simple/token"
  "simple/lexer"
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

func (this *Parser) token() token.Token {
  if this.cursor >= len(this.tokens) {
    return token.Token { token.Error, "Exceeded line length.", this.cursor }
  }
  
  return this.tokens[this.cursor]
}

func (this *Parser) nextStatement() ast.Statement {
  if this.token().Type == token.End {
    return ast.EndStatement {}
  }
  
  if this.token().Type == token.Error {
    return ast.ErrorStatement { this.token().Content }
  }
  
  if len(this.tokens) >= 2 && this.token().Type == token.Identifier && this.tokens[this.cursor + 1].Type == token.Assign {
    return this.parseVarDeclStatement()
  }
  
  if len(this.tokens) >= 1 && (this.token().Type == token.PrintKw || this.token().Type == token.PrintlKw) {
    return this.parsePrintStatement()
  }
  
  return ast.ExpressionStatement { this.parseExpression() }
}

func (this *Parser) parseVarDeclStatement() ast.Statement {
  stat := ast.VarDeclStatement {}
  id := ast.Identifier { Token: this.token(), Value: this.token().Content }
  
  stat.Name = id
  this.advance()
  
  if this.token().Type != token.Assign {
    return ast.ErrorStatement { "Syntax error when declaring a variable. Examples: a = 10; message = 'Hello'." }
  }
  
  this.advance()
  stat.Value = this.parseExpression()
  
  return stat
}

func (this *Parser) parsePrintStatement() ast.Statement {
  stat := ast.PrintStatement {}
  
  tk := this.token()
  this.advance()
  expr := this.parseExpression()
  
  if tk.Type == token.Error || expr == nil {
    return ast.ErrorStatement { "Syntax error in a print statement. Examples: print 'Hello World'; print 1 + 1." }
  }
  
  stat.Token = tk
  stat.BreakLine = tk.Type != token.PrintlKw
  stat.Expression = expr
  
  this.advance()
  
  return stat
}

func (this *Parser) parseExpression() ast.ExpressionNode {
  return this.exp()
}

func (this *Parser) exp() ast.ExpressionNode {
  if this.token().Type == token.Error {
    return nil
  }
  
  res := this.term()
  
  for this.token().Type != token.Error && (this.token().Type == token.Plus || this.token().Type == token.Minus) {
    if this.token().Type == token.Plus {
      this.advance()
      
      res = ast.BinNode { res, this.term(), "+" }
    } else if this.token().Type == token.Minus {
      this.advance()
      
      res = ast.BinNode { res, this.term(), "-" }
    }
  }
  
  return res
}

func (this *Parser) term() ast.ExpressionNode {
  if this.token().Type == token.Error {
    return nil
  }
  
  res := this.factor()
  
  for this.token().Type != token.Error && (this.token().Type == token.Times || this.token().Type == token.Divide) {
    if this.token().Type == token.Times {
      this.advance()
      
      res = ast.BinNode { res, this.term(), "*" }
    } else if this.token().Type == token.Divide {
      this.advance()
      
      res = ast.BinNode { res, this.term(), "/" }
    }
  }
  
  return res
}

func (this *Parser) factor() ast.ExpressionNode {
  tk := this.token()
  
  if tk.Type == token.LParen {
    this.advance()
    res := this.exp()
    
    this.advance()
    return res
  }
  
  if tk.Type == token.Plus {
    this.advance()
    
    return ast.PlusNode { this.factor() }
  }
  
  if tk.Type == token.Minus {
    this.advance()
    
    return ast.MinusNode { this.factor() }
  }
  
  if tk.Type == token.Number {
    this.advance()
    
    value, _ := strconv.ParseFloat(tk.Content, 64)
    return ast.NumberNode { value }
  }
  
  if tk.Type == token.String {
    this.advance()
    
    return ast.StringNode { tk.Content }
  }
  
  if tk.Type == token.Identifier {
    this.advance()
    
    return ast.Identifier { tk, tk.Content }
  }
  
  return nil
}

func Parse(tokens []token.Token) []ast.Statement {
  lines := lexer.SplitLines(tokens)
  stats := []ast.Statement {}
  
  for _, l := range lines {
    p := New(l)
    
    st := p.nextStatement()
    _, ok := st.(ast.EndStatement)
    
    if ok {
      break
    }
    
    stats = append(stats, st)
  }
  
  return stats
}