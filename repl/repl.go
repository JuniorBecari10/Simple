package repl

import (
  "fmt"
  "bufio"
  "os"
  
  "simple/lexer"
)

func Repl() {
  fmt.Println("Simple REPL\n")
  sc := bufio.NewScanner(os.Stdin)
  
  for {
    fmt.Print("> ")
    
  }
}