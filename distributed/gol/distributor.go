package gol

type distributorChannels struct {
	events    chan<- Event
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool
	filename    chan <- string
	input <- chan uint8
	output chan <- uint8
}


// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	/*
		if p.Turns == 0{
			finalWorld = Initialworld
		}
		for turn := 0; turn < p.Turns; turn++ {
			finalWorld = calculateNextState(p, Initialworld, turn, c)
			c.events <- TurnComplete{CompletedTurns: turn + 1}
			if turn == p.Turns-1 {
				c.events <- StateChange{
					CompletedTurns: turn + 1,
					NewState:       Executing,
				}
			} else {
				c.events <- StateChange{
					CompletedTurns: p.Turns,
					NewState:       Quitting,
				}
			}
			for y := 0; y < p.ImageHeight; y++ {
				for x := 0; x < p.ImageWidth; x++ {
					Initialworld[y][x] = finalWorld[y][x]
				}
			}
		}

		/*t := 0 * time.Second
		  t++
		  if t == 2 * time.Second{
		     c.events <- AliveCellsCount{
		        CompletedTurns: turn,
		        CellsCount: len(calculateAliveCells(p, Initialworld)),
		     }
		     //t = 0 *time.Second
		  }


		// TODO: Create a 2D slice to store the world.
		// TODO: For all initially alive cells send a CellFlipped Event.
		// TODO: Execute all turns of the Game of Life.
		// TODO: Send correct Events when required, e.g. CellFlipped, TurnComplete and FinalTurnComplete.
		//     See event.go for a list of all events.
		// Make sure that the Io has finished any output before exiting.
		c.ioCommand <- ioCheckIdle
		<-c.ioIdle
		//c.events <- StateChange{p.Turns, Quitting}
		// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
		close(c.events)*/
}