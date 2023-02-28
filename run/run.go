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

var Error bool = false
var PC int = 0

type Value struct {
  Value Any
}

func Panic(msg string) {
  fmt.Println("ERROR: " + msg)
  Error = true
}

var Variables = map[string]Any {}

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

func Run(stats []ast.Statement) {
  PC = 0
  for PC < len(stats) {
    stat := stats[PC]
    
    _, ok := stat.(ast.EndStatement)
    
    if ok || Error {
      break
    }
    
    RunStat(stat, false)
    
    PC++
  }
}

func RunStat(stat ast.Statement, repl bool) Any {
  fn := GetStatFunc(stat)
  
  if fn == nil {
    Panic("ERROR: Unknown statement.")
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
    
    default:
      return nil
  }
}

// ---

func SolveExpression(ex ast.ExpressionNode) Any {
  fn := GetExprFunc(ex)
  
  if fn == nil {
    Panic("ERROR: Couldn't get function to solve this expression: " + fmt.Sprintf("%q", ex))
  }
  
  return fn(ex)
}

func GetExprFunc(ex ast.ExpressionNode) func(ast.ExpressionNode) Any {
  switch ex.(type) {
    case ast.Identifier:
      return func(ex ast.ExpressionNode) Any {
        value, ok := Variables[ex.(ast.Identifier).Value]
        
        if !ok {
          Panic("Variable " + ex.(ast.Identifier).Value + " doesn't exist.")
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
          Panic("ERROR: Not a number.")
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
          
          default:
            Panic("Unknown operation: " + bin.Op)
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
       Panic("Cannot perform sum on a bool")
     }
     
     return s1 + s2
  }
  
  return n1 + n2
}

func Sub(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    Panic("Cannot perform subtraction on a string or a bool")
  }
  
  return n1 - n2
}

func Mul(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    Panic("Cannot perform multiplication on a string or a bool")
  }
  
  return n1 * n2
}

func Div(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    Panic("Cannot perform division on a string or a bool")
  }
  
  if n2 == 0 {
    Panic("Cannot divide by zero")
  }
  
  return n1 / n2
}

func And(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    Panic("You can only perform and on bools")
  }
  
  return n1 && n2
}

func Or(v1 Any, v2 Any) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    Panic("You can only perform or on bools")
  }
  
  return n1 || n2
}