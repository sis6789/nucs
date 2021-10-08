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
func goRoutineMerger(b *BulkBlock) {
	defer b.goRoutineRequestWG.Done()
	var tempHolder []mongo.WriteModel
	var tempCount int
	var nonOrderedOpt = options.BulkWrite().SetOrdered(false)
	var result *mongo.BulkWriteResult
	for request := range b.goRoutineRequest {
		tempHolder = append(tempHolder, request)
		tempCount++
		if len(tempHolder) >= b.limit {
			log.Println(">>>>>>>>>>>>", b.collectionName, tempCount)
			if result, b.err = b.collection.BulkWrite(context.Background(), tempHolder, nonOrderedOpt); b.err != nil {
				log.Fatalln(caller.Caller(), b, b.err)
			}
			tempHolder = nil
			b.modify += int(result.ModifiedCount)
			b.match += int(result.MatchedCount)
			b.insert += int(result.InsertedCount)
			b.upsert += int(result.UpsertedCount)
			b.delete += int(result.DeletedCount)
		}
	}
	// sender close channel
	if len(tempHolder) > 0 {
		if result, b.err = b.collection.BulkWrite(context.Background(), tempHolder, nonOrderedOpt); b.err != nil {
			log.Fatalln(caller.Caller(), b, b.err)
		}
		tempHolder = nil
		b.modify += int(result.ModifiedCount)
		b.match += int(result.MatchedCount)
		b.insert += int(result.InsertedCount)
		b.upsert += int(result.UpsertedCount)
		b.delete += int(result.DeletedCount)
	}
}

// NewBulk - prepare bulk operation
func (x *KeyDB) NewBulk(dbName, collectionName string, interval int) *BulkBlock {
	var pB *BulkBlock
	var exist bool
	dbCol := dbName + "::" + collectionName
	if pB, exist = x.mapBulk[dbCol]; !exist {
		var b BulkBlock
		pB = &b
		b.dbName = dbName
		b.collectionName = collectionName
		b.collection = x.Col(dbName, collectionName)
		b.limit = interval
		b.goRoutineRequest = make(chan mongo.WriteModel)
		b.goRoutineRequestWG.Add(1)
		x.mapBulk[dbCol] = &b
		go goRoutineMerger(pB)
	}
	return pB
}

// InsertOne - append action InsertOne.
func (bb *BulkBlock) InsertOne(model *mongo.InsertOneModel) {
	bb.goRoutineRequest <- model
}

// UpdateOne - append action UpdateOne.
func (bb *BulkBlock) UpdateOne(model *mongo.UpdateOneModel) {
	bb.goRoutineRequest <- model
}

// Close - send remain accumulated request.
func (bb *BulkBlock) Close() {
	close(bb.goRoutineRequest)
	bb.goRoutineRequestWG.Wait()
}

// String - status message
func (bb *BulkBlock) String() string {
	return fmt.Sprintf("/d:%s/c:%s/(ins:%d mat:%d mod:%d ups:%d del:%d)", bb.dbName, bb.collectionName,
		bb.insert, bb.match, bb.modify, bb.upsert, bb.delete)
}
