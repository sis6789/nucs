package keydb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var myContext = context.Background()
var err error
var dbName string
var mongodbAccess string

var mongoClient *mongo.Client = nil
var mapCollection map[string]*mongo.Collection
var basicColNames []string

func Connect(access string, db string, collectionName ...string) {
	mongodbAccess = access
	dbName = db
	basicColNames = collectionName
	if mongoClient == nil {
		clientOptions := options.Client().ApplyURI(mongodbAccess)
		if client, err := mongo.Connect(myContext, clientOptions); err != nil {
			log.Fatalln(err)
		} else {
			mongoClient = client
		}
		if err = mongoClient.Ping(myContext, nil); err != nil {
			log.Fatalln(err)
		}
	}
	setBaseCollection()
}

func setBaseCollection() {
	// base collection list
	mapCollection = make(map[string]*mongo.Collection)
	for _, name := range basicColNames {
		col := mongoClient.Database(dbName).Collection(name)
		mapCollection[name] = col
	}
}

func Col(name string) *mongo.Collection {
	if collection, exist := mapCollection[name]; exist {
		return collection
	} else {
		log.Fatalln("undefined collection name", name)
		return nil
	}
}

func Add(name string) *mongo.Collection {
	if _, exist := mapCollection[name]; exist {
		Drop(name)
	}
	col := mongoClient.Database(dbName).Collection(name)
	mapCollection[name] = col
	return col
}

func Drop(name string) {
	if col, exist := mapCollection[name]; exist {
		if err = col.Drop(myContext); err != nil {
			log.Fatalln(err)
		}
		delete(mapCollection, name)
	}
}

func DropDb() {
	if err = mongoClient.Database(dbName).Drop(myContext); err != nil {
		log.Fatalln(err)
	}
	setBaseCollection()
}

func ResetCol(name string) *mongo.Collection {
	Drop(name)
	return Add(name)
}

func Index(collectionName string, fieldName ...string) {
	collection := Col(collectionName)
	var vFalse = false
	var keyDef bson.M
	for _, kf := range fieldName {
		keyDef[kf] = 1
	}
	model := mongo.IndexModel{
		Keys: keyDef,
		Options: &options.IndexOptions{
			Unique: &vFalse,
		},
	}
	if _, err := collection.Indexes().CreateOne(myContext, model); err != nil {
		log.Fatalln(err)
	}
}
