package main

import (
	"fmt"
	"os"
)

func main() {
	// var args []string
	args := os.Args

	// check number of arguments
	// slice are of dynamic length, whereas arrays have static length
	if len(args) < 2 {
		fmt.Printf("Usage: ./helloworld <arg1> ... <argN>\n")
		os.Exit(1)
	}

	// fmt.Println("hello world")
	// fmt.Printf("hello world\nArguments: %v\n", args) // %v is for value

	// printing all arguments
	// slices in golang var[n:m] are from n to m-1, starting at 0
	// args[0] is the program name
	fmt.Printf("hello world\nArguments: %v\n", args[1:]) // %v is for value
}
