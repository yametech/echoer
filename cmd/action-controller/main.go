package main

import (
	"flag"
	"fmt"
	"github.com/yametech/echoer/pkg/controller"
	"github.com/yametech/echoer/pkg/storage/mongo"
	"time"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://127.0.0.1:27017/admin", "-storage_uri mongodb://127.0.0.1:27017/admin")
	flag.Parse()

	fmt.Println(fmt.Sprintf("echoer action-controller start... %v", time.Now()))

	stage, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := controller.NewActionController(stage).Run(); err != nil {
			errC <- err
		}
	}()

	panic(<-errC)

}
