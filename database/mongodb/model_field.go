package mongodb

type model struct {
	ID         string
	CreateTime string
	UpdateTime string
}

var ModelFields = model{
	ID:         "id",
	CreateTime: "createTime",
	UpdateTime: "updateTime",
}
