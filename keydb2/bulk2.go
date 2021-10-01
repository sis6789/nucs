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

func (x *KeyDB) NewBulk(dbName, collectionName string, interval int) *BulkBlock {
	var b BulkBlock
	b.dbName = dbName
	b.collectionName = collectionName
	b.collection = x.Col(dbName, collectionName)
	b.limit = interval
	return &b
}

func (b *BulkBlock) InsertOne(model *mongo.InsertOneModel) {
	if len(b.accumulatedAction) >= b.limit {
		b.Apply()
	}
	b.accumulatedAction = append(b.accumulatedAction, model)
}

func (b *BulkBlock) UpdateOne(model *mongo.UpdateOneModel) {
	if len(b.accumulatedAction) >= b.limit {
		b.Apply()
	}
	b.accumulatedAction = append(b.accumulatedAction, model)
}

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

func (b *BulkBlock) Close() {
	b.Apply()
}

func (b *BulkBlock) String() string {
	return fmt.Sprintf("%s-%s(ins:%d mat:%d mod:%d ups:%d del:%d)", b.dbName, b.collectionName,
		b.insert, b.match, b.modify, b.upsert, b.delete)
}
