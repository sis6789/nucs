package keydb2

import (
	"context"
	"log"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sis6789/nucs/caller"
)

var dbMap = make(map[string]*KeyDB)
var mutex sync.Mutex

type KeyDB struct {
	myContext     context.Context //= context.Background()
	err           error
	mongodbAccess string
	mongoClient   *mongo.Client //= nil
	mapCollection map[string]*mongo.Collection
}

// New - prepare mongodb access
func New(access string) *KeyDB {

	mutex.Lock()
	savedKeyDB, exist := dbMap[access]
	mutex.Unlock()
	if exist {
		return savedKeyDB
	}

	var newKeyDB KeyDB
	newKeyDB.myContext = context.Background()
	newKeyDB.connect(access)
	newKeyDB.mapCollection = make(map[string]*mongo.Collection)

	mutex.Lock()
	dbMap[access] = &newKeyDB
	mutex.Unlock()

	return &newKeyDB
}

// GoodBye - disconnect all connection
func GoodBye() {
	for k, kdb := range dbMap {
		if kdb.err = kdb.mongoClient.Disconnect(kdb.myContext); kdb.err != nil {
			log.Println(kdb.mongodbAccess, kdb.err)
		} else {
			log.Println("disconnect", kdb.mongodbAccess, kdb.err)
		}
		delete(dbMap, k)
	}
}

// Open - prepare mongodb access
func Open(access string) *KeyDB {
	return New(access)
}

// connect - access mongodb server and check server availability.
func (x *KeyDB) connect(access string) {
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

// Col - return collection, if not exist make collection and return it.
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

// Add - make new collection, if already exist delete it and make new collection.
func (x *KeyDB) Add(dbName, collectionName string) *mongo.Collection {
	dbCol := dbName + "::" + collectionName
	if _, exist := x.mapCollection[dbCol]; exist {
		x.Drop(dbName, collectionName)
	}
	col := x.mongoClient.Database(dbName).Collection(collectionName)
	x.mapCollection[dbCol] = col
	return col
}

// Drop - delete collection
func (x *KeyDB) Drop(dbName string, collectionNames ...string) {
	for _, colName := range collectionNames {
		dbCol := dbName + "::" + colName
		if col, exist := x.mapCollection[dbCol]; exist {
			if x.err = col.Drop(x.myContext); x.err != nil {
				log.Fatalln(caller.Caller(), x.err)
			}
			delete(x.mapCollection, dbCol)
		}
	}
}

// DropDb - Delete DB and associated collection.
func (x *KeyDB) DropDb(dbName string) {
	if x.err = x.mongoClient.Database(dbName).Drop(x.myContext); x.err != nil {
		log.Fatalln(caller.Caller(), x.err)
	}
	for k := range x.mapCollection {
		dbCol := strings.Split(k, "::")
		if dbCol[0] == dbName {
			delete(x.mapCollection, k)
		}
	}
}

// ResetCol - Delete collection and remake it.
func (x *KeyDB) ResetCol(dbName, CollectionName string) *mongo.Collection {
	x.Drop(dbName, CollectionName)
	return x.Add(dbName, CollectionName)
}

// Index - add index definition. Specify key elements as repeated string.
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
