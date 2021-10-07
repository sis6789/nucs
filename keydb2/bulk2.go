package keydb2

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"

	"github.com/sis6789/nucs/caller"
)

type BulkBlock struct {
	limit              int
	dbName             string
	collectionName     string
	collection         *mongo.Collection
	modify             int
	match              int
	insert             int
	upsert             int
	delete             int
	err                error
	goRoutineRequest   chan mongo.WriteModel
	goRoutineRequestWG sync.WaitGroup
}

// merger - 야러 고루틴에서 보내지는 요구를 모아서 DB에 적용한다.
func (b *BulkBlock) goRoutineMerger() {
	defer b.goRoutineRequestWG.Done()
	var tempHolder []mongo.WriteModel
	var nonOrderedOpt = options.BulkWrite().SetOrdered(false)
	var result *mongo.BulkWriteResult
	for request := range b.goRoutineRequest {
		if len(tempHolder) >= b.limit {
			if result, b.err = b.collection.BulkWrite(context.Background(), tempHolder, nonOrderedOpt); b.err != nil {
				log.Fatalln(caller.Caller(), b.err)
			}
			b.modify += int(result.ModifiedCount)
			b.match += int(result.MatchedCount)
			b.insert += int(result.InsertedCount)
			b.upsert += int(result.UpsertedCount)
			b.delete += int(result.DeletedCount)
			tempHolder = nil
		}
		tempHolder = append(tempHolder, request)
	}
	if result, b.err = b.collection.BulkWrite(context.Background(), tempHolder, nonOrderedOpt); b.err != nil {
		log.Fatalln(caller.Caller(), b.err)
	}
	b.modify += int(result.ModifiedCount)
	b.match += int(result.MatchedCount)
	b.insert += int(result.InsertedCount)
	b.upsert += int(result.UpsertedCount)
	b.delete += int(result.DeletedCount)
	tempHolder = nil
}

// NewBulk - prepare bulk operation
func (x *KeyDB) NewBulk(dbName, collectionName string, interval int) *BulkBlock {
	var b BulkBlock
	b.dbName = dbName
	b.collectionName = collectionName
	b.collection = x.Col(dbName, collectionName)
	b.limit = interval
	b.goRoutineRequest = make(chan mongo.WriteModel)
	b.goRoutineRequestWG.Add(1)
	go b.goRoutineMerger()
	return &b
}

// InsertOne - append action InsertOne.
func (b *BulkBlock) InsertOne(model *mongo.InsertOneModel) {
	b.goRoutineRequest <- model
}

// UpdateOne - append action UpdateOne.
func (b *BulkBlock) UpdateOne(model *mongo.UpdateOneModel) {
	b.goRoutineRequest <- model
}

// Close - send remain accumulated request.
func (b *BulkBlock) Close() {
	close(b.goRoutineRequest)
	b.goRoutineRequestWG.Wait()
}

// String - status message
func (b *BulkBlock) String() string {
	return fmt.Sprintf("/d:%s/c:%s/(ins:%d mat:%d mod:%d ups:%d del:%d)", b.dbName, b.collectionName,
		b.insert, b.match, b.modify, b.upsert, b.delete)
}
