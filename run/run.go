package run

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "bufio"
  "reflect"
  "math"
  "bytes"
  "strings"
  "runtime"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
  "simple/parser"
)

type Any interface {}

type Label struct {
  Name string
  Line int
}

type Value struct {
  Value Any
}

type CodeLine struct {
  Line int
  Statement ast.Statement
}

var Error bool = false
var PC int = 0

var Line int
var LineCode string

var Labels []Label
var Stack []int

var Variables = map[string]Any {}
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

func ShowError(msg, hint string) {
  fmt.Println("\n-------------\n")
  
  fmt.Println("ERROR: On line " + strconv.Itoa(Line + 1) + ".")
  fmt.Println("\n" + msg)
  
  fmt.Println()
  
  if Line > 0 {
    fmt.Printf("%d |\n", Line)
  }
  
  fmt.Printf("%d | %s\n", Line + 1, LineCode)
  fmt.Printf("%d |\n\n", Line + 2)
  
  if hint != "" {
    fmt.Println(hint)
    fmt.Println()
  }
  
  Error = true
}

func ShowWarning(msg, hint string) {
  fmt.Println("\n-------------\n")
  
  fmt.Println("WARNING: On line " + strconv.Itoa(Line + 1) + ".")
  fmt.Println("\n" + msg)
  
  fmt.Println()
  
  if Line > 0 {
    fmt.Printf("%d |\n", Line)
  }
  
  fmt.Printf("%d | %s\n", Line + 1, LineCode)
  fmt.Printf("%d |\n\n", Line + 2)
  
  if hint != "" {
    fmt.Println(hint)
    fmt.Println()
  }
}

func RunCode(code string) {
  lines := strings.Split(strings.TrimSpace(code), "\n")
  codeLines := make([]CodeLine, 0)

  for i, _ := range lines {
    tokens := lexer.Lex(lines[i])
    errs := lexer.CheckErrors(tokens)
    
    if len(errs) > 0 {
      for _, e := range errs {
        // todo: add arrow ^ in hint, getting the position
        LineCode = lines[i]
        ShowError("Error in lexer: " + e, "Consider removing them.")
      }
      
      return
    }

    stats := parser.Parse(tokens)
    errs = parser.CheckErrors(stats)
    
    if len(errs) > 0 {
      for _, e := range errs {
        LineCode = lines[i]
        ShowError("Error in parser: " + e, "")
      }
      
      return
    }

    if len(stats) == 0 {
      continue
    }

    codeLines = append(codeLines, CodeLine { i, stats[0] })
  }
  
  for i, c := range codeLines {
    if ls, ok := c.Statement.(ast.LabelStatement); ok {
      Labels = append(Labels, Label { ":" + ls.Name, i })
    }
  }

  PC = 0

  for PC < len(codeLines) || !Error {
    if PC >= len(codeLines) || Error {
      break
    }

    l := codeLines[PC]

    Line = l.Line
    LineCode = lines[l.Line]

    if _, ok := l.Statement.(ast.LabelStatement); ok {
      PC++
      continue
    }

    RunStat(l.Statement, false)
    PC++
  }
}

func Run(stats []ast.Statement, line int, lineCode string, repl bool) []Any {
  if len(stats) == 0 {
    return nil
  }
  
  ret := []Any {}
  
  Line = line
  LineCode = lineCode
  
  PC = 0
  for PC < len(stats) {
    if err, ok := stats[PC].(ast.ErrorStatement); ok {
      ShowError("Error in parser: " + err.Msg, "Fix it.")
      
      return nil
    }
    
    if _, ok := stats[PC].(ast.LabelStatement); ok {
      PC++
      continue
    }
    
    ret = append(ret, RunStat(stats[PC], repl))
    PC++
  }
  
  return ret
}

func GetTokens(code string) []token.Token {
  tokens := lexer.Lex(code)
  errs := lexer.CheckErrors(tokens)
  
  if len(errs) > 0 {
    for _, e := range errs {
      // todo: add arrow ^ in hint, getting the position
      ShowError("Error in lexer: " + e, "")
    }
    
    return nil
  }
  
  return tokens
}

func GetStatements(code string) []ast.Statement {
  tokens := lexer.Lex(code)
  errs := lexer.CheckErrors(tokens)
  
  if len(errs) > 0 {
    for _, e := range errs {
      ShowError("Error in lexer: " + e, "")
    }
    
    return nil
  }
  
  stats := parser.Parse(tokens)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for _, e := range errs {
      ShowError("Error in parser: " + e, "")
    }
    
    return nil
  }
  
  return stats
}

