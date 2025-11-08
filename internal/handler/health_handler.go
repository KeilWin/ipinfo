package handler

import "net/http"

func NewHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{'a', 'l', 'i', 'v', 'e'})
	}
}
