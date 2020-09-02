package db

type DatabaseI interface {
	StoreItem(string, interface{})
	FetchItem(string) interface{}
	DeleteItem(string)
}

var DB DatabaseI