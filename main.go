package main

import (
  "os"
  
  "simple/repl"
)

func main() {
  if len(os.Args) == 1 {
    repl.Repl()
  }
}