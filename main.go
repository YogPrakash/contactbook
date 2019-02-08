package main

import (
	l "log"
	"net/http"
	"os"
	"io"
	"fmt"
	//"io"
	"gopkg.in/mgo.v2"
	"crypto/tls"
	"time"
	"net"
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
	//http.HandleFunc("/test", helloWorldHandler)
	http.ListenAndServe(":"+port, h)


	//h := curdS.Adapt(curdS.MakeHandler(), auth)
	//fmt.Printf("%+v", h)
	mux := http.NewServeMux()
	mux.Handle("/restservice/", h)

	//// start the server
	//if err := http.ListenAndServe(":" + port, nil); err != nil {
	//	l.Fatal(err)
	//}
}


var tlsConfig = &tls.Config{}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {

	x := mgo.DialInfo{
		Addrs: []string{"cluster0-shard-00-00-awvwh.mongodb.net",
		"cluster0-shard-00-01-awvwh.mongodb.net:",
		"cluster0-shard-00-02-awvwh.mongodb.net"},
		Database:"admin",
		ReplicaSetName:"Cluster0-shard-0",
		Username:"yprakash",
		Password:"test123",
		FailFast:true,
		Timeout:time.Second * 5,
	}
	x.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	db, err := mgo.DialWithInfo(&x)
	if err != nil {
		l.Fatal("unable to connect to mongodb : ", err)
	}
	names, err := db.DatabaseNames()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(names)
	io.WriteString(w, "Hello world!")
}

