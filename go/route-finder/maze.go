// Find the way out of a maze.

package main

import (
	"bytes"
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

const (
	// cells are taller than they are wide, compensate
	gridPieceWidth  = 4
	gridPieceHeight = 2

	windowBackgroundColor = termbox.ColorWhite
	boardBackgroundColor  = termbox.ColorBlack
)

var gridWidth, gridHeight int

// TODO change `int`s to smaller types?

type gridSquare struct {
	// self-awareness/bidirectional ref of position on the grid
	posX, posY int

	// the type/color of the square, coded by letter (see `mapColor`)
	// TODO make an enum of allowed colors
	marker rune

	// calculated distance from end square
	distance int
}

var exitLogBuffer bytes.Buffer

func logOnExit(s string) {
	exitLogBuffer.WriteString(s)
	exitLogBuffer.WriteString("\n")
}

var drawTextQueue []string
func drawText(msg string) {
	drawTextQueue = append(drawTextQueue, msg)
}

func mapColor(r rune) termbox.Attribute {
	colorMap := map[rune]termbox.Attribute{
		'X': termbox.ColorBlue,  // wall
		'S': termbox.ColorRed,   // start
		'E': termbox.ColorGreen, // end

		'Y': termbox.ColorYellow,
		'M': termbox.ColorMagenta,
	}
	color, exists := colorMap[r]
	if exists == true {
		return color
	} else {
		return boardBackgroundColor
	}
}

type gameGrid [][]gridSquare

func initialLayout() *[]string {
	// layout in human-readable format.
	// see `mapColor()` for allowed letters.
	return &[]string{
		"           E           ",
		"                       ",
		"XXXXXXXXXXX XXXXXXXXXXX",
		"X                     X",
		"X   XXXXX     XXXX    X",
		"X                     X",
		"X       XXXXXXX       X",
		"X                     X",
		"X   XXXXX     XXXX    X",
		"X                     X",
		"X          S          X",
		"X                     X",
		"X                     X",
		"XXXXXXXXXXXXXXXXXXXXXXX",
	}
}

// Every spot in the grid should have one.
func addSquareToGrid(gridPtr *gameGrid, posX int, posY int, marker rune) {
	square := gridSquare{
		marker: marker,
		posX:   posX,
		posY:   posY,
		distance: -1,	// unknown
	}
	(*gridPtr)[posX][posY] = square
}

func initGrid() *gameGrid {
	readableGrid := *initialLayout()

	gridHeight = len(readableGrid)
	for _, line := range readableGrid {
		if len(line) > gridWidth {
			gridWidth = len(line)
		}
	}

	grid := make(gameGrid, gridWidth)
	for i := range grid {
		grid[i] = make([]gridSquare, gridHeight)
	}

	for y, line := range readableGrid {
		for x, r := range line {
			addSquareToGrid(&grid, x, y, r)
		}
	}
	return &grid
}

var drawnOnce bool = false
func draw(gridPtr *gameGrid) {
	var outerWidth, outerHeight int = termbox.Size()

	boardWidth := gridWidth * gridPieceWidth
	boardHeight := gridHeight * gridPieceHeight

	if boardWidth > outerWidth || boardHeight > outerHeight {
		panic(fmt.Errorf("The screen is too small (have %v,%v, need %v,%v)",
			outerWidth, outerHeight, boardWidth, boardHeight))
	}

	boardStartX := int((outerWidth - boardWidth) / 2)
	boardStartY := int((outerHeight - boardHeight) / 2)

	termbox.Clear(windowBackgroundColor, windowBackgroundColor)

	if drawnOnce == false {
		drawText(fmt.Sprintf("%vx%v board, %vx%v screen", boardWidth, boardHeight, outerWidth, outerHeight))
	}

	// board
	for x := boardStartX; x < boardStartX+boardWidth; x++ {
		for y := boardStartY; y < boardStartY+boardHeight; y++ {
			termbox.SetCell(x, y, ' ', boardBackgroundColor, boardBackgroundColor)
		}
	}

	getSquareCorners := func(square *gridSquare) (startScreenX, startScreenY, endScreenX, endScreenY int) {
		startScreenX = boardStartX + ((*square).posX * gridPieceWidth)
		startScreenY = boardStartY + ((*square).posY * gridPieceHeight)
		endScreenX = startScreenX + gridPieceWidth - 1
		endScreenY = startScreenY + gridPieceHeight - 1
		return
	}

	// grid squares
	drawGridPiece := func(square *gridSquare) {
		// draw colored box
		startScreenX, startScreenY, endScreenX, endScreenY := getSquareCorners(square)
		color := mapColor((*square).marker)
		for screenX := startScreenX; screenX <= endScreenX; screenX++ {
			for screenY := startScreenY; screenY <= endScreenY; screenY++ {
				termbox.SetCell(screenX, screenY, ' ', color, color)
			}
		}
	}

	drawEmptySpace := func(square *gridSquare) {
		// overlay distance
		startScreenX, startScreenY, _, _ := getSquareCorners(square)
		var distanceStr string
		if (*square).distance != -1 {
			distanceStr = fmt.Sprintf("%v", (*square).distance)
		} else {
			distanceStr = "-"
		}
		for i, r := range distanceStr {
			termbox.SetCell(startScreenX + i, startScreenY, r, termbox.ColorWhite, boardBackgroundColor)
		}
	}

	for _, row := range *gridPtr {
		for _, square := range row {
			if square.marker == ' ' {
				// empty space
				drawEmptySpace(&square)
			} else {
				// wall or piece
				drawGridPiece(&square)
			}
		}
	}

	// sidebar text (logging)
	for y, s := range drawTextQueue {
		for x, c := range s {
			termbox.SetCell(x, y, c, boardBackgroundColor, windowBackgroundColor)
		}
	}

	termbox.Flush()
	drawnOnce = true
}

// Return array of up to 4 pointers to adjacent squares,
// including the start and end squares but excluding walls,
// reflecting where a piece can "move" to - up/down/right/left,
// not diagonal!
func getAdjacentSquares(gridPtr *gameGrid, squarePtr *gridSquare) []*gridSquare {
	posX, posY := (*squarePtr).posX, (*squarePtr).posY

	var adjacent = []*gridSquare{}

	addIfValidAndEmpty := func(x, y int) {
		if x >= 0 && y >= 0 && x < gridWidth && y < gridHeight {
			squarePtr := &(*gridPtr)[x][y]
			if (*squarePtr).marker != 'X' {
				adjacent = append(adjacent, squarePtr)
			}
		}
	}
	/*
	  2
	1 0 3
	  4
	*/
	addIfValidAndEmpty(posX - 1, posY)
	addIfValidAndEmpty(posX, posY - 1)
	addIfValidAndEmpty(posX + 1, posY)
	addIfValidAndEmpty(posX, posY + 1)

	return adjacent
}

func getAdjacentSquareWithShortestRoute(gridPtr *gameGrid, squarePtr, endPtr *gridSquare) (lowestDistanceSquarePtr *gridSquare) {
	adjacentSquarePtrs := getAdjacentSquares(gridPtr, squarePtr)

	//tmp
	origMarkers := make([]rune, 0)

	for _, adjacentSquarePtr := range adjacentSquarePtrs {
		// tmp - flash the adjacent squares
		origMarkers = append(origMarkers, (*adjacentSquarePtr).marker)
		(*adjacentSquarePtr).marker = 'M'

		if areGridSquaresEqual(adjacentSquarePtr, endPtr) {
			// at the end
			lowestDistanceSquarePtr = adjacentSquarePtr
		} else if (*adjacentSquarePtr).distance != -1 &&
			(lowestDistanceSquarePtr == nil || (*lowestDistanceSquarePtr).distance > (*adjacentSquarePtr).distance) {
			lowestDistanceSquarePtr = adjacentSquarePtr
		}
	}

	// tmp
	time.Sleep(100 * time.Millisecond)
	for i, adjacentSquarePtr := range adjacentSquarePtrs {
		(*adjacentSquarePtr).marker = origMarkers[i]
	}

	return
}

func areGridSquaresEqual(squarePtr1, squarePtr2 *gridSquare) bool {
	// TODO why isn't this pointer comparison working? identical squares should have same pointers :-/
	//if squarePtr1 == squarePtr2 {
	if (*squarePtr1).posX == (*squarePtr2).posX && (*squarePtr1).posY == (*squarePtr2).posY {
		return true
	} else {
		return false
	}
}

// Start indicated by 'S', end by 'E'. Should be only 1 of each.
func findMazeStartAndEnd(gridPtr *gameGrid) (startPtr, endPtr *gridSquare) {
	for _, row := range *gridPtr {
		for _, square  := range row {
			// need local var so
			square := square

			switch square.marker {
			case 'S':
				if startPtr != nil {
					panic("Too many start positions ('S') marked!")
				}
				startPtr = &square
				drawText(fmt.Sprintf("Found start at %v,%v (dist:%v)", square.posX, square.posY, square.distance))

			case 'E':
				if endPtr != nil {
					panic("Too many end positions ('E') marked!")
				}
				endPtr = &square
				drawText(fmt.Sprintf("Found end at %v,%v", square.posX, square.posY))
			}
		}
	}
	if startPtr == nil {
		panic("Board is missing start position ('S')!")
	}
	if endPtr == nil {
		panic("Board is missing end position ('E')!")
	}
	logOnExit(fmt.Sprintf("start pos: %v,%v, end pos: %v,%v",
		(*startPtr).posX, (*startPtr).posY, (*endPtr).posX, (*endPtr).posY))
	return
}

func calculateDistancesToExit(gridPtr *gameGrid, startPtr, endPtr *gridSquare) {
	// start at end square
	// get its adjacent squares
	// assign a distance (1) to each of them
	// and add them to queue
	// shift next queued:
	//	get its adjacent squares
	//	assign a distance to each of them
	//	and add them to queue
	// shift next queued ...

	queue := make([]*gridSquare, 0)

	countAssignments := 0
	defer func() {
		logOnExit(fmt.Sprintf("Made %v distance assignments for %v squares", countAssignments, int(gridWidth * gridHeight)))
	}()

	assignDistanceToSquare := func(squarePtr *gridSquare, distance int) {
		countAssignments++
		(*squarePtr).distance = distance
	}

	// TODO can this be simplified to `assignNextSquare` ?

	var assignAndQueueAdjacentSquares func(squarePtr *gridSquare)
	assignAndQueueAdjacentSquares = func(squarePtr *gridSquare) {
		distance := (*squarePtr).distance + 1
		adjacentPtrs := getAdjacentSquares(gridPtr, squarePtr)

		logOnExit(fmt.Sprintf("assigning to %v,%v: distance=%v, %v adjacent",
			(*squarePtr).posX, (*squarePtr).posY, distance, len(adjacentPtrs)))

		for _, adjacentSquarePtr := range adjacentPtrs {
			//logOnExit(fmt.Sprintf("at adjacent square %v,%v", (*adjacentSquarePtr).posX, (*adjacentSquarePtr).posY))

			// if we've reached the starting square, then we've already found the shortest route to the end,
			// so bail.
			if areGridSquaresEqual(adjacentSquarePtr, startPtr) {
				logOnExit("Found start square, enough calculating!")
				return
			}

			// ignore if distance is already known,
			// unless new distance is shorter than previously-calculated distance.
			// (b/c might come from multiple directions)
			if (*adjacentSquarePtr).distance == -1 ||
				((*adjacentSquarePtr).distance > -1 && distance < (*adjacentSquarePtr).distance) {

				assignDistanceToSquare(adjacentSquarePtr, distance)

				queue = append(queue, adjacentSquarePtr)
				//logOnExit(fmt.Sprintf("queued %v,%v", (*adjacentSquarePtr).posX, (*adjacentSquarePtr).posY))
			}
		}
		//logOnExit(fmt.Sprintf("After %v,%v, queue has %v", (*squarePtr).posX, (*squarePtr).posX, len(queue)))

		time.Sleep(10 * time.Millisecond)

		if len(queue) > 0 {
			// shift
			nextSquarePtr := queue[0]
			queue = queue[1:]
			assignAndQueueAdjacentSquares(nextSquarePtr)
		}

	}

	assignDistanceToSquare(endPtr, 0)
	assignAndQueueAdjacentSquares(endPtr)
}


func findShortestRoute(gridPtr *gameGrid) {
	startPtr, endPtr := findMazeStartAndEnd(gridPtr)

	calculateDistancesToExit(gridPtr, startPtr, endPtr)

	logOnExit("Finding route...")
	var routeSquares []*gridSquare
	nextSquarePtr := startPtr
	for {
		previousSquareCopy := *nextSquarePtr

		//logOnExit(fmt.Sprintf("Finding route, at %v,%v", (*nextSquarePtr).posX, (*nextSquarePtr).posY))
		routeSquares = append(routeSquares, nextSquarePtr)

		// the end square should have distance 0,
		// so when we're adjacent to it, it should be the closest/next.
		if areGridSquaresEqual(nextSquarePtr, endPtr) {
			drawText(fmt.Sprintf("Found shortest route in %v steps", len(routeSquares)))
			break
		}

		nextSquarePtr = getAdjacentSquareWithShortestRoute(gridPtr, nextSquarePtr, endPtr)

		// highlight step in yellow
		// and un-highlight previous
		(*gridPtr)[previousSquareCopy.posX][previousSquareCopy.posY].marker = previousSquareCopy.marker
		(*nextSquarePtr).marker = 'Y'

		logOnExit(fmt.Sprintf("-> %v,%v (dist %v)",
			(*nextSquarePtr).posX, (*nextSquarePtr).posY, (*nextSquarePtr).distance))

		if nextSquarePtr == nil {
			logOnExit(fmt.Sprintf("Failed to find shortest route! Stuck at %v steps", len(routeSquares)))
			break
		}

		// TODO is it possible for route to double back over itself and get into infinite loop?
		// Probably not b/c that would mean it deviated from the shortest route at some point ... (?)

		time.Sleep(200 * time.Millisecond)
	}
}


func main() {
	gridPtr := initGrid()

	var err error // ?? what is the right way to keep assigning `err` ??

	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		termbox.Close()
		fmt.Println(exitLogBuffer.String())
	}()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	draw(gridPtr)

	go findShortestRoute(gridPtr)

loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			draw(gridPtr)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
