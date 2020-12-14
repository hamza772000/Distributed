package gol

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}


// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	filename := make(chan string)
	input := make(chan uint8)
	output := make(chan uint8)
//	turns := make(chan int)
//	initialWorld := make(chan [][]uint8)
	//finalWorld := make(chan [][]uint8)
	//cellsCount := make(chan int)
	//completeTurn := make(chan int)

	/*distributorChannels := distributorChannels{
		events,
		ioCommand,
		ioIdle,
		filename,
		input,
		output,
	}
	//go distributor(p, distributorChannels)*/
	//go controller(p, controllerChannels{})



	controllerChannels := controllerChannels{
		events,
		ioCommand,
		ioIdle,
		filename,
		input,
		output,
		keyPresses,
		//turns,
		//initialWorld,
		//finalWorld,
		//cellsCount,
		//completeTurn,
	}
	go controller(p, controllerChannels)




	/*logicChannels := logicChannels{
		turns,
		initialWorld,
		finalWorld,
		cellsCount,
		completeTurn,
	}*/
	//go logic(logicChannels)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: filename,
		output:   output,
		input:    input,
	}
	go startIo(p, ioChannels)

}
