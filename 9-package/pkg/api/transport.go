package api

import (
	"net/http"
)

type MyJWTTransport struct {
	transport http.RoundTripper // so we can use the default Transport in the main function
	token     string
	password  string
	loginURL  string
}

// MyJWTTransport is a pointer because we are making changes to a variable in the struct and do not want to instanciate it every time
func (m *MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.token == "" {
		if m.password != "" {
			token, err := doLoginRequest(http.Client{}, m.loginURL, m.password)
			if err != nil {
				return nil, err
			}
			m.token = token
		}
	}
	if m.token != "" {
		req.Header.Add("Authorization", "Bearer "+m.token)
	}
	return m.transport.RoundTrip(req)
}
