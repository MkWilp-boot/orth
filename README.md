# THE ORTH LANG

Orth is a programming language projected for building simple solutions with a stack-based syntax and
It uses [GO](https://github.com/golang/go) as the mastermind for the entire project.

## UNDER DEVELOPMENT

This language is under development

## What are our goals?

1. I expect a simple language capable of doing things in a stack based way
2. Native code compilation
3. Web suport

## Missing features

* Functions (halfway there)
* Local scope

## Getting started

First of all, you will first need the **GO** for bootstrapping the orth's executable,
The installer can be found [here](https://go.dev/dl/).

### Way of writing

As mentioned, Orth is a stack-based language and the way of reading/writing a program is a bit different than usual. First create a file called **hello.orth** then write this code:

```orth
var hello = s "Hello world in orth!"
hold hello put_string
```

and run using the executable

```console
.\core -c=masm hello.orth && .\output.exe

$ Hello world in orth!
```

Really simple, isn't it? this program just pushes a string value into the stack and prints using the instruction `print` but it can get more complicated as things starts to grow.

```orth
var a = i 0  drop
var b = i 1  drop
var c = i 0  drop
var n = i 15 drop

hold n i 0 == if
    hold a print
else
    i 2 while dup hold n > do
        hold a // read a
        hold b // read b
        dup
        var d call grab_last // backup b
        +
        dup // produce c
        var c call grab_last // c = a + b
        hold d // b backed up
        var a call grab_last // a = b
        var b call grab_last // b = c
        i 1 +
    end
end
s "The number is " hold b call to_string + print
```

So, what do you think this program does? if you said it's the fibonacci sequence, then you are right!</br>
If didn't said it correctly, don't worry, 99% of the people reading this will probably fail, don't worry we will get to the point where this code will look readable.

## Compiled Orth and simulated Orth

Yes, orth has two modes, **Simulation** and **Compilation**. Nowadays (2022-10-12) I am currently migrating all simulation code<br/>
to compilation mode, my plans are to support NASM and MASM but only MASM is working.

## Types

Orth is staticly typed, which means it's operands have types and can not be used in strange situations.</br>
Currently Orth have 6 different super types

* INT
* FLOAT
* BOOL
* STRING
* RNT
* VOID

### Integer

Orth has 5 integer variants

1. `i64` representes a 64 bit number
2. `i32` representes a 32 bit number
3. `i16` representes a 16 bit number
4. `i8` representes a 8 bit number
5. `i` let the compiler decide which integer type will be used, it can be any of the previous mentioned, but it's not sure what it is going to be

### Floats

orth has 2 float variants

1. `f64` representes a 64 bit number
2. `f32` representes a 32 bit number

### Booleans

Boolean in Orth are not different from other languages, it can only be `true` or `false`</br>
and are defined by preceding an variable using _b_

### Strings

Orth strings are literal string, that means what ever you put in a string, it will be used the way it got stored.</br>
**with the exception of "\n" that will always print a new line at the end**</br>
We plan to have other string variants like

* `sl` Will represent a string literal
* `si` Will represent a string that can be interpolated

But for now, wel only have **s** as the only string type available

### RNT

RNT stands for Runtime, this variable's type will calculated at runtime without typechecking

### VOID

Like in C or C++, Orth's **VOID** stands for everything that will not return anything.</br>
They are used mainly by operators such as **"+","-","*"**, functions for side effects and so on

## Variables

As you may guess, orth has variables that store values. To create a variable, use the keyword `var` following by it's name and value:
```orth
var name = s "John"
var age = i 20
```
**orth's variables must be initialized when declared**

To use a variable, use the keyword `hold`
```orth
hold name put_string // John
```
## Conditions

Conditions in Orth are very simple and are made by the `if-end` blocks</br>
first: if you want something to be true or false, the last item on the stack **must** be a bool type.

```orth
i 10 i 10 + i 20 == print
```

the code above will produce a bool type that can be used by if blocks

```orth
i 10 i 10 + i 20 ==
if 
    s "Yes, this is a true statement print
else 
    s "Yes, this is a true statement print
end
```

This code will execute the inner instructions within the if block because 10 + 10 = 20.</br>
Otherwise, the else block would be executed instead of the if block.</br>

## Loops

Orth has only the _while_ loop for now and it's very simple to use</br>
a while loop basic consists of a "until thisis true, then keep doing", that's basiclly what a loop looks in Orth

```orth
i 0 while dup i 10 > do
    dup call to_string s "\n" + print
    i 1 +
end drop
```

So, as you can see we:

1. push an inital value (0 in this case)
2. duplicate it for better usage
3. pushes something to compare (10 in this case)

and until 10 is > 0 we add 1 to the last item on the stack, as simple as this

## Functions (simulation only)

Orth have some bultin functions to help people's life:
* to_string (rnt n)
* size_of (s string)
* length_of (s string)
* make_array (type atype, i capacity)
* free_var (var to_free)
* dump_mem ()
* dump_stack ()
* dump_vars ()
* exit (i code)
* fill (rnt value, x|x+1 rangeable)
* index (i relative_pos)
* grab_at (i absolute_pos)
* grab_last (stack > 1)