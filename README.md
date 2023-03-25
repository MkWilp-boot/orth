# THE ORTH LANG

Orth is a programming language projected for building simple solutions with a stack-based syntax and
It uses [GO](https://github.com/golang/go) as the mastermind for the entire project.

## UNDER DEVELOPMENT

This language is under development so don't expect too much

## What are our goals?

1. I expect a simple language capable of doing things in a stack based way
2. Native code compilation
3. Web suport

## Missing features

* Local scope for if/else/while(DONE)

## Getting started

First of all, you will first need the **GO** for bootstrapping the orth's executable,
The installer can be found [here](https://go.dev/dl/).

### Way of writing

As mentioned, Orth is a stack-based language and the way of reading/writing a program is a bit different than usual. First create a file called **hello.orth** then write this code:

```orth
proc main in
    s "Hello world in orth!" puts
end
```

and run using the executable

```console
.\core -com=masm hello.orth && .\output.exe

$ Hello world in orth!
```

Really simple, isn't it? this program just pushes a string value into the stack and prints using the instruction `puts` but it can get more complicated as things start to grow.

```orth
@define BOARD_CAP i 30
@define BOARD_MINUS_TWO i 28

proc main with 0 out 0 in
    mem BOARD_MINUS_TWO + i 1 .

    i32 0 while dup BOARD_MINUS_TWO > do
        i32 0 while dup BOARD_CAP > do
            dup mem + , if
                dup BOARD_CAP + mem + i8 42 .
            else
                dup BOARD_CAP + mem + i8 32 .
            end
            i32 1 +
        end drop

        mem BOARD_CAP + i8 10 .

        i32 31 mem BOARD_CAP + dump_mem

        mem i 0 + , i 1 lshift mem i32 1 + , lor

        i32 1 while dup BOARD_MINUS_TWO > do
            swap
            i 1 lshift i 7 land
            over mem + i 1 + , lor
            2dup i 110 swap rshift i 1 land
            swap mem + swap .
            swap
            i 1 +
        end drop
        drop

        i32 1 +
    end drop
end
```

So, what do you think this program does? Of course it is the Rule 110! isn't it obvious?!?!?!</br>
Yes, I am joking, but this program is in fact the Rule 110 implementation, you just can't read it (yet)

## Compiled Orth

Yes **Compilation**. You can compile your program to native code by using the "-com=" flag followed by the one of the supported assemblers.<br/>
I have plans to support both NASM and MASM but only MASM is working.

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

1. `i64` representes a 64 bit number (QWORD)
2. `i32` representes a 32 bit number (DWORD)
3. `i16` representes a 16 bit number (WORD)
4. `i8` representes a 8 bit number   (BYTE)
5. `i` let the compiler decide which integer type will be used, it will be a 32 bit integer on a 32 bit machine and a 64 bit integer on a 64 bit machine.</br>
This is usually used for just pushing a number, like a unit

### Floats

orth has 2 float variants

1. `f64` representes a 64 bit number
2. `f32` representes a 32 bit number

### Booleans

Boolean in Orth are not different from other languages, it can only be `true` or `false`</br>
and are defined by preceding an variable using _b_

### Strings

Orth string are defined by using the type _s_ followed by the string literal between _" "_</br>
We plan to have other string variants like

* `si` Will represent a string that can be interpolated

But for now, wel only have _s_ as the only string type available

### RNT

RNT stands for Runtime, this variable's type will calculated at runtime without typechecking

### VOID

Like in C or C++, Orth's **VOID** stands for everything that will not return anything.</br>
They are used mainly by these operations **"putui","puts","invoke","end"**, functions for side effects and so on

## Constants

As you may guess, orth has constants that store values. To create a variable, use the keyword `const` following by it's name, type and value:
```orth
const name = s "John"
const age = i 20
```

## changing a variable

Variables are just like _constants_ but they can be modified at runtime.
```orth
var some_number = i64 57

i 60 # new value for some_number
hold some_number set_number
hold some_number deref putui
```

As you can see, some_number can be changed by using the bultin instruction `set_number` that works for every non decimal number<br/>
It takes two parameters: a pointer to a variable and the new value (they should be on the stack)

## Conditions

Conditions in Orth are very simple and are made by the `if-end` blocks</br>
first: if you want something to be true or false, the last item on the stack **must** be a bool type.

```orth
i 10 i 10 + i 20 == putui
```

the code above will produce a bool type that can be used by if blocks

```orth
i 10 i 10 + i 20 ==
if 
    s "Yes, this is a true statement puts
else 
    s "Yes, this is a true statement puts
end
```

This code will execute the inner instructions within the if block because 10 + 10 = 20.</br>
Otherwise, the else block would be executed instead of the if block.</br>

## Loops

Orth has only the _while_ loop for now and it's very simple to use</br>
a while loop basic consists of a "until thisis true, then keep doing", that's basiclly what a loop looks in Orth

```orth
i 0 while dup i 10 > do
    dup putui s "\n" + puts
    i 1 +
end drop
```

So, as you can see we:

1. push an inital value (0 in this case)
2. duplicate it for better usage
3. pushes something to compare (10 in this case)

and until 10 is > 0 we add 1 to the last item on the stack, as simple as this

## Mem
Orth has a special way of using memory, by default you have a array of 640000 bytes (A LOT) and you can operate in this array by storing or reading values from this arrray.</br>
The _mem_ keyword essentially pushes a pointer pointing to the beginning of the array and can be offseted by adding a number to it:
```orth
mem i32 0 + i8 97 .
mem i32 1 + i8 98 .
mem i32 2 + i8 99 .

i 3 mem i 0 + dump_mem
```
Here we are storing the BYTES</br>
97 at offset 0</br>
98 at offset 1</br>
99 at offset 2</br>
In other words we are storing the string "abc" into the array.</br>
The dump_mem instruction takes 2 values, the amount of elements to print and the start point (which is mem offseted by 0 in this case) </br>
The result is:
```
abc
```