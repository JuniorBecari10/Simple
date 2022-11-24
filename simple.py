import sys
import os

variables = {}

token_types = {          # Examples: (it's the one with brackets)
  "var": "",             # [a] = 10
  "keyword": ["print", "printl", "input", "goto", "exec", "exit", "if"],
  "assign": "=",         # a [=] 10
  "logic": ["==", "!=", ">", ">=", "<", "<="], # if $a [==] a goto 1
  "mathlogic": [">", ">=", "<", "<="],
  "math": ["+=", "-=", "*=", "/="],
  "types": ["num", "str"], # "arr" not yet
  "value": "",           # a = [10]
  "var_ref": "$",        # print [$a]
  "comment": "#",        # [#] comment
  "text": ""             # print [hi]
}

# Extension: .sm

class Token:
  type = ""
  content = ""
  
  def __init__(self, t, c):
    self.type = t
    self.content = c
  
  def __str__(self):
    return "Token | type: " + self.type + ", content: " + self.content
  
  def __repr__(self):
    return "Token | type: " + self.type + ", content: " + self.content

# Source: StackOverflow (Modified)
class Iterator:
  def __init__(self, start, end):
    self.end = end
    self.position = start
  
  def __next__(self):
    if self.position >= self.end:
      raise StopIteration
    else:
      self.position += 1
      return self.position - 1
  
  def __iter__(self):
    return self
  
  def revert(self, n=1):
    self.position = n

# ---

def main():
  if len(sys.argv) != 2:
    print("Usage: python lang.py <file>")
    sys.exit(1)
  
  try:
    with open(sys.argv[1], "r") as f:
      lines = f.read().splitlines()
      tokens = lexer(lines)
      
      run(tokens)
  except FileNotFoundError:
    throw_error_noline(f"The source file '{sys.argv[1]}' doesn't exist.")

# ---

# Return: a list of list of tokens (bidimensional array)
def lexer(lines):
  tokens = []
  
  for i, line in enumerate(lines):
    tk_char = line.split(" ")
    tks = []
    
    for i, ch in enumerate(tk_char):
      if ch == token_types["comment"]:
        break
      
      # skip if is nothing
      if ch == "":
        continue
      elif ch in token_types["keyword"]:
        tks.append(Token("keyword", ch))
      # verify if there's a double equals sign before the current char
      elif ch in token_types["logic"]:
        tks.append(Token("logic", ch))
      # verify if 'ch' is a equals sign
      elif ch == token_types["assign"]:
        tks.append(Token("assign", ch))
      # verify if there's any math sign
      elif ch in token_types["math"]:
        tks.append(Token("math", ch))
       # verify if there's any type keyword
      elif ch in token_types["types"]:
        tks.append(Token("type", ch))
      # if starts with '$'
      elif ch.startswith("$"):
        tks.append(Token("var_ref", ch))
      # verify if there's an equals sign before the current char
      elif token_types["assign"] in tk_char[:i]:
          tks.append(Token("value", ch))
      # verify if is a variable (if it contains an equals sign in the line)
      elif line.__contains__(token_types["assign"]):
        if contains_arr(line, token_types["logic"]) or (contains_arr(line, token_types["math"]) and i > 1):
          if line.__contains__("if") and i == 1:
            tks.append(Token("var", ch))
            continue
          
          tks.append(Token("value", ch))
        else:
          tks.append(Token("var", ch))
      # else it's text
      else:
        if "if" in line:
          if i == 1:
            tks.append(Token("var", ch))
          else:
            tks.append(Token("value", ch))
          continue
        
        tks.append(Token("text", ch))
    
    tokens.append(tks)
  
  return tokens

