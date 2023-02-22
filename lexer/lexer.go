package lexer

import (
  "strings"
  "fmt"
  
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
  
  if this.char() == '\'' {
    pos := this.cursor
    
    this.advance()
    for this.char() != '\'' {
      this.advance()
    }
    this.advance()
    
    return token.Token { token.String, this.chars[pos:this.cursor], pos }
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
  
  if this.char() == '=' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Assign, string(ch), pos }
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
  
  return token.Token { token.Error, "Unknown token: '" + string(ch) + "' char " + fmt.Sprintf("%v", ch) + ", pos " + fmt.Sprintf("%d", pos) + ".", pos }
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

func Lex(chars string) []token.Token {
  lines := strings.Split(chars, "\n")
  tks := []token.Token {}
  
  for i, line := range lines {
    l := New(line)
    
    tk := l.NextToken()
    for tk.Type != token.End {
      tks = append(tks, tk)
      tk = l.NextToken()
    }
    
    // for removing extra NewLine token
    if i == len(lines) - 1 {
      break
    }
    
    tks = append(tks, token.Token { token.NewLine, "", len(line) })
  }
  
  tks = append(tks, token.Token { token.End, "", 0 })
  
  return tks
}

func SplitLines(tokens []token.Token) [][]token.Token {
  toks := make([][]token.Token, 0)
  tks := []token.Token {}
  
  for _, t := range tokens {
    if t.Type != token.NewLine {
      tks = append(tks, t)
      continue
    }
    
    toks = append(toks, tks)
    tks = []token.Token {}
  }
  
  return toks
}