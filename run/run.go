package run

import (
  "fmt"
  "os"
  "strconv"
  "bufio"
  "reflect"
  
  "simple/token"
  "simple/lexer"
  "simple/ast"
  "simple/parser"
)

type Any interface {}

type Label struct {
  Name string
  Line int
}

type Value struct {
  Value Any
}

var Error bool = false
var Labels []Label
var PC int
var Variables = map[string]Any {}
var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

// ---

func Run(code string) {
  tks := lexer.Lex(code)
  errs := lexer.CheckErrors(tks)
  lines := lexer.SplitLines(code)
  
  if len(errs) > 0 {
    for i, e := range errs {
      PrintError(e, "", lines[i], i)
    }
    
    return
  }
  
  stats := parser.Parse(tks, code)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for i, e := range errs {
      PrintError(e, "", lines[i], i) // colocar o numero da linha
    }
    
    return
  }
  
  Exec(stats, lexer.SplitLines(code))
}

func ExecRepl(code string) {
  tks := lexer.Lex(code)
  errs := lexer.CheckErrors(tks)
  
  if len(errs) > 0 {
    for i, e := range errs {
      PrintError(e, "", code, i)
    }
    
    return
  }
  
  stats := parser.Parse(tks, code)
  errs = parser.CheckErrors(stats)
  
  if len(errs) > 0 {
    for i, e := range errs {
      PrintError(e, "", code, i)
    }
    
    return
  }
  
  for _, stat := range stats {
    vl := ExecStat(stat, true, code)
    
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
    
    if !Error {
      fmt.Println("< " + ret)
    }
    
    Error = false
  }
}

func PrintError(msg, hint, code string, line int) {
  fmt.Println("ERROR: On line " + strconv.Itoa(line + 1) + ".")
  fmt.Println("\n" + msg)
  
  fmt.Println()
  
  fmt.Printf("%d | %s\n", line + 1, code)
  
  fmt.Printf("\n%s\n\n", hint);
  
  Error = true
}

// ---

func DetectLabels(stats []ast.Statement) {
  Labels = make([]Label, 0)
  
  for i, v := range stats {
    if ls, ok := v.(ast.LabelStatement); ok {
      Labels = append(Labels, Label { ":" + ls.Name, i })
    }
  }
}

func Exec(stats []ast.Statement, lines []string) {
  DetectLabels(stats)
  
  PC = 0
  for PC < len(stats) {
    stat := stats[PC]
    
    _, ok := stat.(ast.EndStatement)
    
    if ok || Error {
      break
    }
    
    if stat != nil {
      if _, ok := stat.(ast.LabelStatement); ok {
        PC++
        continue
      }
    }
    ExecStat(stat, false, "")
    PC++
  }
}

func ExecStat(stat ast.Statement, repl bool, lineRepl string) Any {
  /*if repl {
    Lines = []string { lineRepl }
  }*/
  
  fn := GetStatFunc(stat)
  
  if st, ok := stat.(ast.LabelStatement); ok && repl {
    PrintError("You cannot declare labels in REPL mode.", "You can only use them when you read an actual script.", st.Code, st.Line)
    return nil
  }
  
  if st, ok := stat.(ast.GotoStatement); ok && repl {
    PrintError("You cannot declare goto statements in REPL mode.", "You can only use them when you read an actual script.", st.Code, st.Line)
    return nil
  }
  
  if st, ok := stat.(ast.IfStatement); ok && repl {
    PrintError("You cannot declare if statements in REPL mode.", "You can only use them when you read an actual script.", st.Code, st.Line)
    return nil
  }
  
  if fn == nil {
    PrintError("Unknown statement.", "Verify if you typed correctly.", "", -1) // corrigir
    return nil
  }
  
  return fn(stat)
}

