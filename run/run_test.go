package run

import (
  "testing"
  
  "simple/token"
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
        ast.StringNode {
          "' world'",
        },
        "+",
      },
    },
    ast.PrintStatement {
      token.Token {
        token.PrintKw,
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
  
  Run(stats)
  
  if Variables["a"] != "'hello'' world'" {
    t.Fatalf("not as expected: expected: %s, got %q", "'hello'' world'", Variables["a"])
  }
}