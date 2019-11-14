package main

import(
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const(
	ITEMS_COUNT = 4
	ROOMS_COUNT = 6
	CHANGE_ROOM = ITEMS_COUNT + 1
)

// ---------- ROOMSTATE ---------- //

// RoomState[0] represents floor.
// RoomState[1] represents windows.
// RoomState[2] represents trash.
// RoomState[3] represents desk.
// Value of 0 represents clean, value > 0 represents degree of dirty.
type RoomState [ITEMS_COUNT]int

func (r *RoomState) initRoomState() {
	for i, _ := range r {
		r[i] = rand.Intn(2)
	}
}

// ---------- LEVEL ----------- //

type Level [ROOMS_COUNT]RoomState

func (l *Level) initLevel() {
	for i, _ := range l {
		l[i].initRoomState()
	}
}

func (l *Level) messUpLevel(degree int, u *LevelRoomUsage) {
	for i := 0; i < degree; i++ {
		var room int = rand.Intn(ROOMS_COUNT)
		var item int = rand.Intn(ITEMS_COUNT)
		l[room][item] += u[room]
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
	return avgDirtness / float64(ROOMS_COUNT)
}

// ---------- LEVELROOMUSAGE ---------- //

type LevelRoomUsage [ROOMS_COUNT]int

func (u *LevelRoomUsage) initLevelRoomUsage(ui int) {
	for i, _ := range u {
		u[i] = rand.Intn(ui) + 1 // +1 to avoid messing up by 0 and represent arg value properly
	}
}

// ---------- BELIEFS ---------- // 

type Beliefs struct {
	currentRoom		int
	action			int
	avgDirtyness	float64
}

// ---------- DESIRES ---------- //

// ---------- INTENTIONS ------- //

// ---------- MCCLEAN ---------- //

type McClean struct {
	beliefs			Beliefs
//	desires			desires
//	intentions		Intentions
}

func (c *McClean) initMcClean() {
	c.beliefs.currentRoom = 0
	c.beliefs.action = 0
	c.beliefs.avgDirtyness = 0
}

// Determine McClean's next action.
func (c *McClean) determineNextAction(l *Level) {

	// determine degree of dirty for whole room
	var dirtyness int = 0
	for _, item := range l[c.beliefs.currentRoom] {
		dirtyness += item
	}

	// update avgDirtyness
	c.beliefs.avgDirtyness += (float64(dirtyness) - c.beliefs.avgDirtyness) / 10
	if c.beliefs.avgDirtyness < 0 {
		c.beliefs.avgDirtyness = 0
	}

	// if above average
	if float64(dirtyness) > c.beliefs.avgDirtyness {

			// determine dirtiest item
			var dirtiest int
			var maxDirtyness int = 0
			for i, item := range l[c.beliefs.currentRoom] {
				if item > maxDirtyness {
					maxDirtyness = item
					dirtiest = i
				}
			}

			// clean that item
			c.beliefs.action = dirtiest
			return
	} else {

	// else change room
	c.beliefs.action = CHANGE_ROOM
	}
}

func (c *McClean) changeRoom(l *Level) {
	c.beliefs.currentRoom = (c.beliefs.currentRoom + 1) % ROOMS_COUNT
}

func (c *McClean) clean(level *Level) {
	level[c.beliefs.currentRoom][c.beliefs.action] = 0
}

func (c *McClean) doAction(l *Level) {
	c.determineNextAction(l)
	if c.beliefs.action == CHANGE_ROOM {
		c.changeRoom(l)
	} else {
		c.clean(l)
	}
}

// ---------- GENERAL FUNCTIONS ---------- //

func printState(l *Level, m *McClean, csv bool) {
	if csv {
		for i := 0; i < ROOMS_COUNT; i++ {
			for j := 0; j < ITEMS_COUNT; j++ {
				fmt.Printf("%d,", l[i][j])
			}
		}
		fmt.Printf("%d,%d,%f,%f\n", m.beliefs.currentRoom, m.beliefs.action, m.beliefs.avgDirtyness, l.getAverageDirtyness())
	} else {
		for i := 0; i < ROOMS_COUNT; i++ {
			fmt.Printf("Room %d: ", i)
			for j := 0; j < ITEMS_COUNT; j++ {
				fmt.Printf("%2d ", l[i][j])
			}
		}
		fmt.Printf("McClean: currentRoom = %d, action = %d, avgDirtness = %2.2f Average Dirtyness: %2.2f\n", m.beliefs.currentRoom, m.beliefs.action, m.beliefs.avgDirtyness, l.getAverageDirtyness())
	}
}

func main() {

	// Parse command line args
	var messUpArg int
	flag.IntVar(&messUpArg, "messup", 1, "Specify how many items to mess up in every iteration.")

	var speedArg int
	flag.IntVar(&speedArg, "speed", 1000, "Specify the pause in each iteration in milliseconds.")

	var seedArg int
	flag.IntVar(&seedArg, "seed", 123, "Specify random number generator seed.")

	var roomUsageIntervalArg int
	flag.IntVar(&roomUsageIntervalArg, "usage", 3, "Specify maximal room usage. Usage is determined randomly, depending on seed.")

	// TODO: Is this idiomatic?
	var csvArg bool
	b := flag.Bool("csv", false, "Output state in csv format.")
	flag.Parse()
	csvArg = *b

	// Seed random number generator
	rand.Seed(int64(seedArg))

	// Level
	var level Level
	level.initLevel()

	// Room usage
	var usage LevelRoomUsage
	usage.initLevelRoomUsage(roomUsageIntervalArg)
	log.Print("Room usage: ", usage)

	// McClean
	var mcClean McClean
	mcClean.initMcClean()

	// Main Loop
	for true {
		printState(&level, &mcClean, csvArg)
		mcClean.doAction(&level)
		level.messUpLevel(messUpArg, &usage)
		time.Sleep(time.Duration(speedArg) * time.Millisecond)
	}
}
