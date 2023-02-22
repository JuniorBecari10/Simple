package run

import (
  "simple/ast"
)

var Variables = map[string]string {}

func Run(stats []ast.Statement, repl bool) string {
  pc := 0
  for pc < len(stats) {
    stat := stats[pc]
    
    fn := GetStatFunc(stat)
    
    if fn == nil {
      return "ERROR: Unknown statement."
    }
    
    fn(stat)
    
    pc++
  }
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
  
}

func RunPrintStatement(st ast.PrintStatement) string {
  
}