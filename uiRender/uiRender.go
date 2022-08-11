/*
Module for rendering frequently needed UI-components on picture
Renders and provides functions for manipulating state of
*/

package uirender

import (
	"image"

	"github.com/hjkoskel/gomonochromebitmap"
)

type ScrollVerticalSelectMenu struct {
	Bitmaps         []gomonochromebitmap.MonoBitmap
	SelectedIndex   int
	Scroll          int //where this thing is scrolled
	InvertSelection bool
	Arrow           *gomonochromebitmap.MonoBitmap //Moves on left side
	ScrollBar       int                            //How many pixels
}

//Helper function for returning bitmaps generated from strings
func GetStringBitmaps(arr []string, font map[rune]gomonochromebitmap.MonoBitmap, w int, h int, lineSpacing, gap int) []gomonochromebitmap.MonoBitmap {
	result := make([]gomonochromebitmap.MonoBitmap, len(arr))
	workArea := gomonochromebitmap.NewMonoBitmap(w, h, false)
	for i := 0; i < len(arr); i++ {
		usedRect := workArea.Print(arr[i], font, lineSpacing, gap, workArea.Bounds(), true, true, false, true)
		usedRect.Max.X = w
		usedRect.Max.Y = lineSpacing
		//Shrink on vertical
		result[i] = gomonochromebitmap.NewMonoBitmap(usedRect.Dx(), usedRect.Dy(), false)
		result[i].DrawBitmap(workArea, usedRect, image.Point{X: 0, Y: 0}, true, true, false)
		workArea.Fill(usedRect, false)
	}
	return result
}

func (p *ScrollVerticalSelectMenu) Render(w int, h int) gomonochromebitmap.MonoBitmap {
	if p.SelectedIndex < 0 {
		p.SelectedIndex = 0
	}
	if p.SelectedIndex >= len(p.Bitmaps) {
		p.SelectedIndex = len(p.Bitmaps) - 1 //Clamp or rotate? Calling software decides. This just clamps
	}

	//calc scroll value so that selected menu fits nicely on screen (try keep selection inside center of 1/3 of height
	heightCounter := 0
	totalHeightCounter := 0
	for i := 0; i < len(p.Bitmaps); i++ {
		totalHeightCounter += p.Bitmaps[i].H
		if i < p.SelectedIndex {
			heightCounter += p.Bitmaps[i].H
		}
	}
	p.Scroll = intMax(0, heightCounter-h/3)
	p.Scroll = intMin(totalHeightCounter-h, p.Scroll)
	p.Scroll = intMax(0, p.Scroll)

	barHeight := h * h / totalHeightCounter
	barStart := h * p.Scroll / totalHeightCounter

	result := gomonochromebitmap.NewMonoBitmap(w, h, false)

	leftMargin := 0 //TODO set Arrow Width
	//Ok, render bitmaps
	totalHeightCounter = 0
	for i := 0; i < len(p.Bitmaps); i++ {
		drawPos := totalHeightCounter - p.Scroll
		if drawPos > h {
			break //Others are over
		}
		if 0 < (drawPos + p.Bitmaps[i].H) {
			invert := (i == p.SelectedIndex) && p.InvertSelection
			corner := image.Point{X: leftMargin, Y: drawPos}

			result.DrawBitmap(
				p.Bitmaps[i],
				p.Bitmaps[i].Bounds(),
				corner, true, true, invert)
		}
		totalHeightCounter += p.Bitmaps[i].H
	}

	//do we need scroll bar
	if 0 < p.ScrollBar {
		//Black background
		sidebarMargin := 6
		sbA := image.Point{X: w - sidebarMargin - p.ScrollBar, Y: 0}
		sbB := image.Point{X: w - p.ScrollBar, Y: h - 1}

		result.Fill(image.Rectangle{Min: sbA, Max: sbB}, false)

		//fmt.Printf("ScrollBar height=%v start=%v\n", barHeight, barStart)
		//result.Vline(w-2, barStart, barStart+barHeight, true)
		for x := w - p.ScrollBar; x < w; x++ {
			result.Vline(x, barStart, barStart+barHeight, true)
		}

	}
	return result
}

//Private Utils
func intMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func intMin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
