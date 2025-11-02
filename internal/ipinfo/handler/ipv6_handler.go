package handler

import "net/http"

func NewIpV6Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{'i', 'p', 'v', '6'})
	}
}
