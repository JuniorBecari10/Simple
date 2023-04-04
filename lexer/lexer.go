package lexer

import (
  "strings"
  "fmt"
  "unicode"
  
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

func (this *Lexer) peekChar() byte {
  if this.cursor + 1 >= len(this.chars) {
    return 0
  }
  
  return this.chars[this.cursor + 1]
}

func (this *Lexer) NextToken() token.Token {
  if this.cursor >= len(this.chars) {
    return token.Token { token.End, "", this.cursor }
  }
  
  for (unicode.IsSpace(rune(this.char())) && this.char() != '\n') || this.char() == 0 {
    if this.cursor >= len(this.chars) {
      return token.Token { token.End, "", this.cursor }
    }

    this.advance()
  }
  
  if this.char() == '#' {
    return token.Token { token.End, "", this.cursor }
  }
  
  if this.char() == '\'' {
    pos := this.cursor
    
    this.advance()
    for this.char() != '\'' {
      this.advance()
    }
    this.advance()
    
    return token.Token { token.String, this.chars[pos + 1:this.cursor - 1], pos }
  }
  
  if this.char() == '"' {
    pos := this.cursor
    
    this.advance()
    for this.char() != '"' {
      this.advance()
    }
    this.advance()
    
    return token.Token { token.String, this.chars[pos + 1:this.cursor - 1], pos }
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
    
    for IsLetter(this.char()) || IsDigit(this.char()) {
      this.advance()
    }
    
    txt := this.chars[pos:this.cursor]
    
    if IsKeyword(txt) {
      return token.Token { token.KeyTokens[txt], txt, pos }
    }
    
    if IsType(txt) {
      return token.Token { token.TypeTokens[txt], txt, pos }
    }
    
    return token.Token { token.Identifier, txt, pos }
  }
  
  if this.char() == ':' {
    pos := this.cursor
    
    this.advance()
    
    for IsLetter(this.char()) || IsDigit(this.char()) {
      this.advance()
    }
    
    txt := this.chars[pos:this.cursor]
    
    if IsKeyword(txt) || IsType(txt) {
      return token.Token { token.Error, "Cannot use neither a keyword nor a type as a label name.", pos }
    }
    
    return token.Token { token.Label, txt, pos }
  }
  
  if this.char() == '=' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.Equals, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '!' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.Different, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '>' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.GreaterEq, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '<' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.LessEq, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '>' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Greater, string(ch), pos }
  }
  
  if this.char() == '<' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Less, string(ch), pos }
  }
  
  if this.char() == '=' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Assign, string(ch), pos }
  }
  
  if this.char() == '+' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.PlusAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '-' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.MinusAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '*' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.TimesAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '/' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.DivideAssign, this.chars[pos:pos + 2], pos }
  }

  if this.char() == '^' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.PowerAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '%' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.ModAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '&' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.AndAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '|' && this.peekChar() == '=' {
    pos := this.cursor
    
    this.advance()
    this.advance()
    
    return token.Token { token.OrAssign, this.chars[pos:pos + 2], pos }
  }
  
  if this.char() == '%' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Mod, string(ch), pos }
  }
  
  if this.char() == '&' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.And, string(ch), pos }
  }
  
  if this.char() == '|' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Or, string(ch), pos }
  }
  
  // xor will be added later
  
  if this.char() == '!' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Bang, string(ch), pos }
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

  if this.char() == '^' {
    pos := this.cursor
    ch := this.char()
    this.advance()
    
    return token.Token { token.Power, string(ch), pos }
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
  
  return token.Token { token.Error, "Unknown token: '" + string(ch) + "', char " + fmt.Sprintf("%d", ch) + ", pos " + fmt.Sprintf("%d", pos) + ".", pos }
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

func IsType(s string) bool {
  for _, t := range token.Types {
    if s == t {
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

func SplitTokens(tokens []token.Token) [][]token.Token {
  toks := make([][]token.Token, 0)
  tks := []token.Token {}
  
  for _, t := range tokens {
    if t.Type != token.NewLine && t.Type != token.End {
      tks = append(tks, t)
      continue
    }
    
    toks = append(toks, tks)
    tks = []token.Token {}
  }
  
  return toks
}

func SplitLines(s string) []string {
  lines := []string {}
  str := ""
  
  for _, c := range s {
    if c != '\n' && c != ';' {
      str += string(c)
      continue
    }
    
    lines = append(lines, str)
    str = ""
  }
  
  if str != "" {
    lines = append(lines, str)
    str = ""
  }
  
  return lines
}

func CheckErrors(tks []token.Token) []string {
  errs := []string {}
  
  for _, t := range tks {
    if t.Type == token.Error {
      errs = append(errs, t.Content)
    }
  }
  
  return errs
}