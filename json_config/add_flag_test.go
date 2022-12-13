package json_config

import (
	"flag"
	"log"
	"testing"
)

func TestAddFlag(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Put("jfjdklsfjsdklfjd", "djfklsdklfjdksl")
	var myFlag = flag.Int("I-Love", 12345, "aaaa")
	AddFlag()
	v := Int("I-love")
	log.Printf("'%#v' '%#v'", *myFlag, v)
	Write("config.all")
	WriteUsed("config.used")
}
