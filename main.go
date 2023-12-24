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

// WINDOW CONSTANTS
const (
	SCREEN_WIDTH  = 1400
	SCREEN_HEIGHT = 640
)

// TOOL STATUS
var (
	activeTool Tool
	drawing    bool
)

// UI ELEMENTS
var (
	canvasA, canvasB Canvas

	buttonPencil, buttonEraser, buttonFlagStart, buttonFlagEnd Button
	buttonClearPath, buttonClearCanvas                         Button
	buttonGenerateTerrain                                      Button
	buttonTerrainSizeS, buttonTerrainSizeM, buttonTerrainSizeL Button
	buttonPlay                                                 Button
	buttonMsMinus, buttonMsPlus                                Button

	categoryTools, categoryClear, categoryTerrainSize, categoryCooldown string
)

// OTHERS
var (
	canvasSize int
	cellSize   int

	mononokiFFace font.Face
	stopSignal    chan struct{}

	iterationCooldownMS int
)

// MOUSE TOOLS
type Tool string

const (
	PENCIL     Tool = "PENCIL"
	ERASER     Tool = "ERASER"
	FLAG_START Tool = "FLAG_START"
	FLAG_END   Tool = "FLAG_END"
)

// CANVAS SIZES
const (
	SIZE_S int = 22
	SIZE_M int = 55
	SIZE_L int = 110
)

type Game struct {
	sc *ebiten.Image
}

