package serial_number

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sis6789/nucs/keydb2"
)

func TestSerialNumber_Next(t *testing.T) {

	mongoClient := keydb2.New("mongodb://localhost:27017")
	mongoClient.DropDb("bulktest")
	bulk01 := mongoClient.NewBulk("bulktest", "bulktest01", 100)
	bulk02 := mongoClient.NewBulk("bulktest", "bulktest02", 100)
	n01 := New()
	n02 := New()
	var func01, func02 sync.WaitGroup
	func01.Add(1)
	go func() {
		defer func01.Done()
		var wg01 sync.WaitGroup
		for jj := 0; jj < 10; jj++ {
			wg01.Add(1)
			go func(goSubID int) {
				defer wg01.Done()
				for ix := 0; ix < 1000; ix++ {
					idValue := n01.Next()
					t1 := time.Now().Format("2006-01-02 15:04:05.999999")
					model1 := mongo.NewInsertOneModel().SetDocument(bson.D{
						{"_id", idValue},
						{"goSub", goSubID},
						{"val1", t1},
					})
					bulk01.InsertOne(model1)
				}
			}(jj)
		}
		wg01.Wait()
	}()
	func02.Add(1)
	go func() {
		defer func02.Done()
		var wg02 sync.WaitGroup
		for jj := 0; jj < 10; jj++ {
			wg02.Add(1)
			go func(goSubID int) {
				defer wg02.Done()
				for ix := 0; ix < 1000; ix++ {
					idValue := n02.Next()
					t1 := time.Now().Format("2006-01-02 15:04:05.999999")
					model1 := mongo.NewInsertOneModel().SetDocument(bson.D{
						{"_id", idValue},
						{"goSub", goSubID},
						{"val1", t1},
					})
					bulk02.InsertOne(model1)
				}
			}(jj)
		}
		wg02.Wait()
	}()
	func01.Wait()
	func02.Wait()
	bulk01.Close()
	bulk02.Close()
	fmt.Println("EOJ", bulk01)
	fmt.Println("EOJ", bulk02)
}
