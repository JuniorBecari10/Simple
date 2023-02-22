package run

import (
  "testing"
  
  "simple/ast"
)

func TestRun(t *testing.T) {
  stats := []ast.Statement {
    ast.VarDeclStatement {
      ast.Identifier {
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
      token.Token {
        token.PrintlKw,
        "printl",
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
            17,
          },
          "a",
        },
        "+",
      },
    },
    ast.EndStatement {},
  }
  
  ret := Run(stats)
  
  
}