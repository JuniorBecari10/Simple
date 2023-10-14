package main

import (
	"fmt"
	"os"
	"strings"
)

const (
  Version = "Release v1.0"

  SourceExt = ".sm"
  AssemblyExt = ".sma"
  BytecodeExt = ".smb"

  UsageMsg = "Usage: csimple compile / build <file> | run <file> | assemble <file> | -v / --version | -h / --help"
)

func main() {
  switch len(os.Args) {
    // -v / --version | -h / --help
    case 2:
      option := os.Args[1]

      switch option {
        case "-v":
        case "--version":
          version()

        case "-h":
        case "--help":
          help()
      }
    
    // run | compile | assemble <file>
    case 3:
      command := os.Args[1]
      file := os.Args[2]
      
      content, err := os.ReadFile(file)
      str := string(content)
    
      if err != nil {
        fmt.Println("File '" + file + "' doesn't exist.")
        fmt.Println("Check if you typed the name correctly.")
        
        return
      }

      if !strings.HasSuffix(file, SourceExt) {
        fmt.Printf("Warning: Simple/CSimple source files have the extension '%s'.\n", SourceExt)
      }

      switch command {
        case "run":
          Run(str)
        
        case "assemble":
          Write(strings.Split(file, ".")[0] + BytecodeExt, Assemble(str))
        
        case "compile":
          Write(strings.Split(file, ".")[0] + BytecodeExt, Compile(str))
      }
    
    default:
      fmt.Println(UsageMsg)
  }
}

func Compile(content string) string {
  return ""
}

func Run(content string) {

}

func Assemble(content string) string {
  return ""
}

func Write(file string, content string) {

}

func version() {
  fmt.Println("CSimple - " + Version)
  fmt.Println("Made by JuniorBecari10")
}

func help() {
  fmt.Printf("CSimple - %s\n\n", Version)
  
  fmt.Println("A simple, (now compiled!) programming language.\nIt's very easy to use.")

  fmt.Printf("%s\n\n", UsageMsg)

  fmt.Println("csimple compile / build <file> | compile source code in 'file' to bytecode")
  fmt.Println("csimple run <file>             | run bytecode in 'file'")
  fmt.Println("csimple assemble <file>        | assemble code in 'file' to bytecode")
  fmt.Println("csimple -v / --version         | show the version number")
  fmt.Println("csimple -h / --help            | show this help message")
  
  fmt.Println("\nMade by JuniorBecari10.")
  fmt.Println("https://github.com/JuniorBecari10/CSimple")
}
