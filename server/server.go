package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"simiulee.io/settlers/server/impl"
	settlers "simiulee.io/settlers/server/proto/service"
)

func main() {
	fmt.Println("Initializing..")

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v, err")
	}

	s := grpc.NewServer()
	settlers.RegisterSettlersServer(s, &impl.Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v, err")
	}

}
