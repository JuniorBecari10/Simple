package lexer

import (
  "strings"
  
  "simple/token"
)

type Lexer struct {
  chars  string
  cursor int
}

func New(chars string) *Lexer {
  return &Lexer { chars, 0 }
}

func (this *Lexer) advance() {
  this.cursor++
}

func (this *Lexer) char() byte {
  if this.cursor >= len(this.chars) {
    return 0
  }
  
  return this.chars[this.cursor]
}

func (this *Lexer) NextToken() token.Token {
  if this.cursor >= len(this.chars) {
    return token.Token { token.End, "", this.cursor }
  }
  
  for this.char() == ' ' {
    this.advance()
  }
  
  if IsDigit(this.char()) {
    pos := this.cursor
    
    for IsDigit(this.char()) {
      this.advance()
    }
    
    nb := this.chars[pos:this.cursor]
    
    if strings.Count(nb, ".") > 1 {
      return token.Token { token.Error, "A Number cannot have multiple dots!", pos }
    }
    
    return token.Token { token.Number, nb, pos }
  }
  
  if IsLetter(this.char()) {
    pos := this.cursor
    
    for IsLetter(this.char()) {
      this.advance()
    }
    
    txt := this.chars[pos:this.cursor]
    
    if IsKeyword(txt) {
      return token.Token { token.Keyword, txt, pos }
    }
    
    return token.Token { token.Identifier, txt, pos }
  }
  
  if this.char() == '+' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Plus, string(ch), pos }
  }
  
  if this.char() == '-' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Minus, string(ch), pos }
  }
  
  if this.char() == '*' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Times, string(ch), pos }
  }
  
  if this.char() == '/' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Divide, string(ch), pos }
  }
  
  if this.char() == '(' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.LParen, string(ch), pos }
  }
  
  if this.char() == ')' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.RParen, string(ch), pos }
  }
  
  pos := this.cursor
  ch := this.char()
  this.advance()
  
  return token.Token { token.Error, "Unknown token: '" + string(ch) + "'.", pos }
}

// -- Helper -- //

func IsDigit(b byte) bool {
  return (b >= '0' && b <= '9') || b == '.'
}

func IsLetter(b byte) bool {
  return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}

func IsKeyword(s string) bool {
  for _, k := range token.Keywords {
    if s == k {
      return true
    }
  }
  
  return false
}

// ---

func Lex(chars string) [][]token.Token {
  lines := strings.Split(chars, "\n")
  tks := [][]token.Token {}
  lineTks := []token.Token {}
  
  for _, line := range lines {
    l := New(line)
    lineTks := []token.Token {}
    
    tk := l.NextToken() 
    for tk.Type != token.End {
      lineTks = append(lineTks, tk)
      tk = l.NextToken()
    }
    
    tks = append(tks, lineTks)
  }
  
  return tks
}