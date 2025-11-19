package handler

import (
	"encoding/json"
	"net/http"
)

func NewHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		healthData := NewHealthData()
		okResponse := NewOkResponse(healthData)
		res, err := json.Marshal(okResponse)
		if err != nil {
			WriteInternalServerError(w)
			return
		}
		w.Write(res)
	}
}
