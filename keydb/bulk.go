package keydb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type BulkBlock struct {
	limit             int
	collectionName    string
	collection        *mongo.Collection
	accumulatedAction []mongo.WriteModel
	modify            int
	match             int
	insert            int
	upsert            int
	delete            int
}

func NewBulk(collection string, interval int) *BulkBlock {
	var x BulkBlock
	x.collectionName = collection
	x.collection = Col(x.collectionName)
	x.limit = interval
	return &x
}

func (x *BulkBlock) InsertOne(model *mongo.InsertOneModel) {
	if len(x.accumulatedAction) > x.limit {
		x.Apply()
		x.accumulatedAction = []mongo.WriteModel{model}
	} else {
		x.accumulatedAction = append(x.accumulatedAction, model)
	}
}

func (x *BulkBlock) UpdateOne(model *mongo.UpdateOneModel) {
	if len(x.accumulatedAction) > x.limit {
		x.Apply()
		x.accumulatedAction = []mongo.WriteModel{model}
	} else {
		x.accumulatedAction = append(x.accumulatedAction, model)
	}
}

func (x *BulkBlock) Apply() {
	var nonOrderedOpt = options.BulkWrite().SetOrdered(false)
	if len(x.accumulatedAction) == 0 {
		return
	}
	var result *mongo.BulkWriteResult
	if result, err = x.collection.BulkWrite(context.Background(), x.accumulatedAction, nonOrderedOpt); err != nil {
		log.Fatalln(err)
	}
	x.accumulatedAction = nil
	x.modify += int(result.ModifiedCount)
	x.match += int(result.MatchedCount)
	x.insert += int(result.InsertedCount)
	x.upsert += int(result.UpsertedCount)
	x.delete += int(result.DeletedCount)
}

func (x *BulkBlock) String() string {
	return fmt.Sprintf("insert:%d match:%d modify:%d upsert:%d delete:%d",
		x.insert, x.match, x.modify, x.upsert, x.delete)
}
