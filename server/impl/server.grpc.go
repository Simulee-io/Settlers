package impl

import (
	"context"
	"fmt"
	"math/rand"

	settlers "simiulee.io/settlers/server/proto/service"
)

type Server struct{}

// Function: RollDice
//  In: *settlers.DiceRequest {empty}
//  Out: *setttlers.DiceResponse {Dice1: int, Dice2: int}, error
//  Description:
//   Example & reference function for gRPC. Takes in an empty DiceRequest, generates two random numbers
//   between 1 and 6 each representing a dice roll. The function then returns these values in a DiceResponse.
func (*Server) RollDice(ctx context.Context, in *settlers.DiceRequest) (*settlers.DiceResponse, error) {

	fmt.Println("Rolling dice..")

	d1 := (rand.Intn(6) + 1)
	d2 := (rand.Intn(6) + 1)

	fmt.Println("dice 1: %v", d1)
	fmt.Println("dice 2: %v", d2)
	fmt.Println()

	response := &settlers.DiceResponse{
		Dice1: int32(d1),
		Dice2: int32(d2),
	}

	//temp call
	//call function
	MakeMapFromCSV("mini.csv")

	return response, nil
}