func (g *Game) Update() error {
	// GET MOUSE POSITION
	posX, posY := ebiten.CursorPosition()

	// BUTTON HOVER STATES
	buttonPencil.hover(posX, posY)
	buttonEraser.hover(posX, posY)
	buttonFlagStart.hover(posX, posY)
	buttonFlagEnd.hover(posX, posY)
	buttonClearPath.hover(posX, posY)
	buttonClearCanvas.hover(posX, posY)
	buttonGenerateTerrain.hover(posX, posY)
	buttonTerrainSizeS.hover(posX, posY)
	buttonTerrainSizeM.hover(posX, posY)
	buttonTerrainSizeL.hover(posX, posY)
	buttonPlay.hover(posX, posY)
	buttonMsMinus.hover(posX, posY)

	buttonMsPlus.hover(posX, posY)

	// BUTTON SELECTION STATES
	switch canvasSize {
	case SIZE_S:
		buttonTerrainSizeS.active = true
		buttonTerrainSizeM.active = false
		buttonTerrainSizeL.active = false
	case SIZE_M:
		buttonTerrainSizeS.active = false
		buttonTerrainSizeM.active = true
		buttonTerrainSizeL.active = false
	case SIZE_L:
		buttonTerrainSizeS.active = false
		buttonTerrainSizeM.active = false
		buttonTerrainSizeL.active = true
	}

	if canvasA.grid.Status == STATUS_PATHING || canvasB.grid.Status == STATUS_PATHING {
		buttonPlay.active = true
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
		buttonPlay.active = false
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

	// BUTTON CLICKS
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if buttonPencil.hovered {
			activeTool = PENCIL
			buttonPencil.active = true
			buttonEraser.active = false
			buttonFlagStart.active = false
			buttonFlagEnd.active = false
		} else if buttonEraser.hovered {
			activeTool = ERASER
			buttonPencil.active = false
			buttonEraser.active = true
			buttonFlagStart.active = false
			buttonFlagEnd.active = false
		} else if buttonFlagStart.hovered {
			activeTool = FLAG_START
			buttonPencil.active = false
			buttonEraser.active = false
			buttonFlagStart.active = true
			buttonFlagEnd.active = false
		} else if buttonFlagEnd.hovered {
			activeTool = FLAG_END
			buttonPencil.active = false
			buttonEraser.active = false
			buttonFlagStart.active = false
			buttonFlagEnd.active = true
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
				func() {
					close(stopSignal)
					for canvasA.grid.Status == STATUS_PATHING || canvasB.grid.Status == STATUS_PATHING {
						// Just wait until both closed to prevent closing the closed channel. Can happen in high cooldown, when pressing
						// the stop button multiple times.
					}
				}()
			}
		} else if buttonMsMinus.hovered {
			if iterationCooldownMS <= 10 {
				if iterationCooldownMS > 0 {
					iterationCooldownMS--
				}
			} else if iterationCooldownMS <= 20 {
				iterationCooldownMS -= 5
			} else if iterationCooldownMS > 20 && iterationCooldownMS <= 100 {
				iterationCooldownMS -= 10
			} else if iterationCooldownMS > 100 {
				iterationCooldownMS -= 100
			}
		} else if buttonMsPlus.hovered {
			if iterationCooldownMS < 10 {
				iterationCooldownMS++
			} else if iterationCooldownMS < 20 {
				iterationCooldownMS += 5
			} else if iterationCooldownMS >= 20 && iterationCooldownMS < 100 {
				iterationCooldownMS += 10
			} else if iterationCooldownMS >= 100 && iterationCooldownMS < 1000 {
				iterationCooldownMS += 100
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		canvasA.grid.Status != STATUS_PATHING && canvasB.grid.Status != STATUS_PATHING {
		if i, j, canvas := mousePosCoords(&canvasA, &canvasB, posX, posY); canvas != nil {
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
		if i, j, canvas := mousePosCoords(&canvasA, &canvasB, posX, posY); canvas != nil {
			if !canvas.grid.Start.Coord.Eq(pair.New(i, j)) && !canvas.grid.End.Coord.Eq(pair.New(i, j)) {
				canvasA.grid.Cells[i][j].IsWall = activeTool == PENCIL
				canvasB.grid.Cells[i][j].IsWall = activeTool == PENCIL
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.sc == nil {
		g.sc = screen
	}

	// BUTTON DRAWING
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

	// LEFT TEXTS DRAWING
	textColor := color.RGBA{255, 255, 255, 255}
	if canvasA.grid.Status == STATUS_PATHING || canvasB.grid.Status == STATUS_PATHING {
		textColor = color.RGBA{0x4b, 0x4b, 0x4b, 255}
	}

	text.Draw(screen, categoryTools, mononokiFFace, 15, 55, textColor)
	text.Draw(screen, categoryClear, mononokiFFace, 15, 210+40, textColor)
	text.Draw(screen, "Canvas size", mononokiFFace, 15, 210+40+120, textColor)
	text.Draw(screen, categoryCooldown, mononokiFFace, 15, SCREEN_HEIGHT-135+30, color.White)
	text.Draw(screen, fmt.Sprintf("%dms", iterationCooldownMS), mononokiFFace, 80, SCREEN_HEIGHT-105+30, color.White)

	// CANVAS DRAWING
	canvasA.Draw(screen)
	canvasB.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

//go:embed all:assets/**
var assets embed.FS

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("pathfinding - keelus")
	ebiten.SetWindowIcon([]image.Image{loadImage("assets/icons/greenFlag.png")})

	mononokiFFace = getFont("assets/fonts/mononoki.ttf", 18)

	activeTool = PENCIL

	iterationCooldownMS = 10

	// CREATE CANVAS & SET GRID (default: Medium)
	canvasA = NewCanvas(550, 550, 200, 40, "Dijkstra")
	canvasB = NewCanvas(550, 550, 800, 40, "A*")
	canvasA.SetGrid(NewGrid(SIZE_M, pair.New(SIZE_M-1, 0), pair.New(0, SIZE_M-1)))
	canvasB.SetGrid(NewGrid(SIZE_M, pair.New(SIZE_M-1, 0), pair.New(0, SIZE_M-1)))

	// LEFT BUTTONS
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

	buttonMsMinus = NewButton(30, 30, (200-150)/2, SCREEN_HEIGHT-125+30, "-", false, nil)
	buttonMsPlus = NewButton(30, 30, (200-150)/2+120, SCREEN_HEIGHT-125+30, "+", false, nil)

	// LEFT TEXTS
	categoryTools = "Tools"
	categoryClear = "Clear"
	categoryCooldown = "Cooldown"

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// getFont returns the font located at fpath, read from embed 'assets'. Path should start with 'assets/...'.
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

// getImage returns the *ebiten.Image located at fpath, read from embed 'assets'. Path should start with 'assets/...'.
func getImage(fpath string) *ebiten.Image {
	return ebiten.NewImageFromImage(loadImage(fpath))
}

// getImage returns the *ebiten.Image located at fpath, read from embed 'assets'. Path should start with 'assets/...'.
func loadImage(fpath string) image.Image {
	imgBytes, err := assets.ReadFile(fpath)
	if err != nil {
		log.Fatalf("Error when opening the icon.")
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Fatal(err)
	}

	return img
}
