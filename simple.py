import sys
import os
import shlex

# Version number
version = "v1.1"

# Declare the dictionaries for variables and labels
variables = {}
labels = {}

# Dictionary for all types of tokens
token_types = {          # Examples: (it's the one with brackets)
  "var": "",             # [a] = 10
  "keyword": ["print", "printl", "input", "goto", "exec", "exit", "if", "emptystr"],
  "assign": "=",         # a [=] 10
  "logic": ["==", "!=", ">", ">=", "<", "<="], # if $a [==] a goto 1
  "mathlogic": [">", ">=", "<", "<="],
  "math": ["+=", "-=", "*=", "/=", "%="],
  "types": ["num", "str"], # "arr" not yet
  "value": "",           # a = [10]
  "var_ref": "$",        # print [$a]
  "comment": "#",        # [#] comment
  "label": ":",
  "text": ""             # print [hi]
}

# Extension: .sm

# Create the class Token, with two values: type and content.
# type: The type of the token, one among token_types.
# content: the literal text the token carries
class Token:
  type = ""
  content = ""
  
  def __init__(self, t, c):
    self.type = t
    self.content = c
  
  def __repr__(self):
    return "Token | type: " + str(self.type) + ", content: " + str(self.content)

# Source: StackOverflow (Modified)
# A custom iterator that allows the modification of its count (for goto statements)
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

# The main function. The program starts here.
def main():
  # Verify the CLI args
  if len(sys.argv) != 2:
    print("Usage: simple <file> | [-v | --version]")
    sys.exit(1)
  
  if sys.argv[1] == "-v" or sys.argv[1] == "--version":
    print(f"Simple {version}")
    print("Made by Antonio Carlos (JuniorBecari10).")
    sys.exit(0)
  
  # Run the program
  try:
    # Open the source file
    with open(sys.argv[1], "r") as f:
      # Read all the file and split into lines
      lines = f.read().splitlines()
      
      # Run the Lexer
      tokens = lexer(lines)
      
      # Add labels for goto statemente
      add_labels(tokens)
      
      # Run the code with the tokens.
      run(tokens)
  except FileNotFoundError:
    throw_error_noline(f"The source file '{sys.argv[1]}' doesn't exist.")

# ---

# The lexer. It will split the file into tokens.
# Return: a list of list of tokens (bidimensional array)
def lexer(lines):
  # Define the 2d array
  tokens = []
  
  # Read the file, line by line.
  for i, line in enumerate(lines):
    # Split the line by spaces, because the tokens will be separated by spaces.
    # Code snippet to join quoted strings and split by spaces (Source: StackOverflow)
    tk_char = shlex.split(line, posix=False)
    
    # Declare list
    tks = []
    
    # Read word by word
    for i, ch in enumerate(tk_char):
      # Remove quotes from final text
      ch = ch.replace("\"", "")
      
      # Don't add the comments to the token list
      if ch == token_types["comment"]:
        break
      
      # skip if there's nothing
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
      elif ch.startswith(token_types["var_ref"]):
        tks.append(Token("var_ref", ch))
      # if starts with ':'
      elif ch.startswith(token_types["label"]): # atentar se n é o index 1
        tks.append(Token("label", ch))
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
        # If the lexer can't define correctly, there's a fallback
        if "if" in line:
          if i == 1:
            tks.append(Token("var", ch))
          else:
            tks.append(Token("value", ch))
          continue
        
        tks.append(Token("text", ch))
    
    # append the list to the tokens 2d list
    tokens.append(tks)
  
  return tokens

# Add the labels by reading the entire file first
def add_labels(tokens):
  for lc, line in enumerate(tokens):
    for i, t in enumerate(line):
      # if it's a label
      if t.type == "label":
        # Labels must be the first token in line. See the error message.
        if i == 0:
          # If yes, append to the labels dictionary
          labels[t.content[1:]] = lc + 1
        else:
          line_str = " ".join(token_to_str(line))
          
          # If not, it must be with a goto statement.
          if not line_str.__contains__("goto"):
            # If not, throw an error.
            throw_error("A label must be the first token in line!", lc + 1)

