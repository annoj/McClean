package main

import(
	"bufio"
	"log"
	"math/rand"
	"os"
)

type RoomState [4]int

func (r *RoomState) initRoomState() {
	for i, _ := range r {
		r[i] = rand.Intn(2)
	}
}

type Level [5]RoomState

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

	// Seed random number generator
	rand.Seed(123)

	// Level
	var level Level
	level.initLevel()

	// McClean
	var mcClean McClean
	mcClean.initMcClean()

	// Main Loop
	input := bufio.NewScanner(os.Stdin)
	for true {
		log.Print(level, mcClean)
		mcClean.doAction(&level)
		level.messUpLevel(3)
		input.Scan()
	}
}
