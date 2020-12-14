package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)


/*type logicChannels struct {
	turns <-chan int
	initialWorld <-chan [][]uint8
	finalWorld chan<- [][]uint8
	cellsCount chan<- int
	completeTurn chan<- int
}*/



const alive = 255
const dead = 0
func mod(x, m int) int {
	return (x + m) % m
}
var turnno int
var alivecellno int
var KeySignal int


func calculateNeighbours(ImageHeight,ImageWidth, x, y int, world [][]uint8) int {
	neighbours := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i != 0 || j != 0 {
				if world[mod(y+i, ImageHeight)][mod(x+j, ImageWidth)] == alive {
					neighbours++
				}
			}
		}
	}
	return neighbours
}
func makeWorld(ImageHeight, ImageWidth int) [][]uint8{
	world := make([][]uint8, ImageHeight)
	for i := range world{
		world[i] = make([]uint8, ImageWidth)
	}
	return world
}
/*func makeImmutableWorld(world [][]uint8) func(y, x int) uint8{
   return func(y, x int) uint8 {
      return world[y][x]
   }
}*/
func calculateNextState(ImageHeight, ImageWidth int,world [][]uint8) [][]uint8 {
	newWorld := makeWorld(ImageHeight, ImageWidth)
	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			neighbours := calculateNeighbours(ImageHeight, ImageWidth, x, y, world)
			if world[y][x] == alive {
				if neighbours == 2 || neighbours == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
					/*c.events <- CellFlipped{
						CompletedTurns: turns+1,
						Cell:           util.Cell{x,y},
					}*/
				}
			} else {
				if neighbours == 3 {
					newWorld[y][x] = alive
					/*c.events <- CellFlipped{
						CompletedTurns: turns+1,
						Cell:           util.Cell{x,y},
					}*/
				} else {
					newWorld[y][x] = dead
				}
			}
		}
	}
	return newWorld
}

/*func worker(p Params, world [][]uint8, turns int, c distributorChannels, out chan<- [][]uint8){
   imagePart := calculateNextState(p, world, turns, c)
   out <- imagePart
}*/
func calculateAliveCells(ImageHeight,ImageWidth int, world [][]uint8) []util.Cell {
	aliveCells := []util.Cell{}
	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			if world[y][x] == alive {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

type LogicOp struct {}






func (s *LogicOp) LogicTwo(req stubs.Request,res *stubs.Response)(err error) {
	//fmt.Print("bAJBAKJdb")

//	ticker := time.NewTicker(2*time.Second)
//	quit := make(chan struct{})


	turn := req.Turns
	initialWorld := req.World

	ImageHeight := len(initialWorld)
	ImageWidth := ImageHeight
	finalWorld := makeWorld(ImageHeight, ImageWidth)
	if turn == 0{
		finalWorld = initialWorld
	}
	//for t:=0; t<turn; t++{
	if turn != 0 {
		for i := 0; i < turn; i++ {
			if KeySignal == 1{
				switch  {
				case KeySignal == 1:
					time.Sleep(time.Second)
				default:
					fmt.Printf("")

				}
				}


			finalWorld = calculateNextState(ImageHeight, ImageWidth, initialWorld)
			initialWorld = finalWorld
			turnno = i
			alivecellno = len(calculateAliveCells(ImageHeight, ImageWidth, finalWorld))


			/*go func() {
				for {
					select {
					case <-ticker.C:
						res.CompletedTurns = i
						turnno = i
						alivecellno = len(calculateAliveCells(ImageHeight, ImageWidth, finalWorld))

					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()*/
		}
	}
	//}
	//close(quit)
	res.World = finalWorld
	//fmt.Println(finalWorld)
	res.FinalAlive = calculateAliveCells(ImageHeight,ImageWidth,finalWorld)

	//res.CellCount = len(calculateAliveCells(ImageHeight,ImageWidth,finalWorld))

	return

}


func (s *LogicOp) SendCells(req stubs.Request,res *stubs.Response)(err error){
	res.Turns = turnno
	res.CellCount = alivecellno
	return
}

func (s *LogicOp) SendKey(req stubs.Keys, res * stubs.ResKeys)(err error){
	KeySignal = req.Val
	return
}



func main() {
	pAddr := flag.String("port","8030","Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&LogicOp{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)

}
