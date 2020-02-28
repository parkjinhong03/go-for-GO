package handlers

import (
	"../data"
	"encoding/json"
	"net/http"
)

// 새끼 고양이 정보를 반환하는 기능의 핸들러이다.
// 필드에 검색 기능을 가지고 있는 데이터 저장소 객체를 가지고 있다.
type SearchHandler struct {
	DataStore data.Store
}

type searchRequest struct {
	Query string `json:"query"`
}

type searchResponse struct {
	Kittens []data.Kitten `json:"kittens"`
}

func (h SearchHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(searchRequest)
	err := decoder.Decode(request)

	if err != nil || len(request.Query) < 1 {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	kittens := h.DataStore.Search(request.Query)

	encoder := json.NewEncoder(rw)
	encoder.Encode(searchResponse{Kittens: kittens})
}