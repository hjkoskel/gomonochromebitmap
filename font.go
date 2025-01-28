/*
This is more "modern" font solution supporting embed directive (avail since go 1.16)

Fonts are generated beforehand and loaded with embed directive.

Font data is like this
- font width uint16
- font height uint16

Then repeating part
- flags, now 0,  1=RLE encoded
- first 8bit char code
- number of chars  first,first+1,first+2
- byte  w*h*n bits of data. Alligned to byte on each char

# Rendering supports rotation

Run length coding on large fonts
Option to scale down font when rendering (embed largest font, scale down for smaller fonts)
*/
package gomonochromebitmap

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
)

// Must be just []byte in type definition so this works
type Font []byte

type FontFileHeader struct {
	Width  uint16
	Height uint16
}

type FontFileBlockHeader struct {
	FirstRune     rune // rune is 32bit
	NumberOfCodes byte //
}

func (p *Font) GetHeader() (FontFileHeader, error) {
	if len(*p) < 4 {
		return FontFileHeader{}, fmt.Errorf("invalid size of font data %v", len(*p))
	}
	a := []byte(*p)
	return FontFileHeader{
		Width:  uint16(a[0]) | uint16(a[1])<<8,
		Height: uint16(a[2]) | uint16(a[3])<<8,
	}, nil
}

func (p *Font) GetRune(c rune) ([]byte, error) {
	h, errH := p.GetHeader()
	if errH != nil {
		return nil, errH
	}
	pixelsPerchar := h.Width * h.Height
	bytesPerchar := int(math.Ceil(float64(pixelsPerchar) / 8))
	//Avoidin bytes package
	i := 4
	var bh FontFileBlockHeader
	arr := []byte(*p)
	for i < len(*p) { //Loop thru
		bh.FirstRune = rune(uint32(arr[i+0]) | uint32(arr[i+1])<<8 | uint32(arr[i+2])<<16 | uint32(arr[i+3])<<24)
		bh.NumberOfCodes = arr[i+4]
		if bh.FirstRune <= c && (bh.FirstRune+rune(bh.NumberOfCodes)) <= c {
			startIndex := i + 5 + int(c-bh.FirstRune)*bytesPerchar
			endIndex := startIndex + bytesPerchar
			if len(arr) <= endIndex {
				return nil, fmt.Errorf("internal error, not enough data on font")
			}
			return arr[startIndex:endIndex], nil
		}
		i += 5 + bytesPerchar
	}
	return nil, fmt.Errorf("rune %s not found in font", strconv.QuoteRune(c))
}

// On grayscale
type Typesetter struct {
	DrawTrue  bool
	DrawFalse bool
	//TODO SUPPORT LATER SourceJumpX int //0 or 1 = full size, 2=half, 3= 1/3
	//TODO SUPPORT LATER SourceJumpY int //0 or 1 = full size, 2=half, 3= 1/3
	//TargetRepeats int //how many repeated per pixel (remember float calc expensive on microcontroller)

	Gap int

	Typeface *Font
}
