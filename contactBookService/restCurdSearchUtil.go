package contactBookService

import "gopkg.in/mgo.v2/bson"

func (r *recipe) updateQ() bson.M {
	updateQ := bson.M{}

	if r.PrepTime != nil {
		updateQ["prepTime"] = *r.PrepTime
	}
	if len(r.Name) != 0 {
		updateQ["name"] = r.Name
	}
	if r.Vegetarian != nil {
		updateQ["vegetarian"] = *r.Vegetarian
	}
	if r.Difficulty != nil {
		updateQ["difficulty"] = *r.Difficulty
	}
	if r.Rating != nil {
		updateQ["rating"] = *r.Rating
	}
	return bson.M{"$set": updateQ}
}
