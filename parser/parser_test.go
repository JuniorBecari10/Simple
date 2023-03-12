package parser

import (
  "testing"
  "reflect"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
)

func TestParser(t *testing.T) {
  input := `exec "cat " + "file"`
  
  tokens := lexer.Lex(input)
  checkLexerErrors(t, tokens)
  
  stats := Parse(tokens)
  
  expect := []ast.Statement {
    ast.ExpressionStatement {
      ast.BinNode {
        ast.NumberNode {
          1,
        },
        ast.NumberNode {
          1,
        },
        "!=",
      },
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