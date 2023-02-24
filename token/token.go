package token

const (
  Number     = "Number"
  String     = "String"
  Identifier = "Identifier"
  Label      = "Label"
  PrintlnKw  = "PrintlnKw"
  PrintKw    = "PrintKw"
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
  "println",
  "print",
}

var KeyTokens = map[string]TokenType {
  "println": PrintlnKw,
  "print":   PrintKw,
}