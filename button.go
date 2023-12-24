package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Button struct {
	x, y           float64
	w, h           int
	rect           *ebiten.Image
	op             ebiten.DrawImageOptions
	title          string
	titleW, titleH int
	grid           Grid
	buttonIcon     *ebiten.Image

	hovered  bool
	selected bool
	disabled bool
}

func NewButton(w, h int, x, y float64, title string, selected bool, buttonIcon *ebiten.Image) Button {
	rect := ebiten.NewImage(w, h)
	rect.Fill(color.RGBA{255, 0, 0, 255})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	return Button{
		rect:  rect,
		op:    op,
		title: title,
		x:     x, y: y, w: w, h: h,
		titleW:     text.BoundString(loadedFont, title).Dx(),
		titleH:     text.BoundString(loadedFont, title).Dy(),
		selected:   selected,
		buttonIcon: buttonIcon,
	}
}

func (c *Button) Draw(screen *ebiten.Image) {
	screen.DrawImage(c.rect, &c.op)

	if c.buttonIcon == nil {
		textColor := color.RGBA{255, 255, 255, 255}
		if c.disabled {
			textColor = color.RGBA{0x4b, 0x4b, 0x4b, 255}
		}
		text.Draw(screen, c.title, loadedFont, int(c.x)+c.w/2-c.titleW/2, int(c.y)+c.h/2+18/2-5, textColor)
	} else {
		iconOps := ebiten.DrawImageOptions{}
		iconOps.GeoM.Translate(c.x+2, c.y+2)
		if c.disabled {
			iconOps.ColorScale.ScaleAlpha(0.3)
		}
		screen.DrawImage(c.buttonIcon, &iconOps)
	}

	rowSize := c.w * 4
	bytes := make([]byte, c.w*c.h*4)

	IDLE_WIDTH := 1
	ACTIVE_WIDTH := 2

	var bColor byte = 0x87
	var bColorSelected byte = 0xff
	if c.hovered {
		bColor = 0xff
	} else if c.disabled {
		bColor = 0x4b
		bColorSelected = 0x5b
	}

	for i := 0; i < rowSize; i++ {
		for h := 0; h < IDLE_WIDTH; h++ {
			bytes[i+rowSize*h] = bColor
			bytes[i+rowSize*(c.h-1)-rowSize*h] = bColor
		}
	}

	for i := 0; i < c.h; i++ {
		for h := 0; h < IDLE_WIDTH; h++ {
			bytes[rowSize*i+h*4] = bColor
			bytes[rowSize*i+(h*4)+1] = bColor
			bytes[rowSize*i+(h*4)+2] = bColor
			bytes[rowSize*i+(h*4)+3] = bColor

			bytes[rowSize*(i+1)-(h*4)-1] = bColor
			bytes[rowSize*(i+1)-(h*4)-2] = bColor
			bytes[rowSize*(i+1)-(h*4)-3] = bColor
			bytes[rowSize*(i+1)-(h*4)-4] = bColor
		}
	}

	if c.selected {
		for i := 0; i < rowSize; i++ {
			for h := 0; h < ACTIVE_WIDTH; h++ {
				bytes[i+rowSize*h] = bColorSelected
				bytes[i+rowSize*(c.h-1)-rowSize*h] = bColorSelected
			}
		}

		for i := 0; i < c.h; i++ {
			for h := 0; h < ACTIVE_WIDTH; h++ {
				bytes[rowSize*i+h*4] = bColorSelected
				bytes[rowSize*i+(h*4)+1] = bColorSelected
				bytes[rowSize*i+(h*4)+2] = bColorSelected
				bytes[rowSize*i+(h*4)+3] = bColorSelected

				bytes[rowSize*(i+1)-(h*4)-1] = bColorSelected
				bytes[rowSize*(i+1)-(h*4)-2] = bColorSelected
				bytes[rowSize*(i+1)-(h*4)-3] = bColorSelected
				bytes[rowSize*(i+1)-(h*4)-4] = bColorSelected
			}
		}
	}

	c.rect.WritePixels(bytes)
}

func (b *Button) hover(x, y int) {
	b.hovered = !b.disabled && x >= int(b.x) && x <= int(b.x)+b.w && y >= int(b.y) && y <= int(b.y)+b.h
}
