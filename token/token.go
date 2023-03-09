package token

const (
  Number       = "Number"
  String       = "String"
  Identifier   = "Identifier"
  Label        = "Label"
  
  PrintlnKw    = "PrintlnKw"
  PrintKw      = "PrintKw"
  
  InputKw      = "InputKw"
  IfKw         = "IfKw"
  GotoKw       = "GotoKw"
  ExitKw       = "ExitKw"
  
  Assign       = "Assign"
  PlusAssign   = "PlusAssign"
  MinusAssign  = "MinusAssign"
  TimesAssign  = "TimesAssign"
  DivideAssign = "DivideAssign"
  AndAssign    = "AndAssign"
  OrAssign     = "OrAssign"
  
  Plus         = "Plus"
  Minus        = "Minus"
  Times        = "Times"
  Divide       = "Divide"
  
  Power        = "Power"
  Bang         = "Bang"
  
  Equals       = "Equals"
  Different    = "Different"
  
  Greater      = "Greater"
  GreaterEq    = "GreaterEq"
  
  Less         = "Less"
  LessEq       = "LessEq"
  
  And          = "And"
  Or           = "Or"
  Xor          = "Xor"
  Not          = "Not"
  
  TrueKw       = "True"
  FalseKw      = "False"
  
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
  "true",
  "false",
  "if",
  "goto",
  "exit",
}

var Types   = []string {
  "num",
  "str",
  "bool",
}

var KeyTokens = map[string]TokenType {
  "println": PrintlnKw,
  "print":   PrintKw,
  "input":   InputKw,
  "true":    TrueKw,
  "false":   FalseKw,
  "if":      IfKw,
  "goto":    GotoKw,
  "exit":    ExitKw,
}

var TypeTokens = map[string]TokenType {
  "num":  TypeNum,
  "str":  TypeStr,
  "bool": TypeBool,
}