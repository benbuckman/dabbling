// Find the way out of a maze.

package main

import (
	"bytes"
	"fmt"
	"github.com/nsf/termbox-go"
	"math"
	"time"
)

const (
	gridWidth  = 30
	gridHeight = 15
	// cells are taller than they are wide, compensate
	gridPieceWidth  = 4
	gridPieceHeight = 2

	windowBackgroundColor = termbox.ColorWhite
	boardBackgroundColor  = termbox.ColorBlack
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

func mapColor(r rune) termbox.Attribute {
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

	printStr(0, 0, boardBackgroundColor, windowBackgroundColor,
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
				termbox.SetCell(screenX, screenY, 'X', color, color)
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
	readableGrid := []string{
		"GRB           BRG",
		"",
		"       BRG       ",
		"",
		"G R B G R B G R B",
		"",
		"       GRB       ",
		"",
		"GRB           BRG",
	}

	// center the human-readable grid on the full grid

	padWidth, padHeight := func() (int, int) {
		height := len(readableGrid)
		width := 0
		for _, line := range readableGrid {
			if len(line) > width {
				width = len(line)
			}
		}
		if width > gridWidth || height > gridHeight {
			panic("Defined grid is too large")
		}
		return int(math.Floor(float64(gridWidth - width) / 2)),
			int(math.Floor(float64(gridHeight - height) / 2))
	}()

	for y, line := range readableGrid {
		for x, r := range line {
			if r != ' ' {
				(*grid)[x+padWidth][y+padHeight] = r
			} else {
				(*grid)[x+padWidth][y+padHeight] = ' '
			}
		}
	}
}

func main() {
	grid := newGameGrid(gridWidth, gridHeight)
	initGrid(grid)

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
