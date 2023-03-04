package main

import (
  "fmt"
  "os"
  
  "simple/repl"
  "simple/lexer"
  "simple/parser"
  "simple/run"
)

const (
  Version = "v0.1 Alpha"
)

func main() {
  if len(os.Args) == 1 {
    repl.Repl()
    return
  }
  
  if len(os.Args) == 2 {
    if os.Args[1] == "-v" || os.Args[1] == "--version" {
      fmt.Println("Simple " + Version)
      fmt.Println("Made by JuniorBecari10")
      return
    }
    
    if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
      help()
    }
    
    content, err := os.ReadFile(os.Args[1])
    
    if err != nil {
      fmt.Println("File " + os.Args[1] + " doesn't exist.")
    }
    
    repl.Run(string(content), false)
    
    return
  }
  
  if len(os.Args) == 3 {
    if os.Args[1] == "run" {
      Run(os.Args[2])
    }
    
    return
  }
  
  help()
}

func help() {
  fmt.Println("Usage: simple [file] | [-v | --version] | run [code]\n")
  
  fmt.Println("Run 'simple' to open the REPL;")
  fmt.Println("Run 'simple run [code]' to automatically run the code you typed;")
  fmt.Println("Run 'simple [file] to run code from file;'")
  fmt.Println("Run 'simple -v' or 'simple --version' to show up the version number.")
  
  fmt.Println("\nhttps://github.com/JuniorBecari10/Simple")
}

func Run(code string) {
  tks := lexer.Lex(code)
  errs := lexer.CheckErrors(tks)
  
  if len(errs) > 0 {
    for _, e := range errs {
      fmt.Println("ERROR: " + e)
    }
    
    return
  }
  
  stats := parser.Parse(tks)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for _, e := range errs {
      fmt.Println("ERROR: " + e)
    }
    
    return
  }
  
  run.Run(stats, lexer.SplitLines(code))
}
