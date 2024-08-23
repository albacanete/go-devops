package main

import (
	"fmt"
	"io"
	"log"
)

type MySlowReader struct {
	Contents string
}

// MySlowReader needs a Read function in other to be a Reader and be able to pass it to io.ReadAll
func (m MySlowReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func main() {

	mySlowReaderInstance := MySlowReader{
		Contents: "hello world!",
	}

	out, err := io.ReadAll(mySlowReaderInstance)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("output: %s\n", out)
}
