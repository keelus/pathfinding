package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"pathfinding/pair"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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

	buttonPencil, buttonEraser, buttonFlagStart, buttonFlagEnd, buttonCleanState, buttonCleanCanvas, buttonPlay, buttonMsMinus, buttonMsPlus Button
	categoryTools, categoryClear, categoryCooldown                                                                                           string
)

type Tool string

const (
	PENCIL     Tool = "PENCIL"
	ERASER     Tool = "ERASER"
	FLAG_START Tool = "FLAG_START"
	FLAG_END   Tool = "FLAG_END"
)

type Game struct {
	count int
}

func (g *Game) Update() error {

	x, y := ebiten.CursorPosition()
	buttonPencil.hover(x, y)
	buttonEraser.hover(x, y)
	buttonFlagStart.hover(x, y)
	buttonFlagEnd.hover(x, y)
	buttonCleanState.hover(x, y)
	buttonCleanCanvas.hover(x, y)
	buttonPlay.hover(x, y)
	buttonMsMinus.hover(x, y)
	buttonMsPlus.hover(x, y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if buttonPencil.hovered {
			activeTool = PENCIL
			buttonPencil.selected = true
			buttonEraser.selected = false
			buttonFlagStart.selected = false
			buttonFlagEnd.selected = false
		} else if buttonEraser.hovered {
			activeTool = ERASER
			buttonPencil.selected = false
			buttonEraser.selected = true
			buttonFlagStart.selected = false
			buttonFlagEnd.selected = false
		} else if buttonFlagStart.hovered {
			activeTool = FLAG_START
			buttonPencil.selected = false
			buttonEraser.selected = false
			buttonFlagStart.selected = true
			buttonFlagEnd.selected = false
		} else if buttonFlagEnd.hovered {
			activeTool = FLAG_END
			buttonPencil.selected = false
			buttonEraser.selected = false
			buttonFlagStart.selected = false
			buttonFlagEnd.selected = true
		} else if buttonCleanState.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(true)
				canvasB.grid.Restart(true)
			}
		} else if buttonCleanCanvas.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(false)
				canvasB.grid.Restart(false)
			}
		} else if buttonPlay.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(true)
				canvasB.grid.Restart(true)
				go canvasA.grid.DoDijkstra()
				go canvasB.grid.DoAStar()
				buttonPlay.selected = true
			}
		} else if buttonMsMinus.hovered {
			if MS_COOLDOWN > 0 {
				MS_COOLDOWN -= 10
			}
		} else if buttonMsPlus.hovered {
			if MS_COOLDOWN < 1000 {
				MS_COOLDOWN += 10
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		for i, row := range canvasA.grid.Cells {
			for j := range row {
				nodeA := &canvasA.grid.Cells[i][j]
				nodeB := &canvasB.grid.Cells[i][j]
				if nodeA != canvasA.grid.Start && nodeB != canvasA.grid.End {
					isWall := rand.Intn(100) < 30
					nodeA.IsWall = isWall
					nodeB.IsWall = isWall
				}
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		pos_x, pos_y := ebiten.CursorPosition()

		if i, j, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			if activeTool == PENCIL || activeTool == ERASER {
				drawing = true
			} else if activeTool == FLAG_START {
				if !canvas.grid.Cells[i][j].IsWall && !canvas.grid.Start.Coord.Eq(pair.New(i, j)) {
					canvasA.grid.Start = &canvasA.grid.Cells[i][j]
					canvasB.grid.Start = &canvasB.grid.Cells[i][j]
				}
			} else if activeTool == FLAG_END {
				if !canvas.grid.Cells[i][j].IsWall && !canvas.grid.Start.Coord.Eq(pair.New(i, j)) {
					canvasA.grid.End = &canvasA.grid.Cells[i][j]
					canvasB.grid.End = &canvasB.grid.Cells[i][j]
				}
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		pos_x, pos_y := ebiten.CursorPosition()

		if _, _, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			if activeTool == PENCIL || activeTool == ERASER {
				drawing = false
			}
		}
	}

	if drawing {
		pos_x, pos_y := ebiten.CursorPosition()

		if i, j, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			if !canvas.grid.Start.Coord.Eq(pair.New(i, j)) && !canvas.grid.End.Coord.Eq(pair.New(i, j)) {
				canvasA.grid.Cells[i][j].IsWall = activeTool == PENCIL
				canvasB.grid.Cells[i][j].IsWall = activeTool == PENCIL
			}
		}
	}

	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	buttonPencil.Draw(screen)
	buttonEraser.Draw(screen)
	buttonFlagStart.Draw(screen)
	buttonFlagEnd.Draw(screen)
	buttonCleanState.Draw(screen)
	buttonCleanCanvas.Draw(screen)
	buttonPlay.Draw(screen)
	buttonMsMinus.Draw(screen)
	buttonMsPlus.Draw(screen)

	text.Draw(screen, categoryTools, loadedFont, 15, 55, color.White)
	text.Draw(screen, categoryClear, loadedFont, 15, 210, color.White)
	text.Draw(screen, categoryCooldown, loadedFont, 15, screenHeight-135, color.White)
	text.Draw(screen, fmt.Sprintf("%dms", MS_COOLDOWN), loadedFont, 80, screenHeight-105, color.White)

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

	canvasA.SetGrid(NewGrid(30, pair.New(29, 0), pair.New(0, 29)))
	canvasB.SetGrid(NewGrid(30, pair.New(29, 0), pair.New(0, 29)))

	buttonPencil = NewButton(50, 50, (200-100)/2, 100-35, "P", true)
	buttonEraser = NewButton(50, 50, (200-100)/2+50, 100-35, "E", false)
	buttonFlagStart = NewButton(50, 50, (200-100)/2, 150-35, "F1", false)
	buttonFlagEnd = NewButton(50, 50, (200-100)/2+50, 150-35, "F2", false)
	buttonCleanState = NewButton(150, 40, (200-150)/2, 250-30, "Clear state", false)
	buttonCleanCanvas = NewButton(150, 40, (200-150)/2, 290-30, "Clear canvas", false)
	buttonPlay = NewButton(150, 40, (200-150)/2, 430, "Play", false)
	buttonMsMinus = NewButton(30, 30, (200-150)/2, screenHeight-125, "-", false)
	buttonMsPlus = NewButton(30, 30, (200-150)/2+120, screenHeight-125, "+", false)

	categoryTools = "Tools"
	categoryClear = "Clear"
	categoryCooldown = "Cooldown"

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
