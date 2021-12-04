package impl

import (
	"errors"
	"log"
	"os"
	"strconv"

	board "simiulee.io/settlers/server/proto/board"
)

//           //
// USE CASES //
//           //

func GenerateMapFromCSV(fileName string) (*board.Point, error) {

	//open & read map csv file
	file, err := os.Open("server\\res\\rawmaps\\" + fileName)
	if err != nil {
		log.Fatalf("Error opening map file: %v", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalf("File error: %v", err)
		}
	}()

	// content, err := ioutil.ReadAll(file)

	// //process
	// rows := strings.Split(string(content), string(10))
	// var IDIter int32 = 0
	// var head *board.Point
	// horizHexCount := 0
	// prevHexs := make([]*board.Hexagon, 25)
	// currHexs := make([]*board.Hexagon, 25)
	// //prevX := 0
	// var currHex *board.Hexagon
	// var prevHex *board.Hexagon
	// currX := 0
	// prevHexs = nil
	// for i, r := range rows {
	// 	elements := strings.Split(string(r), ",")
	// 	prevHexs = currHexs
	// 	prevHexs = append(prevHexs, nil)
	// 	currHexs = nil
	// 	//println("i", i)
	// 	//prevX = currX
	// 	currX = 0
	// 	horizHexCount = 0

	// 	for y, c := range elements {

	// 		char := string(c)
	// 		//println("y", y)
	// 		//println(char)
	// 		_, numeric := strconv.Atoi(char)
	// 		if char == "x" {
	// 			if y < 3 {
	// 				currX++

	// 			}
	// 			currHexs = append(currHexs, nil)
	// 		} else if char == "h" {
	// 			currHex = makeHex(IDIter, "")
	// 			IDIter++

	// 			if horizHexCount > 0 {
	// 				connectHex(prevHex, currHex, "right", "")
	// 			}
	// 			if i == 0 && horizHexCount == 0 {
	// 				head = currHex.P[0]
	// 			}
	// 			horizHexCount++

	// 			currHexs = append(currHexs, currHex)
	// 			currHexs = append(currHexs, currHex)

	// 			println(i, y)
	// 			if prevHexs != nil {
	// 				if prevHexs[y] != nil {
	// 					if prevHexs[y+1] == prevHexs[y] {
	// 						connectHex(prevHexs[y], currHex, "down", "left")
	// 					} else {
	// 						connectHex(prevHexs[y], currHex, "down", "right")
	// 					}
	// 				}
	// 			}

	// 		} else if numeric == nil { //is a numeric
	// 			populatePorts(currHex, char)
	// 			prevHex = currHex

	// 			println(i, y)
	// 			if prevHexs != nil {
	// 				if prevHexs[y] != nil {
	// 					if prevHexs[y+1] == prevHexs[y] {
	// 						connectHex(prevHexs[y], currHex, "down", "left")
	// 					} else {
	// 						connectHex(prevHexs[y], currHex, "down", "right")
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// copy := head.Down.End.Down.End
	// counter := 0
	// for copy.Right != nil {
	// 	copy = copy.Right.End
	// 	counter += 1
	// }
	// println("found points ", counter)

	//temp return
	return nil, nil
}

//                  //
// HELPER FUNCTIONS //
//                  //

// HEX & POINT MANIPULATION:

//create hex and link all points
func makeHex(id int32, ports string) *board.Hexagon {

	//create and link all edges & points
	this := new(board.Hexagon)
	this.Id = id
	this.P = make([]*board.Point, 6)
	this.E = make([]*board.Edge, 6)
	this.Block = new(board.Block)
	this.Block.Resource = board.Resource_DESERT

	this.P[0] = new(board.Point)
	this.P[1] = new(board.Point)
	this.E[0] = new(board.Edge)
	connectPoints(this.P[0], this.E[0], this.P[1], "right")

	this.P[2] = new(board.Point)
	this.E[1] = new(board.Edge)
	connectPoints(this.P[1], this.E[1], this.P[2], "right")

	this.P[3] = new(board.Point)
	this.E[2] = new(board.Edge)
	connectPoints(this.P[2], this.E[2], this.P[3], "down")

	this.P[4] = new(board.Point)
	this.E[3] = new(board.Edge)
	connectPoints(this.P[4], this.E[3], this.P[3], "right")

	this.P[5] = new(board.Point)
	this.E[4] = new(board.Edge)
	connectPoints(this.P[5], this.E[4], this.P[4], "right")

	this.E[5] = new(board.Edge)
	connectPoints(this.P[0], this.E[5], this.P[5], "down")

	//populate ports
	populatePorts(this, ports)

	return this

}