# Run the code
def run(tokens):
  # Define the iterator
  it = Iterator(0, len(tokens))
  
  # Loop through it
  for line_count in it:
    # Define as 'line' the current line it's iterating over
    line = tokens[line_count]
    
    # If there's nothing, continue
    if len(line) == 0:
      continue
    
    # replace all var_ref's by the values
    # And also 'emptystr' by ""
    #
    # Ex:
    # a = 10
    # b = $a -> b = 10
    for i, t in enumerate(line):
      if t.type == "var_ref":
        try:
          line[i].content = variables[t.content[1:]]
          line[i].type = "value"
        except Exception:
          throw_error(f"Variable '{line[i].content[1:]}' doesn't exist.", i + 1)
      
      # Try to convert to number. If not, pass.
      try:
        line[i].content = float(line[i].content)
      except Exception:
        pass
      
      if t.content == "emptystr":
        t.content = ""
    
    # Verify if the current line is a variable declaration
    if is_var_decl(line):
      # The value is the third (or second, counting up from 0) element.
      # var = [value]
      #  0  1    2
      value = line[2].content
      
      # Verify if the 'value' is a keyword
      if line[2].type == "keyword":
        # If it's 'input'
        if line[2].content == "input":
          # Loop until the user types the desired type
          while True:
            value = input("")
            
            # If the programmer only types 'input', break, because the type doesn't matter
            if len(line) == 3:
              break
            
            # if not, proceed.
            if len(line) == 4:
              # Verify each one of the types
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
      
      # Add the variable.
      add_variable(line[0].content, value)
    # Verify if the programmer typed something wrong
    elif line[0].type == "var" and not any_has_type(line, "math"):
      throw_error("Syntax error while declaring a variable.", line_count + 1)
    
    # Verify if it's a mathematical operation
    if is_var_math(line):
      # Get variables
      var1 = variables[line[0].content]
      var2 = line[2].content
      
      # Convert to num if possible
      try:
        var1 = float(var1)
      except Exception:
        var1 = variables[line[0].content]
      
      # Prevention to not do math operations with distinct variable types besides plus
      if type(var1) != type(var2) and not line[1].content.startswith("+"):
        throw_error("Variable types don't match.", line_count + 1)
      
      # Variable to check if both of them are num's
      is_num = type(var1) is float and type(var2) is float
      
      # Plus. With str's you can do, the effect will be a concatenation
      try:
        if line[1].content.startswith("+"):
          if type(var1) != type(var2):
            var1 = str(var1)
            var2 = str(var2)
          
          variables[line[0].content] = var1 + var2
        # Minus and others. These cannot be done with str's.
        elif line[1].content.startswith("-") and is_num:
          variables[line[0].content] = var1 - var2
        elif line[1].content.startswith("*") and is_num:
          variables[line[0].content] = var1 * var2
        elif line[1].content.startswith("/") and is_num:
          if var2 == 0:
            throw_error("Cannot divide by zero.", line_count + 1)
          
          variables[line[0].content] = var1 / var2
        elif line[1].content.startswith("%") and is_num:
          if var2 == 0:
            throw_error("Cannot divide by zero.", line_count + 1)
          
          variables[line[0].content] = var1 % var2
        
      except Exception:
        throw_error(f"Variable '{line[0].content}' doesn't exist.", line_count + 1)
    elif len(line) > 1 and any_has_type(line, "math"):
      throw_error("Syntax error while doing a math operation.", line_count + 1)
    
    # Operations with keywords
    if line[0].type == "keyword":
      # Print statement
      if line[0].content.startswith("print"):
        # Loop all through the line and print each one of them
        for i, t in enumerate(line):
          if i == 0:
            continue
          
          ch = t.content
          
          try:
            ch = float(ch)
            
            if round(ch) == ch:
              ch = int(ch)
          except Exception:
            pass
          
          print(ch, end="")
        
        # If the keyword is printl, don't break the line
        if line[0].content != "printl":
          print()
      # Goto statement
      elif line[0].content == "goto":
        try:
          if len(line) != 2:
            throw_error("Syntax error on a goto statement.", line_count + 1)
          
          # Check if the user typed a line number or label
          if line[1].type == "value":
            # must be int
            line_go = int(line[1].content)
          elif line[1].type == "label":
            line_go = labels[line[1].content[1:]]
          else:
            raise Exception()
          
          # Check if the line is inside the bounds of the file
          if line_go < 0 or line_go > len(tokens):
            throw_error(f"Line out of bounds: {line_go}", line_count + 1)
          
          # Do the 'goto'
          it.revert(line_go - 1)
        except Exception:
          throw_error("Couldn't parse the line number or label to go. Value read: " + line[1].content, line_count + 1)
      # Exec statement
      elif line[0].content == "exec":
        if len(line) == 1:
          throw_error("No commands to run.", line_count + 1)
        
        com = []
        
        # Loop through the line and add to the command list
        for i, t in enumerate(line):
          if i == 0:
            continue
          
          com.append(t.content)
        
        # Run
        os.system(" ".join(com))
      # Exit
      elif line[0].content == "exit":
        # Get status code
        if len(line) == 2:
          try:
            # must be int
            status = int(line[1].content)
          except Exception:
            throw_error("Invalid status code.", line_count + 1)
          
          # Exit
          sys.exit(status)
      # Condition: if statement
      elif is_condition(line):
        try:
          # Get values
          var = variables[line[1].content]
          value = line[3].content
          oper = line[2].content
          line_go = line[5].content
          
          # Check if the value typed for goto is a line or a label
          try:
            if line[5].type == "value":
              # must be int
              line_go = int(line[5].content)
            elif line[5].type == "label":
              line_go = labels[line[5].content[1:]]
            else:
              raise Exception()
          except Exception:
            throw_error("Couldn't parse the line number or label to go. Value read: " + line_go, line_count + 1)
          
          go = False
          
          var_num = True
          value_num = True
          
          # Convert the values to num
          try:
            var = float(var)
          except Exception:
            var_num = False
          
          try:
            value = float(value)
          except Exception:
            value_num = False
          
          is_num = var_num and value_num
          
          # Do the checks, based on the operator the programmer has typed
          if oper == "==":
            go = str(var) == str(value)
          elif oper == "!=":
            go = str(var) != str(value)
          
          # Below here, only num's can perform these logic operations
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
          
          # If succeeded, go to the specified line
          if go:
            if line_go < 0 or line_go > len(tokens):
              throw_error(f"Line out of bounds: {line_go}", line_count + 1)
            
            it.revert(line_go - 1)
          
        except Exception:
          throw_error(f"Variable '{line[1].content}' doesn't exist.", line_count + 1)
      elif any_has_type(line, "logic"):
        throw_error("Syntax error on a if statement.", line_count + 1)