func RunStat(stat ast.Statement, repl bool) Any {
  fn := GetStatFunc(stat)

  if _, ok := stat.(ast.LabelStatement); ok && repl {
    ShowError("You cannot declare labels in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if _, ok := stat.(ast.GotoStatement); ok && repl {
    ShowError("You cannot declare goto statements in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if _, ok := stat.(ast.IfStatement); ok && repl {
    ShowError("You cannot declare if statements in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if _, ok := stat.(ast.RetStatement); ok && repl {
    ShowError("You cannot declare ret statements in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if fn == nil {
    ShowError("Unknown statement: " + fmt.Sprintf("%T", stat), "Verify if you typed correctly.")
    return nil
  }
  
  return fn(stat)
}

func GetStatFunc(st ast.Statement) func(ast.Statement) Any {
  switch st.(type) {
    case ast.VarDeclStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.VarDeclStatement)
        
        vl := SolveExpression(s.Value)

        if vl == nil {
          return nil
        }

        Variables[s.Name.Value] = vl
        
        return vl
      }
    
    case ast.OperationStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.OperationStatement)

        vl := SolveExpression(s.Value)
        _, ok := Variables[s.Name.Value]

        if vl == nil {
          return nil
        }

        if !ok {
          ShowError("The variable " + s.Name.Value + " doesn't exist.", "Create one declaring it, like: a = 10, b = 'Hello', c = true.")
          return nil
        }

        switch s.Op {
          case "+":
            vl := Sum(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "-":
            vl := Sub(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "*":
            vl := Mul(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "/":
            vl := Div(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl

          case "^":
            vl := Pow(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "%":
            vl := Mod(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "&":
            vl := And(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "|":
            vl := Or(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
        }

        return vl
      }
    
    case ast.PrintStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.PrintStatement)
        
        exp := SolveExpression(s.Expression)

        if exp == nil {
          return nil
        }

        fmt.Print(exp)
        
        if s.BreakLine {
          fmt.Println()
        }
        
        return exp
      }
    
    case ast.ExpressionStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.ExpressionStatement)
        
        return SolveExpression(s.Expression)
      }
    
    case ast.GotoStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.GotoStatement)
        
        label := s.Label
        
        for _, l := range Labels {
          if l.Name == label {
            Stack = append(Stack, PC)
            PC = l.Line - 1
            return ""
          }
        }
        
        ShowError("Couldn't find label '" + label + "'.", "Verify if you typed the name correctly.")
        return nil
      }
    
    case ast.IfStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.IfStatement)
        
        res := SolveExpression(s.Expression)
        
        label := s.Label
        pc := 0
        
        for _, l := range Labels {
          if l.Name == label {
            pc = l.Line
            
            vl, ok := res.(bool)
            if ok {
              if vl {
                Stack = append(Stack, PC)
                PC = pc - 1
                return ""
              }
              
              // in case returning false, also return and don't print the error
              return ""
            }
            
            ShowError("Cannot use non-boolean expressions inside an if statement.", "You must use only boolean expressions.")
            return nil
          }
        }
        
        ShowError("Couldn't find label '" + label + "'.", "Verify if you typed the name correctly.")
        return nil
      }
    
    case ast.ExitStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.ExitStatement)
        
        code, ok := SolveExpression(s.Code).(float64)
        i := int(code)
        
        if !ok {
          ShowError("The exit code provided must be a positive integer.", "Examples: exit 0, exit 1 + 1, exit a + b.")
          return nil
        }
        
        if i < 0 {
          i = 0
        }
        
        os.Exit(i)
        
        return nil
      }
    
    case ast.RetStatement:
      return func(st ast.Statement) Any {
        if len(Stack) == 0 {
          ShowError("Cannot return in call stack because it's empty.", "The call stack is empty.")
          return nil
        }
        
        pc := Stack[len(Stack) - 1]
        Stack = Stack[:len(Stack) - 1]
        
        PC = pc
        
        return pc
      }
    
    default:
      return nil
  }
}

// ---

