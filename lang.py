import sys

token_types = {          # Examples: (that's the one with brackets)
  "var": "",             # [a] = 10
  "keyword": ["print"],  # [print]
  "assign": "=",         # a [=] 10
  "value": "",           # a = [10]
  "var_ref": "$",        # print [$a]
}

class Token:
  type = ""
  content = ""
  
  def __init__(self, t, c):
    self.type = t
    self.content = c
  
  def __str__(self):
    return "Token: type: " + self.type + ", content: " + self.content
  
  def __repr__(self):
    return "Token: type: " + self.type + ", content: " + self.content

# ---

def main():
  if len(sys.argv) != 2:
    print("Usage: python lang.py <file>")
    sys.exit(1)
  
  try:
    with open(sys.argv[1], "r") as f:
      lines = f.readlines()
      tokens = lexer(lines)
      
      print(tokens)
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
    
    for ch in tk_char:
      # verify if 'ch' is a equals sign
      if ch == token_types["assign"]:
        tks.append(Token("assign", ch))
      # verify if is a variable (if it contains an equals sign in the line)
      elif line.__contains__(token_types["assign"]):
        tks.append(Token("var", ch))
      
    
    tokens.append(tks)
  
  return tokens

# ---

if __name__ == "__main__":
  main()