package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}

	// check if arg is a valid url
	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is invalid: %s\n", err)
		os.Exit(1)
	}

	// make http request
	// res is *http.Response
	res, err := http.Get(args[1])
	if err != nil {
		log.Fatal(err)
	}

	// close the response body, since it is streamed on demand
	defer res.Body.Close()

	// now we're ready to read the body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTTP Status code: %d\nBody: %s", res.StatusCode, body)

	// parse the body
	// {"page":"words","input":"","words":["word1"]}

}
