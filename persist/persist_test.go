package persist

import (
	"log"
	"testing"
)

func TestSave(t *testing.T) {
	var x = map[int]int{
		1: 1, 2: 2, 3: 3,
	}
	var y = map[int]int{
		11: 1, 12: 2, 13: 3,
	}
	err := Save("c:/tmp/map.save", x)
	if err != nil {
		log.Printf("%v", err)
	}

	err = Load("c:/tmp/map.save", &y)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("%#v\n", y)
}
