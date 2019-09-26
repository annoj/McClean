package main

import(
	"flag"
	"log"
	"math/rand"
	"time"
)

// Value of 0 represents clean, 1 represents dirty.
// RoomState[0] represents floor.
// RoomState[1] represents windows.
// RoomState[2] represents trash.
// RoomState[3] represents desk.
type RoomState [4]int

func (r *RoomState) initRoomState() {
	for i, _ := range r {
		r[i] = rand.Intn(2)
	}
}

type Level [6]RoomState

func (l *Level) initLevel() {
	for i, _ := range l {
		l[i].initRoomState()
	}
}

func (l *Level) messUpLevel(degree int) {
	for i := 0; i < degree; i++ {
		l[rand.Intn(len(l))][rand.Intn(len(l[0]))] = 1
	}
}

type McClean struct {
	currentRoom	int
	action		int
}

func (c *McClean) initMcClean() {
	c.currentRoom = 0
	c.action = 0
}

// ISSUE: McClean starts with last action, could result in infinite loop
// 		  should start with i := c.action + 1 and handle i == 5
func (c *McClean) determineNextAction(l *Level) {
	for i := c.action; i < len(l[c.currentRoom]); i++ {
		if l[c.currentRoom][i] == 1 {
			c.action = i
			return
		}
	}
	c.action = 5
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
	c.action = 0
}

func (c *McClean) clean(level *Level) {
	level[c.currentRoom][c.action] = 0
}

func main() {

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
		log.Print(level, mcClean)
		mcClean.doAction(&level)
		level.messUpLevel(messUpArg)
		time.Sleep(time.Duration(speedArg) * time.Millisecond)
	}
}
