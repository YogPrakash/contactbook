package contactBookService

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	l "log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	contactBookID          = "1234u51234567"
	incorrectContactBookID = "@&ASDFGE@3244"
)

func SetupData(session *mgo.Session) {
	var cb contactBook
	cb.ID = contactBookID
	cb.FirstName = "Yog"
	cb.LastName = "Prakash"
	cb.UserID = "U111"
	cb.GroupID = "G11"
	email := "testing@gmail.com"
	cb.Email = &email
	createTime := time.Now()
	cb.CreateDateTime = &createTime
	cb.LastUpdatedDateTime = &createTime
	cb.DocumentVersion = &docVersion
	cb.IsActive = &trueVal
	// insert it into the database
	if err := session.DB(dbName).C(collectionName).Insert(&cb); err != nil {
		l.Printf("error inserting in db error:%v", err)
	}

}
func ClearData(session *mgo.Session) {
	findQ := bson.M{"_id": contactBookID}
	if err := session.DB(dbName).C(collectionName).Remove(findQ); err != nil {
		l.Printf("error inserting in db error:%v", err)
	}
}

//************************** insertH starts *************************************
func TestInserContactCorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	res := httptest.NewRecorder()
	correctPayload := `{
	"userID": "U2",
  	"groupID": "G2",
  	"firstName": "Yog1",
  	"lastName": "Prakash1",
  	"email": "test1@gmail.com",
  	"contact": [
    	{
     	 "type": 0,
	     "number": "2234567890",
         "countryCode": "91",
         "isPrimary": true
    	}
  	],
  	"notes": "prime user 1",
  	"lastUpdatedByUser": "Yog"
	}`
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodPost, "/cb_service/contact_book", strings.NewReader(correctPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx := context.WithValue(req.Context(), "database", db)
	req = req.WithContext(ctx)
	insertH(res, req)
	a.Equal(res.Code, http.StatusOK)
}

func TestInserContactIncorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	res := httptest.NewRecorder()
	incorrectPayload := `{
    			"userID": "U2",
  	"groupID": 1,
  	"firstName": "Yog1",
  	"lastName": "Prakash1",
  	"email": "test1@gmail.com",
  	"contact": [
    	{
     	 "type": 0,
	     "number": "2234567890",
         "countryCode": "91",
         "isPrimary": true
    	}
  	],
  	"notes": "prime user 1",
  	"lastUpdatedByUser": "Yog"
	}`
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodPost, "/cb_service/contact_book", strings.NewReader(incorrectPayload))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	insertH(res, req)
	a.Equal(res.Code, http.StatusBadRequest)
}

//************************** insertH ends *************************************

//************************** readOneH starts *************************************
func TestReadOnecorrectID(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done
	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/cb_service/contact_book/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), "database", db)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", encodedAuthToken)
		req = req.WithContext(ctx)
		readOneH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/cb_service/contact_book/" + contactBookID
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}

	a.Equal(res.StatusCode, http.StatusOK)
	ClearData(db)
}

func TestReadOneIncorrectID(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done
	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/cb_service/contact_book/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", encodedAuthToken)
		req = req.WithContext(ctx)
		readOneH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/cb_service/contact_book/" + incorrectContactBookID
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}

	a.Equal(res.StatusCode, http.StatusNotFound)
	ClearData(db)
}

//************************** readOneH ends *************************************

//************************** readAll starts *************************************
func TestReadAllCorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)
	res := httptest.NewRecorder()
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodGet, "/cb_service/contact_book", nil)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx)
	readAllH(res, req)
	a.Equal(res.Code, http.StatusOK)
	ClearData(db)
}

func TestReadAllIncorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)
	res := httptest.NewRecorder()
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodGet, "/cb_service/contact_book?fromPage='!@#$'", nil)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	readAllH(res, req)
	a.Equal(res.Code, http.StatusBadRequest)
	ClearData(db)
}

//************************** readAll ends *************************************
//************************** Search starts *************************************
func TestSearchCorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)
	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/cb_service/contact_book/search/{query}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", encodedAuthToken)
		req = req.WithContext(ctx)
		searchH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/cb_service/contact_book/search/Yog"
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}

	a.Equal(res.StatusCode, http.StatusOK)
	ClearData(db)
}

func TestSearchIncorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/cb_service/contact_book/search/{query}?fromPage='!@#$'", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", encodedAuthToken)
		req = req.WithContext(ctx)
		searchH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/cb_service/contact_book/search/ggghhg?fromPage='!@#$'"
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	a.Equal(res.StatusCode, http.StatusNotFound)
	ClearData(db)
}

//************************** Search ends *************************************

//************************** readOneH starts *************************************
func TestDeleteHCorrectID(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done
	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/cb_service/contact_book/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", encodedAuthToken)
		req = req.WithContext(ctx)
		deleteH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/cb_service/contact_book/" + contactBookID
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}

	a.Equal(res.StatusCode, http.StatusOK)
	ClearData(db)
}

func TestDeleteHIncorrectID(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done
	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/cb_service/contact_book/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", encodedAuthToken)
		req = req.WithContext(ctx)
		deleteH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/cb_service/contact_book/" + incorrectContactBookID
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}

	a.Equal(res.StatusCode, http.StatusNotFound)
	ClearData(db)
}

//************************** readOneH ends *************************************
