package main

import (
  "fmt"
  "os"
  
  "simple/repl"
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
    }
    
    if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
      help()
      return
    }
    
    
  } else {
    help()
  }
}

func help() {
  fmt.Println("Usage: simple [file] | [-v | --version] | [code]\n")
  
  fmt.Println("Execute 'simple' to open the REPL;")
  fmt.Println("Execute 'simple [code]' to automatically run this code.\n")
  
  fmt.Println("https://github.com/JuniorBecari10/Simple")
}