func SolveExpression(ex ast.ExpressionNode) Any {
  fn := GetExprFunc(ex)
  
  if fn == nil {
    if ex == nil {
      ShowError("The expression is incomplete.", "Certify that you completed it correctly.")
      return nil
    }
    
    ShowError("Couldn't get function to solve this expression: " + fmt.Sprintf("%q", ex), "This happens when you use an operator the wrong way or the operator isn't supported.")
    return nil
  }
  
  return fn(ex)
}

func GetExprFunc(ex ast.ExpressionNode) func(ast.ExpressionNode) Any {
  switch ex.(type) {
    case ast.Identifier:
      return func(ex ast.ExpressionNode) Any {
        value, ok := Variables[ex.(ast.Identifier).Value]
        
        if !ok {
          ShowError("Variable '" + ex.(ast.Identifier).Value + "' doesn't exist.", "Verify if you typed the name correctly.")
          return nil
        }
        
        return value
      }
    
    case ast.NumberNode:
      return func(ex ast.ExpressionNode) Any {
        return ex.(ast.NumberNode).Value
      }
    
    case ast.StringNode:
      return func(ex ast.ExpressionNode) Any {
        return ex.(ast.StringNode).Value
      }
    
    case ast.PlusNode:
      return func(ex ast.ExpressionNode) Any {
        nb, _ := SolveExpression(ex.(ast.PlusNode).Value).(float64)
        
        return nb
      }
    
    case ast.MinusNode:
      return func(ex ast.ExpressionNode) Any {
        nb, ok := SolveExpression(ex.(ast.MinusNode).Value).(float64)
        
        if !ok {
          ShowError("You can only use numbers with the operator '-'.", "Examples: -10, -25.5, -a.")
          return nil
        }
        
        return -nb
      }
    
    case ast.BinNode:
      return func(ex ast.ExpressionNode) Any {
        bin := ex.(ast.BinNode)
        
        v1, v2 := SolveExpression(bin.NodeA), SolveExpression(bin.NodeB)
        
        switch bin.Op {
          case "+":
            return Sum(v1, v2)
          
          case "-":
            return Sub(v1, v2)
          
          case "*":
            return Mul(v1, v2)
          
          case "/":
            return Div(v1, v2)
          
          case "^":
            return Pow(v1, v2)

          case "%":
            return Mod(v1, v2)
          
          case "&":
            return And(v1, v2)
          
          case "|":
            return Or(v1, v2)
          
          case "==":
            return Eq(v1, v2)
          
          case "!=":
            return Diff(v1, v2)
          
          case ">":
            return Greater(v1, v2)
          
          case ">=":
            return GreaterEq(v1, v2)
          
          case "<":
            return Less(v1, v2)
          
          case "<=":
            return LessEq(v1, v2)
          
          default:
            ShowError("Unknown operation: " + bin.Op, "")
            return nil
        }
      }
    
    case ast.InputNode:
      return func(exp ast.ExpressionNode) Any {
        inp := exp.(ast.InputNode)
        vl := ""
        
        for {
          scanner.Scan()
          vl = scanner.Text()
          
          if inp.Type == "" {
            return vl
          }
          
          value, err := strconv.ParseFloat(vl, 64)
          
          if inp.Type == token.TypeNum {
            if err != nil {
              fmt.Println("Please enter a valid num!")
              continue
            }
            
            return value
          }
          
          if inp.Type == token.TypeStr {
            if err != nil {
              
              if vl == "true" || vl == "false" {
                fmt.Println("Please enter a valid str!")
                continue
              }
              
              return vl
            }
            
            fmt.Println("Please enter a valid str!")
            continue
          }
          
          if inp.Type == token.TypeBool {
            if vl == "true" {
              return true
            }
            
            if vl == "false" {
              return false
            }
            
            fmt.Println("Please enter a valid bool!")
            continue
          }
          
          // fun fact: this error will never happen
          ShowError("Unknown type used on input expressions.", "Verify if you typed correctly.")
          return nil
        }
        
        return vl
      }
    
    case ast.BoolNode:
      return func(ex ast.ExpressionNode) Any {
        bo := ex.(ast.BoolNode)
        
        return bo.Type == ast.TrueType
      }
    
    case ast.FactorialNode:
      return func(ex ast.ExpressionNode) Any {
        f := ex.(ast.FactorialNode)
        
        n, ok := SolveExpression(f.Node).(float64)
        
        if !ok {
          ShowError("Can only perform factorial on a number.", "Examples: 5!, 10.5!, a!")
          return nil
        }
        
        if n < 0 {
          ShowError("Cannot calculate factorial of a negative number.", "You cannot calculate it.")
          return nil
        }
        
        return Factorial(n)
      }
    
    case ast.ExecNode:
      return func(ex ast.ExpressionNode) Any {
        e := ex.(ast.ExecNode)
        
        c := SolveExpression(e.Command).(string)
        var cmd *exec.Cmd
        
        if runtime.GOOS == "windows" {
          cmd = exec.Command("cmd", "/c", c)
        } else {
          cmd = exec.Command("bash", "-c", c)
        }
        
        var out bytes.Buffer
        cmd.Stdout = &out

        var stderr bytes.Buffer
        cmd.Stderr = &stderr
        
        err := cmd.Run()
        
        if err != nil {
          ShowError("An error occurred while executing the command '" + c + "':\n" + stderr.String(), "Fix the error and try again.")
          return nil
        }
        
        return strings.TrimSpace(out.String())
      }
    
    default:
      return nil
  }
}

