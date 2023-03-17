package main

import (
  "fmt"
  "os"
  "strings"
  "reflect"
  
  "simple/repl"
  "simple/run"
)

const (
  Version = "Beta v1.3"
  
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
      fmt.Println("Simple - " + Version)
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
      
      return
    }
    
    if Mode == ModeTokens {
      fmt.Println("Tokens:\n")
      tks := run.GetTokens(string(content))
      
      for _, t := range tks {
        fmt.Printf("%+v\n", t)
      }
      
      return
    } else if Mode == ModeStatements {
      fmt.Println("Statements:\n")
      stats := run.GetStatements(string(content))
      
      for _, s := range stats {
        fmt.Printf("%s | %+v\n", reflect.TypeOf(s), s)
      }
      
      return
    }
    
    lines := strings.Split(string(content), "\n")
    
    run.GetLabels(string(content))
    for i, line := range lines {
      run.Run(run.GetStatements(line), i, line, false)
    }
    
    return
  }
  
  if len(os.Args) >= 3 {
    if os.Args[1] == "run" {
      if Mode == ModeTokens {
        fmt.Println("Tokens:\n")
        tks := run.GetTokens(os.Args[2])
        
        for _, t := range tks {
          fmt.Printf("%+v\n", t)
        }
        
        return
      } else if Mode == ModeStatements {
        fmt.Println("Statements:\n")
        stats := run.GetStatements(os.Args[2])
        
        for _, s := range stats {
          fmt.Printf("%s | %+v\n", reflect.TypeOf(s), s)
        }
        
        return
      }
      
      run.GetLabels(os.Args[2])
      run.Run(run.GetStatements(os.Args[2]), 1, os.Args[2], false)
      return
    }
    
    content, err := os.ReadFile(os.Args[1])
    
    if err != nil {
      fmt.Println("File '" + os.Args[1] + "' doesn't exist.")
      fmt.Println("Verify if you typed the name correctly.")
      
      return
    }
    
    if Mode == ModeTokens {
      fmt.Println("Tokens:\n")
      tks := run.GetTokens(string(content))
      
      for _, t := range tks {
        fmt.Printf("%+v\n", t)
      }
      
      return
    } else if Mode == ModeStatements {
      fmt.Println("Statements:\n")
      stats := run.GetStatements(string(content))
      
      for _, s := range stats {
        fmt.Printf("%s | %+v\n", reflect.TypeOf(s), s)
      }
      
      return
    }
    
    // not ; because you can use it inside a string
    lines := strings.Split(string(content), "\n")
    
    run.GetLabels(string(content))
    for i, line := range lines {
      run.Run(run.GetStatements(strings.TrimSpace(line)), i, strings.TrimSpace(line), false)
    }
    
    return
  }
  
  help()
}

func help() {
  fmt.Println("Simple - " + Version)
  
  fmt.Println("\nA simple, interpreted programming language. It's very easy to use.\n")
  
  fmt.Println("Usage: simple [file] | [-v | --version] | [-h | --help] | run [code] [-t | --tokens | -s | --statements]\n")
  
  fmt.Println("Run 'simple' to open the REPL;")
  fmt.Println("Run 'simple run [code]' to automatically run the code you typed;")
  fmt.Println("Run 'simple [file] to run code from file;'")
  fmt.Println("Run 'simple -v' or 'simple --version' to show the version number;")
  fmt.Println("Run 'simple -h' or 'simple --help' to show this help message.")
  
  fmt.Println("\nhttps://github.com/JuniorBecari10/Simple")
}