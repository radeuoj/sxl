# SXL Programming Language

This is an almost identical interpreter (at least at the time of writing this readme) with the one from [Thorsten Ball's book](https://interpreterbook.com/) which I read in order to build this.

```rs
println("Hello world");
```

## Future plans for the language
- Imports, `import "file.sxl"` syntax
- Rewrite it in Rust
- Add a stack memory model
- Compile to C


## Language features

More examples in the `/sxl` folder

```rs
// Variables
let a = 10;
a = 20 * 39;
inspect(a + 130); // this prints the argument's value
// it's basically an export of the (*Node).Inspect() function

// If statements
let a = 1;
let b = 2;

if a > b {
    inspect(a);
} else if a == b {
    inspect(0);
} else {
    inspect(b);
}

// User defined functions
fn fib(n) {
    if n <= 1 {
        return n;
    } else {
        return fib(n - 2) + fib(n - 1);
    }
}

fib(35); // this is ridiculously slow (15s) vs sub 100ms on rust

// Strings
fn hello(who) {
    "Hello " + who;
}

println(hello("world"));
println(hello("Marcus"));
```

## How to run a file
```
$ go run . sxl/println.sxl
```

---

That's about it, I hope I'll have time to improve the language further.
