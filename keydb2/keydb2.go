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
var dbMapMutex sync.Mutex

type KeyDB struct {
	myContext     context.Context //= context.Background()
	err           error
	mongodbAccess string
	mongoClient   *mongo.Client //= nil
	mapCollection map[string]*mongo.Collection
	mapBulk       map[string]*BulkBlock
	colMapMutex   sync.Mutex
}

// New - prepare mongodb access
func New(access string) *KeyDB {
	var err error
	dbMapMutex.Lock()
	savedKeyDB, exist := dbMap[access]
	if !exist {
		var newKeyDB KeyDB
		newKeyDB.myContext = context.Background()
		//newKeyDB.connect(access)
		newKeyDB.mapCollection = make(map[string]*mongo.Collection)
		newKeyDB.mongodbAccess = access
		clientOptions := options.Client().ApplyURI(newKeyDB.mongodbAccess)
		if newKeyDB.mongoClient, newKeyDB.err = mongo.Connect(newKeyDB.myContext, clientOptions); err != nil {
			log.Fatalln(caller.Caller(), newKeyDB.err)
		}
		if newKeyDB.err = newKeyDB.mongoClient.Ping(newKeyDB.myContext, nil); newKeyDB.err != nil {
			log.Fatalln(caller.Caller(), newKeyDB.err)
		}
		newKeyDB.mapCollection = make(map[string]*mongo.Collection)
		newKeyDB.mapBulk = make(map[string]*BulkBlock)
		dbMap[access] = &newKeyDB
		savedKeyDB = &newKeyDB
	}
	dbMapMutex.Unlock()
	return savedKeyDB
}

// GoodBye - disconnect all connection
func GoodBye() {
	dbMapMutex.Lock()
	for k, kdb := range dbMap {
		if kdb.err = kdb.mongoClient.Disconnect(kdb.myContext); kdb.err != nil {
			log.Print(kdb.mongodbAccess, kdb.err)
		} else {
			log.Print("disconnect", kdb.mongodbAccess, kdb.err)
		}
		delete(dbMap, k)
	}
	dbMapMutex.Unlock()
}

// Col - return collection, if not exist make collection and return it.
func (x *KeyDB) Col(dbName, collectionName string) *mongo.Collection {
	dbCol := dbName + "::" + collectionName
	x.colMapMutex.Lock()
	var collection *mongo.Collection
	var exist bool
	if collection, exist = x.mapCollection[dbCol]; !exist {
		collection = x.mongoClient.Database(dbName).Collection(collectionName)
		x.mapCollection[dbCol] = collection
	}
	x.colMapMutex.Unlock()
	return collection
}

// Drop - delete collection
func (x *KeyDB) Drop(dbName string, collectionNames ...string) {
	x.colMapMutex.Lock()
	for _, colName := range collectionNames {
		dbCol := dbName + "::" + colName
		if col, exist := x.mapCollection[dbCol]; exist {
			if x.err = col.Drop(x.myContext); x.err != nil {
				log.Fatalln(caller.Caller(), x.err)
			}
			delete(x.mapCollection, dbCol)
			delete(x.mapBulk, dbCol)
		}
	}
	x.colMapMutex.Unlock()
}

// DropDb - Delete DB and associated collection.
func (x *KeyDB) DropDb(dbName string) {
	x.colMapMutex.Lock()
	if x.err = x.mongoClient.Database(dbName).Drop(x.myContext); x.err != nil {
		log.Fatalln(caller.Caller(), x.err)
	}
	for k := range x.mapCollection {
		dbCol := strings.Split(k, "::")
		if dbCol[0] == dbName {
			delete(x.mapCollection, k)
		}
	}
	for k := range x.mapBulk {
		dbCol := strings.Split(k, "::")
		if dbCol[0] == dbName {
			delete(x.mapBulk, k)
		}
	}
	x.colMapMutex.Unlock()
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
