package main

import (
	"encoding/json"
	"flag"
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
	return fmt.Sprintf("Words: %s", strings.Join(w.Words, ", "))
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
	return fmt.Sprintf("Words: %s", strings.Join(out, ", "))
}

func main() {
	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	// Alternative to flag: https://github.com/spf13/cobra (allows optional params)
	flag.StringVar(&requestURL, "url", "", "url to access")
	flag.StringVar(&password, "password", "", "password to access our api")

	flag.Parse()

	// check if arg is a valid url
	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Unvalid URL error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	client := http.Client{}

	if password != "" {
		token, err := doLoginRequest(client, parsedURL.Scheme+"://"+parsedURL.Host+"/login", password)
		if err != nil {
			if requestErr, ok := err.(RequestError); ok {
				fmt.Printf("Error: %s (HTTP Code: %d, Body: %s\n)", requestErr.Error(), requestErr.HTTPCode, requestErr.Body)
				os.Exit(1)
			}
			fmt.Printf("Login failed: %s\n", err)
			os.Exit(1)
		}
		client.Transport = MyJWTTransport{
			transport: http.DefaultTransport,
			token:     token,
		}
	}

	res, err := doRequest(client, parsedURL.String())
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s\n)", requestErr.Error(), requestErr.HTTPCode, requestErr.Body)
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

func doRequest(client http.Client, requestURL string) (Response, error) {
	// make http request
	res, err := client.Get(requestURL)
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
		return nil, fmt.Errorf("invalid output (HTTP Code: %d): %s", res.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: res.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("invalid JSON body error"),
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
