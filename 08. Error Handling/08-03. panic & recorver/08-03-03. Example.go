package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok:=recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

func handle(w http.ResponseWriter, r * http.Request) {
	v := r.URL.Query()
	a, _ := strconv.Atoi(v.Get("dividend"))
	b, _ := strconv.Atoi(v.Get("divisor"))
	fmt.Fprintf(w, "%d / %d = %d", a, b, a/b)
}

func main() {
	http.HandleFunc("/divide", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(handle)(w, r)
	})
	http.ListenAndServe(":8080", nil)
}