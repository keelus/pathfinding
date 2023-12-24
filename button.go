package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Button struct {
	x, y       float64
	w, h       int
	rect       *ebiten.Image
	op         ebiten.DrawImageOptions
	title      string
	titleW     int
	grid       Grid
	buttonIcon *ebiten.Image

	hovered  bool
	active   bool
	disabled bool
}

func NewButton(w, h int, x, y float64, title string, active bool, buttonIcon *ebiten.Image) Button {
	rect := ebiten.NewImage(w, h)
	rect.Fill(color.RGBA{255, 0, 0, 255})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	return Button{
		rect:  rect,
		op:    op,
		title: title,
		x:     x, y: y, w: w, h: h,
		titleW:     text.BoundString(mononokiFFace, title).Dx(),
		active:     active,
		buttonIcon: buttonIcon,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	if b.buttonIcon == nil {
		textColor := color.RGBA{255, 255, 255, 255}
		if b.disabled {
			textColor = color.RGBA{0x4b, 0x4b, 0x4b, 255}
		}
		text.Draw(screen, b.title, mononokiFFace, int(b.x)+b.w/2-b.titleW/2, int(b.y)+b.h/2+18/2-5, textColor)
	} else {
		iconOps := ebiten.DrawImageOptions{}
		iconOps.GeoM.Translate(b.x+2, b.y+2)
		if b.disabled {
			iconOps.ColorScale.ScaleAlpha(0.3)
		}
		screen.DrawImage(b.buttonIcon, &iconOps)
	}

	rowSize := b.w * 4
	bytes := make([]byte, b.w*b.h*4)

	var bColor byte = 0x87
	var bColorSelected byte = 0xff

	if b.hovered {
		bColor = 0xff
	} else if b.disabled {
		bColor = 0x4b
		bColorSelected = 0x5b
	}

	if !b.active {
		for i := 0; i < rowSize; i++ {
			for h := 0; h < 1; h++ { // Border width = 1
				bytes[i+rowSize*h] = bColor                 // Top border
				bytes[i+rowSize*(b.h-1)-rowSize*h] = bColor // Bottom border
			}
		}

		for i := 0; i < b.h; i++ {
			for h := 0; h < 1; h++ { // Border width = 1
				for j := 0; j < 4; j++ {
					bytes[rowSize*i+(h*4)+j] = bColor         // Left border
					bytes[rowSize*(i+1)+(h*4)-(j+1)] = bColor // Right border
				}
			}
		}
	} else {
		for i := 0; i < rowSize; i++ {
			for h := 0; h < 2; h++ { // Border width = 2
				bytes[i+rowSize*h] = bColorSelected                 // Top border
				bytes[i+rowSize*(b.h-1)-rowSize*h] = bColorSelected // Bottom border
			}
		}

		for i := 0; i < b.h; i++ {
			for h := 0; h < 2; h++ { // Border width = 2
				for j := 0; j < 4; j++ {
					bytes[rowSize*i+(h*4)+j] = bColorSelected         // Left border
					bytes[rowSize*(i+1)-(h*4)-(j+1)] = bColorSelected // Right border
				}
			}
		}
	}

	b.rect.WritePixels(bytes)
	screen.DrawImage(b.rect, &b.op)
}

func (b *Button) hover(x, y int) {
	b.hovered = !b.disabled && x >= int(b.x) && x <= int(b.x)+b.w && y >= int(b.y) && y <= int(b.y)+b.h
}
