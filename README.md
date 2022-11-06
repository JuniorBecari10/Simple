# Lang

A simple programming language made in Python.

## Keywords

`print` - Prints text to the screen. <br>
`printl` - Prints text to the screen, but doesn't break the line. <br>
`input` - For use when declaring a variable. _See variables_. <br>

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

## Variables

A variable can hold any value, for a while its type will always be `string`.

```
var = value
```

In this case, the variable `var` will have the value `value`. <br>
You can use the `input` keyword to get user's input.

```
var = input
```

The user will be prompted for a value, and the value that'll be typed in will be the value of the variable.

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