import sys

variables = {}

token_types = {          # Examples: (that's the one with brackets)
  "var": "",             # [a] = 10
  "keyword": ["print"],  # [print]
  "assign": "=",         # a [=] 10
  "value": "",           # a = [10]
  "var_ref": "$",        # print [$a]
  "text": ""             # print [olá]
}

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

# ---

def main():
  if len(sys.argv) != 2:
    print("Usage: python lang.py <file>")
    sys.exit(1)
  
  try:
    with open(sys.argv[1], "r") as f:
      lines = f.readlines()
      tokens = lexer(lines)
      
      run(tokens)
      
      print(tokens)
      print(variables)
  except FileNotFoundError:
    print(f"The file '{sys.argv[1]}' doesn't exist.")
    sys.exit(1)

# ---

# Return: a list of list of tokens (bidimensional array)
def lexer(lines):
  tokens = []
  
  for line in lines:
    tk_char = line.split(" ")
    tks = []
    
    for i, ch in enumerate(tk_char):
      # skip if is nothing
      if ch == "":
        continue
      # verify if there's an equals sign before the current char
      elif token_types["assign"] in tk_char[:i]:
        tks.append(Token("value", ch))
      # verify if 'ch' is a equals sign
      elif ch == token_types["assign"]:
        tks.append(Token("assign", ch))
      # verify if is a variable (if it contains an equals sign in the line)
      elif line.__contains__(token_types["assign"]):
        tks.append(Token("var", ch))
      # if it's a keyword
      elif ch in token_types["keyword"]:
        tks.append(Token("keyword", ch))
      # if starts with '$'
      elif ch[0] == "$":
        tks.append(Token("var_ref", ch))
      # else it's text
      else:
        tks.append(Token("text", ch))
    
    tokens.append(tks)
  
  return tokens

def run(tokens):
  for line in tokens:
    if is_var_decl(line):
      add_variable(line[0].content, line[2].content)

def is_var_decl(tokens):
  return len(tokens) == 3 and tokens[0].type == "var" and tokens[1].type == "assign" and tokens[2].type == "value"

def add_variable(name, value):
  variables[name] = value

# ---

if __name__ == "__main__":
  main()