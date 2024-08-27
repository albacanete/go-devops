package main

import "net/http"

type MyJWTTransport struct {
	transport http.RoundTripper // so we can use the default Transport in the main function
	token     string
}

func (m MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.token != "" {
		req.Header.Add("Authorization", "Bearer "+m.token)
	}
	return m.transport.RoundTrip(req)
}
