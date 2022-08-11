package gomonochromebitmap

import (
	"fmt"
	"image"
	"testing"
)

//Testin visually, rendering
func TestBlockRender(t *testing.T) {
	powerOfTwo := NewMonoBitmap(16, 16, false)
	powerOfTwo.Line(image.Point{X: 0, Y: 0}, image.Point{X: 15, Y: 15}, true)

	//Non power of two dimensions
	nontwo := NewMonoBitmap(17, 15, false)
	nontwo.Line(image.Point{X: 0, Y: 0}, image.Point{X: 16, Y: 14}, true)

	//17*len(runelist),7+1+8+1+16+1+6+1+10+1+12+1+5,false)
	simpleRender := BlockGraphics{Clear: false, HaveBorder: true}
	fmt.Printf("%s", simpleRender.ToFullBlockChars(&powerOfTwo))
	fmt.Printf("%s", simpleRender.ToHalfBlockChars(&powerOfTwo))
	fmt.Printf("%s", simpleRender.ToQuadBlockChars(&powerOfTwo))

	fmt.Printf("%s", simpleRender.ToFullBlockChars(&nontwo))
	fmt.Printf("%s", simpleRender.ToHalfBlockChars(&nontwo))
	fmt.Printf("%s", simpleRender.ToQuadBlockChars(&nontwo))

}
