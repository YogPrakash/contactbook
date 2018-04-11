package curdSearchApisService

import (
	"bytes"
	"context"
	"encoding/json"
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
	recipeID          = "1234u51234567"
	incorrectRecipeID = "@&ASDFGE@3244"
)

func SetupData(session *mgo.Session) {
	var rec recipe
	rec.ID = recipeID
	rec.Name = "chicken"
	diff := 2
	rec.Difficulty = &diff
	prepTime := time.Now()
	rec.PrepTime = &prepTime
	veg := false
	rec.Vegetarian = &veg
	// insert it into the database
	if err := session.DB(dbName).C(collectionName).Insert(&rec); err != nil {
		l.Printf("error inserting in db error:%v", err)
	}

}
func ClearData(session *mgo.Session) {
	findQ := bson.M{"_id": recipeID}
	if err := session.DB(dbName).C(collectionName).Remove(findQ); err != nil {
		l.Printf("error inserting in db error:%v", err)
	}
}

//************************** handleInsert starts *************************************
func TestInserRecipeCorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	res := httptest.NewRecorder()
	correctPayload := `{
    		"name":"dal makhani",
			"difficulty":1,
			"vegetarian":true
	}`
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodPost, "/restservice/recipe", strings.NewReader(correctPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx := context.WithValue(req.Context(), "database", db)
	req = req.WithContext(ctx)
	handleInsert(res, req)
	a.Equal(res.Code, http.StatusOK)
}

func TestInserRecipeIncorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	res := httptest.NewRecorder()
	incorrectPayload := `{
    		"name":"dal makhani",
			"difficulty":"1",
			"vegetarian":true
	}`
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodPost, "/restservice/recipe", strings.NewReader(incorrectPayload))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	handleInsert(res, req)
	a.Equal(res.Code, http.StatusBadRequest)
}

//************************** handleInsert ends *************************************

//************************** readOneH starts *************************************
func TestReadOnecorrectID(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done
	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/restservice/recipe/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), "database", db)
		req.Header.Set("Content-Type", contentType)
		req = req.WithContext(ctx)
		readOneH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/restservice/recipe/" + recipeID
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
	router.HandleFunc("/restservice/recipe/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req = req.WithContext(ctx)
		readOneH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/restservice/recipe/" + incorrectRecipeID
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
	req, err := http.NewRequest(http.MethodGet, "/restservice/recipe", nil)
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
	req, err := http.NewRequest(http.MethodGet, "/restservice/recipe?fromPage='!@#$'", nil)
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
	res := httptest.NewRecorder()
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodGet, "/restservice/search/?query=ch", nil)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	searchH(res, req)
	a.Equal(res.Code, http.StatusOK)
	ClearData(db)
}

func TestSearchIncorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)
	res := httptest.NewRecorder()
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodGet, "/restservice/search/?fromPage='!@#$'", nil)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	searchH(res, req)
	a.Equal(res.Code, http.StatusBadRequest)
	ClearData(db)
}

//************************** Search ends *************************************

//************************** rating starts *************************************
func TestRatingCorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)
	ratingReq := recipeRatingReq{}
	rate := 5
	ratingReq.Rating = &rate
	payload, err := json.Marshal(ratingReq)
	if err != nil {
		l.Printf("failed to marshal the payload :%v ", err)
	}
	res := httptest.NewRecorder()
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodPost, "/restservice/recipe/"+recipeID+"/rating", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	ratingH(res, req)
	a.Equal(res.Code, http.StatusNotFound)
	ClearData(db)
}

func TestRatingIncorrectPayload(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done

	SetupData(db)
	ratingReq := recipeRatingReq{}
	rate := 5
	ratingReq.Rating = &rate
	payload, err := json.Marshal(ratingReq)
	if err != nil {
		l.Printf("failed to marshal the payload :%v ", err)
	}
	res := httptest.NewRecorder()
	a := assert.New(t)
	req, err := http.NewRequest(http.MethodPost, "/restservice/recipe/"+incorrectRecipeID+"/rating", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", encodedAuthToken)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}
	// save it in the request context
	ctx1 := context.WithValue(req.Context(), dbSessionKey, db)
	req = req.WithContext(ctx1)
	ratingH(res, req)
	a.Equal(res.Code, http.StatusNotFound)
	ClearData(db)
}

//************************** rating ends *************************************

//************************** readOneH starts *************************************
func TestDeleteHCorrectID(t *testing.T) {
	db := DBSession()
	defer db.Close() // clean up when we’re done
	SetupData(db)

	a := assert.New(t)
	router := mux.NewRouter()
	router.HandleFunc("/restservice/recipe/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req = req.WithContext(ctx)
		deleteH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/restservice/recipe/" + recipeID
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
	router.HandleFunc("/restservice/recipe/{id}", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// save it in the request context
		ctx := context.WithValue(req.Context(), dbSessionKey, db)
		req.Header.Set("Content-Type", contentType)
		req = req.WithContext(ctx)
		deleteH(res, req)
	}))

	server := httptest.NewServer(router)
	defer server.Close()
	reqURL := server.URL + "/restservice/recipe/" + incorrectRecipeID
	res, err := http.Get(reqURL)
	if err != nil {
		l.Printf("Cannot Make Request :%v ", err)
		a.Error(err)
	}

	a.Equal(res.StatusCode, http.StatusNotFound)
	ClearData(db)
}

//************************** readOneH ends *************************************
