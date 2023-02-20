package parser

import (
  "testing"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
)

func TestParser(t *testing.T) {
  input = ""
  
  tokens := lexer.Lex(input)
  checkLexerErrors(tokens)
  
  stats := Parse(tokens)
  
  expect := []ast.Statement {
    EndStatement {},
  }
  
  for _, st := range stats {
    if st != expect[i] {
      t.Fatalf("parser error: expected type %T, got %T.", st, expect[i])
    }
  }
  
}

func checkLexerErrors(tokens []token.Token) {
  for _, tk := range tokens {
    if tk.Type == token.Error {
      t.Fatalf("lexer error: %s", tk.Content)
    }
  }
}