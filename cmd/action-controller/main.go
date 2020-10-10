package main

import (
	"flag"
	"fmt"
	"github.com/yametech/echoer/pkg/action"
	"github.com/yametech/echoer/pkg/controller"
	"github.com/yametech/echoer/pkg/storage/mongo"
	"time"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://127.0.0.1:27017/admin", "-storage_uri mongodb://127.0.0.1:27017/admin")
	flag.Parse()

	fmt.Println(fmt.Sprintf("echoer action-controller start... %v", time.Now()))
	stage, err := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}
	hc := action.NewHookClient()
	server := controller.NewActionController(stage, hc)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
