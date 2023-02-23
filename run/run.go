package run

import (
  "fmt"
  "os"
  "strconv"
  
  "simple/ast"
)

type Any interface {}

type Value struct {
  Value Any
}

func Panic(msg string) {
  fmt.Println("ERROR: " + msg)
  os.Exit(1)
}

var Variables = map[string]Any {}

func Run(stats []ast.Statement) {
  pc := 0
  for pc < len(stats) {
    stat := stats[pc]
    
    _, ok := stat.(ast.EndStatement)
    
    if ok {
      break
    }
    
    RunStat(stat, false)
    
    pc++
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
        
        n1, ok1 := v1.(float64)
        n2, ok2 := v2.(float64)
        
        switch bin.Op {
          case "+":
            if !ok1 || !ok2 {
              s1 := ""
              s2 := ""
              
              if ok1 {
                s1 = strconv.FormatFloat(n1, 'f', -1, 64)
              } else {
                s1 = v1.(string)
              }
              
              if ok2 {
                s2 = strconv.FormatFloat(n2, 'f', -1, 64)
              } else {
                s2 = v2.(string)
              }
              
              return s1 + s2
            }
            
            return n1 + n2
          
          case "-":
            if !ok1 || !ok2 {
              Panic("Cannot perforn subtraction on a string")
            }
            
            return n1 - n2
          
          case "*":
            if !ok1 || !ok2 {
              Panic("Cannot perforn multiplication on a string")
            }
            
            return n1 * n2
          
          case "/":
            if !ok1 || !ok2 {
              Panic("Cannot perforn division on a string")
            }
            
            if n2 == 0 {
              Panic("Cannot divide by 0")
            }
            
            return n1 / n2
          
          default:
            Panic("Unknown operation: " + bin.Op)
            return ""
        }
      }
    
    default:
      return nil
  }
}