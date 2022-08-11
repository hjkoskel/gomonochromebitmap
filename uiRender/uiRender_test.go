/*
Package for testing/demonstrating ui rendering capabilities
*/
package uirender

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/hjkoskel/gomonochromebitmap"
)

func TestSimple(t *testing.T) {
	fmt.Printf("--- Testing uiRender ---\n")

	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	testfont1 := gomonochromebitmap.GetFont_8x8()

	textArr := []string{"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrott", "Golf", "Hotel", "India", "Juliet", "Kilo", "Lima", "Mike", "November", "Oscar", "Papa", "Quebec", "Romeo", "Sierra", "Tango"}

	menu1 := ScrollVerticalSelectMenu{
		Bitmaps:         GetStringBitmaps(textArr, testfont1, 127, 32, 8, 1),
		SelectedIndex:   2,
		Scroll:          0,
		InvertSelection: true,
		Arrow:           nil,
		ScrollBar:       1,
	}

	menu2 := ScrollVerticalSelectMenu{
		Bitmaps:         GetStringBitmaps(textArr, testfont1, 127, 32, 8, 1),
		SelectedIndex:   5,
		Scroll:          0,
		InvertSelection: true,
		Arrow:           nil,
		ScrollBar:       1,
	}

	test1 := menu1.Render(128, 64)
	out1, _ := os.Create("test1.png")
	png.Encode(out1, test1.GetImage(colTrue, colFalse))
	out1.Close()

	aaa := gomonochromebitmap.BlockGraphicsSettings{
		Clear:       false,
		HaveBorder:  true,
		BorderColor: gomonochromebitmap.FGANSI_BLUE + gomonochromebitmap.BGANSI_YELLOW,
		TextColor:   gomonochromebitmap.FGANSI_BRIGHT_RED + gomonochromebitmap.BGANSI_BRIGHT_GREEN}
	//fmt.Printf("\n%s\n", test1.ToFullBlockChars(aaa))
	//fmt.Printf("\n%s\n", test1.ToHalfBlockChars(aaa))
	fmt.Printf("\n%s\n", test1.ToQuadBlockChars(aaa))

	aaa = gomonochromebitmap.BlockGraphicsSettings{
		Clear:       false,
		HaveBorder:  true,
		BorderColor: gomonochromebitmap.FGANSI_BLUE + gomonochromebitmap.BGANSI_YELLOW,
		TextColor:   ""}
	//fmt.Printf("\n%s\n", test1.ToFullBlockChars(aaa))
	//fmt.Printf("\n%s\n", test1.ToHalfBlockChars(aaa))
	fmt.Printf("\n%s\n", test1.ToQuadBlockChars(aaa))

	test2 := menu2.Render(128, 64)
	out2, _ := os.Create("test2.png")
	png.Encode(out2, test2.GetImage(colTrue, colFalse))
	out2.Close()

	//Large version, is this scalable
	testfont2 := gomonochromebitmap.GetFont_11x16()

	largemenu1 := ScrollVerticalSelectMenu{
		Bitmaps:         GetStringBitmaps(textArr, testfont2, 127, 32, 16, 1),
		SelectedIndex:   2,
		Scroll:          0,
		InvertSelection: true,
		Arrow:           nil,
		ScrollBar:       7,
	}

	largemenu2 := ScrollVerticalSelectMenu{
		Bitmaps:         GetStringBitmaps(textArr, testfont2, 127, 32, 16, 1),
		SelectedIndex:   5,
		Scroll:          0,
		InvertSelection: true,
		Arrow:           nil,
		ScrollBar:       7,
	}

	largetest1 := largemenu1.Render(128, 64)
	largeout1, _ := os.Create("largetest1.png")
	png.Encode(largeout1, largetest1.GetImage(colTrue, colFalse))
	largeout1.Close()

	largetest2 := largemenu2.Render(128, 64)
	largeout2, _ := os.Create("largetest2.png")
	png.Encode(largeout2, largetest2.GetImage(colTrue, colFalse))
	largeout2.Close()

}
