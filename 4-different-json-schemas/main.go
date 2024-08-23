package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Page struct {
	Name string `json:"page"`
}

// {"page":"words","input":"","words":["word1"]}
// output of "$go run main.go http://localhost:8080/words?input=word1"
type Words struct {
	// to add metadata that parses the JSON you need ``
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

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

	// parse data if response code is 200
	if res.StatusCode != 200 {
		fmt.Printf("Invalid output (HTTP Code: %d): %s\n", res.StatusCode, body)
		os.Exit(1)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}

		// strings.Join separates elements of an array/slice with the separator you want
		fmt.Printf("JSON Parsed\nPage: %s\nWords: %v\n", page.Name, strings.Join(words.Words, ", "))
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			log.Fatal(err)
		}

		// print occurence map (MAPS ARE NOT ORDERED, UNLIKE ARRAYS: ouput will be different every time)
		for word, occ := range occurrence.Words {
			fmt.Printf("Word: %s; Occurrence: %d\n", word, occ)
		}

	default:
		fmt.Printf("Page not found!\n")
	}

}
