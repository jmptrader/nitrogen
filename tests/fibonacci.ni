/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates recursion using the Fibonacci sequence.
 */

func fib(x) {
     if (x == 0) {
        return 0;
     }

     if (x == 1) {
        return 1;
     }

     return fib(x-1) + fib(x-2);
}

func main() {
    let fibTest = fib(10);
    if (fibTest != 55) {
        println("Fibonacci is broken!");
        print("Got: ");
        print(fibTest);
        print(", Expected: 55");
    }
}

main();
