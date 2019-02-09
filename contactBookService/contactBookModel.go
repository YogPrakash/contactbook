package contactBookService

import (
	"time"
)

type contactBook struct {
	ID             string      `bson:"_id" json:"id"`
	UserID         string      `bson:"userID" json:"userID"`
	GroupID        string      `bson:"groupID" json:"groupID"`
	FirstName      string      `bson:"firstName" json:"firstName"`
	LastName       string      `bson:"lastName" json:"lastName"`
	Email          *string     `bson:"email" json:"email"`
	Contact        []phoneInfo `json:"contact"`
	Notes          string      `bson:"notes" json:"notes"`
	CreatedByUser  string      `bson:"createdByUser" json:"createdByUser"`
	CreateDateTime *time.Time  `bson:"createDateTime" json:"createDateTime"`
	// IsActive T or F (TRUE or FALSE) -- DEFAULT 'T'
	IsActive *bool `bson:"isActive" json:"isActive"`
	// LastUpdatedByUser Data updated by who
	LastUpdatedByUser   string     `bson:"lastUpdatedByUser" json:"lastUpdatedByUser"`
	LastUpdatedDateTime *time.Time `bson:"lastUpdatedDateTime" json:"lastUpdatedDateTime"`
	// DocumentVersion to keep track of the changes - DEFAULT 1.0
	DocumentVersion *float32 `bson:"documentVersion" json:"documentVersion"`
}

type phoneInfo struct {
	Type        *ContactInfoType `bson:"type" json:"type"`
	Number      string           `bson:"number" json:"number"`
	CountryCode string           `bson:"countryCode" json:"countryCode"`
	IsPrimary   *bool            `bson:"isPrimary" json:"isPrimary"`
}
