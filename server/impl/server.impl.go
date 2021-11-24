package impl

import (
	"context"
	"fmt"
	"math/rand"

	settlers "simiulee.io/settlers/server/proto/service"
)

type Server struct{}

func (*Server) RollDice(ctx context.Context, in *settlers.DiceRequest) (*settlers.DiceResponse, error) {

	fmt.Printf("Rolling dice..")

	d1 := (rand.Intn(6) + 1)
	d2 := (rand.Intn(6) + 1)

	fmt.Printf("dice 1: %v", d1)
	fmt.Printf("dice 2: %v", d2)

	response := &settlers.DiceResponse{
		Dice1: int32(d1),
		Dice2: int32(d2),
	}

	return response, nil
}