func GetStatFunc(st ast.Statement) func(ast.Statement) Any {
  switch st.(type) {
    case ast.VarDeclStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.VarDeclStatement)
        
        vl := SolveExpression(s.Value)
        Variables[s.Name.Value] = vl
        
        return vl
      }
    
    case ast.OperationStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.OperationStatement)
        
        vl := SolveExpression(s.Value)
        
        
        switch s.Op {
          case "+":
            vl := Sum(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "-":
            vl := Sub(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "*":
            vl := Mul(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "/":
            vl := Div(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "&":
            vl := And(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
          case "|":
            vl := Or(Variables[s.Name.Value], vl)
            Variables[s.Name.Value] = vl
            
            return vl
          
        }
        
        return vl
      }
    
    case ast.PrintStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.PrintStatement)
        
        exp := SolveExpression(s.Expression)
        fmt.Print(exp)
        
        if s.BreakLine {
          fmt.Println()
        }
        
        return exp
      }
    
    case ast.ExpressionStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.ExpressionStatement)
        
        return SolveExpression(s.Expression)
      }
    
    case ast.GotoStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.GotoStatement)
        
        label := s.Label
        
        for _, l := range Labels {
          if l.Name == label {
            PC = l.Line
            return ""
          }
        }
        
        PrintError("Couldn't find label '" + label + "'.", "Verify if you typed the name correctly.", s.Code, s.Line)
        return nil
      }
    
    case ast.IfStatement:
      return func(st ast.Statement) Any {
        s := st.(ast.IfStatement)
        
        res := SolveExpression(s.Expression)
        
        label := s.Label
        pc := 0
        
        for _, l := range Labels {
          if l.Name == label {
            pc = l.Line
            
            if vl, ok := res.(bool); ok {
              if vl {
                PC = pc
                return ""
              }
            }
            
            PrintError("Cannot use non-boolean expressions inside an if statement.", "You should use only boolean expressions.", s.Code, s.Line)
            return nil
          }
        }
        
        PrintError("Couldn't find label '" + label + "'.", "Verify if you typed the name correctly.", s.Code, s.Line)
        return nil
      }
    
    default:
      return nil
  }
}

// ---

func SolveExpression(ex ast.ExpressionNode, code string, line int) Any {
  fn := GetExprFunc(ex)
  
  if fn == nil {
    PrintError("Couldn't get function to solve this expression: " + fmt.Sprintf("%q", ex), "This happens when you use an operator the wrong way or the operator isn't supported.", code, line)
    return nil
  }
  
  return fn(ex)
}

func GetExprFunc(ex ast.ExpressionNode) func(ast.ExpressionNode) Any {
  switch ex.(type) {
    case ast.Identifier:
      return func(ex ast.ExpressionNode) Any {
        value, ok := Variables[ex.(ast.Identifier).Value]
        
        if !ok {
          PrintError("Variable '" + ex.(ast.Identifier).Value + "' doesn't exist.", "Verify if you typed the name correctly.", s.Code, s.Line)
        }
        
        return value
      }
    
    case ast.NumberNode:
      return func(ex ast.ExpressionNode) Any {
        return ex.(ast.NumberNode).Value
      }
    
    case ast.StringNode:
      return func(ex ast.ExpressionNode) Any {
        return ex.(ast.StringNode).Value
      }
    
    case ast.PlusNode:
      return func(ex ast.ExpressionNode) Any {
        return ex.(ast.PlusNode).Value
      }
    
    case ast.MinusNode:
      return func(ex ast.ExpressionNode) Any {
        s := ex.(ast.MinusNode)
        nb, ok := SolveExpression(s.Value).(float64)
        
        if !ok {
          PrintError("You can only use numbers with the operator '-'.", "Examples: -10, -a, -25.5.", s.Code, s.Line)
        }
        
        return -nb
      }
    
    case ast.BinNode:
      return func(ex ast.ExpressionNode) Any {
        bin := ex.(ast.BinNode)
        
        v1, v2 := SolveExpression(bin.NodeA), SolveExpression(bin.NodeB)
        
        switch bin.Op {
          case "+":
            return Sum(v1, v2)
          
          case "-":
            return Sub(v1, v2)
          
          case "*":
            return Mul(v1, v2)
          
          case "/":
            return Div(v1, v2)
          
          case "&":
            return And(v1, v2)
          
          case "|":
            return Or(v1, v2)
          
          case "==":
            return Eq(v1, v2)
          
          case "!=":
            return Diff(v1, v2)
          
          case ">":
            return Greater(v1, v2)
          
          case ">=":
            return GreaterEq(v1, v2)
          
          case "<":
            return Less(v1, v2)
          
          case "<=":
            return LessEq(v1, v2)
          
          default:
            PrintError("Unknown operation: " + bin.Op, "", bin.Code, bin.Line)
            return ""
        }
      }
    
    case ast.InputNode:
      return func(exp ast.ExpressionNode) Any {
        inp := exp.(ast.InputNode)
        vl := ""
        
        for {
          scanner.Scan()
          vl = scanner.Text()
          
          if inp.Type == "" {
            return vl
          }
          
          value, err := strconv.ParseFloat(vl, 64)
          
          if inp.Type == token.TypeNum {
            if err != nil {
              continue
            }
            
            return value
          }
          
          if inp.Type == token.TypeStr {
            if err != nil {
              return vl
            }
            
            continue
          }
        }
        
        return vl
      }
    
    case ast.BoolNode:
      return func(ex ast.ExpressionNode) Any {
        bo := ex.(ast.BoolNode)
        
        return bo.Type == ast.TrueType
      }
    
    case ast.FactorialNode:
      return func(ex ast.ExpressionNode) Any {
        f := ex.(ast.FactorialNode)
        
        n, ok := SolveExpression(f.Node).(float64)
        
        if !ok {
          PrintError("Can only perform factorial on a number.", "Examples: 5!, 10.5!, a!", f.Code, f.Line)
        }
        
        if n < 0 {
          PrintError("Cannot calculate factorial of a negative number.", "You cannot calculate it.", f.Code, f.Line)
        }
        
        return Factorial(n)
      }
    
    default:
      return nil
  }
}

