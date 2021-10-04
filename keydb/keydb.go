package keydb

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sis6789/nucs/keydb2"
)

var err error
var dbName string

var dKDB2 *keydb2.KeyDB

func Connect(access string, db string) {
	dKDB2 = keydb2.New(access)
	dbName = db
}

func ReConnect(access string, db string) {
	Connect(access, db)
}

func Col(name string) *mongo.Collection {
	return dKDB2.Col(dbName, name)
}

func Add(name string) *mongo.Collection {
	return dKDB2.Add(dbName, name)
}

func Drop(name string) {
	dKDB2.Drop(dbName, name)
}

func DropDb() {
	dKDB2.DropDb(dbName)
}

func ResetCol(name string) *mongo.Collection {
	return dKDB2.ResetCol(dbName, name)
}

func Index(collectionName string, fieldName ...string) {
	dKDB2.Index(dbName, collectionName, fieldName...)
}

func GoodBye() {
	keydb2.GoodBye()
}
