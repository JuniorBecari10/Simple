package run

import (
  "fmt"
  
  "simple/ast"
)

type Any interface {}

type Value struct {
  Value Any
}

var Variables = map[Any]Any {}

func Run(stats []ast.Statement) {
  pc := 0
  for pc < len(stats) {
    stat := stats[pc]
    
    RunStat(stat)
    
    pc++
  }
}

func RunStat(stat ast.Statement, repl bool) Any {
  fn := GetStatFunc(stat)
  
  if fn == nil {
    return "ERROR: Unknown statement."
  }
  
  return fn(stat)
}

func GetStatFunc(st ast.Statement) func(ast.Statement) Any {
  switch st.(type) {
    case VarDeclStatement:
      return RunVarDeclStatement
    
    case PrintStatement:
      return RunPrintStatement
    
    default:
      return nil
  }
}

func RunVarDeclStatement(st ast.VarDeclStatement) Any {
  Variables[st.Name.Value] = SolveExpression(st.Value)
}

func RunPrintStatement(st ast.PrintStatement) Any {
  
}

// ---

func SolveExpression(ex ast.ExpressionNode) Any {
  fn := GetExprFunc(ex)
  
  if fn == nil {
    return "ERROR: Couldn't get function to solve this expression: " + fmt.Sprintf("%q", ex)
  }
  
  return fn(ex)
}

func GetExprFunc(ex ast.ExpressionNode) func(ast.ExpressionNode) Any {
  switch ex.(type) {
    case Identifier:
      return func(ex ast.ExpressionNode) Any {
        return ex.(Identifier).Value
      }
    
    case NumberNode:
      return func(ex ast.ExpressionNode) Any {
        return ex.(NumberNode).Value
      }
    
    case Identifier:
      return func(ex ast.ExpressionNode) Any {
        return ex.(Identifier).Value
      }
    
  }
}