//populate hex's ports
func populatePorts(h *board.Hexagon, ports string) {
	if ports == "0" {
		return
	}

	//populate ports
	pSlice := []rune(ports)
	for x := range pSlice {
		value, _ := strconv.Atoi(string(pSlice[x]))
		h.P[value-1].Port = new(board.TradePort)
	}
}

//connect 2 points through an edge
func connectPoints(p1 *board.Point, edge *board.Edge, p2 *board.Point, dir string) error {
	edge.Start = p1
	edge.End = p2

	if dir == "right" {
		p1.Right = edge
		p2.Left = edge
	} else if dir == "left" {
		p1.Left = edge
		p2.Right = edge
	} else if dir == "down" {
		p1.Down = edge
		p2.Up = edge
	} else if dir == "up" {
		p1.Up = edge
		p2.Down = edge
	} else {
		return (errors.New("connectPoints: invalid direction received"))
	}

	p1.Progression = 1
	p2.Progression = 1

	return nil
}

//connect 2 hexs (replacing 2nd hex's colliding points & edges with the first hex's)
func connectHex(h1 *board.Hexagon, h2 *board.Hexagon, dir string, offset string) error {
	if dir == "right" {
		if h2.P[0].Port != nil {
			h1.P[2].Port = h2.P[0].Port
		}
		if h2.P[5].Port != nil {
			h1.P[3].Port = h2.P[5].Port
		}

		h1.P[2].Right = h2.P[0].Right
		h1.P[3].Right = h2.P[5].Right

		h2.P[0] = h1.P[2]
		h2.P[5] = h1.P[3]
		h2.E[5] = h1.E[2]
	} else if dir == "left" {
		if h2.P[2].Port != nil {
			h1.P[0].Port = h2.P[2].Port
		}
		if h2.P[3].Port != nil {
			h1.P[5].Port = h2.P[3].Port
		}

		h1.P[0].Left = h2.P[2].Left
		h1.P[5].Left = h2.P[3].Left

		h2.P[2] = h1.P[0]
		h2.P[3] = h1.P[5]
		h2.E[2] = h1.E[5]
	} else if dir == "down" {
		if offset == "left" {
			if h2.P[1].Port != nil {
				h1.P[5].Port = h2.P[1].Port
			}
			if h2.P[2].Port != nil {
				h1.P[4].Port = h2.P[2].Port
			}

			h1.P[5].Down = h2.P[1].Down
			h1.P[4].Down = h2.P[2].Down

			h2.P[1] = h1.P[5]
			h2.P[2] = h1.P[4]
			h2.E[1] = h1.E[4]
		} else if offset == "right" {
			if h2.P[0].Port != nil {
				h1.P[4].Port = h2.P[0].Port
			}
			if h2.P[1].Port != nil {
				h1.P[3].Port = h2.P[1].Port
			}

			h1.P[4].Down = h2.P[0].Down
			h1.P[3].Down = h2.P[1].Down

			h2.P[0] = h1.P[4]
			h2.P[1] = h1.P[3]
			h2.E[0] = h1.E[3]
		} else {
			return (errors.New("connectHex: invalid offset received"))
		}
	} else if dir == "up" {
		if offset == "left" {
			if h2.P[5].Port != nil {
				h1.P[1].Port = h2.P[5].Port
			}
			if h2.P[4].Port != nil {
				h1.P[2].Port = h2.P[4].Port
			}

			h1.P[1].Up = h2.P[5].Up
			h1.P[2].Up = h2.P[4].Up

			h2.P[5] = h1.P[1]
			h2.P[4] = h1.P[2]
			h2.E[4] = h1.E[1]
		} else if offset == "right" {
			if h2.P[4].Port != nil {
				h1.P[0].Port = h2.P[4].Port
			}
			if h2.P[3].Port != nil {
				h1.P[1].Port = h2.P[3].Port
			}

			h1.P[0].Up = h2.P[4].Up
			h1.P[1].Up = h2.P[3].Up

			h2.P[4] = h1.P[0]
			h2.P[3] = h1.P[1]
			h2.E[3] = h1.E[0]
		} else {
			return (errors.New("connectHex: invalid offset received"))
		}
	} else {
		return (errors.New("connectHex: invalid direction received"))
	}

	return nil
}
