package lexer

import (
  "testing"
  
  "simple/token"
)

func TestNumber(t *testing.T) {
  inp := "10 20.2"
  tks := Lex(inp)
  
  res := []token.Token {
    {token.Number, "10", 1},
    {token.Number, "20.2", 1},
  }
  
  for i, tk := range tks {
    if tk.Type == token.Error {
      t.Fatalf("lexer error: %s", tk.Content)
    }
    
    if tk.Type != res[i].Type {
      t.Fatalf("type wrong. expected %s, got %s.", res[i].Type, tk.Type)
    }
    
    if tk.Content != res[i].Content {
      t.Fatalf("content wrong. expected %s, got %s.", res[i].Content, tk.Content)
    }
  }
}