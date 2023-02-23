package repl

import (
  "fmt"
  "bufio"
  "os"
  "strconv"
  
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
    
    text := sc.Text()
    
    if text == "" {
      continue
    }
    
    if text == "exit" {
      os.Exit(0)
    }
    
    fmt.Println("< " + Perform(text))
  }
}

func Perform(q string) string {
  tks := lexer.Lex(q)
  stats := parser.Parse(tks)
  vl := run.RunStat(stats[0], true) // it's going to return only one statement (I think)
  
  value, ok := vl.(float64)
  
  if ok {
    return strconv.FormatFloat(value, 'f', -1, 64)
  }
  
  return vl.(string)
}