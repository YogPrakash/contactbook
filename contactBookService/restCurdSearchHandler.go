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
	"time"
)

func handleInsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(dbSessionKey).(*mgo.Session)
	// decode the request body
	var recipe recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	// give the recipe a unique ID and set the time
	recipe.ID = uuid.NewV1().String()
	prepTime := time.Now()
	recipe.PrepTime = &prepTime
	// insert it into the database
	if err := db.DB(dbName).C(collectionName).Insert(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}
	HTTPResponse(w, http.StatusOK, "inserted recipe successfully", recipe)
}

func readOneH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	var recipe *recipe
	vars := mux.Vars(r)
	id := vars["id"]

	if err := db.DB(dbName).C(collectionName).
		Find(bson.M{"_id": id}).One(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "recipe by id response", recipe)
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

	var recipe []*recipe

	if err := db.DB(dbName).C(collectionName).
		Find(nil).Sort("-prepTime").Skip(s).Limit(e).All(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "recipe list response", recipe)
}

func patchH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	var recipe *recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	findQ := bson.M{"_id": id}
	updateQ := recipe.updateQ()

	// update database for given recipe id
	if err := db.DB(dbName).C(collectionName).Update(findQ, updateQ); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	//fetch the updated doc from db to return as response
	if err := db.DB(dbName).C(collectionName).
		Find(findQ).One(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "updated recipe", recipe)
}

func deleteH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	vars := mux.Vars(r)
	id := vars["id"]
	findQ := bson.M{"_id": id}

	//remove the doc from db for given recipe id
	if err := db.DB(dbName).C(collectionName).Remove(findQ); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}
	msg := fmt.Sprintf("deleted the  recipe document with ID: %v", id)
	HTTPResponse(w, http.StatusOK, msg, nil)
}

func ratingH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)
	vars := mux.Vars(r)
	id := vars["id"]

	var payload *recipeRatingReq
	var recipe *recipe

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if payload.Rating == nil {
		err := fmt.Errorf("invalid request rating can't be nil")
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	findQ := bson.M{"_id": id}
	if err := db.DB(dbName).C(collectionName).
		Find(findQ).One(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	if recipe.Rating != nil {
		errMsg := fmt.Errorf("rating for recipe with id : %v is already done", id)
		HTTPErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	recipe.Rating = payload.Rating
	updateQ := recipe.updateQ()
	if err := db.DB(dbName).C(collectionName).Update(findQ, updateQ); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}
	HTTPResponse(w, http.StatusOK, "rated  recipe successfully", recipe)
}

func searchH(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value("database").(*mgo.Session)

	query := r.URL.Query().Get("query")
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
		err := fmt.Errorf("failed to convert string to int")
		HTTPErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var recipe []*recipe
	findQ := bson.M{}
	findQ["name"] = bson.M{"$regex": query + ".*", "$options": "i"}
	if err := db.DB(dbName).C(collectionName).
		Find(findQ).Sort("-prepTime").Skip(s).Limit(e).All(&recipe); err != nil {
		HTTPErrorResponse(w, http.StatusNotFound, err)
		return
	}

	HTTPResponse(w, http.StatusOK, "recipe search list response", recipe)
}

func ffff(w http.ResponseWriter, r *http.Request) {
	HTTPResponse(w, http.StatusOK, "recipe search list response", "hello")
}
