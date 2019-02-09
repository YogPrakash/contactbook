package contactBookService

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (cb *contactBook) updateQ() bson.M {
	updateQ := bson.M{}

	if len(cb.UserID) != 0 {
		updateQ["userID"] = cb.UserID
	}
	if len(cb.GroupID) != 0 {
		updateQ["groupID"] = cb.GroupID
	}
	if len(cb.FirstName) != 0 {
		updateQ["firstName"] = cb.FirstName
	}
	if len(cb.LastName) != 0 {
		updateQ["lastName"] = cb.LastName
	}
	if cb.Email != nil {
		updateQ["email"] = cb.Email
	}
	if len(cb.Contact) != 0 {
		updateQ["contact"] = cb.Contact
	}
	if len(cb.Notes) != 0 {
		updateQ["notes"] = cb.Notes
	}
	if len(cb.LastUpdatedByUser) != 0 {
		updateQ["lastUpdatedByUser"] = cb.LastUpdatedByUser
	}

	return bson.M{"$set": updateQ}
}

func emailExists(email string, db *mgo.Session) bool {
	findQ := bson.M{"email": email}
	selectQ := bson.M{"_id": 1}

	var cb contactBook
	err := db.DB(dbName).C(collectionName).Find(findQ).Select(selectQ).One(&cb)
	if err != nil && err == mgo.ErrNotFound {
		return false
	}

	if err != nil {
		return false
	}

	return true
}
