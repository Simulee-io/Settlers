package impl

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	board "simiulee.io/settlers/server/proto/board"
)

//           //
// USE CASES //
//           //

func MakeMapFromCSV(fileName string) (*board.Board, error) {

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

	hexagons := [][]*board.Hexagon{}

	content, err := ioutil.ReadAll(file)
	rows := strings.Split(string(content), string(10))
	fmtd := [][]string{}

	//cycle through rows of text, generate readable array
	for x, r := range rows {
		row := strings.TrimRight(string(r), "\r\n")
		row = strings.TrimRight(string(row), string(10))
		items := strings.Split(row, ",")
		prev := "#"
		iter := 0

		//cycle through ',' delimetered items in row
		for _, rawVal := range items {
			fmtd = append(fmtd, []string{})
			val := string(rawVal)
			if iter > 1 {
				iter = 0
			}

			if val == "" {
				if prev != "" {
					iter = 0
				}
				if prev == "" && iter == 1 {
					fmtd[x] = append(fmtd[x], "_")
				}
				iter += 1
			} else if val == "w" {
				if prev != "w" {
					iter = 0
				}
				if prev == "w" && iter == 1 {
					fmtd[x] = append(fmtd[x], "w")
				}
				iter += 1
			} else { // port positions
				if prev == "h" {
					fmtd[x] = append(fmtd[x], val)
				}
			}

			prev = val

		}
	}

	//create hexagon array from fmtd
	for x, xv := range fmtd {

		hexagons = append(hexagons, []*board.Hexagon{})

		for y, yv := range xv {

			val := string(yv)

			if val == "_" {
				//empty becomes water
				hexagons[x] = append(hexagons[x], makeHex(strconv.Itoa(x)+strconv.Itoa(y), "", true))
			} else if val == "w" {
				hexagons[x] = append(hexagons[x], makeHex(strconv.Itoa(x)+strconv.Itoa(y), "", true))
			} else { //numeric indicating land hex with port positions
				//hexagons[x][y] = makeHex(strconv.Itoa(x)+strconv.Itoa(y), val, false)
				hexagons[x] = append(hexagons[x], makeHex(strconv.Itoa(x)+strconv.Itoa(y), val, false))
			}
		}
	}

	//link hexagons horizontally
	for x, xv := range hexagons {
		for y, _ := range xv {
			if y > 0 {
				connectHexRight(hexagons[x][y-1], hexagons[x][y])
			}
		}
	}

	//link hexagons vertically
	for x, xv := range hexagons {
		if x > 0 {
			for y, _ := range xv {
				//print(strconv.Itoa(x) + " " + strconv.Itoa(y) + "\n")

				linePos := (x + 1) % 2
				if linePos == 0 {
					connectHexUp(hexagons[x][y], hexagons[x-1][y], "right")
					if y < len(hexagons[x-1]) {
						connectHexUp(hexagons[x][y], hexagons[x-1][y+1], "left")
					}
				} else {
					if y < len(hexagons[x-1]) {
						connectHexUp(hexagons[x][y], hexagons[x-1][y], "left")
					}
					if y != 0 {
						connectHexUp(hexagons[x][y], hexagons[x-1][y-1], "right")
					}
				}
			}
		}
	}

	//populate vtables
	vPoints := []*board.Point{}
	vEdges := []*board.Edge{}
	vHexs := []*board.Hexagon{}
	x := 0
	for x = range hexagons {
		if len(hexagons[x]) == 0 {
			break
		}

		if x%2 != 0 {
			continue
		}

		//init ptr to first of upper points of hexagon in row
		ptr := hexagons[x][0].P[0]

		//iterate along the row
		for ptr.Right != nil {
			//append point
			vPoints = append(vPoints, ptr)

			//append present edges
			vEdges = append(vEdges, ptr.Right)
			if ptr.Down != nil {
				vEdges = append(vEdges, ptr.Down)
			}

			//iterate
			ptr = ptr.Right.End
		}
		//append final point in row
		vPoints = append(vPoints, ptr)
		if ptr.Down != nil {
			vEdges = append(vEdges, ptr.Down)
		}

		//set ptr to first of lower points of hexagon row
		ptr = hexagons[x][0].P[5]

		//iterate along the row
		for ptr.Right != nil {
			//append point
			vPoints = append(vPoints, ptr)

			//append present edges
			vEdges = append(vEdges, ptr.Right)
			if ptr.Down != nil {
				vEdges = append(vEdges, ptr.Down)
			}

			//iterate
			ptr = ptr.Right.End
		}
		//append final point in row
		vPoints = append(vPoints, ptr)
		if ptr.Down != nil {
			vEdges = append(vEdges, ptr.Down)
		}
	}

	for _, x := range hexagons {
		for _, y := range x {
			vHexs = append(vHexs, y)
		}
	}

	//create board and return
	_board := new(board.Board)
	_board.Edges = vEdges
	_board.Points = vPoints
	_board.Hexs = vHexs
	bRows := []*board.Row{}
	for _, xv := range hexagons {
		r := new(board.Row)
		r.Hexagons = xv
		bRows = append(bRows, r)
	}
	_board.Rows = bRows

	return _board, nil
}

//                  //
// HELPER FUNCTIONS //
//                  //

// HEX & POINT MANIPULATION:

//create hex and link all points
func makeHex(id string, ports string, ocean bool) *board.Hexagon {

	//create and link all edges & points
	this := new(board.Hexagon)
	this.Id = id
	this.P = make([]*board.Point, 6)
	this.E = make([]*board.Edge, 6)
	this.Block = new(board.Block)
	this.Block.Resource = board.Resource_NONE

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

	//ocean
	if ocean {
		this.Ocean = true
	} else {
		this.Ocean = false
	}

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

	return nil
}

func connectHexRight(h1 *board.Hexagon, h2 *board.Hexagon) error {
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

	return nil
}

func connectHexUp(h1 *board.Hexagon, h2 *board.Hexagon, offset string) error {
	if offset == "left" {
		if h2.P[5].Port != nil {
			h1.P[1].Port = h2.P[5].Port
		}
		if h2.P[4].Port != nil {
			h1.P[2].Port = h2.P[4].Port
		}

		h1.P[1].Up = h2.P[5].Up
		h1.P[2].Up = h2.P[4].Up
		connectPoints(h1.P[2], h2.E[3], h2.P[3], "right")

		h2.P[5] = h1.P[1]
		h2.P[4] = h1.P[2]
		connectPoints(h2.P[5], h2.E[4], h2.P[4], "right")

	} else if offset == "right" {
		if h2.P[4].Port != nil {
			h1.P[0].Port = h2.P[4].Port
		}
		if h2.P[3].Port != nil {
			h1.P[1].Port = h2.P[3].Port
		}

		h1.P[0].Up = h2.P[4].Up
		h1.P[1].Up = h2.P[3].Up
		connectPoints(h2.P[5], h2.E[4], h1.P[0], "right")

		h2.P[4] = h1.P[0]
		h2.P[3] = h1.P[1]
		connectPoints(h2.P[4], h2.E[3], h2.P[3], "right")

	} else {
		return (errors.New("connectHex: invalid offset received"))
	}

	return nil
}
