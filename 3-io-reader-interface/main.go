package main

import (
	"fmt"
	"io"
	"log"
)

type MySlowReader struct {
	// variables can be lowercase because we don't want to export them
	contents string
	pos      int
}

// MySlowReader needs a Read function in other to be a Reader and be able to pass it to io.ReadAll
// need to pass var by reference (pointer) because if not, m.pos++ takes no effect
func (m *MySlowReader) Read(p []byte) (n int, err error) {
	if m.pos+1 <= len(m.contents) { // Pos starts at 0
		elems := copy(p, m.contents[m.pos:m.pos+1])
		m.pos++
		return elems, nil
	}
	return 0, io.EOF
}

func main() {

	mySlowReaderInstance := &MySlowReader{
		contents: "hello world!",
		pos:      0,
	}

	out, err := io.ReadAll(mySlowReaderInstance)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("output: %s\n", out)
}
