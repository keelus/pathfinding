package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1400
	screenHeight = 640

	RECTANGLE_AMOUNT = 50
	RECTANGLE_WIDTH  = screenWidth/RECTANGLE_AMOUNT - 1
	RECTANGLE_MARGIN = 1

	RECTANGLE_HEIGHT_MULT = screenHeight / RECTANGLE_AMOUNT
)

var (
	loadedFont       font.Face
	canvasA, canvasB Canvas

	activeTool Tool
	drawing    bool
)

type Tool string

const (
	PENCIL Tool = "PENCIL"
	ERASER Tool = "ERASER"
)

type Game struct {
	count int
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		drawing = true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		drawing = false
	}

	if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			canvasA.grid.Restart(false)
			canvasB.grid.Restart(false)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			canvasA.grid.Restart(true)
			canvasB.grid.Restart(true)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			pos_x, pos_y := ebiten.CursorPosition()

			if pos_x >= int(canvasA.x) && pos_x <= int(canvasA.x)+canvasA.w && pos_y >= int(canvasA.y) && pos_y <= int(canvasA.y)+canvasA.h {
				cellSize := (canvasA.w) / len(canvasA.grid.Cells)
				x, y := pos_x-int(canvasA.x), pos_y-int(canvasA.y)
				j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))

				if !canvasA.grid.Cells[i][j].IsWall && !canvasA.grid.End.Eq(image.Point{i, j}) {
					canvasA.grid.Start = image.Point{i, j}
					canvasB.grid.Start = image.Point{i, j}
				}
			}

			if pos_x >= int(canvasB.x) && pos_x <= int(canvasB.x)+canvasB.w && pos_y >= int(canvasB.y) && pos_y <= int(canvasB.y)+canvasB.h {
				cellSize := (canvasB.w) / len(canvasB.grid.Cells)
				x, y := pos_x-int(canvasB.x), pos_y-int(canvasB.y)
				j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))

				if !canvasA.grid.Cells[i][j].IsWall && !canvasA.grid.End.Eq(image.Point{i, j}) {
					canvasA.grid.Start = image.Point{i, j}
					canvasB.grid.Start = image.Point{i, j}
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			pos_x, pos_y := ebiten.CursorPosition()

			if pos_x >= int(canvasA.x) && pos_x <= int(canvasA.x)+canvasA.w && pos_y >= int(canvasA.y) && pos_y <= int(canvasA.y)+canvasA.h {
				cellSize := (canvasA.w) / len(canvasA.grid.Cells)
				x, y := pos_x-int(canvasA.x), pos_y-int(canvasA.y)
				j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))

				if !canvasA.grid.Cells[i][j].IsWall && !canvasA.grid.Start.Eq(image.Point{i, j}) {
					canvasA.grid.End = image.Point{i, j}
					canvasB.grid.End = image.Point{i, j}
				}
			}

			if pos_x >= int(canvasB.x) && pos_x <= int(canvasB.x)+canvasB.w && pos_y >= int(canvasB.y) && pos_y <= int(canvasB.y)+canvasB.h {
				cellSize := (canvasB.w) / len(canvasB.grid.Cells)
				x, y := pos_x-int(canvasB.x), pos_y-int(canvasB.y)
				j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))
				if !canvasA.grid.Cells[i][j].IsWall && !canvasA.grid.Start.Eq(image.Point{i, j}) {
					canvasA.grid.End = image.Point{i, j}
					canvasB.grid.End = image.Point{i, j}
				}
			}
		}

	}

	if canvasA.grid.Status == STATUS_IDLE && canvasB.grid.Status == STATUS_IDLE {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			go canvasA.grid.DoDijkstra()
			go canvasB.grid.DoAStar()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		if activeTool == PENCIL {
			activeTool = ERASER
		} else {
			activeTool = PENCIL
		}
	}

	if drawing {
		pos_x, pos_y := ebiten.CursorPosition()

		if pos_x >= int(canvasA.x) && pos_x <= int(canvasA.x)+canvasA.w && pos_y >= int(canvasA.y) && pos_y <= int(canvasA.y)+canvasA.h {
			cellSize := (canvasA.w) / len(canvasA.grid.Cells)
			x, y := pos_x-int(canvasA.x), pos_y-int(canvasA.y)
			j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))

			if !canvasA.grid.Start.Eq(image.Point{i, j}) && !canvasA.grid.End.Eq(image.Point{i, j}) {
				canvasA.grid.Cells[i][j].IsWall = activeTool == PENCIL
				canvasB.grid.Cells[i][j].IsWall = activeTool == PENCIL
			}
		}

		if pos_x >= int(canvasB.x) && pos_x <= int(canvasB.x)+canvasB.w && pos_y >= int(canvasB.y) && pos_y <= int(canvasB.y)+canvasB.h {
			cellSize := (canvasB.w) / len(canvasB.grid.Cells)
			x, y := pos_x-int(canvasB.x), pos_y-int(canvasB.y)
			j, i := int(math.Floor(float64(x)/float64(cellSize))), int(math.Floor(float64(y)/float64(cellSize)))

			if !canvasA.grid.Start.Eq(image.Point{i, j}) && !canvasA.grid.End.Eq(image.Point{i, j}) {
				canvasA.grid.Cells[i][j].IsWall = activeTool == PENCIL
				canvasB.grid.Cells[i][j].IsWall = activeTool == PENCIL
			}
		}
	}

	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	canvasA.Draw(screen)
	canvasB.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	fontData, err := ioutil.ReadFile("./mononoki_bold.ttf")
	if err != nil {
		log.Fatalf("Error opening the font.")
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	loadedFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dijkstra vs A*")

	activeTool = PENCIL

	canvasA = NewCanvas(550, 550, 200, 40, "Dijkstra")
	canvasB = NewCanvas(550, 550, 800, 40, "A*")

	canvasA.SetGrid(NewGrid(15, image.Point{3, 10}, image.Point{0, 12}))
	canvasB.SetGrid(NewGrid(15, image.Point{3, 10}, image.Point{0, 12}))

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
