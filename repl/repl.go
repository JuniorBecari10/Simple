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
    
    Run(text, true)
  }
}

func Run(code string, printRet bool) {
  tks := lexer.Lex(code)
  errs := lexer.CheckErrors(tks)
  
  if len(errs) > 0 {
    for _, e := range errs {
      fmt.Println("ERROR: " + e)
    }
    
    return
  }
  
  stats := parser.Parse(tks)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for _, e := range errs {
      fmt.Println("ERROR: " + e)
    }
    
    return
  }
  
  for _, stat := range stats {
    vl := run.RunStat(stat, true)
    
    if vl == nil {
      continue
    }
    
    value, ok := vl.(float64)
    ret := ""
    
    if ok {
      ret = strconv.FormatFloat(value, 'f', -1, 64)
    } else {
      b, ok := vl.(bool)
      
      if !ok {
        ret = vl.(string)
      }
      
      ret = fmt.Sprintf("%t", b)
    }
    
    if printRet && !run.Error {
      fmt.Println("< " + ret)
    }
    
    run.Error = false
  }
}