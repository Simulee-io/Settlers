package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	settlers "simiulee.io/settlers/server/service"
)

func main() {
	fmt.Println("pseudo client initialized")

	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close() //calls close function once main function completes

	c := settlers.NewSettlersClient(cc)

	fmt.Printf("Created client %f", c)

}
