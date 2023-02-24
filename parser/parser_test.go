package parser

import (
  "testing"
  "reflect"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
)

func TestParser(t *testing.T) {
  input := `a = input`
  
  tokens := lexer.Lex(input)
  checkLexerErrors(t, tokens)
  
  stats := Parse(tokens)
  
  expect := []ast.Statement {
    ast.VarDeclStatement {
      ast.Identifier {
        token.Token {
          token.Identifier,
          "a",
          0,
        },
        "a",
      },
      ast.InputNode {
        "",
      },
    },
    ast.EndStatement {},
  }
  
  for i, s := range stats {
    if !reflect.DeepEqual(s, expect[i]) {
      t.Fatalf("not equal. len %d.\nexpected %+v\ngot %+v", len(stats), expect[i], s)
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