# ---

# Function to check if a value is a num or not
def is_num(var):
  try:
    _ = float(variables[var])
    return True
  except Exception:
    return False

# Function to convert a list of Tokens to a list of strs (t.content)
def token_to_str(tokens):
  strs = []
  
  for i, n in enumerate(tokens):
    strs.append(n.content)
  
  return strs

# Function to ckeck if an array is inside another
def contains_arr(line, arr):
  for n in arr:
    if line.__contains__(n):
      return True
  
  return False

def any_has_type(tokens, type):
  for t in tokens:
    if t.type == type:
      return True
  
  return False

# ---

# Functions to check certain situations
def is_var_math(tokens):
  return len(tokens) == 3 and tokens[0].type == "var" and tokens[1].type == "math" and (tokens[2].type == "value")

def is_var_decl(tokens):
  return (len(tokens) == 3 and tokens[0].type == "var" and tokens[1].type == "assign" and (tokens[2].type == "value" or tokens[2].type == "keyword")) or (len(tokens) == 4 and tokens[0].type == "var" and tokens[1].type == "assign" and (tokens[2].type == "value" or tokens[2].type == "keyword") and tokens[3].type == "type")

def is_condition(tokens):
  return len(tokens) == 6 and (tokens[0].type == "keyword" and tokens[0].content == "if") and tokens[1].type == "var" and tokens[2].type == "logic" and tokens[3].type == "value" and (tokens[4].type == "keyword" and tokens[4].content == "goto") and (tokens[5].type == "value" or tokens[5].type == "label")

# ---

# Function to add a variable
def add_variable(name, value):
  variables[name] = value

# Function used to throw errors
def throw_error(message, line_number):
  line_number = str(line_number)
  
  print("-----")
  print("ERROR")
  print("On line " + line_number + "\n")
  print(message)
  print("-----")
  
  sys.exit(1)

# Function to throw error without specifing the line
def throw_error_noline(message):
  print("-----")
  print("ERROR\n")
  print(message)
  print("-----")
  
  sys.exit(1)

# ---

# Main code
if __name__ == "__main__":
  main()