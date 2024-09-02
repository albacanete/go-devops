package main

import "fmt"

func main() {
	a := "hola"
	testPointer(a)
	fmt.Printf("a: %s\n", a) // prints "hola"

	testPointer2(&a)
	fmt.Printf("a: %v\n", a)  // prints "adeu"
	fmt.Printf("a: %v\n", &a) // prints memory direction

	b := []string{"quetal"}
	testPointerSlice(b)
	fmt.Printf("b: %v\n", b) // prints "quepasa" SLICES WORK AS POINTERS

	testPointerAppend(&b)
	fmt.Printf("b: %v\n", b) // prints "quepasa" if not passed by value

	c := testPointerAppend2(b)
	fmt.Printf("c: %v\n", c) // prints "quepasa moltbe"

	d := make(map[string]string)
	d["a"] = "reveure"
	testPointerMap(d)
	fmt.Printf("d: %v\n", d)
}

func testPointer(a string) {
	a = "adeu"
}

func testPointer2(a *string) {
	*a = "adeu"
}

func testPointerSlice(a []string) {
	a[0] = "quepasa"
}

func testPointerAppend(a *[]string) {
	*a = append(*a, "moltbe")
}

func testPointerAppend2(a []string) []string {
	a = append(a, "estupendo")
	return a
}

func testPointerMap(a map[string]string) {
	a["fins"] = "aviat"
}