func Sum(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    s1 := ""
    s2 := ""
    
    if ok1 {
      s1 = strconv.FormatFloat(n1, 'f', -1, 64)
    } else {
      s1, ok1 = v1.(string)
    }
    
    if ok2 {
      s2 = strconv.FormatFloat(n2, 'f', -1, 64)
    } else {
      s2, ok2 = v2.(string)
    }
    
    if !ok1 || !ok2 {
      ShowError("You can only sum numbers or strings, or the variables don't exist.", "Examples: 1 + 1, 'hello ' + 'world', a + 1, 3 + 'hi'.")
      return nil
    }
     
    return s1 + s2
  }
  
  return n1 + n2
}

func Sub(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only subtract numbers, or the variables don't exist.", "Examples: 10 - 4, a - 4, c - f.")
    return nil
  }
  
  return n1 - n2
}

func Mul(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only multiply numbers, or the variables don't exist.", "Examples: 5 * 5, 3 * b, a * c.")
    return nil
  }
  
  return n1 * n2
}

func Div(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only divide numbers, or the variables don't exist.", "Examples: 10 / 5, 20 / a, a / b.")
    return nil
  }
  
  if n2 == 0 {
    ShowError("Cannot divide by zero.", "The divisor is equal to zero.")
    return nil
  }
  
  return n1 / n2
}

func Pow(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only apply power on numbers, or the variables don't exist.", "Examples: 10 ^ 5, 20 ^ a, a ^ b.")
    return nil
  }
  
  return math.Pow(n1, n2)
}

func Mod(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform modulo on numbers, or the variables don't exist.", "Examples: 10 % 5, 20 % a, a % b.")
    return nil
  }
  
  if n2 == 0 {
    ShowError("Cannot divide by zero.", "The divisor is equal to zero.")
    return nil
  }
  
  return math.Mod(n1, n2)
}

func And(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform And on bools, or the variables don't exist.", "Examples: a & b, true & false, false & d.")
    return nil
  }
  
  return n1 && n2
}

func Or(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform Or on bools, or the variables don't exist.", "Examples: a | b, true | false, false | d.")
    return nil
  }
  
  return n1 || n2
}

func Eq(v1 Any, v2 Any) Any {
  return reflect.DeepEqual(v1, v2)
}

func Diff(v1 Any, v2 Any) Any {
  return !reflect.DeepEqual(v1, v2)
}

func Greater(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform Greater on numbers, or the variables don't exist.", "Examples: a > b, 1 > 2, 2 > c.")
    return nil
  }
  
  return n1 > n2
}

func GreaterEq(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform Greater or Equal on numbers, or the variables don't exist.", "Examples: a >= b, 1 >= 2, 2 >= c.")
    return nil
  }
  
  return n1 >= n2
}

func Less(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform Less on numbers, or the variables don't exist.", "Examples: a < b, 1 < 2, 2 < c.")
    return nil
  }
  
  return n1 < n2
}

func LessEq(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    ShowError("You can only perform Less or Equal on numbers, or the variables don't exist.", "Examples: a <= b, 1 <= 2, 2 <= c.")
    return nil
  }
  
  return n1 <= n2
}

func Factorial(n float64) float64 {
  if n == 0 {
    return 1
  }
  
  if n <= 1 && n > 0 {
    return n
  }
  
  return n * Factorial(n - 1)
}