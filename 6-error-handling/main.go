package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Response can either be Words or Occurrence
type Response interface {
	GetResponse() string
}

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

// add GetResponse function to Words struct
func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

// add GetResponse function to Occurrence struct
func (o Occurrence) GetResponse() string {
	out := []string{}
	for word, occ := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occ))
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}

	res, err := doRequest(args[1])
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s\n)", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error occurred: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}

func doRequest(requestURL string) (Response, error) {
	// check if arg is a valid url
	if _, err := url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("Unvalid URL error: %s", err)
	}

	// make http request
	res, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %s", err)
	}

	// close the response body, since it is streamed on demand
	defer res.Body.Close()

	// now we're ready to read the body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	// parse data if response code is 200
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output (HTTP Code: %d): %s\n", res.StatusCode, body)
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: res.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Invalid JSON body error"),
		}
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: res.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page unmarsall error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: res.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words unmarsall error: %s", err),
			}
		}
		return words, nil

	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: res.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Unmarsall error: %s", err),
			}
		}
		return occurrence, nil

	default:
		return nil, nil
	}

}
