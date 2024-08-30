package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Assignment1 struct {
	Name         string             `json:"page"`
	Words        []string           `json:"words"`
	Percentatges map[string]float32 `json:"percentages"`
	Special      []string           `json:"special"`
	ExtraSpecial []any              `json:"extraSpecial"`
}

func main() {
	var (
		requestURL string
		parsedURL  *url.URL
		err        error
	)

	flag.StringVar(&requestURL, "url", "", "url to access")
	flag.Parse()

	// check if arg is a valid url
	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Unvalid URL error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	res, err := http.Get(parsedURL.String())
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		fmt.Printf("Invalid output (HTTP Code: %d): %s\n", res.StatusCode, body)
		os.Exit(1)
	}

	var assignment1 Assignment1
	err = json.Unmarshal(body, &assignment1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON Parsed for Page: %s\n", assignment1.Name)
	// print words
	fmt.Printf("Words: %v\n", strings.Join(assignment1.Words, ", "))

	// print percentages
	for word, per := range assignment1.Percentatges {
		fmt.Printf("Word: %s; Percentage: %f\n", word, per)
	}

	// print special
	fmt.Printf("Special words: %v\n", strings.Join(assignment1.Special, ", "))

	// print extraSpecial
	fmt.Printf("ExtraSpecial words: %v\n", fmt.Sprintf("%v", assignment1.ExtraSpecial))
}
