package main

import(
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

const(
	ITEMS_COUNT = 4
	ROOMS_COUNT = 6
)

// Value of 0 represents clean, value > 0 represents degree of dirty.
// RoomState[0] represents floor.
// RoomState[1] represents windows.
// RoomState[2] represents trash.
// RoomState[3] represents desk.
type RoomState [ITEMS_COUNT]int

func (r *RoomState) initRoomState() {
	for i, _ := range r {
		r[i] = rand.Intn(2)
	}
}

type Level [ROOMS_COUNT]RoomState

func (l *Level) initLevel() {
	for i, _ := range l {
		l[i].initRoomState()
	}
}

func (l *Level) messUpLevel(degree int) {
	for i := 0; i < degree; i++ {
		l[rand.Intn(len(l))][rand.Intn(len(l[0]))]++
	}
}

func (l *Level) getAverageDirtyness() float64 {
	var avgDirtness float64 = 0
	for _, room := range l {
		var dirtyness int = 0
		for _, item := range room {
			dirtyness += item
		}
		avgDirtness += float64(dirtyness)
	}
	return avgDirtness / float64(len(l))
}

type McClean struct {
	currentRoom		int
	action			int
	avgDirtyness	float64
}

func (c *McClean) initMcClean() {
	c.currentRoom = 0
	c.action = 0
	c.avgDirtyness = 0
}

// Determine McClean's next action.
func (c *McClean) determineNextAction(l *Level) {

	// determine degree of dirty for whole room
	var dirtyness int = 0
	for _, item := range l[c.currentRoom] {
		dirtyness += item
	}

	// update avgDirtyness
	c.avgDirtyness += (float64(dirtyness) - c.avgDirtyness) / 10
	if c.avgDirtyness < 0 {
		c.avgDirtyness = 0
	}

	// if above average
	if float64(dirtyness) > c.avgDirtyness {

			// determine dirtiest item
			var dirtiest int
			var maxDirtyness int = 0
			for i, item := range l[c.currentRoom] {
				if item > maxDirtyness {
					maxDirtyness = item
					dirtiest = i
				}
			}

			// clean that item
			c.action = dirtiest
			return
	} else {

	// else change room
	c.action = 5
	}
}

func (c *McClean) doAction(l *Level) {
	c.determineNextAction(l)
	if c.action == 5 {
		c.changeRoom(l)
	} else {
		c.clean(l)
	}
}

func (c *McClean) changeRoom(l *Level) {
	c.currentRoom = (c.currentRoom + 1) % len(l)
}

func (c *McClean) clean(level *Level) {
	level[c.currentRoom][c.action] = 0
}

func main() {

	// Set log output to stdout
	log.SetOutput(os.Stdout)

	// Parse command line args
	var messUpArg int
	flag.IntVar(&messUpArg, "messup", 1, "Specify how many items to mess up in every iteration.")

	var speedArg int
	flag.IntVar(&speedArg, "speed", 1000, "Specify the pause in each iteration in milliseconds.")

	var seedArg int
	flag.IntVar(&seedArg, "seed", 123, "Specify random number generator seed.")
	flag.Parse()

	// Seed random number generator
	rand.Seed(123)

	// Level
	var level Level
	level.initLevel()

	// McClean
	var mcClean McClean
	mcClean.initMcClean()

	// Main Loop
	for true {
		log.Print(level, mcClean, level.getAverageDirtyness())
		mcClean.doAction(&level)
		level.messUpLevel(messUpArg)
		time.Sleep(time.Duration(speedArg) * time.Millisecond)
	}
}
