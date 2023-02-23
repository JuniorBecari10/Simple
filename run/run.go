package run

import (
  "simple/ast"
)

var Variables = map[string]string {}

func Run(stats []ast.Statement) {
  pc := 0
  for pc < len(stats) {
    stat := stats[pc]
    
    RunStat(stat)
    
    pc++
  }
}

func RunStat(stat ast.Statement, repl bool) string {
  fn := GetStatFunc(stat)
  
  if fn == nil {
    return "ERROR: Unknown statement."
  }
  
  return fn(stat)
}

func GetStatFunc(st ast.Statement) func(ast.Statement) string {
  switch st.(type) {
    case VarDeclStatement:
      return RunVarDeclStatement
    
    case PrintStatement:
      return RunPrintStatement
    
    default:
      return nil
  }
}

func RunVarDeclStatement(st ast.VarDeclStatement) string {
  Variables[st.Name.Value] = 
}

func RunPrintStatement(st ast.PrintStatement) string {
  
}

// ---

func SolveExpression(ex ast.ExpressionNode) string {
  
}