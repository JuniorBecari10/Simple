package run

import (
  "fmt"
  "os"
  "strconv"
  "bufio"
  
  "simple/token"
  "simple/ast"
)

type Any interface {}

type Label struct {
  Name string
  Line int
}

type Value struct {
  Value Any
}

var Error bool = false
var PC int = 0
var Labels []Label
var Lines []string

func Panic(msg, hint string) {
  fmt.Println("ERROR: On statement " + strconv.Itoa(PC + 1) + ".")
  fmt.Println("\n" + msg)
  
  fmt.Println()
  
  if PC > 0 {
    fmt.Printf("%d |\n", PC)
  }
  
  fmt.Printf("%d | %s\n", PC + 1, Lines[PC])
  fmt.Printf("%d |\n\n", PC + 2)
  
  if hint != "" {
    fmt.Println(hint)
  }
  
  fmt.Println()
  
  Error = true
}

var Variables = map[string]Any {}

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

func DetectLabels(stats []ast.Statement) {
  Labels = make([]Label, 0)
  
  for i, v := range stats {
    if ls, ok := v.(ast.LabelStatement); ok {
      Labels = append(Labels, Label { ":" + ls.Name, i })
    }
  }
}

func Run(stats []ast.Statement, lines []string) {
  DetectLabels(stats)
  
  PC = 0
  Lines = lines
  for PC < len(stats) {
    stat := stats[PC]
    
    _, ok := stat.(ast.EndStatement)
    
    if ok || Error {
      break
    }
    
    if stat != nil {
      if _, ok := stat.(ast.LabelStatement); ok {
        PC++
        continue
      }
    }
    RunStat(stat, false, "")
    PC++
  }
}

func RunStat(stat ast.Statement, repl bool, s string) Any {
  if repl {
    Lines = []string { s }
  }
  
  fn := GetStatFunc(stat)
  
  if _, ok := stat.(ast.LabelStatement); ok && repl {
    Panic("You cannot declare labels in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if _, ok := stat.(ast.GotoStatement); ok && repl {
    Panic("You cannot declare goto statements in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if _, ok := stat.(ast.IfStatement); ok && repl {
    Panic("You cannot declare if statements in REPL mode.", "You can only use them when you read an actual script.")
    return nil
  }
  
  if fn == nil {
    Panic("Unknown statement.", "Verify if you typed correctly.")
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
        Variables[s.Name.Value] = vl
        
        return vl
      }
    
    case ast.OperationStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.OperationStatement)
        
        vl := SolveExpression(s.Value)
        
        
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
            PC = l.Line
            return ""
          }
        }
        
        Panic("Couldn't find label '" + label + "'.", "Verify if you typed the name correctly.")
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
            
            if vl, ok := res.(bool); ok {
              if vl {
                PC = pc
                return ""
              }
            }
            
            Panic("Cannot use non-boolean expressions inside an if statement.", "You should use only boolean expressions.")
            return nil
          }
        }
        
        Panic("Couldn't find label '" + label + "'.", "Verify if you typed the name correctly.")
        return nil
      }
    
    default:
      return nil
  }
}

// ---

func SolveExpression(ex ast.ExpressionNode) Any {
  fn := GetExprFunc(ex)
  
  if fn == nil {
    Panic("Couldn't get function to solve this expression: " + fmt.Sprintf("%q", ex), "This happens when you use an operator the wrong way or the operator isn't supported.")
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
          Panic("Variable '" + ex.(ast.Identifier).Value + "' doesn't exist.", "Verify if you typed the name correctly.")
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
        return ex.(ast.PlusNode).Value
      }
    
    case ast.MinusNode:
      return func(ex ast.ExpressionNode) Any {
        nb, ok := SolveExpression(ex.(ast.MinusNode).Value).(float64)
        
        if !ok {
          Panic("You can only use numbers with the operator '-'.", "Examples: -10, -a, -25.5.")
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
          
          case "&":
            return And(v1, v2)
          
          case "|":
            return Or(v1, v2)
          
          case "==":
            return Eq(v1, v2)
          
          case "!=":
            return Diff(v1, v2)
          
          default:
            Panic("Unknown operation: " + bin.Op, "")
            return ""
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
              continue
            }
            
            return value
          }
          
          if inp.Type == token.TypeStr {
            if err != nil {
              return vl
            }
            
            continue
          }
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
          Panic("Can only perform factorial on a number.", "Examples: 5!, 10.5!, a!")
        }
        
        if n < 0 {
          Panic("Cannot calculate factorial of a negative number.", "You cannot calculate it.")
        }
        
        return Factorial(n)
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
       Panic("Cannot perform sum on a bool.", "You can only add numbers and strings.")
     }
     
     return s1 + s2
  }
  
  return n1 + n2
}

func Sub(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    Panic("Cannot perform subtraction on a string or a bool", "Examples: 10 - 4, a - 4, c - f.")
  }
  
  return n1 - n2
}

func Mul(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    Panic("You can only multiply numbers.", "Examples: 5 * 5, 3 * b, a * c.")
  }
  
  return n1 * n2
}

func Div(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    Panic("You can only divide numbers.", "Examples: 10 / 5, 20 / a, a / b.")
  }
  
  if n2 == 0 {
    Panic("Cannot divide by zero.", "Self explanatory.")
  }
  
  return n1 / n2
}

func And(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    Panic("You can only perform AND on bools.", "Examples: a & b, true & false, false & d.")
  }
  
  return n1 && n2
}

func Or(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    Panic("You can only perform OR on bools.", "Examples: a | b, true | false, false | d.")
  }
  
  return n1 || n2
}

func Eq(v1 Any, v2 Any) Any {
  
}

func Diff(v1 Any, v2 Any) Any {
  
}

func Factorial(n float64) float64 {
  if n == 0 {
    return 1
  }
  
  if n == 1 {
    return n
  }
  
  return n * Factorial(n - 1)
}