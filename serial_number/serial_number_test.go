package serial_number

import (
	"fmt"
	"github.com/sis6789/nucs/keydb2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"testing"
	"time"
)

func TestSerialNumber_Next(t *testing.T) {

	db := keydb2.New("mongodb://localhost:27017")
	db.DropDb("bulktest")
	bulk := db.NewBulk("bulktest", "bulktest", 21)
	n := New()
	var wg2 sync.WaitGroup
	for jj := 0; jj < 100; jj++ {
		wg2.Add(1)
		go func(goSubID int) {
			defer wg2.Done()
			for ix := 0; ix < 10; ix++ {
				idValue := n.Next()
				t1 := time.Now().Format("2006-01-02 15:04:05.999999")
				model1 := mongo.NewInsertOneModel().SetDocument(bson.D{
					{"_id", idValue},
					{"goSub", goSubID},
					{"val1", t1},
				})
				bulk.InsertOne(model1)
			}
		}(jj)
	}
	wg2.Wait()
	bulk.Close()
	fmt.Println("EOJ", bulk)
}
