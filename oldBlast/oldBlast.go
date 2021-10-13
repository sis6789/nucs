package oldBlast

import (
	"context"
	"github.com/sis6789/nucs/keydb2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OldBlast struct {
	B           string `bson:"_id"`
	BT          string `bson:"bt"`
	C           int    `bson:"c"`
	L           int    `bson:"l"`
	QS          int    `bson:"qs"`
	SF          int    `bson:"sf"`
	SS          int    `bson:"ss"`
	ST          int    `bson:"st"`
	mongoAccess string
	dbName      string
	colName     string
	client      *keydb2.KeyDB
	collection  *mongo.Collection
	result      *mongo.SingleResult
	err         error
}

func New(mongoAccess, dbName, colName string) *OldBlast {
	var x OldBlast
	x.mongoAccess = mongoAccess
	x.dbName = dbName
	x.colName = colName
	x.client = keydb2.New(mongoAccess)
	x.collection = x.client.Col(dbName, colName)
	return &x
}

// Query - 서열이 존재하면 해당 값을 블록에 저장하고 true를 반환한다.
func (x *OldBlast) Query(query string) bool {
	x.result = x.collection.FindOne(context.TODO(), bson.M{"_id": query})
	if x.result.Err() != nil {
		// 만족하는 서열 없음
		return false
	}
	if x.err = x.result.Decode(x); x.err != nil {
		return false
	}
	return true
}

// Close - 모든 처리를 종료한다.
func Close() {
	//bulkOldBlast.Close()
}