func Sum(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    s1 := ""
    s2 := ""
    
     if ok1 {
       s1 = strconv.FormatFloat(n1, 'f', -1, 64)
     } else {
       s1, ok1 = v1.(string)
     }
     
     if ok2 {
       s2 = strconv.FormatFloat(n2, 'f', -1, 64)
     } else {
       s2, ok2 = v2.(string)
     }
     
     if !ok1 || !ok2 {
       PrintError("Cannot perform sum on a bool.", "You can only add numbers and strings.", code, line)
     }
     
     return s1 + s2
  }
  
  return n1 + n2
}

func Sub(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("Cannot perform subtraction on a string or a bool", "Examples: 10 - 4, a - 4, c - f.", code, line)
  }
  
  return n1 - n2
}

func Mul(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("You can only multiply numbers.", "Examples: 5 * 5, 3 * b, a * c.", code, line)
  }
  
  return n1 * n2
}

func Div(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("You can only divide numbers.", "Examples: 10 / 5, 20 / a, a / b.", code, line)
  }
  
  if n2 == 0 {
    PrintError("Cannot divide by zero.", "Self explanatory.")
  }
  
  return n1 / n2
}

func And(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    PrintError("You can only perform AND on bools.", "Examples: a & b, true & false, false & d.", code, line)
  }
  
  return n1 && n2
}

func Or(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(bool)
  n2, ok2 := v2.(bool)
  
  if !ok1 || !ok2 {
    PrintError("You can only perform OR on bools.", "Examples: a | b, true | false, false | d.", code, line)
  }
  
  return n1 || n2
}

func Eq(v1 Any, v2 Any) Any {
  return reflect.DeepEqual(v1, v2)
}

func Diff(v1 Any, v2 Any) Any {
  return !reflect.DeepEqual(v1, v2)
}

func Greater(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("You can only perform Greater on numbers.", "Examples: a > b, 1 > 2, 2 > c.", code, line)
  }
  
  return n1 > n2
}

func GreaterEq(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("You can only perform Greater or Equal on numbers.", "Examples: a >= b, 1 >= 2, 2 >= c.", code, line)
  }
  
  return n1 >= n2
}

func Less(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("You can only perform Less on numbers.", "Examples: a < b, 1 < 2, 2 < c.", code, line)
  }
  
  return n1 < n2
}

func LessEq(v1 Any, v2 Any, code string, line int) Any {
  n1, ok1 := v1.(float64)
  n2, ok2 := v2.(float64)
  
  if !ok1 || !ok2 {
    PrintError("You can only perform Less or Equal on numbers.", "Examples: a <= b, 1 <= 2, 2 <= c.", code, line)
  }
  
  return n1 <= n2
}

func Factorial(n float64) float64 {
  if n == 0 {
    return 1
  }
  
  if n == 1 {
    return n
  }
  
  return n * Factorial(n - 1)
}