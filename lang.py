import sys

token_types = [  # Examples: (that's the one with brackets)
  "var",         # [a] = 10
  "keyword",     # [print]
  "assign",      # a [=] 10
  "value",       # a = [10]
  "var_ref"      # print [$a]
]

class Token:
  def __init__(self, t, c):
    self.type = t
    self.content = c

# ---

def main():
  if len(sys.argv) != 2:
    print("Usage: python lang.py <file>")
    sys.exit(1)
  
  try:
    f = open(sys.argv[1], "r")
  except FileNotFoundError:
    print(f"The file '{sys.argv[1]}' doesn't exist.")
    sys.exit(1)
  
  lines = f.readlines()
  tokens = lexer(lines)

# ---

# Return: a list of tokens
def lexer(lines):
  

# ---

if __name__ == "__main__":
  main()