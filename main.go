package main

import (
	l "log"
	"net/http"
	"os"
	"io"
	//"fmt"
	//"io"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		l.Fatal("$PORT must be set")
	}
	
	l.Println("restCURDSearchApis service started on port", port)
	//
	//db := curdS.DBSession()
	//defer db.Close() // clean up when we’re done
	//
	//// Adapt our handle function using withDB and BasicAuth
	//withDB := curdS.WithDB(db)
	//auth := curdS.BasicAuth()
	////h := curdS.Adapt(curdS.MakeHandler(), withDB, auth)
	http.HandleFunc("/test", helloWorldHandler)
	//http.ListenAndServe(":80", nil)


	//h := curdS.Adapt(curdS.MakeHandler(), auth)
	//fmt.Printf("%+v", h)
	//mux := http.NewServeMux()
	//mux.Handle("/restservice/", h)

	// start the server
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		l.Fatal(err)
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

