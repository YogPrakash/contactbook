package contactBookService

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/cb_service/contact_book", insertH).Methods(http.MethodPost)
	r.HandleFunc("/cb_service/contact_book/{id}", readOneH).Methods(http.MethodGet)
	r.HandleFunc("/cb_service/contact_book", readAllH).Methods(http.MethodGet)
	r.HandleFunc("/cb_service/contact_book/{id}", patchH).Methods(http.MethodPatch)
	r.HandleFunc("/cb_service/contact_book/{id}", deleteH).Methods(http.MethodDelete)
	r.HandleFunc("/cb_service/contact_book/search/{query}", searchH).Methods(http.MethodGet)
	return r
}