def run(tokens):
  it = Iterator(0, len(tokens))
  
  for line_count in it:
    line = tokens[line_count]
    
    if len(line) == 0:
      continue
    
    if is_var_decl(line):
      value = line[2].content
      
      if line[2].type == "keyword" and line[2].content == "input":
        while True:
            value = input("")
            
            if len(line) == 3:
              break
            
            if len(line) == 4:
              if line[3].content == "num":
                try:
                  _ = float(value)
                  
                  break
                except Exception:
                  continue
              elif line[3].content == "str":
                try:
                  _ = float(value)
              
                  continue
                except Exception:
                  break
              else:
                throw_error("Type '" + line[3].content + "' is not allowed for input.", line_count + 1)
      
      add_variable(line[0].content, value)
    elif line[0].type == "var" and line[1].type != "math":
      throw_error("Syntax error while declaring a variable.", line_count + 1)
    
    if is_var_math(line):
      var1 = variables[line[0].content]
      var2 = line[2].content
      
      if var2.startswith(token_types["var_ref"]):
        try:
          var2 = variables[line[2].content[1:]]
        except Exception:
          throw_error(f"Variable '{line[2].content[1:]}' doesn't exist.", line_count + 1)
      
      try:
        var1 = float(var1)
      except Exception:
        var1 = variables[line[0].content]
      
      if not var2.startswith(token_types["var_ref"]):
        try:
          var2 = float(var2)
        except Exception:
          var2 = line[2].content
      
      if type(var1) != type(var2) and not line[1].content.startswith("+"):
        throw_error("Variable types don't match.", line_count + 1)
      
      is_num = type(var1) is float
      
      try:
        if line[1].content.startswith("+"):
          if type(var1) != type(var2):
            var1 = str(var1)
            var2 = str(var2)
          
          variables[line[0].content] = var1 + var2
        elif line[1].content.startswith("-") and is_num:
          variables[line[0].content] = var1 - var2
        elif line[1].content.startswith("*") and is_num:
          variables[line[0].content] = var1 * var2
        elif line[1].content.startswith("/") and is_num:
          if var2 == 0:
            throw_error("Cannot divide by zero.", line_count + 1)
          
          variables[line[0].content] = var1 / var2
        
      except Exception:
        throw_error(f"Variable '{line[0].content}' doesn't exist.", line_count + 1)
    elif len(line) > 1 and line[1].type == "math":
      throw_error("Syntax error while doing a math operation.", line_count + 1)
    
    if line[0].type == "keyword":
      if line[0].content.startswith("print"):
        for i, t in enumerate(line):
          if i == 0:
            continue
          
          if t.type == "var_ref" and t.content.startswith(token_types["var_ref"]):
            try:
              print(variables[t.content[1:]], end=" ")
            except Exception:
              throw_error(f"Variable '{t.content[1:]}' doesn't exist.", i)
          else:
            print(t.content, end=" ")
        
        if line[0].content != "printl":
          print()
      elif line[0].content == "goto":
        try:
          if len(line) != 2:
            throw_error("Syntax error on a goto statement.", line_count + 1)
          
          # must be int
          line_go = int(line[1].content)
          
          if line_go < 0 or line_go > len(tokens):
            throw_error(f"Line out of bounds: {line_go}", line_count + 1)
          
          it.revert(line_go - 1)
        except Exception:
          throw_error("Couldn't parse the line number to go. Value read: " + line[1].content, line_count + 1)
      elif line[0].content == "exec":
        if len(line) == 1:
          throw_error("No commands to run.", line_count + 1)
        
        com = []
        
        for i, t in enumerate(line):
          if i == 0:
            continue
          
          if t.type == "var_ref":
            com.append(variables[t.content[1:]])
            continue
          
          com.append(t.content)
        
        os.system(" ".join(com))
      elif line[0].content == "exit":
        if len(line) == 2:
          try:
            # must be int
            status = int(line[1].content)
          except Exception:
            throw_error("Invalid status code.", line_count + 1)
          
          sys.exit(status)
      elif is_condition(line):
        try:
          var = variables[line[1].content]
          value = line[3].content
          oper = line[2].content
          line_go = line[5].content
          
          try:
            # must be int
            line_go = int(line_go)
          except Exception:
            throw_error("Couldn't parse the line number to go. Value read: " + line_go, line_count + 1)
          
          go = False
          
          if line[3].type == "var_ref":
            try:
              value = variables[value[1:]]
            except Exception:
              throw_error(f"Variable '{value[1:]}' doesn't exist.", line_count + 1)
          
          var_num = True
          value_num = True
          
          try:
            var = float(var)
          except Exception:
            var_num = False
          
          try:
            value = float(value)
          except Exception:
            value_num = False
          
          is_num = var_num and value_num
          
          if oper == "==":
            go = str(var) == str(value)
          elif oper == "!=":
            go = str(var) != str(value)
          
          if not is_num and oper in token_types["mathlogic"]:
            throw_error(f"Cannot perform math logical operations on strings.", line_count + 1)
          
          if oper == ">":
            go = var > value
          elif oper == ">=":
            go = var >= value
          elif oper == "<":
            go = var < value
          elif oper == "<=":
            go = var <= value
          
          if go:
            if line_go < 0 or line_go > len(tokens):
              throw_error(f"Line out of bounds: {line_go}", line_count + 1)
            
            it.revert(line_go - 1)
          
        except Exception:
          throw_error(f"Variable '{line[1].content}' doesn't exist.", line_count + 1)

def is_num(var):
  try:
    _ = float(variables[var])
    return True
  except Exception:
    return False

def token_to_str(tokens):
  strs = []
  
  for i, n in enumerate(tokens):
    strs.append(n.content)
  
  return strs

def contains_arr(line, arr):
  for n in arr:
    if line.__contains__(n):
      return True
  
  return False

def is_var_math(tokens):
  return len(tokens) == 3 and tokens[0].type == "var" and tokens[1].type == "math" and (tokens[2].type == "value" or tokens[2].type == "var_ref")

def is_var_decl(tokens):
  return (len(tokens) == 3 and tokens[0].type == "var" and tokens[1].type == "assign" and (tokens[2].type == "value" or tokens[2].type == "keyword")) or (len(tokens) == 4 and tokens[0].type == "var" and tokens[1].type == "assign" and (tokens[2].type == "value" or tokens[2].type == "keyword") and tokens[3].type == "type")

def is_condition(tokens):
  return len(tokens) == 6 and tokens[0].type == "keyword" and tokens[1].type == "var" and tokens[2].type == "logic" and (tokens[3].type == "value" or tokens[3].type == "var_ref") and tokens[4].type == "keyword" and tokens[5].type == "value"

def add_variable(name, value):
  variables[name] = value

def throw_error(message, line_number):
  line_number = str(line_number)
  
  print("-----")
  print("ERROR")
  print("On line " + line_number + "\n")
  print(message)
  print("-----")
  
  sys.exit(1)

def throw_error_noline(message):
  print("-----")
  print("ERROR\n")
  print(message)
  print("-----")
  
  sys.exit(1)

# ---

if __name__ == "__main__":
  main()
