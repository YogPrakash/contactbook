package contactBookService

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	l "log"
	"net"
	"net/http"
	"strings"
	"time"
)

var tlsConfig = &tls.Config{}

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func WithDB(db *mgo.Session) Adapter {
	// return the Adapter
	return func(next http.Handler) http.Handler {
		// the adapter (when called) should return a new handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// copy the database session
			dbsession := db.Copy()
			defer dbsession.Close() // clean up
			// save it in the request context
			ctx := context.WithValue(r.Context(), "database", dbsession)
			r = r.WithContext(ctx)
			//// pass execution to the original handler
			next.ServeHTTP(w, r)
		})
	}
}
func BasicAuth() Adapter {
	// return the Adapter
	return func(next http.Handler) http.Handler {
		//the adapter (when called) should return a new handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//compare the request url and method type
			if strings.Contains(strings.TrimSpace(r.URL.Path), "/cb_service/") {
				//generic error message for invalid autherization
				err := fmt.Errorf("User Not Authorized")
				username, password, authOK := r.BasicAuth()
				if authOK == false {
					HTTPErrorResponse(w, http.StatusUnauthorized, err)
					return
				}

				if username != "username" || password != "password" {
					HTTPErrorResponse(w, http.StatusUnauthorized, err)
					return
				}
			}
			// pass execution to the original handler
			next.ServeHTTP(w, r)
		})
	}
}

// create a db session
func DBSession() *mgo.Session {
	// connect to the database
	x := mgo.DialInfo{
		Addrs: []string{"cluster0-shard-00-00-awvwh.mongodb.net",
			"cluster0-shard-00-01-awvwh.mongodb.net:",
			"cluster0-shard-00-02-awvwh.mongodb.net"},
		Database:       "admin",
		ReplicaSetName: "Cluster0-shard-0",
		Username:       "yprakash",
		Password:       "test123",
		FailFast:       true,
		Timeout:        time.Second * 30,
	}
	x.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	db, err := mgo.DialWithInfo(&x)
	if err != nil {
		l.Fatal("unable to connect to mongodb : ", err)
	}
	return db
}

// MetaData of HTTP API response
type MetaData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//Response -  structure of  HTTP response Meta + Data
type Response struct {
	Meta MetaData    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// HTTPResponse writes the HTTPResponse and renders the json
func HTTPResponse(w http.ResponseWriter, statusCode int, msg string, data interface{}) {
	renderer := render.New()
	res := Response{}
	res.Meta.Code = statusCode
	res.Meta.Msg = msg
	res.Data = data
	renderer.JSON(w, statusCode, res)
}

// HTTPErrorResponse writes the HTTPErrorResponse and renders the json
func HTTPErrorResponse(w http.ResponseWriter, errorCode int, err error) {
	renderer := render.New()
	res := Response{}
	res.Meta.Code = errorCode
	res.Meta.Msg = err.Error()
	l.Printf("error:%v", err)
	renderer.JSON(w, errorCode, res)
}
