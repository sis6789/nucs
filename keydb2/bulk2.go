package keydb2

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sis6789/nucs/caller"
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

// NewBulk - prepare bulk operation
func (x *KeyDB) NewBulk(dbName, collectionName string, interval int) *BulkBlock {
	var b BulkBlock
	b.dbName = dbName
	b.collectionName = collectionName
	b.collection = x.Col(dbName, collectionName)
	b.limit = interval
	return &b
}

// InsertOne - append action InsertOne.
func (b *BulkBlock) InsertOne(model *mongo.InsertOneModel) {
	if len(b.accumulatedAction) >= b.limit {
		b.Apply()
	}
	b.accumulatedAction = append(b.accumulatedAction, model)
}

// UpdateOne - append action UpdateOne.
func (b *BulkBlock) UpdateOne(model *mongo.UpdateOneModel) {
	if len(b.accumulatedAction) >= b.limit {
		b.Apply()
	}
	b.accumulatedAction = append(b.accumulatedAction, model)
}

// Apply - send accumulated request.
func (b *BulkBlock) Apply() {
	var nonOrderedOpt = options.BulkWrite().SetOrdered(false)
	if len(b.accumulatedAction) == 0 {
		return
	}
	var result *mongo.BulkWriteResult
	if result, b.err = b.collection.BulkWrite(context.Background(), b.accumulatedAction, nonOrderedOpt); b.err != nil {
		log.Fatalln(caller.Caller(), b.err)
	}
	b.accumulatedAction = nil
	b.modify += int(result.ModifiedCount)
	b.match += int(result.MatchedCount)
	b.insert += int(result.InsertedCount)
	b.upsert += int(result.UpsertedCount)
	b.delete += int(result.DeletedCount)
}

// Close - send remain accumulated request.
func (b *BulkBlock) Close() {
	b.Apply()
}

// String - status message
func (b *BulkBlock) String() string {
	return fmt.Sprintf("%s-%s(ins:%d mat:%d mod:%d ups:%d del:%d)", b.dbName, b.collectionName,
		b.insert, b.match, b.modify, b.upsert, b.delete)
}
