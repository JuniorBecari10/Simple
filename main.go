package main

import (
	"fmt"
	"os"
)

const (
  Version = "Release v1.0"
)

func main() {
  switch len(os.Args) {
    // -v | --version | -h | --help | <file>
    case 2:
      args1 := os.Args[1]

      switch args1 {
        case "-v":
        case "--version":
          fmt.Println("CSimple - " + Version)
          fmt.Println("Made by JuniorBecari10")
          return

        case "-h":
        case "--help":
          help()
          return
      }
      
      content, err := os.ReadFile(args1)
      
      if err != nil {
        fmt.Println("File '" + args1 + "' doesn't exist.")
        fmt.Println("Verify if you typed the name correctly.")
        
        return
      }
      
      Compile(string(content))
    
    // run | assemble <file>
    case 3:
      command := os.Args[1]
      file := os.Args[2]
      
      content, err := os.ReadFile(file)
    
      if err != nil {
        fmt.Println("File '" + file + "' doesn't exist.")
        fmt.Println("Verify if you typed the name correctly.")
        
        return
      }

      switch command {
        case "run":
          Run(string(content))
        
        case "assemble":
          Assemble(string(content))
      }
    
    default:
      help()
  }
}

func help() {
  fmt.Println("CSimple - " + Version)
  
  fmt.Println("\nA simple, interpreted programming language.\nIt's very easy to use.\n")
  
  fmt.Println("Usage: csimple <file> | <-v | --version> | <-h | --help>\n")
  
  fmt.Println("csimple <file>                 | compile source code in 'file' to bytecode")
  fmt.Println("csimple run <file>             | run bytecode in 'file'")
  fmt.Println("csimple assemble <file>        | assemble the bytecode in 'file'")
  fmt.Println("csimple -v | csimple --version | show the version number")
  fmt.Println("csimple -h | csimple --help    | show this help message")
  
  fmt.Println("\nMade by JuniorBecari10.")
  fmt.Println("https://github.com/JuniorBecari10/CSimple")
}

func Compile(content string) {

}

func Run(content string) {

}

func Assemble(content string) {

}
