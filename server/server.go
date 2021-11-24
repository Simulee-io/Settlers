package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"
	settlers "simiulee.io/settlers/server/proto/service"
)

type server struct{}

func (*server) RollDice(ctx context.Context, in *settlers.DiceRequest) (*settlers.DiceResponse, error) {
	d1 := (rand.Intn(6) + 1)
	d2 := (rand.Intn(6) + 1)

	response := &settlers.DiceResponse{
		Dice1: int32(d1),
		Dice2: int32(d2),
	}

	return response, nil
}

func main() {
	fmt.Println("Initializing..")

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v, err")
	}

	s := grpc.NewServer()
	settlers.RegisterSettlersServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v, err")
	}

}
