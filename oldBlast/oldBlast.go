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
	bulk        *keydb2.BulkBlock
	err         error
}

func New(mongoAccess, dbName, colName string) *OldBlast {
	var x OldBlast
	x.mongoAccess = mongoAccess
	x.dbName = dbName
	x.colName = colName
	x.client = keydb2.New(mongoAccess)
	x.collection = x.client.Col(dbName, colName)
	x.bulk = x.client.NewBulk(dbName, colName, 10000)
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

// Save - 블록에 데이터를 설정하고 save를 불러 저장한다. 오루가 없으면 true를 반환한다.
func (x *OldBlast) Save() bool {
	_, x.err = x.collection.InsertOne(context.TODO(), bson.D{
		{"_id", x.B},
		{"bt", x.BT},
		{"c", x.C},
		{"l", x.L},
		{"qs", x.QS},
		{"sf", x.SF},
		{"ss", x.SS},
		{"st", x.ST},
	})
	if x.err != nil {
		return false
	} else {
		return true
	}
}

// SaveBulk - 10000개 기준으로 저장한다. 최종 Close를 불러 저장해야 한다.
func (x *OldBlast) SaveBulk() {
	model1 := mongo.NewInsertOneModel().SetDocument(bson.D{
		{"_id", x.B},
		{"bt", x.BT},
		{"c", x.C},
		{"l", x.L},
		{"qs", x.QS},
		{"sf", x.SF},
		{"ss", x.SS},
		{"st", x.ST},
	})
	x.bulk.InsertOne(model1)
}

// Close - 모든 처리를 종료한다.
func (x *OldBlast) Close() {
	x.bulk.Close()
}
