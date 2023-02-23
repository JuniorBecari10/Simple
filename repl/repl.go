package repl

import (
  "fmt"
  "bufio"
  "os"
  
  "simple/lexer"
  "simple/parser"
  "simple/run"
)

func Repl() {
  fmt.Println("Simple REPL\n")
  sc := bufio.NewScanner(os.Stdin)
  
  for {
    fmt.Print("> ")
    
    sc.Scan()
    fmt.Println("< " + Perform(sc.Text()))
  }
}

func Perform(q string) string {
  tks := lexer.Lex(q)
  stats := parser.Parse(tks)
  vl := run.RunStat(stats[0], true) // it's going to return only one statement
  
  return vl.(string)
}