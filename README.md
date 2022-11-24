<img src="logo-big.png">

# Simple

A simple, interpreted programming language. <br>
It's very easy to use.

## Keywords

`print` - Prints text to the screen; <br>
`printl` - Prints text to the screen, but doesn't break the line; <br>
`input` - For use when declaring a variable _See Variables_; <br>
`goto` - Go to a specified line; <br>
`exec` - Run a system command; <br>
`exit` - Exit the program;
`if` - Checks a condition and go to a line if it's true _See Conditions_.

## Print

For printing things on the screen, you can use the keywords `print` and `printl`. <br>

The difference between them is: <br>
`print` prints the text and break the line. <br>
`printl` prints the text and **DOES NOT** break the line. <br>

### Printing Examples:

```
print Hello World
```

```
printl Hello
print World
```

You can also print a variable.

```
a = 10
print a is $a
```

## Variables

A variable can hold any value, its type can be `str` or `num`.

```
var = <value | [input [str | num]]>
```

In this case, the variable `var` will have the value `value`. <br>
You can use the `input` keyword to get user's input.

A variable type will be `num` if the value is a number. Simple. <br>
So there's no `str`s with numbers inside, unless there's a character that's not a number inside the variable.

```
var = input
```

The user will be prompted for a value, and the value that'll be typed in will be the value of the variable.

## User Input

You can request user's input with the `input` keyword.

```
var = input
```

Within this example, the user can type anything, and the variable will store what the user has typed. <br>
Also, you can force a type with:

```
var = input num
```

The language will keep prompting the user until type a number. <br>
You can do the same with `str`.

```
var = input str
```

In this case, the language will not accept numbers, but strings.

## Conditions

A condition is done with the `if` keyword.

```
if variable <logic> <value | $var> goto <line>
```

## Examples

Print Hello World:

```
print Hello World
```

Set a variable:

```
var = value
```

Get user input and store it in a variable:

```
var = input
```

```
var = input num
```

Print a variable:

```
var = 10
print var is $var
```

------

Program to get user's name and age:

```
printl Type your name:
name = input

printl Type your age:
age = input

print Your name is $name and your age is $age
```

## 