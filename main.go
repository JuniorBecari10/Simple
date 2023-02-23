package main

import (
  "fmt"
  "os"
  
  "simple/repl"
)

func main() {
  if len(os.Args) == 0 {
    repl.Repl()
  }
}