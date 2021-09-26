package keydb2

import (
	"context"
	"github.com/sis6789/nucs/caller"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type KeyDB struct {
	myContext     context.Context //= context.Background()
	err           error
	mongodbAccess string
	mongoClient   *mongo.Client //= nil
	mapCollection map[string]*mongo.Collection
}

func New(access string) *KeyDB {
	var keyDB KeyDB
	keyDB.myContext = context.Background()
	keyDB.Connect(access)
	return &keyDB
}

func (x *KeyDB) Close() {
	if x.mongoClient == nil {
		return
	}
	if x.err = x.mongoClient.Disconnect(x.myContext); x.err != nil {
		log.Fatalln(x.err)
	}
	x.mongoClient = nil
	x.mapCollection = nil
}

func (x *KeyDB) Connect(access string) {
	if x.mongoClient == nil {
		x.mapCollection = make(map[string]*mongo.Collection)
		x.mongodbAccess = access
		clientOptions := options.Client().ApplyURI(x.mongodbAccess)
		if client, err := mongo.Connect(x.myContext, clientOptions); err != nil {
			log.Fatalln(caller.Caller(), err)
		} else {
			x.mongoClient = client
		}
		if x.err = x.mongoClient.Ping(x.myContext, nil); x.err != nil {
			log.Fatalln(caller.Caller(), x.err)
		}
	}
}

func (x *KeyDB) ReConnect(access string) {
	if x.mongoClient != nil {
		_ = x.mongoClient.Disconnect(x.myContext)
	}
	x.Connect(access)
}

func (x *KeyDB) Col(dbName, collectionName string) *mongo.Collection {
	dbCol := dbName + "::" + collectionName
	if collection, exist := x.mapCollection[dbCol]; exist {
		return collection
	} else {
		col := x.mongoClient.Database(dbName).Collection(collectionName)
		x.mapCollection[dbCol] = col
		return col
	}
}

func (x *KeyDB) Add(dbName, collectionName string) *mongo.Collection {
	dbCol := dbName + "::" + collectionName
	if _, exist := x.mapCollection[dbCol]; exist {
		x.Drop(dbName, collectionName)
	}
	col := x.mongoClient.Database(dbName).Collection(collectionName)
	x.mapCollection[dbCol] = col
	return col
}

func (x *KeyDB) Drop(dbName, collectionName string) {
	dbCol := dbName + "::" + collectionName
	if col, exist := x.mapCollection[dbCol]; exist {
		if x.err = col.Drop(x.myContext); x.err != nil {
			log.Fatalln(caller.Caller(), x.err)
		}
		delete(x.mapCollection, dbCol)
	}
}

func (x *KeyDB) DropDb(dbName string) {
	if x.err = x.mongoClient.Database(dbName).Drop(x.myContext); x.err != nil {
		log.Fatalln(caller.Caller(), x.err)
	}
}

func (x *KeyDB) ResetCol(dbName, CollectionName string) *mongo.Collection {
	x.Drop(dbName, CollectionName)
	return x.Add(dbName, CollectionName)
}

func (x *KeyDB) Index(dbName, collectionName string, fieldName ...string) {
	collection := x.Col(dbName, collectionName)
	var vFalse = false
	var keyDef bson.D
	for _, kf := range fieldName {
		keyDef = append(keyDef, bson.E{Key: kf, Value: 1})
	}
	model := mongo.IndexModel{
		Keys: keyDef,
		Options: &options.IndexOptions{
			Unique: &vFalse,
		},
	}
	if _, err := collection.Indexes().CreateOne(x.myContext, model); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
}
