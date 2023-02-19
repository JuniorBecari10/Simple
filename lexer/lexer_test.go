package lexer

import (
  "testing"
  
  "simple/token"
)

func TestLexer(t *testing.T) {
  inp := "10 20.2 ola_+-*/()print"
  tks := Lex(inp)
  
  res := []token.Token {
    {token.Number, "10", 0},
    {token.Number, "20.2", 0},
    {token.Identifier, "ola_", 0},
    {token.Plus, "+", 0},
    {token.Minus, "-", 0},
    {token.Times, "*", 0},
    {token.Divide, "/", 0},
    {token.LParen, "(", 0},
    {token.RParen, ")", 0},
    {token.Keyword, "print", 0},
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