package gol

import (
	"fmt"
	"net/rpc"
	"os"
	"time"

	//"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)
func makeWorld(p Params) [][]uint8{
	world := make([][]uint8, p.ImageHeight)
	for i := range world{
		world[i] = make([]uint8,p.ImageWidth)
	}
	return world
}

type controllerChannels struct {
	events    chan<- Event
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool
	filename  chan<- string
	input <- chan uint8
	output chan <- uint8
	keyPresses <- chan rune


}


var tempWorld [][]uint8
var CCount int
var CTurns int
var FinalAlive []util.Cell
var GetTurns int
var GetCells int


func makeCall(client rpc.Client, world [][]uint8,cellcount, Cturns, turns int, Cellalive []util.Cell){


	request := stubs.Request{
		World:         world,
		CellCount:      cellcount,
		CompletedTurns: Cturns,
		Turns:          turns,
		FinalAlive:     nil,

	}
	response:= new(stubs.Response)
	//response = &stubs.Request{Message:turns}

	client.Call(stubs.LogicOne,request, response)
	//return response.FinalWorld
CTurns = response.CompletedTurns
CCount = response.CellCount
tempWorld = response.World
FinalAlive = response.FinalAlive
	//fmt.Println(response.Turns)
}

func TickCall(client rpc.Client,turns int, cellcount int){
	request := stubs.Ticks{Turns:turns, CellCount:cellcount}
	response := new(stubs.Ticks)
	client.Call(stubs.Ticker,request,response)
	GetTurns = response.Turns
	GetCells = response.CellCount
}

func Keys(client rpc.Client, val int){
	request := stubs.Keys{val}
	response := new(stubs.ResKeys)
	client.Call(stubs.Key,request,response)


}


func calculateAliveCells(ImageHeight,ImageWidth int, world [][]uint8) []util.Cell {
	aliveCells := []util.Cell{}
	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			if world[y][x] == 255 {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}


func controller(p Params, c controllerChannels) {

	ticker := time.NewTicker(2*time.Second)
	quit := make(chan struct{})
	Initialworld := makeWorld(p)
	var test  []util.Cell
	client, _ := rpc.Dial("tcp", "52.87.238.29:8030")
	client1, _ := rpc.Dial("tcp", "52.87.238.29:8030")


	c.ioCommand <- ioInput
	c.filename <-  fmt.Sprintf("%vx%v",p.ImageHeight,p.ImageWidth)

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			Initialworld[y][x] = <- c.input
		}
	}

	defer client.Close()
//	defer client1.Close()

	go func() {
		for {
			select {
			case <-ticker.C:
				TickCall(*client1,0,0)
				c.events <- AliveCellsCount{
					GetTurns+1,
					GetCells,
				}
			case key := <-c.keyPresses:
				switch key {
				case 's':
					c.ioCommand <- ioOutput
					c.filename <-  fmt.Sprintf("%vx%v",p.ImageHeight,p.ImageWidth)
					for y := 0; y < p.ImageHeight; y++ {
						for x := 0; x < p.ImageWidth; x++ {
							c.output <- Initialworld[y][x]
						}
					}
				case 'q':
					c.events <- StateChange{GetTurns+1,Quitting}
					os.Exit(2)
				case 'p':
					c.events <-StateChange{GetTurns+1, Paused}
					Keys(*client,1)
					fmt.Println(GetTurns+1)
					select{
					case x := <-c.keyPresses:
						switch x {
						case 'p':
							c.events <- StateChange{
								GetTurns+1,
								Executing,
							}
							Keys(*client1,0)
							fmt.Println("Continuing")
						default:
							fmt.Println("Press p to continue")
						}
					}



				default:
					fmt.Print("")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()


	if p.Turns == 0{
		CCount = len(calculateAliveCells(p.ImageHeight,p.ImageWidth,Initialworld))
		FinalAlive = calculateAliveCells(p.ImageHeight,p.ImageWidth,Initialworld)
	}else{
			makeCall(*client, Initialworld, 0, 0, p.Turns, test)
			Initialworld = tempWorld
	}

	c.events <- FinalTurnComplete{
		CompletedTurns: p.Turns	,
		Alive:       FinalAlive,
	}

	c.ioCommand <- ioOutput
	c.filename <-  fmt.Sprintf("%vx%v",p.ImageHeight,p.ImageWidth)

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.output <- Initialworld[y][x]
		}
	}

	c.events <- ImageOutputComplete{
		CompletedTurns: p.Turns,
		Filename:       fmt.Sprintf("%vx%vx%v", p.ImageHeight, p.ImageWidth,p.Turns),
	}


	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	close(quit)
	close(c.events)

}







