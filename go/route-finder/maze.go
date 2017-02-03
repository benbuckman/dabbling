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

type gridSquare struct {
	// self-awareness/bidirectional ref of position on the grid
	posX, posY int

	// the type/color of the square, coded by letter (see `mapColor`)
	// TODO make an enum of allowed colors
	marker rune
}

var exitLogBuffer bytes.Buffer

func logOnExit(s string) {
	exitLogBuffer.WriteString(s)
	exitLogBuffer.WriteString("\n")
}

func drawText(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func mapColor(r rune) termbox.Attribute {
	colorMap := map[rune]termbox.Attribute{
		'X': termbox.ColorBlue,  // wall
		'S': termbox.ColorRed,   // start
		'F': termbox.ColorGreen, // finish

		'Y': termbox.ColorYellow,
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
		"           F           ",
		"                       ",
		"XXXXXXXXXXX XXXXXXXXXXX",
		"X                     X",
		"X       XXXXXXX       X",
		"X                     X",
		"X          S          X",
		"X                     X",
		"X                     X",
		"X                     X",
		"X                     X",
		"X                     X",
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

	drawText(0, 0, boardBackgroundColor, windowBackgroundColor,
		fmt.Sprintf("%vx%v board, %vx%v screen", boardWidth, boardHeight, outerWidth, outerHeight))

	// board
	for x := boardStartX; x < boardStartX+boardWidth; x++ {
		for y := boardStartY; y < boardStartY+boardHeight; y++ {
			termbox.SetCell(x, y, ' ', boardBackgroundColor, boardBackgroundColor)
		}
	}

	// grid squares
	drawGridPiece := func(square gridSquare) {
		color := mapColor(square.marker)
		startScreenX := boardStartX + (square.posX * gridPieceWidth)
		startScreenY := boardStartY + (square.posY * gridPieceHeight)
		endScreenX := startScreenX + gridPieceWidth - 1
		endScreenY := startScreenY + gridPieceHeight - 1
		for screenX := startScreenX; screenX <= endScreenX; screenX++ {
			for screenY := startScreenY; screenY <= endScreenY; screenY++ {
				termbox.SetCell(screenX, screenY, ' ', color, color)
			}
		}
	}

	for _, row := range *gridPtr {
		for _, square := range row {
			if square.marker != ' ' {
				drawGridPiece(square)
			}
		}
	}

	termbox.Flush()
}

// Return array of up to 8 pointers to surrounding *empty* squares.
func getSurroundingSpaces(gridPtr *gameGrid, sPtr *gridSquare) []*gridSquare {
	posX, posY := (*sPtr).posX, (*sPtr).posY
	/*
		1 2 3
		4 X 5
		6 7 8

		1 = x-1, y-1
		2 = x,   y-1
		3 = x+1, y-1
		...
	*/
	var surrounding = []*gridSquare{}

	getValidEmptySquare := func(x, y int) (*gridSquare) {
		if x >= 0 && y >= 0 &&
			x < gridWidth && y < gridHeight &&
			!(x == posX && y == posY) {
				otherSquarePtr := &(*gridPtr)[x][y]
				if (*otherSquarePtr).marker == ' ' {
					return otherSquarePtr
				}
		}
		return nil
	}

	for y := posY - 1; y <= posY+1; y++ {
		for x := posX - 1; x <= posX+1; x++ {
			otherSquarePtr := getValidEmptySquare(x, y)
			if otherSquarePtr != nil {
				surrounding = append(surrounding, otherSquarePtr)
			}
		}
	}

	logOnExit(fmt.Sprintf("surrounding pieces for %v,%v:", posX, posY))
	for _, square := range surrounding {
		logOnExit(fmt.Sprintf("- %v,%v", (*square).posX, (*square).posY))
	}

	return surrounding
}

func highlightSpaces(ss []*gridSquare) {
	for _, sPtr := range ss {
		(*sPtr).marker = 'Y'
	}
}

func calculateDistancesToExit(gridPtr *gameGrid) {
	highlightSpaces(getSurroundingSpaces(gridPtr, &((*gridPtr)[0][0])))
	highlightSpaces(getSurroundingSpaces(gridPtr, &((*gridPtr)[7][5])))
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
		fmt.Printf("log:\n%s\n", exitLogBuffer.String())
	}()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	draw(gridPtr)

	calculateDistancesToExit(gridPtr)

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
