package keydb2

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sis6789/nucs/caller"
)

type BulkBlock struct {
	limit              int
	dbName             string
	collectionName     string
	collection         *mongo.Collection
	goRoutineRequest   chan mongo.WriteModel
	goRoutineRequestWG sync.WaitGroup
	tempCount          int
}

// merger - 야러 고루틴에서 보내지는 요구를 모아서 DB에 적용한다.
func goRoutineMerger(b *BulkBlock) {
	defer b.goRoutineRequestWG.Done()
	var tempHolder []mongo.WriteModel
	var nonOrderedOpt = options.BulkWrite().SetOrdered(false)
	var wgAsync sync.WaitGroup
	for request := range b.goRoutineRequest {
		tempHolder = append(tempHolder, request)
		b.tempCount++
		if len(tempHolder) >= b.limit {
			wgAsync.Add(1)
			go func(models []mongo.WriteModel) {
				defer wgAsync.Done()
				if _, err := b.collection.BulkWrite(context.Background(), models, nonOrderedOpt); err != nil {
					log.Fatalln(caller.Caller(), b, err)
				}
			}(tempHolder)
			tempHolder = []mongo.WriteModel{}
		}
	}
	// sender close channel
	if len(tempHolder) > 0 {
		wgAsync.Add(1)
		go func(models []mongo.WriteModel) {
			defer wgAsync.Done()
			if _, err := b.collection.BulkWrite(context.Background(), models, nonOrderedOpt); err != nil {
				log.Fatalln(caller.Caller(), b, err)
			}
		}(tempHolder)
		tempHolder = []mongo.WriteModel{}
	}
	wgAsync.Wait()
	log.Printf("close bulk %v", b)
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
	return fmt.Sprintf("%s.%s", bb.dbName, bb.collectionName)
}
