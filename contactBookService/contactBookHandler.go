package contactBookService

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func insertH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(dbSessionKey).(*mgo.Session)
	// decode the request body
	var cb contactBook
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	//normalize email to be saved only in lowercase
	if cb.Email != nil {
		*cb.Email = strings.ToLower(*cb.Email)
	} else {
		HTTPErrorResponse(w, http.StatusBadRequest, errors.New("email can not be empty"))
		return
	}

	// check if email exists to maintain the unique email
	if !emailExists(*cb.Email, db) {
		// give the contact info  a unique ID and set the time
		cb.ID = uuid.NewV1().String()
		crtDate := time.Now()
		cb.CreateDateTime = &crtDate
		cb.LastUpdatedDateTime = &crtDate
		// insert it into the database
		if err := db.DB(dbName).C(collectionName).Insert(&cb); err != nil {
			HTTPErrorResponse(w, http.StatusNotFound, err)
			return
		}
	} else {
		HTTPErrorResponse(w, http.StatusNotFound, errors.New("duplicate email, email already exists"))
		return
	}

	if err := db.DB(dbName).C(collectionName).Find(bson.M{"_id": cb.ID}).One(&cb); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}
	HTTPResponse(w, http.StatusOK, "inserted contact book  successfully", cb)
}

func readOneH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	var cb *contactBook
	vars := mux.Vars(r)
	id := vars["id"]

	if err := db.DB(dbName).C(collectionName).
		Find(bson.M{"_id": id}).One(&cb); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "contact book by id response", cb)
}

func readAllH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	fromPage := r.URL.Query().Get("fromPage")
	toPage := r.URL.Query().Get("toPage")

	//variable to store starting and end page for pagination
	var s, e int
	var err error
	if len(fromPage) > 0 {
		s, err = strconv.Atoi(fromPage)
		if err != nil {
			HTTPErrorResponse(w, http.StatusBadRequest, errors.New("failed to convert string to int"))
		}
	}

	if len(toPage) > 0 {
		e, err = strconv.Atoi(toPage)
		if err != nil {
			HTTPErrorResponse(w, http.StatusBadRequest, errors.New("failed to convert string to int"))
		}
	}

	var cbs []contactBook

	if s == 0 && e == 0 {
		e = 20
	}

	if err := db.DB(dbName).C(collectionName).
		Find(nil).Sort("-lastUpdatedDateTime").Skip(s).Limit(e).All(&cbs); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "contact book list response", cbs)
}

func patchH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	var cb *contactBook
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	findQ := bson.M{"_id": id}
	updateQ := cb.updateQ()

	// update database for given contact book id
	if err := db.DB(dbName).C(collectionName).Update(findQ, updateQ); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	//fetch the updated doc from db to return as response
	if err := db.DB(dbName).C(collectionName).
		Find(findQ).One(&cb); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "updated contact book", cb)
}

func deleteH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	vars := mux.Vars(r)
	id := vars["id"]
	findQ := bson.M{"_id": id}

	//remove the doc from db for given contact book id
	if err := db.DB(dbName).C(collectionName).Remove(findQ); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}
	msg := fmt.Sprintf("deleted the  contact book document with ID: %v", id)
	HTTPResponse(w, http.StatusOK, msg, nil)
}

func searchH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)

	vars := mux.Vars(r)
	query := vars["query"]
	fromPage := r.URL.Query().Get("fromPage")
	toPage := r.URL.Query().Get("toPage")

	//variable to store starting and end page for pagination
	var s, e int
	var err error
	if len(fromPage) > 0 {
		s, err = strconv.Atoi(fromPage)
		if err != nil {
			HTTPErrorResponse(w, http.StatusBadRequest, errors.New("failed to convert string to int"))
		}
	}

	if len(toPage) > 0 {
		e, err = strconv.Atoi(toPage)
		if err != nil {
			HTTPErrorResponse(w, http.StatusBadRequest, errors.New("failed to convert string to int"))
		}
	}

	if len(query) == 0 {
		err := fmt.Errorf("query can not be empty")
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if s == 0 && e == 0 {
		e = 10
	}

	var cbs []contactBook
	findQ := bson.M{}
	orQ := []bson.M{}
	orQ = append(orQ, bson.M{"lastName": bson.M{"$regex": query + ".*", "$options": "i"}},
		bson.M{"email": bson.M{"$regex": query + ".*.", "$options": "i"}},
		bson.M{"firstName": bson.M{"$regex": query + ".*", "$options": "i"}})

	findQ["$or"] = orQ
	if err := db.DB(dbName).C(collectionName).
		Find(findQ).Sort("-lastUpdatedDateTime").Skip(s).Limit(e).All(&cbs); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "contact book search list response", cbs)
}
