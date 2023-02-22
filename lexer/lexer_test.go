package lexer

import (
  "testing"
  
  "simple/token"
)

func TestLexer(t *testing.T) {
  inp := "print 'sei lá 1-24423'10 :aaa"
  tks := Lex(inp)
  
  res := []token.Token {
    {token.Keyword, "print", 0},
    {token.String, "'sei lá 1-24423'", 0},
    {token.Number, "10", 0},
    {token.Label, ":aaa", 0},
    {token.End, "", 0},
  }
  
  for i, tk := range tks {
    if tk.Type == token.Error {
      t.Fatalf("lexer error: %s", tk.Content)
    }
    
    if tk.Type != res[i].Type {
      t.Fatalf("type wrong. expected %s, got %s.", res[i].Type, tk.Type)
    }
    
    if tk.Content != res[i].Content {
      t.Fatalf("content wrong. expected %s, got %s. pos: %d", res[i].Content, tk.Content, tk.Pos)
    }
  }
}