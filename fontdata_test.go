package gomonochromebitmap_test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/hjkoskel/gomonochromebitmap"
)

func TestFont_all(t *testing.T) {
	runelist := []rune{
		' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		':', ';', '<', '=', '>', '?', '@',
		'€',
		'A', 'Ä', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'Ö', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'[', '\\', ']', '^', '_', '`',
		'a', 'ä', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'ö', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'{', '|', '}', '~'}

	testfont1 := gomonochromebitmap.GetFont_5x7()
	testfont2 := gomonochromebitmap.GetFont_8x8()
	testfont3 := gomonochromebitmap.GetFont_11x16()
	testfont4 := gomonochromebitmap.GetFont_3x6()
	testfont5 := gomonochromebitmap.GetFont_6x10()
	testfont6 := gomonochromebitmap.GetFont_8x12()
	testfont7 := gomonochromebitmap.GetFont_4x5()

	//TODO https://github.com/BaronWilliams/Vertical-Fonts

	bm := gomonochromebitmap.NewMonoBitmap(17*len(runelist), 7+1+8+1+16+1+6+1+10+1+12+1+5, false)
	for i, r := range runelist {
		x := 17 * i

		fb, hazfb := testfont1[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font1", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 0}, true, true, false)

		fb, hazfb = testfont2[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font2", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 7 + 1}, true, true, false)

		fb, hazfb = testfont3[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 7 + 1 + 8}, true, true, false)

		fb, hazfb = testfont4[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 7 + 1 + 8 + 1 + 16}, true, true, false)

		fb, hazfb = testfont5[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 7 + 1 + 8 + 1 + 16 + 1 + 6}, true, true, false)

		fb, hazfb = testfont6[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 7 + 1 + 8 + 1 + 16 + 1 + 6 + 1 + 10}, true, true, false)

		fb, hazfb = testfont7[r]
		if !hazfb {
			t.Errorf("Char %v not defined in font", r)
		}
		bm.DrawBitmap(fb, fb.Bounds(), image.Point{X: x, Y: 7 + 1 + 8 + 1 + 16 + 1 + 6 + 1 + 10 + 1 + 12}, true, true, false)

	}

	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	out, _ := os.Create("testFonts.png")
	png.Encode(out, bm.GetImage(colTrue, colFalse))
	out.Close()

	/*
		testfont2:=gomonochromebitmap.GetFont_5x7()
		testfont3:=gomonochromebitmap.GetFont_11x16()
	*/
}
