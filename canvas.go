package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Canvas struct {
	x, y   float64
	w, h   int
	rect   *ebiten.Image
	op     ebiten.DrawImageOptions
	title  string
	titleW int
	grid   Grid
}

func (c Canvas) TopLeftX() float64 {
	return c.x
}

func NewCanvas(w, h int, x, y float64, title string) Canvas {
	rect := ebiten.NewImage(w, h)
	rect.Fill(color.RGBA{25, 25, 25, 255})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	return Canvas{
		rect:  rect,
		op:    op,
		title: title,
		x:     x, y: y, w: w, h: h,
		titleW: text.BoundString(loadedFont, title).Dx(),
	}
}

func (c *Canvas) SetGrid(grid Grid) {
	c.grid = grid
}

func (c *Canvas) Draw(screen *ebiten.Image) {
	screen.DrawImage(c.rect, &c.op)
	text.Draw(screen, c.title, loadedFont, int(c.x)+c.w/2-c.titleW/2, 33, color.White)

	text.Draw(screen, fmt.Sprintf("Path length: %d | Iterations: %d", c.grid.PathLength, c.grid.Iterations), loadedFont, int(c.x)+c.w/2-150, int(c.y)+c.h+22, color.White)

	cellSize := (c.w - len(c.grid.Cells)) / len(c.grid.Cells)

	rowSize := c.w * 4
	bytes := make([]byte, c.w*c.w*4)

	bytes[0] = 255

	for i := range bytes {
		bytes[i] = 0x32
	}

	for i, row := range c.grid.Cells {
		for j, node := range row {
			nodeColor := color.RGBA{100, 100, 100, 255}

			if node.IsWall {
				nodeColor = color.RGBA{0, 0, 0, 255}
			} else if node.Coord == c.grid.Start.Coord {
				nodeColor = color.RGBA{0, 255, 0, 255}
			} else if node.Coord == c.grid.End.Coord {
				nodeColor = color.RGBA{255, 0, 0, 255}
			} else if node.IsPath {
				nodeColor = color.RGBA{243, 240, 90, 255}
			} else if node.Visited {
				nodeColor = color.RGBA{66, 135, 245, 255}
			} else if node.Added {
				nodeColor = color.RGBA{50, 168, 82, 255}
			}

			drawNodePixels(i, j, cellSize, rowSize, &bytes, nodeColor)
		}
	}

	c.rect.WritePixels(bytes)
}

func drawNodePixels(cellI, cellJ int, cellSize int, rowSize int, bytes *[]byte, cellColor color.RGBA) {
	for i := 0; i < cellSize; i++ {
		for j := 0; j < cellSize; j++ {
			// index := base........ + vertical displa. + horizontal displaceme. + 1 row margin??

			index := i * rowSize                // Vertical displacement
			index += rowSize * cellI            // One pixel margin (between rows)
			index += cellI * rowSize * cellSize // Vertical specific displacement

			index += 4 * j                  // Horizontal displacement
			index += 4 * (cellJ * cellSize) // Horizontal specific displacement
			index += 4 * cellJ              // One pixel margin (between cols)

			(*bytes)[index] = cellColor.R
			(*bytes)[index+1] = cellColor.G
			(*bytes)[index+2] = cellColor.B
			(*bytes)[index+3] = cellColor.A
		}
	}
}

func mousePosCoords(canvasA, canvasB *Canvas, pos_x, pos_y int) (int, int, *Canvas) {
	var clickedCanvas *Canvas = nil

	if pos_x >= int(canvasA.x) && pos_x <= int(canvasA.x)+canvasA.w && pos_y >= int(canvasA.y) && pos_y <= int(canvasA.y)+canvasA.h {
		clickedCanvas = canvasA
	} else if pos_x >= int(canvasB.x) && pos_x <= int(canvasB.x)+canvasB.w && pos_y >= int(canvasB.y) && pos_y <= int(canvasA.y)+canvasB.h {
		clickedCanvas = canvasB
	}

	if clickedCanvas == nil {
		return -1, -1, nil
	}

	cellSize := (clickedCanvas.w) / len(clickedCanvas.grid.Cells)
	x, y := pos_x-int(clickedCanvas.x), pos_y-int(clickedCanvas.y)
	j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))

	if i >= 0 && i < len(clickedCanvas.grid.Cells) && j >= 0 && j < len(clickedCanvas.grid.Cells[0]) {
		return i, j, clickedCanvas
	}

	return -1, -1, nil
}
