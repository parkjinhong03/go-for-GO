package handlers

import (
	"net/http"
)

type bangHandler struct {}

func (h *bangHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	panic("Somethings gone wrong again")
}

func NewBangHandler() *bangHandler {
	return &bangHandler{}
}