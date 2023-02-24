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
  errs := lexer.CheckErrors(tks)
  
  if len(errs) > 0 {
    for _, e := range errs {
      fmt.Println(e)
    }
    
    return
  }
  
  stats := parser.Parse(tks)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for _, e := range errs {
      fmt.Println(e)
    }
    
    return
  }
  
  for _, stat := range stats {
    vl := run.RunStat(stat, true)
    
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