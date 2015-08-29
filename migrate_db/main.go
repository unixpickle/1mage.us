package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"os"
)

type Image struct {
	MIME string `json:"mimeType" bson:"mime"`
	Seq  int    `json:"id" bson:"sequence"`
}

func die(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		die("Usage: migrate_db <images path>")
	}

	session, err := mgo.Dial("mongodb://127.0.0.1:27017/1mage")
	if err != nil {
		die(err)
	}
	collection := session.DB("1mage").C("images")
	var result []Image
	if err := collection.Find(nil).All(&result); err != nil {
		die(err)
	}

	// TODO: parse the dates from the path argument.
	fmt.Println(result)
}
