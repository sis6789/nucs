package keydb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type BulkBlock struct {
	limit             int
	collectionName    string
	collection        *mongo.Collection
	accumulatedAction []mongo.WriteModel
}

func NewBulk(collection string, applyLimit int) *BulkBlock {
	var x BulkBlock
	x.collectionName = collection
	x.collection = Col(x.collectionName)
	x.limit = applyLimit
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
	if len(x.accumulatedAction) == 0 {
		return
	}
	issueLen := len(x.accumulatedAction)
	nonOrderedOpt := options.BulkWrite().SetOrdered(false)
	var result *mongo.BulkWriteResult
	if result, err = x.collection.BulkWrite(context.TODO(), x.accumulatedAction, nonOrderedOpt); err != nil {
		log.Fatalln(err)
	}
	x.accumulatedAction = nil
	log.Printf("%v bulk issue:%d match:%d modify:%d\n", x.collectionName, issueLen,
		result.MatchedCount, result.ModifiedCount)
}
