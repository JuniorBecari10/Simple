package token

const (
  Number     = "Number"
  String     = "String"
  Identifier = "Identifier"
  Label      = "Label"
  PrintKw    = "PrintKw"
  PrintlKw   = "PrintlKw"
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
  "printl",
}

var KeyTokens = map[string]TokenType {
  "print":  PrintKw,
  "printl": PrintlKw,
}