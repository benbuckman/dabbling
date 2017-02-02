// Find the way out of a maze.

package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
	"bytes"
	"time"
)

const (
	gridWidth = 30
	gridHeight = 15
	// cells are taller than they are wide, compensate
	gridPieceWidth = 4
	gridPieceHeight = 2

	windowBackgroundColor = termbox.ColorWhite
	boardBackgroundColor = termbox.ColorBlack
)

var exitLogBuffer bytes.Buffer
func logOnExit(s string) {
	exitLogBuffer.WriteString(s)
	exitLogBuffer.WriteString("\n")
}

func printStr(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func mapColor(r rune) (termbox.Attribute) {
	colorMap := map[rune]termbox.Attribute{
		'B': termbox.ColorBlue,
		'R': termbox.ColorRed,
		'G': termbox.ColorGreen,
	}
	color, exists := colorMap[r]
	if exists == true {
		return color
	} else {
		return boardBackgroundColor
	}
}

var drawCounter int = 0

func draw(grid *gameGrid) (error) {
	drawCounter++

	var outerWidth, outerHeight int = termbox.Size()

	boardWidth := gridWidth * gridPieceWidth
	boardHeight := gridHeight * gridPieceHeight

	if boardWidth > outerWidth || boardHeight > outerHeight {
		return fmt.Errorf("The screen is too small (have %v,%v, need %v,%v)",
			outerWidth, outerHeight, boardWidth, boardHeight)
	}

	boardStartX := int((outerWidth - boardWidth) / 2)
	boardStartY := int((outerHeight - boardHeight) / 2)

	termbox.Clear(windowBackgroundColor, windowBackgroundColor)

	printStr(0, 0, boardBackgroundColor, windowBackgroundColor,
		fmt.Sprintf("%vx%v board, %vx%v screen", boardWidth, boardHeight, outerWidth, outerHeight))

	// board
	for x := boardStartX; x < boardStartX + boardWidth; x++ {
		for y := boardStartY; y < boardStartY + boardHeight; y++ {
			termbox.SetCell(x, y, ' ', boardBackgroundColor, boardBackgroundColor)
		}
	}

	// grid pieces
	drawGridPiece := func(r rune, gridX, gridY int, color termbox.Attribute) {
		startScreenX := boardStartX + (gridX * gridPieceWidth)
		startScreenY := boardStartY + (gridY * gridPieceHeight)
		endScreenX := startScreenX + gridPieceWidth - 1
		endScreenY := startScreenY + gridPieceHeight - 1
		if drawCounter == 1 {
			logOnExit(fmt.Sprintf("%v,%v from %v,%v to %v,%v",
				gridX, gridY, startScreenX, startScreenY,
				endScreenX, endScreenY))
		}
		for screenX := startScreenX; screenX <= endScreenX; screenX++ {
			for screenY := startScreenY; screenY <= endScreenY; screenY++ {
				termbox.SetCell(screenX, screenY, 'X', color, color)
			}
		}
	}

	for x, row := range *grid {
		for y, r := range row {
			if drawCounter == 1 {
				logOnExit(fmt.Sprintf("%v,%v: [%s]", x, y, string(r)))
			}
			if r != ' ' {
				drawGridPiece(r, x, y, mapColor(r))
			}
		}
	}

	termbox.Flush()
	return nil
}

type gameGrid [][]rune

func newGameGrid(width, height int) *gameGrid {
	grid := make(gameGrid, width)
	for i := range grid {
		grid[i] = make([]rune, height)
	}
	return &grid
}

func initGrid(grid *gameGrid) {
	// layout in human-readable format.
	// rows here don't have to be as long as actual grid (but can't be longer).
	// X's indicate a piece, space indicates no piece.
	layout := []string{
		"BRG",
		"",
		"G R B",
		"GRB",

		//"  X  X  X  ",
		//"            ",
		//" B  R  B  R",
	}

	for y, line := range layout {
		for x, r := range line {
			if r != ' ' {
				(*grid)[x][y] = r
			} else {
				(*grid)[x][y] = ' '
			}
		}
	}
}

func main() {
	grid := newGameGrid(gridWidth, gridHeight)
	initGrid(grid)

	err := termbox.Init()
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
			err := draw(grid)
			if err != nil {
				logOnExit(err.Error())
				break loop
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
