package main

import (
	curdS "github.com/yogprakash/contactbook/contactBookService"
	l "log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	l.Println("contact book service started on port", port)

	db := curdS.DBSession()
	defer db.Close() // clean up when we’re done

	// Adapt our handle function using withDB and BasicAuth
	withDB := curdS.WithDB(db)
	auth := curdS.BasicAuth()
	h := curdS.Adapt(curdS.MakeHandler(), withDB, auth)
	http.ListenAndServe(":"+port, h)
	mux := http.NewServeMux()
	mux.Handle("/cb_service/", h)

	// start the server
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		l.Fatal(err)
	}
}
