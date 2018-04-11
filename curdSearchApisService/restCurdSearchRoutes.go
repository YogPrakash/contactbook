package curdSearchApisService

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/restservice/recipe", handleInsert).Methods(http.MethodPost)
	r.HandleFunc("/restservice/recipe/{id}", readOneH).Methods(http.MethodGet)
	r.HandleFunc("/restservice/recipe", readAllH).Methods(http.MethodGet)
	r.HandleFunc("/restservice/recipe/{id}", patchH).Methods(http.MethodPatch)
	r.HandleFunc("/restservice/recipe/{id}", deleteH).Methods(http.MethodDelete)
	r.HandleFunc("/restservice/recipe/{id}/rating", ratingH).Methods(http.MethodPost)
	r.HandleFunc("/restservice/recipe/search/", searchH).Methods(http.MethodGet)
	return r
}
