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
    
    Perform(text)
  }
}

func Perform(q string) {
  tks := lexer.Lex(q)
  stats := parser.Parse(tks)
  
  for _, stat := range stats {
    vl := run.RunStat(stat, true) // it's going to return only one statement (I think)
    
    value, ok := vl.(float64)
    ret := ""
    
    if ok {
      ret = strconv.FormatFloat(value, 'f', -1, 64)
    } else {
      ret = vl.(string)
    }
    
    fmt.Println("< " + ret)
  }
}