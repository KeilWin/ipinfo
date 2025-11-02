package handler

import "net/http"

func newHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{'a', 'l', 'i', 'v', 'e'})
	}
}
