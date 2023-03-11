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
    for i, e := range errs {
      Panic(e, code, i)
    }
    
    return
  }
  
  
  stats := parser.Parse(tks)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for i, e := range errs {
      Panic(e, code, i)
    }
    
    return
  }
  
  
  for _, stat := range stats {
    vl := run.RunStat(stat, true, code)
    
    if vl == nil {
      continue
    }
    
    ret := ""
    
    num, ok1 := vl.(float64)
    str, ok2 := vl.(string)
    boo, ok3 := vl.(bool)
    
    if ok1 {
      ret = strconv.FormatFloat(num, 'f', -1, 64)
    } else if ok2 {
      ret = str
    } else if ok3 {
      ret = fmt.Sprintf("%t", boo)
    }
    
    if printRet && !run.Error {
      fmt.Println("< " + ret)
    }
    
    run.Error = false
  }
}

func Panic(msg, lineStr string, line int) {
  fmt.Println("ERROR: On line " + strconv.Itoa(line + 1) + ".")
  fmt.Println("\n" + msg)
  
  fmt.Println()
  
  if line > 0 {
    fmt.Printf("%d |\n", line)
  }
  
  fmt.Printf("%d | %s\n", line + 1, lineStr)
  fmt.Printf("%d |\n", line + 2)
  
  fmt.Println()
}