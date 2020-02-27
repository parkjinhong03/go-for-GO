package handlers

import (
	"encoding/json"
	"net/http"
)

type SearchHandler struct {}

type searchRequest struct {
	Query string `json:"query"`
}

func (h SearchHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(searchRequest)
	err := decoder.Decode(request)

	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}
}