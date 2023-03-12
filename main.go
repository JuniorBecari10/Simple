package main

import (
  "fmt"
  "os"
  "reflect"
  
  "simple/repl"
  "simple/lexer"
  "simple/parser"
  "simple/run"
)

const (
  Version = "v1.2 Beta"
  
  ModeTokens     = "Tokens"
  ModeStatements = "Statements"
)

var (
  Mode = ""
)

func main() {
  if len(os.Args) == 1 {
    repl.Repl()
    return
  }
  
  for _, a := range os.Args {
    if a == "-t" || a == "--tokens" {
      Mode = ModeTokens
      break
    }
    
    if a == "-s" || a == "--statements" {
      Mode = ModeStatements
      break
    }
  }
  
  if len(os.Args) == 2 {
    if os.Args[1] == "-v" || os.Args[1] == "--version" {
      fmt.Println("Simple " + Version)
      fmt.Println("Made by JuniorBecari10")
      return
    }
    
    if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
      help()
      return
    }
    
    content, err := os.ReadFile(os.Args[1])
    
    if err != nil {
      fmt.Println("File '" + os.Args[1] + "' doesn't exist.")
      fmt.Println("Verify if you typed the name correctly.")
    }
    
    Run(string(content))
    
    return
  }
  
  if len(os.Args) >= 3 {
    if os.Args[1] == "run" {
      Run(os.Args[2])
    }
    
    return
  }
  
  help()
}

func help() {
  fmt.Println("Simple")
  fmt.Println("Version " + Version)
  
  fmt.Println("\nA simple, interpreted programming language. It's very easy to use.\n")
  
  fmt.Println("Usage: simple [file] | [-v | --version] | run [code] [-t | --tokens | -s | --statements]\n")
  
  fmt.Println("Run 'simple' to open the REPL;")
  fmt.Println("Run 'simple run [code]' to automatically run the code you typed;")
  fmt.Println("Run 'simple [file] to run code from file;'")
  fmt.Println("Run 'simple -v' or 'simple --version' to show up the version number.")
  
  fmt.Println("\nhttps://github.com/JuniorBecari10/Simple")
}

func Run(code string) {
  tks := lexer.Lex(code)
  errs := lexer.CheckErrors(tks)
  lines := lexer.SplitLines(code)
  
  if len(errs) > 0 {
    for i, e := range errs {
      repl.Panic(e, lines[i], i)
    }
    
    return
  }
  
  if Mode == ModeTokens {
    fmt.Println("Tokens:\n")
    
    for _, t := range tks {
      fmt.Println(t)
    }
    return
  }
  
  stats := parser.Parse(tks)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for i, e := range errs {
      repl.Panic(e, lines[i], i) // o 'i' não reflete a linha, mas o indice dos erros, adicionar numero da linha nos statements
    }
    
    return
  }
  
  if Mode == ModeStatements {
    fmt.Println("Statements:\n")
    
    for _, s := range stats {
      fmt.Printf("%s | %+v\n", reflect.TypeOf(s), s)
    }
    return
  }
  
  run.Run(stats, lexer.SplitLines(code))
}