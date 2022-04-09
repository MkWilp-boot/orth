[o] new flag `-d`:
    `-d` should dump the generated VM code
[o] new intruction `print`:
    Dump should dump the entire object while Print will print the operand only
[o] new instruction `while`:
    while should iterate to a certain amount of times decided by the last item in the stack
[o] orth should only parse `.orth` files and warn if the read file is not of type `orth`
[o] maybe os constant values for orthtypes like `orthtypes.RNT, orthtypes.I64, orthtypes.f32` etc..
[o] new instruction `drop`:
    drops the last element in the stack
[o] new instruction `swap`:
    swap the position of the last 2 elements in the stack
[x] new storage `memory`:
    creates a `heap` like storage for more control
[x] new instruction `mem x ,`:
    reads an element at position `x` and puts them in the stack
[x] new instruction `mem x obj .`:
    stores an element `obj` at position `x` and puts them in the `heap`
[o] new instruction `%`:
    takes the modulo of the last 2 items in the stack
[x] when compiling to VM code, in case of unknow tokens, print them all and exit with error