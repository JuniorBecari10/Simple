package token

const (
  Number     = "Number"
  Identifier = "Identifier"
  Keyword    = "Keyword"
  Assign     = "Assign"
  Plus       = "Plus"
  Minus      = "Minus"
  Times      = "Times"
  Divide     = "Divide"
  LParen     = "LParen"
  RParen     = "RParen"
  NewLine    = "NewLine"
  End        = "End"
  Error      = "Error"
)

type TokenType string

type Token struct {
  Type    TokenType
  Content string
  Pos     int
}

var Keywords = []string {
  "print",
}