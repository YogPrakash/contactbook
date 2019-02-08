package main

import (
	curdS "github.com/yogprakash/contactbook/contactBookService"
	l "log"
	"net/http"
)

func main() {
	l.Println("restCURDSearchApis service started on port :8080")

	db := curdS.DBSession()
	defer db.Close() // clean up when we’re done

	// Adapt our handle function using withDB and BasicAuth
	withDB := curdS.WithDB(db)
	auth := curdS.BasicAuth()
	h := curdS.Adapt(curdS.MakeHandler(), withDB, auth)

	mux := http.NewServeMux()
	mux.Handle("/restservice/", h)

	// start the server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		l.Fatal(err)
	}
}
