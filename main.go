package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
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

	buttonPencil, buttonEraser, buttonFlagStart, buttonFlagEnd Button
	buttonClearPath, buttonClearCanvas                         Button
	buttonGenerateTerrain                                      Button
	buttonTerrainSizeS, buttonTerrainSizeM, buttonTerrainSizeL Button
	buttonPlay                                                 Button

	buttonMsMinus, buttonMsPlus Button

	categoryTools, categoryClear, categoryCooldown, categoryTerrainSize string

	stopSignal chan struct{}
)

type Tool string

const (
	PENCIL     Tool = "PENCIL"
	ERASER     Tool = "ERASER"
	FLAG_START Tool = "FLAG_START"
	FLAG_END   Tool = "FLAG_END"
)

const (
	SIZE_S int = 22
	SIZE_M int = 55
	SIZE_L int = 110
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
	buttonClearPath.hover(x, y)
	buttonClearCanvas.hover(x, y)
	buttonGenerateTerrain.hover(x, y)
	buttonTerrainSizeS.hover(x, y)
	buttonTerrainSizeM.hover(x, y)
	buttonTerrainSizeL.hover(x, y)
	buttonPlay.hover(x, y)
	buttonMsMinus.hover(x, y)
	buttonMsPlus.hover(x, y)

	if len(canvasA.grid.Cells) == SIZE_S {
		buttonTerrainSizeS.selected = true
		buttonTerrainSizeM.selected = false
		buttonTerrainSizeL.selected = false
	} else if len(canvasA.grid.Cells) == SIZE_M {
		buttonTerrainSizeS.selected = false
		buttonTerrainSizeM.selected = true
		buttonTerrainSizeL.selected = false
	} else {
		buttonTerrainSizeS.selected = false
		buttonTerrainSizeM.selected = false
		buttonTerrainSizeL.selected = true
	}

	if canvasA.grid.Status == STATUS_PATHING || canvasB.grid.Status == STATUS_PATHING {
		buttonPlay.selected = true
		buttonPlay.title = "Stop"

		buttonPencil.disabled = true
		buttonEraser.disabled = true
		buttonFlagStart.disabled = true
		buttonFlagEnd.disabled = true
		buttonClearPath.disabled = true
		buttonClearCanvas.disabled = true
		buttonGenerateTerrain.disabled = true
		buttonTerrainSizeS.disabled = true
		buttonTerrainSizeM.disabled = true
		buttonTerrainSizeL.disabled = true
	} else if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
		buttonPlay.selected = false
		buttonPlay.title = "Play"

		buttonPencil.disabled = false
		buttonEraser.disabled = false
		buttonFlagStart.disabled = false
		buttonFlagEnd.disabled = false
		buttonClearPath.disabled = false
		buttonClearCanvas.disabled = false
		buttonGenerateTerrain.disabled = false
		buttonTerrainSizeS.disabled = false
		buttonTerrainSizeM.disabled = false
		buttonTerrainSizeL.disabled = false
	}

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
		} else if buttonClearPath.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(true)
				canvasB.grid.Restart(true)
			}
		} else if buttonClearCanvas.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(false)
				canvasB.grid.Restart(false)
			}
		} else if buttonGenerateTerrain.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(false)
				canvasB.grid.Restart(false)
				for i, row := range canvasA.grid.Cells {
					for j := range row {
						nodeA := &canvasA.grid.Cells[i][j]
						nodeB := &canvasB.grid.Cells[i][j]
						if !canvasA.grid.Start.Coord.Eq(pair.New(i, j)) && !canvasA.grid.End.Coord.Eq(pair.New(i, j)) {
							isWall := rand.Intn(100) < 20
							nodeA.IsWall = isWall
							nodeB.IsWall = isWall
						}
					}
				}
			}
		} else if buttonTerrainSizeS.hovered {
			canvasA.SetGrid(NewGrid(SIZE_S, pair.New(SIZE_S-1, 0), pair.New(0, SIZE_S-1)))
			canvasB.SetGrid(NewGrid(SIZE_S, pair.New(SIZE_S-1, 0), pair.New(0, SIZE_S-1)))
		} else if buttonTerrainSizeM.hovered {
			canvasA.SetGrid(NewGrid(SIZE_M, pair.New(SIZE_M-1, 0), pair.New(0, SIZE_M-1)))
			canvasB.SetGrid(NewGrid(SIZE_M, pair.New(SIZE_M-1, 0), pair.New(0, SIZE_M-1)))
		} else if buttonTerrainSizeL.hovered {
			canvasA.SetGrid(NewGrid(SIZE_L, pair.New(SIZE_L-1, 0), pair.New(0, SIZE_L-1)))
			canvasB.SetGrid(NewGrid(SIZE_L, pair.New(SIZE_L-1, 0), pair.New(0, SIZE_L-1)))
		} else if buttonPlay.hovered {
			if canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
				canvasA.grid.Restart(true)
				canvasB.grid.Restart(true)
				go canvasA.grid.DoDijkstra()
				go canvasB.grid.DoAStar()
				stopSignal = make(chan struct{})
			} else {
				close(stopSignal)
			}
		} else if buttonMsMinus.hovered {
			if MS_COOLDOWN <= 10 {
				if MS_COOLDOWN > 0 {
					MS_COOLDOWN--
				}
			} else if MS_COOLDOWN <= 20 {
				MS_COOLDOWN -= 5
			} else if MS_COOLDOWN > 20 && MS_COOLDOWN <= 100 {
				MS_COOLDOWN -= 10
			} else if MS_COOLDOWN > 100 {
				MS_COOLDOWN -= 100
			}
		} else if buttonMsPlus.hovered {
			if MS_COOLDOWN < 10 {
				MS_COOLDOWN++
			} else if MS_COOLDOWN < 20 {
				MS_COOLDOWN += 5
			} else if MS_COOLDOWN >= 20 && MS_COOLDOWN < 100 {
				MS_COOLDOWN += 10
			} else if MS_COOLDOWN >= 100 && MS_COOLDOWN < 1000 {
				MS_COOLDOWN += 100
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
		pos_x, pos_y := ebiten.CursorPosition()

		if i, j, canvas := mousePosCoords(&canvasA, &canvasB, pos_x, pos_y); canvas != nil {
			switch activeTool {
			case PENCIL, ERASER:
				drawing = true
			case FLAG_START:
				if !canvas.grid.Cells[i][j].IsWall && !canvas.grid.Start.Coord.Eq(pair.New(i, j)) {
					canvasA.grid.Start = &canvasA.grid.Cells[i][j]
					canvasB.grid.Start = &canvasB.grid.Cells[i][j]
				}
			case FLAG_END:
				if !canvas.grid.Cells[i][j].IsWall && !canvas.grid.Start.Coord.Eq(pair.New(i, j)) {
					canvasA.grid.End = &canvasA.grid.Cells[i][j]
					canvasB.grid.End = &canvasB.grid.Cells[i][j]
				}
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if activeTool == PENCIL || activeTool == ERASER {
			drawing = false
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
	buttonClearPath.Draw(screen)
	buttonClearCanvas.Draw(screen)
	buttonGenerateTerrain.Draw(screen)
	buttonTerrainSizeS.Draw(screen)
	buttonTerrainSizeM.Draw(screen)
	buttonTerrainSizeL.Draw(screen)
	buttonPlay.Draw(screen)
	buttonMsMinus.Draw(screen)
	buttonMsPlus.Draw(screen)

	textColor := color.RGBA{255, 255, 255, 255}
	if canvasA.grid.Status == STATUS_PATHING || canvasB.grid.Status == STATUS_PATHING {
		textColor = color.RGBA{0x4b, 0x4b, 0x4b, 255}
	}
	text.Draw(screen, categoryTools, loadedFont, 15, 55, textColor)
	text.Draw(screen, categoryClear, loadedFont, 15, 210+40, textColor)
	text.Draw(screen, "Canvas size", loadedFont, 15, 210+40+120, textColor)
	text.Draw(screen, categoryCooldown, loadedFont, 15, screenHeight-135+30, color.White)
	text.Draw(screen, fmt.Sprintf("%dms", MS_COOLDOWN), loadedFont, 80, screenHeight-105+30, color.White)

	canvasA.Draw(screen)
	canvasB.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed all:assets/**
var assets embed.FS

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dijkstra vs A*")

	loadedFont = getFont("assets/fonts/mononoki.ttf", 18)

	activeTool = PENCIL

	canvasA = NewCanvas(550, 550, 200, 40, "Dijkstra")
	canvasB = NewCanvas(550, 550, 800, 40, "A*")

	canvasA.SetGrid(NewGrid(SIZE_M, pair.New(SIZE_M-1, 0), pair.New(0, SIZE_M-1)))
	canvasB.SetGrid(NewGrid(SIZE_M, pair.New(SIZE_M-1, 0), pair.New(0, SIZE_M-1)))

	buttonPencil = NewButton(50, 50, (200-100)/2, 100-35, "P", true, getImage("assets/icons/pencil.png"))
	buttonEraser = NewButton(50, 50, (200-100)/2+50, 100-35, "E", false, getImage("assets/icons/eraser.png"))
	buttonFlagStart = NewButton(50, 50, (200-100)/2, 150-35, "F1", false, getImage("assets/icons/greenFlag.png"))
	buttonFlagEnd = NewButton(50, 50, (200-100)/2+50, 150-35, "F2", false, getImage("assets/icons/redFlag.png"))

	buttonClearPath = NewButton(150, 40, (200-150)/2, 250-30+40, "Clear path", false, nil)
	buttonClearCanvas = NewButton(150, 40, (200-150)/2, 290-30+40, "Clear canvas", false, nil)

	buttonGenerateTerrain = NewButton(170, 40, (200-170)/2, 180, "Generate terrain", false, nil)

	buttonTerrainSizeS = NewButton(40, 40, (200-150)/2, 250-30+40+120, "S", false, nil)
	buttonTerrainSizeM = NewButton(40, 40, (200-150)/2+40+15, 250-30+40+120, "M", false, nil)
	buttonTerrainSizeL = NewButton(40, 40, (200-150)/2+40+15+40+15, 250-30+40+120, "L", false, nil)

	buttonPlay = NewButton(150, 40, (200-150)/2, 430+40, "Play", false, nil)

	buttonMsMinus = NewButton(30, 30, (200-150)/2, screenHeight-125+30, "-", false, nil)
	buttonMsPlus = NewButton(30, 30, (200-150)/2+120, screenHeight-125+30, "+", false, nil)

	categoryTools = "Tools"
	categoryClear = "Clear"
	categoryCooldown = "Cooldown"

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func getFont(fpath string, size float64) font.Face {
	fontData, err := assets.ReadFile("assets/fonts/mononoki.ttf")
	if err != nil {
		log.Fatalf("Error opening the font.")
	}
	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	fontType, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	return fontType
}

func getImage(fpath string) *ebiten.Image {
	iconBytes, err := assets.ReadFile(fpath)
	if err != nil {
		log.Fatalf("Error when opening the icon.")
	}

	img, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}
