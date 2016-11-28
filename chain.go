package main

import "net/http"

// Adapter function to be used on middlewares
type Adapter func(http.Handler) http.Handler

func chain(handler http.Handler, adapters ...Adapter) http.Handler {

	for i := len(adapters) - 1; i >= 0; i-- {
		handler = adapters[i](handler)
	}

	return handler
}
