package curdSearchApisService

import "time"

type recipe struct {
	ID         string     `bson:"_id" json:"recipeID"`
	Name       string     `bson:"name" json:"name"`
	PrepTime   *time.Time `bson:"prepTime" json:"prepTime"`
	Difficulty *int       `bson:"difficulty" json:"difficulty"`
	Vegetarian *bool      `bson:"vegetarian" json:"vegetarian"`
	Rating     *int       `bson:"rating" json:"rating"`
}

type recipeRatingReq struct {
	Rating *int `bson:"rating" json:"rating"`
}
