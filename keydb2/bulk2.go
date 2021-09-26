package keydb2

import (
	"context"
	"fmt"
	"github.com/sis6789/nucs/caller"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type BulkBlock struct {
	limit             int
	dbName            string
	collectionName    string
	collection        *mongo.Collection
	accumulatedAction []mongo.WriteModel
	modify            int
	match             int
	insert            int
	upsert            int
	delete            int
	err               error
}

func (x *KeyDB) NewBulk(dbName, collectionName string, interval int) *BulkBlock {
	var y BulkBlock
	y.dbName = dbName
	y.collectionName = collectionName
	y.collection = x.Col(dbName, collectionName)
	y.limit = interval
	return &y
}

func (x *BulkBlock) InsertOne(model *mongo.InsertOneModel) {
	if len(x.accumulatedAction) >= x.limit {
		x.Apply()
	}
	x.accumulatedAction = append(x.accumulatedAction, model)
}

func (x *BulkBlock) UpdateOne(model *mongo.UpdateOneModel) {
	if len(x.accumulatedAction) >= x.limit {
		x.Apply()
	}
	x.accumulatedAction = append(x.accumulatedAction, model)
}

func (x *BulkBlock) Apply() {
	var nonOrderedOpt = options.BulkWrite().SetOrdered(false)
	if len(x.accumulatedAction) == 0 {
		return
	}
	var result *mongo.BulkWriteResult
	if result, x.err = x.collection.BulkWrite(context.Background(), x.accumulatedAction, nonOrderedOpt); x.err != nil {
		log.Fatalln(caller.Caller(), x.err)
	}
	x.accumulatedAction = nil
	x.modify += int(result.ModifiedCount)
	x.match += int(result.MatchedCount)
	x.insert += int(result.InsertedCount)
	x.upsert += int(result.UpsertedCount)
	x.delete += int(result.DeletedCount)
}

func (x *BulkBlock) Close() {
	x.Apply()
}

func (x *BulkBlock) String() string {
	return fmt.Sprintf("%s(ins:%d mat:%d mod:%d ups:%d del:%d)", x.collectionName,
		x.insert, x.match, x.modify, x.upsert, x.delete)
}
