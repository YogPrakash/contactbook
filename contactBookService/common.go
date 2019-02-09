package contactBookService

const (
	dbName           = "ContactBookService"
	collectionName   = "ContactBook"
	encodedAuthToken = "Basic dXNlcm5hbWU6cGFzc3dvcmQ="
	dbSessionKey     = "database"
	contentType      = "application/json"
)

//User contact number types
type ContactInfoType int

const (
	Personal ContactInfoType = iota
	Home
	Work
)
