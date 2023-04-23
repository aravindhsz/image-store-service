package main

import (
	"fmt"
	"image-store-service/pkg/mongoserver"
	"image-store-service/pkg/server"
)

func main() {
	fmt.Println("starting image-store-service service")
	// var wg sync.WaitGroup
	// wg.Add(2)

	client, err := mongoserver.InitMongoClient()
	if err != nil {
		fmt.Println(err)
	}

	err = server.Start(client)

	if err != nil {
		fmt.Println(err)
	}

	// wg.Wait()
}
