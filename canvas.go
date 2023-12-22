package main

import (
	"fmt"
	"image/color"

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

	// startCell := ebiten.NewImage(int(cellSize), int(cellSize))
	// startCell.Fill(color.RGBA{255, 0, 0, 255})
	// startCellOp := ebiten.DrawImageOptions{}
	// startCellOp.GeoM.Translate(c.x+float64(c.grid.Start.Y)*cellSize+1, c.y+float64(c.grid.Start.X)*cellSize+1)
	// screen.DrawImage(startCell, &startCellOp)

	// endCell := ebiten.NewImage(int(cellSize), int(cellSize))
	// endCell.Fill(color.RGBA{0, 255, 0, 255})
	// endCellOp := ebiten.DrawImageOptions{}
	// endCellOp.GeoM.Translate(c.x+float64(c.grid.End.Y)*cellSize+1, c.y+float64(c.grid.End.X)*cellSize+1)
	// screen.DrawImage(endCell, &endCellOp)

	rowSize := c.w * 4
	bytes := make([]byte, c.w*c.w*4)

	for i := range bytes {
		bytes[i] = 0x32
	}

	// for i := 0; i < iCellSize; i++ {
	// 	for j := 0; j < iCellSize; j++ {
	// 		bytes[i*rowSize+4*j] = 0xff
	// 		bytes[i*rowSize+4*j+1] = 0xff
	// 		bytes[i*rowSize+4*j+2] = 0xff
	// 		bytes[i*rowSize+4*j+3] = 0xff
	// 	}
	// }

	for i, row := range c.grid.Cells {
		for j, node := range row {
			nodeColor := color.RGBA{100, 100, 100, 255}

			if node.IsWall {
				nodeColor = color.RGBA{0, 0, 0, 255}
			} else if node.Coord == c.grid.Start {
				nodeColor = color.RGBA{0, 255, 0, 255}
			} else if node.Coord == c.grid.End {
				nodeColor = color.RGBA{255, 0, 0, 255}
			} else if node.IsPath {
				nodeColor = color.RGBA{243, 240, 90, 255}
			} else if node.Visited {
				nodeColor = color.RGBA{66, 135, 245, 255}
			}

			drawNodePixels(i, j, cellSize, rowSize, &bytes, nodeColor)
		}
	}

	c.rect.WritePixels(bytes)

	// for i := -1; i < len(c.grid.Cells); i++ {
	// 	line := ebiten.NewImage(1, c.h)
	// 	line.Fill(color.RGBA{255, 0, 0, 255})
	// 	op := ebiten.DrawImageOptions{}
	// 	op.GeoM.Translate(float64((i+1)*cellSize)+float64(i)+1, 0)
	// 	c.rect.DrawImage(line, &op)
	// }

	// for j := -1; j < len(c.grid.Cells); j++ {
	// 	line := ebiten.NewImage(c.w, 1)
	// 	line.Fill(color.RGBA{255, 0, 0, 255})
	// 	op := ebiten.DrawImageOptions{}
	// 	op.GeoM.Translate(0, float64((j+1)*cellSize)+float64(j)+1)
	// 	c.rect.DrawImage(line, &op)
	// }

	// borderL := ebiten.NewImage(1, c.h)
	// borderL.Fill(color.RGBA{50, 50, 50, 255})
	// borderLOp := ebiten.DrawImageOptions{}
	// borderLOp.GeoM.Translate(0, 0)
	// c.rect.DrawImage(borderL, &borderLOp)

	// borderR := ebiten.NewImage(1, c.h)
	// borderR.Fill(color.RGBA{50, 50, 50, 255})
	// borderROp := ebiten.DrawImageOptions{}
	// borderROp.GeoM.Translate(float64(c.w)-1, 0)
	// c.rect.DrawImage(borderR, &borderROp)

	// borderU := ebiten.NewImage(c.w, 1)
	// borderU.Fill(color.RGBA{50, 50, 50, 255})
	// borderUOp := ebiten.DrawImageOptions{}
	// borderUOp.GeoM.Translate(0, 0)
	// c.rect.DrawImage(borderU, &borderUOp)

	// borderD := ebiten.NewImage(c.w, 1)
	// borderD.Fill(color.RGBA{50, 50, 50, 255})
	// borderDOp := ebiten.DrawImageOptions{}
	// borderDOp.GeoM.Translate(0, float64(c.h)-1)
	// c.rect.DrawImage(borderD, &borderDOp)
}

func drawNodePixels(cellI, cellJ int, cellSize int, rowSize int, bytes *[]byte, cellColor color.RGBA) {
	for i := 0; i < cellSize; i++ {
		for j := 0; j < cellSize; j++ {
			// index := base........ + vertical displa. + horizontal displaceme. + 1 row margin??
			index := i*rowSize + 4*j + 4*cellJ*cellSize + cellI*rowSize*cellSize + rowSize*cellI + 4*cellJ + rowSize + 4

			(*bytes)[index] = cellColor.R
			(*bytes)[index+1] = cellColor.G
			(*bytes)[index+2] = cellColor.B
			(*bytes)[index+3] = cellColor.A
		}
	}
}
