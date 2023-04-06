package repl

import (
  "fmt"
  "bufio"
  "os"
  "strconv"
  "strings"
  "reflect"
  
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
    
    Run(text)
  }
}

func Run(code string) {
  split := strings.Split(code, " ")

  if split[0] == "!t" {
    fmt.Println("Tokens:\n")
    tks := run.GetTokens(code[2:])
    
    for _, t := range tks {
      fmt.Printf("%+v\n", t)
    }
    
    return
  } else if split[0] == "!s" {
    fmt.Println("Statements:\n")
    stats := run.GetStatements(code[2:])
    
    for _, s := range stats {
      fmt.Printf("%s | %+v\n", reflect.TypeOf(s), s)
    }
    
    return
  }

  tokens := lexer.Lex(code)
  errs := lexer.CheckErrors(tokens)
  
  if len(errs) > 0 {
    for _, e := range errs {
      // todo: add arrow ^ in hint, getting the position
      run.LineCode = code
      run.ShowError("Error in lexer: " + e, "")
    }
    
    return
  }
  
  stats := parser.Parse(tokens)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for _, e := range errs {
      run.LineCode = code
      run.ShowError("Error in parser: " + e, "")
    }
    
    return
  }
  
  vls := run.Run(stats, 0, code, true)
  
  for _, vl := range vls {
    if vl == nil {
      return
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
        
        fmt.Println("< " + ret)
      }
      
      run.Error = false
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