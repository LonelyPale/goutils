package mongodb

// Field Constant
const (
	// db bson field name
	IDKey         = "_id" // MongoDB ObjectID Field Name
	CreateTimeKey = "createTime"
	UpdateTimeKey = "updateTime"

	// go struct field name
	IDField         = "ID"
	CreateTimeField = "CreateTime"
	UpdateTimeField = "UpdateTime"
)
