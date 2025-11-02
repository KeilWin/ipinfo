package handler

import "net/http"

func NewIpV4Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{'i', 'p', 'v', '4'})
	}
}
