package parser

import (
  "fmt"
  "testing"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
)

func TestParser(t *testing.T) {
  input := "a = 10 + 1"
  
  tokens := lexer.Lex(input)
  checkLexerErrors(t, tokens)
  
  stats := Parse(tokens)
  
  expect := []ast.Statement {
    ast.VarDeclStatement {
      &ast.Identifier {
        token.Token {
          token.Identifier,
          "a",
          0,
        },
        "a",
      },
      ast.NumberNode {
        10,
      },
    },
    ast.EndStatement {},
  }
  
  for i, s := range stats {
    str1 := fmt.Sprintf("%+v", s)
    str2 := fmt.Sprintf("%+v", expect[i])
    
    if str1 != str2 {
      t.Fatalf("error. expect %s, got %s", str1, str2)
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