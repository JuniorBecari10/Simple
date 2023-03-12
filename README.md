<img src='logo.png'>

# Simple

A simple, interpreted programming language. <br>
It's very easy to use.

## Keywords

`println` - Prints text to the screen; <br>
`print` - Prints text to the screen, but doesn't break the line; <br>
`input` - For use when declaring a variable | _See Variables_; <br>
`goto` - Go to a specified line or label | _See labels_; <br>
`ret` - Return to the previous goto statement; <br>
`exit` - Exit the program; <br>
`if` - Checks a condition and go to a line if it's true | _See Conditions_; <br>

## Print

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

A variable can hold any value, its type can be `str` or `num`.

```
var = <value | [input [str | num | bool]]>
```

In this case, the variable `var` will have the value `value`. <br>
You can use the `input` keyword to get user's input.

```py
var = input
```

The user will be prompted for a value, and the value that'll be typed in will be the value of the variable.

You can force a type with `input` too. Just type the type in front of the `input` keyword.

```
age = input num
```

The language will keep prompting until you type a value that satisfies the condiftion; <br>
The same is for `str`; in thic case numbers will not be accepted, and vice versa.

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

### Returning

You can return to the last goto statement that has been executed.

```
goto :first
println 'After'
exit 0

:first
println 'Before'
ret

```

As you can see, the `exit 0` command won't be executed.

## User Input

You can request user's input with the `input` keyword.

```py
var = input
```

Within this example, the user can type anything, and the variable will store what the user has typed. <br>
Also, you can force a type with:

```py
var = input num
```

The language will keep prompting the user until type a number. <br>
You can do the same with `str`.

```py
var = input str
```

In this case, the language won't accept numbers, but strings.

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
```

```py
var = input num
```

```py
var = input str
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