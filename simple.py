import sys
import os

variables = {}

token_types = {          # Examples: (that's the one with brackets)
  "var": "",             # [a] = 10
  "keyword": ["print", "printl", "input", "goto", "exec", "exit"],
  "assign": "=",         # a [=] 10
  "value": "",           # a = [10]
  "var_ref": "$",        # print [$a]
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
      # skip if is nothing
      if ch == "":
        continue
      # if it's a keyword
      elif ch in token_types["keyword"]:
        tks.append(Token("keyword", ch))
      # verify if there's an equals sign before the current char
      elif token_types["assign"] in tk_char[:i]:
        tks.append(Token("value", ch))
      # verify if 'ch' is a equals sign
      elif ch == token_types["assign"]:
        tks.append(Token("assign", ch))
      # verify if is a variable (if it contains an equals sign in the line)
      elif line.__contains__(token_types["assign"]):
        tks.append(Token("var", ch))
      # if starts with '$'
      elif ch.startswith("$"):
        tks.append(Token("var_ref", ch))
      # else it's text
      else:
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
        value = input("")
      
      add_variable(line[0].content, value)
    elif line[0].type == "var":
      throw_error("Syntax error while declaring a variable.", line_count + 1)
    
    elif line[0].type == "keyword":
      if line[0].content.startswith("print"):
        for i, t in enumerate(line):
          if i == 0:
            continue
          
          if t.type == "var_ref" and t.content.startswith("$"):
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
          
          line_go = int(line[1].content)
          
          if line_go < 0 or line_go > len(tokens):
            throw_error("Line out of bounds.", line_count + 1)
          
          it.revert(line_go - 1)
        except Exception:
          throw_error("Couldn't parse the line number to go. Value: " + line[1].content, line_count + 1)
      elif line[0].content == "exec":
        if len(line) == 1:
          throw_error("No commands to run.", line_count + 1)
        
        os.system(" ".join(token_to_str(line[1:])))
      elif line[0].content == "exit":
        if len(line) == 2:
          try:
            status = int(line[1].content)
          except Exception:
            throw_error("Invalid status code.", line_count + 1)
          
          sys.exit(status)

def token_to_str(tokens):
  strs = []
  
  for i, n in enumerate(tokens):
    strs.append(n.content)
  
  return strs

def is_var_decl(tokens):
  return len(tokens) == 3 and tokens[0].type == "var" and tokens[1].type == "assign" and (tokens[2].type == "value" or tokens[2].type == "keyword")

def add_variable(name, value):
  variables[name] = value

def throw_error(message, line_number):
  line_number = str(line_number)
  
  print("----")
  print("ERROR")
  print("On line " + line_number + "\n")
  print(message)
  print("----")
  
  sys.exit(1)

def throw_error_noline(message):
  print("----")
  print("ERROR\n")
  print(message)
  print("----")
  
  sys.exit(1)

# ---

if __name__ == "__main__":
  main()