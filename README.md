<img src='logo.png'>

# CSimple

A simple, interpreted programming language. <br>
It's very easy to use.

## How to Execute

```
$ simple [file] | [-v | --version] | [-h | --help] [-t | --tokens | -s | --statements]
```

`file` - The script file you want to execute; <br>
`-v` or `--version` - Show the version number; <br>
`-h` or `--help` - Show the help message; <br>
`-t` or `--tokens` - Show the tokens read and don't execute the code; <br>
`-s` or `--statements` - Show the statements read and don't execute the code. <br>

## Keywords

`println` - Prints text to the screen; <br>
`print` - Prints text to the screen, but doesn't break the line; <br>
`input` - For use when declaring a variable | _See Variables_; <br>
`goto` - Go to a specified label | _See Labels_; <br>
`ret` - Return to the previous goto statement | _See Labels_; <br>
`exec` - Run a system command | _See Executing Commands_; <br>
`exit` - Exit the program; <br>
`if` - Checks a condition and go to a label if it's true | _See Conditions_; <br>

## REPL

If you execute the interpreter without any arguments (`$ simple`) the REPL will start. <br>
It's the same thing as running the code from a file, but here you have the response immediately.

You can also request only the tokens or the statements using `!t` or `!s` before your code.

Examples:
```
> !t a = 10
> !s println "a is " + a
```

## Printing

For printing things on the screen, you can use the keywords `print` and `println`. <br>

The difference between them is: <br>
`println` prints the text and break the line. <br>
`print` prints the text and **DOES NOT** break the line. <br>

### Printing Examples:

```py
println 'Hello World'
```

```py
print 'Hello '
println 'World'
```

You can also print a variable.

```py
a = 10
println 'a is ' + a
```

## Variables

A variable can hold any value, its type can be `str`, `num` or `bool`.

```
var = <value | [input [str | num | bool]]>
```

Examples:
```
counter = 1
name = 'John'
exp = true
```

### User Input

You can use the `input` keyword to get user's input.

```py
var = input
```

The user will be prompted for a value, and the value that'll be typed in will be the value of the variable. <br>
You can force a type with `input` too. Just type the type in front of the `input` keyword.

```
age = input num
name = input str
exp = input bool
```

The language will keep prompting until you type a value that satisfies the condition; <br>
The same is for `str` and `bool`; in this case numbers (or anything else but `true` and `false` if the type is `bool`) will not be accepted, and vice versa.

## Labels and Jumps

Labels are used to identify a part of the code. <br>
They are used exclusively with the `goto` keyword.

Example:

```py
goto :menu
exit 0

:menu
# menu...
```

As you can see, the `exit 0` command won't be executed.

### Returning

You can return to the last goto statement that has been executed, because its line gets push onto the Call Stack.

```
goto :first
println 'After'
exit 0

:first
println 'Before'
ret
```

The `ret` keyword will return to the `goto :first` line, and continue executing, without the need to jump there using `goto` again. <br>
But if you return and the call stack is empty, you'll get an error.

## Conditions

A condition is done with the `if` keyword.

```
if <expression> goto <label>
```

```py
x = 10

if x == 10 goto :end
exit 0

:end
println 'x is 10!'
```

## Comments

You can comment using `#`.

```py
# Get user's name
name = input str

# Print user's name
print 'Your name is ' + name
```

## Executing Commands

In Simple, you can also execute system commands, with the `exec` keyword. <br>
It runs the command provided, and returns the value as an expression, so you can store the standard output in a variable! <br>
It doesn't print the command on the screen, to do so you need to assign it to a variable and then print the variable. <br>

```
exec 'command'
```

It automatically detects the current OS, so you don't need to specify OS specific leading commands, such as `cmd /c` or `bash -c`. <br>

Example:
```
print 'Enter a file name: '
file = input str

cont = exec 'cat ' + file

println cont
```

## Errors

In Simple, error messages are well explained and detailed.

Consider the following program:
```py
1 | a = 1
2 | b = 0
3 | c = a / b
4 |
5 | println c
```

The following error message will be printed:
```
-------------

ERROR: On line 3.

Cannot divide by zero.

2 |
3 | c = a / b
4 |

The divisor is equal to zero.
```

When the program runs into an error, it gets terminated.

## Examples

Print Hello World:

```py
print 'Hello World'
```

Set a variable:

```py
var = 'value'
```

Get user input and store it in a variable:

```py
var = input
var = input num
var = input str
var = input bool
```

Print a variable:

```py
var = 10
print 'var is ' + var
```

------

Program to get user's name and age:

```py
printl 'Type your name: '
name = input str

printl 'Type your age: '
age = input num

print 'Your name is ' + name + ' and your age is ' + age + '.'
```

Basic Calculator:
```
print 'Type a number: '

na = input num

print 'Type another number: '

nb = input num

:ope
print 'Type the operation (+ - * /): '

op = input str

if op == '+' goto :plus
if op == '-' goto :minus
if op == '*' goto :times
if op == '/' goto :divide

println 'Invalid operation.'
goto :ope

:plus
res = na + nb
goto :res

:minus
res = na - nb
goto :res

:times
res = na * nb
goto :res

:divide
res = na / nb
goto :res

:res
print 'The result is ' + res
```

Basic Terminal:
```
:start
print "$ "

com = input str
res = exec com

println res
if com == "exit" goto :end
goto :start

:end
```

### Snippets

More snippets on this [Gist](https://gist.github.com/JuniorBecari10/c756b122ebfd35b0e9114d759d84142f). <br>
They were used to test the language, so feel free to submit more!