package lexer

import (
  "testing"
  
  "simple/token"
)

func TestLexer(t *testing.T) {
  inp := "&|!true false"
  tks := Lex(inp)
  
  res := []token.Token {
    {token.And, "&", 0},
    {token.Or, "|", 0},
    {token.Bang, "!", 0},
    {token.TrueKw, "true", 0},
    {token.FalseKw, "false", 0},
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