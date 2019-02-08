package main

import (
	l "log"
	"net/http"
	"os"
	curdS "github.com/yogprakash/contactbook/contactBookService"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		l.Fatal("$PORT must be set")
	}
	
	l.Println("restCURDSearchApis service started on port", port)

	db := curdS.DBSession()
	defer db.Close() // clean up when we’re done

	// Adapt our handle function using withDB and BasicAuth
	withDB := curdS.WithDB(db)
	auth := curdS.BasicAuth()
	h := curdS.Adapt(curdS.MakeHandler(), withDB, auth)

	mux := http.NewServeMux()
	mux.Handle("/restservice/", h)

	// start the server
	if err := http.ListenAndServe(":" + port, mux); err != nil {
		l.Fatal(err)
	}
}
