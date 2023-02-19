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
  
  return token.Token { token.End, "", this.cursor }
}

// -- Helper -- //

func IsDigit(b byte) bool {
  return (b >= '0' && b <= '9') || b == '.'
}

func IsLetter(b byte) bool {
  return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}

// ---

func Lex(chars string) []token.Token {
  l := New(chars)
  tks := []token.Token {}
  
  tk := l.NextToken()
  for tk.Type != token.End {
    tks = append(tks, tk)
    tk = l.NextToken()
  }
  
  return tks
}