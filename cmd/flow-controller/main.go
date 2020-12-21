package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/yametech/echoer/pkg/controller"
	"github.com/yametech/echoer/pkg/storage/mongo"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://127.0.0.1:27017/admin", "-storage_uri mongodb://127.0.0.1:27017/admin")
	flag.Parse()

	fmt.Println(fmt.Sprintf("echoer flow-controller start... %v", time.Now()))

	stage, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := controller.NewFlowController(stage).Run(); err != nil {
			errC <- err
		}
	}()

	panic(<-errC)
}
