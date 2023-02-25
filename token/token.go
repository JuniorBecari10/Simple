package token

const (
  Number       = "Number"
  String       = "String"
  Identifier   = "Identifier"
  Label        = "Label"
  
  PrintlnKw    = "PrintlnKw"
  PrintKw      = "PrintKw"
  InputKw      = "InputKw"
  
  Assign       = "Assign"
  PlusAssign   = "PlusAssign"
  MinusAssign  = "MinusAssign"
  TimesAssign  = "TimesAssign"
  DivideAssign = "DivideAssign"
  
  Plus         = "Plus"
  Minus        = "Minus"
  Times        = "Times"
  Divide       = "Divide"
  
  LParen       = "LParen"
  RParen       = "RParen"
  
  NewLine      = "NewLine"
  End          = "End"
  Error        = "Error"
  
  TypeNum      = "Num"
  TypeStr      = "Str"
  TypeBool     = "Bool"
)

type TokenType string

type Token struct {
  Type    TokenType
  Content string
  Pos     int
}

var Keywords   = []string {
  "println",
  "print",
  "input",
}

var Types   = []string {
  "num",
  "str",
  "bool",
}

var KeyTokens   = map[string]TokenType {
  "println": PrintlnKw,
  "print":   PrintKw,
  "input":   InputKw,
}

var TypeTokens   = map[string]TokenType {
  "num":  TypeNum,
  "str":  TypeStr,
  "bool": TypeBool,
}