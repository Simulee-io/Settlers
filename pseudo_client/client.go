package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	settlers "simiulee.io/settlers/server/proto/service"
)

func main() {
	fmt.Println("pseudo client initialized")

	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close() //calls close function once main function completes

	c := settlers.NewSettlersClient(cc)

	//fmt.Printf("Created client %f", c)
	doUnary(c)

}

func doUnary(c settlers.SettlersClient) {
	req := &settlers.DiceRequest{}
	resp, err := c.RollDice(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling dice roll: %v", err)
	}

	log.Println("Response from dice roll:")
	log.Printf("dice 1: %v", resp.Dice1)
	log.Printf("dice 2: %v", resp.Dice2)
}
