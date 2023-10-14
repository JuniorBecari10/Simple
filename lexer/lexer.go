package lexer

import (
	"csimple/token"
)

type Lexer struct {
  input  string

  start  int
  cursor int

  startPos token.Position
  pos      token.Position
}

func New(input string) Lexer {
  return Lexer {
    input: input,
    
    start: 0,
    cursor: 0,

    startPos: token.Position {
      Line: 0,
      Col: 0,
    },

    pos: token.Position {
      Line: 0,
      Col: 0,
    },
  }
}

func (l *Lexer) Lex() []token.Token {
  tokens := []token.Token {}

}
