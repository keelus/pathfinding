package main

import (
	_ "image/png"
	"io/ioutil"
	"log"
	"pathfinding/pair"

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
		pos_x, pos_y := ebiten.CursorPosition()

		if _, _, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			drawing = true
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		pos_x, pos_y := ebiten.CursorPosition()

		if _, _, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			drawing = false
		}
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

			if i, j, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
				if !canvas.grid.Cells[i][j].IsWall && !canvas.grid.End.Eq(pair.Pair{i, j}) {
					canvasA.grid.Start = pair.Pair{i, j}
					canvasB.grid.Start = pair.Pair{i, j}
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			pos_x, pos_y := ebiten.CursorPosition()

			if i, j, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
				if !canvas.grid.Cells[i][j].IsWall && !canvas.grid.Start.Eq(pair.Pair{i, j}) {
					canvasA.grid.End = pair.Pair{i, j}
					canvasB.grid.End = pair.Pair{i, j}
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

		if i, j, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			if !canvas.grid.Start.Eq(pair.Pair{i, j}) && !canvas.grid.End.Eq(pair.Pair{i, j}) {
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
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dijkstra vs A*")

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

	activeTool = PENCIL

	canvasA = NewCanvas(550, 550, 200, 40, "Dijkstra")
	canvasB = NewCanvas(550, 550, 800, 40, "A*")

	canvasA.SetGrid(NewGrid(30, pair.New(4, 5), pair.New(1, 2)))
	canvasB.SetGrid(NewGrid(30, pair.New(4, 5), pair.New(1, 2)))

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
