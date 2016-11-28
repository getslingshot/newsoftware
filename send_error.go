package main

// E is an error type to respond with
type E struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors"`
}

func sendError() E {
	m := make(map[string]string)
	return E{"", m}
}
