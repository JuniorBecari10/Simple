package parser

import (
  "testing"
  "reflect"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
)

func TestParser(t *testing.T) {
  input := `a = 'hello' + b
print 'hello' + a
`
  
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
      ast.BinNode {
        ast.StringNode {
          "'hello'",
        },
        ast.Identifier {
          token.Token {
            token.Identifier,
            "b",
            14,
          },
          "b",
        },
        "+",
      },
    },
    ast.PrintStatement {
      &token.Token {
        token.Keyword,
        "print",
        0,
      },
      ast.BinNode {
        ast.StringNode {
          "'hello'",
        },
        ast.Identifier {
          token.Token {
            token.Identifier,
            "a",
            16,
          },
          "a",
        },
        "+",
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