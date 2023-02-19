package token

const (
  Number     = "Number"
  Identifier = "Identifier"
  Plus       = "Plus"
  Minus      = "Minus"
  Times      = "Times"
  Divide     = "Divide"
  LParen     = "LParen"
  RParen     = "RParen"
)

type TokenType string

type Token struct {
  Type    TokenType
  Content string
  Pos     int
}