package mongodb

// Field Constant
const (
	// go struct field name
	IDField         = "ID"
	CreateTimeField = "CreateTime"
	UpdateTimeField = "UpdateTime"

	// db bson field name
	IDBson         = "_id" // MongoDB ObjectID Field Name
	CreateTimeBson = "createTime"
	UpdateTimeBson = "updateTime"

	// web json field name
	IDJson         = "id"
	CreateTimeJson = "createTime"
	UpdateTimeJson = "updateTime"
)
