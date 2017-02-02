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
		'X': termbox.ColorBlue,		// wall
		'S': termbox.ColorRed,		// start
		'F': termbox.ColorGreen,	// finish
	}
	color, exists := colorMap[r]
	if exists == true {
		return color
	} else {
		return boardBackgroundColor
	}
}

type gameGrid [][]rune

func initGrid() *gameGrid {
	// layout in human-readable format.
	// rows here don't have to be as long as actual grid (but can't be longer).
	// X's indicate a piece, space indicates no piece.
	readableGrid := []string{
		"              F               ",
		"",
		"XXXXXXXXXXXXXX XXXXXXXXXXXXXXX",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"X                            X",
		"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	gridHeight = len(readableGrid)
	for _, line := range readableGrid {
		if len(line) > gridWidth {
			gridWidth = len(line)
		}
	}

	grid := make(gameGrid, gridWidth)
	for i := range grid {
		grid[i] = make([]rune, gridHeight)
	}

	for y, line := range readableGrid {
		for x, r := range line {
			if r != ' ' {
				grid[x][y] = r
			} else {
				grid[x][y] = ' '
			}
		}
	}
	return &grid
}

func draw(grid *gameGrid) {
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

	// grid pieces
	drawGridPiece := func(r rune, gridX, gridY int, color termbox.Attribute) {
		startScreenX := boardStartX + (gridX * gridPieceWidth)
		startScreenY := boardStartY + (gridY * gridPieceHeight)
		endScreenX := startScreenX + gridPieceWidth - 1
		endScreenY := startScreenY + gridPieceHeight - 1
		for screenX := startScreenX; screenX <= endScreenX; screenX++ {
			for screenY := startScreenY; screenY <= endScreenY; screenY++ {
				termbox.SetCell(screenX, screenY, ' ', color, color)
			}
		}
	}

	for x, row := range *grid {
		for y, r := range row {
			if r != ' ' {
				drawGridPiece(r, x, y, mapColor(r))
			}
		}
	}

	termbox.Flush()
}

func main() {
	grid := initGrid()

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

	draw(grid)

loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			draw(grid)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
