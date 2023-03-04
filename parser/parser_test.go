package parser

import (
  "testing"
  "reflect"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
)

func TestParser(t *testing.T) {
  input := `goto :label`
  
  tokens := lexer.Lex(input)
  checkLexerErrors(t, tokens)
  
  stats := Parse(tokens)
  
  expect := []ast.Statement {
    ast.GotoStatement {
      token.Token {
        token.GotoKw,
        "goto",
        0,
      },
      ":label",
    },
    ast.EndStatement {},
  }
  
  for i, s := range stats {
    if !reflect.DeepEqual(s, expect[i]) {
      t.Fatalf("not equal. len %d.\nexpected %+v\ngot      %+v", len(stats), expect[i], s)
    }
  }
}

func checkLexerErrors(t *testing.T, tokens []token.Token) {
  for _, tk := range tokens {
    if tk.Type == token.Error {
      t.Fatalf("lexer error: %s", tk.Content)
    }
  }
}