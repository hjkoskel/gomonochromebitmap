// Package provides functions for operating monochrome images
package gomonochromebitmap

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type MonoBitmap struct {
	Pix []uint32 //using byte vs uint16 vs uint32 vs uint64...  32bit shoud suit well for raspi1/2
	W   int
	H   int
}

// NewMonoBitmap initializes empty bitmap fill is default value
func NewMonoBitmap(w int, h int, fill bool) MonoBitmap {
	result := MonoBitmap{W: w, H: h, Pix: make([]uint32, w*h/32+1)}
	if fill {
		for index := range result.Pix {
			result.Pix[index] = 0xFFFFFFFF
		}
	}
	return result
}

// NewMonoBitmapFromImage initializes bitmap from image. Color conversion: if any Red,Green or Blue value is over threshold then pixel is true
func NewMonoBitmapFromImage(img image.Image, area image.Rectangle, threshold byte, invert bool) MonoBitmap {
	b := img.Bounds()
	w := b.Max.X
	h := b.Max.Y
	result := NewMonoBitmap(w, h, false)
	for x := 0; x <= w; x++ {
		for y := 0; y < h; y++ {
			vr, vg, vb, _ := img.At(x, y).RGBA()
			v := byte((max(int(vr), max(int(vg), int(vb)))) >> 8)
			if v > threshold {
				result.SetPix(x, y, !invert)
			} else {
				result.SetPix(x, y, invert)
			}
		}
	}
	return result
}

// Bounds returns W,H in Rect struct
func (p *MonoBitmap) Bounds() image.Rectangle {
	return image.Rect(0, 0, p.W, p.H)
}

// RLEdecode decodes run length compressed bitmap data
func (p *MonoBitmap) RLEdecode(activeFirst bool, data []byte) error {
	//TODO line drawing... LESS naive solution
	activeNow := activeFirst

	for y := 0; y < p.H; y++ {
		for x := 0; x < p.W; x++ {
			for data[0] == 0 { //remove zeros
				data = data[1:]
				if len(data) == 0 {
					return fmt.Errorf("runned out of RLE data")
				}
				//fmt.Printf("data len =%v\n", len(data))
				activeNow = !activeNow
			}
			data[0]--
			p.SetPix(x, y, activeNow)
		}
	}
	return nil
}

// RLEencodes bitmap in runlength compressed format
func (p *MonoBitmap) RLEencode(activeFirst bool) []byte {
	counter := byte(0)
	activeNow := activeFirst
	result := []byte{}

	for y := 0; y < p.H; y++ {
		for x := 0; x < p.W; x++ {
			if activeNow == p.GetPixNoCheck(x, y) {
				if counter < 254 {
					counter++ //Nothing changed increase
				} else {
					//overflow
					result = append(result, 255) //add maximum..this is pixel by pixel
					activeNow = !activeNow
					counter = 0
				}
			} else {
				activeNow = !activeNow
				result = append(result, counter) //write previous value
				counter = 1
			}
		}
	}
	result = append(result, counter)
	return result
}

// BenchmarkGetBlankImage-16    	     750	   1588057 ns/op	 1925196 B/op	       4 allocs/op
// GetImage Creates RGBA image from bitmap
func (p *MonoBitmap) GetImage(trueColor color.Color, falseColor color.Color) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, p.W, p.H))

	trueR, trueG, trueB, trueA := trueColor.RGBA()
	falseR, falseG, falseB, falseA := falseColor.RGBA()

	offset := 0
	//fmt.Printf("STRIDE %v\n", result.Stride)
	n := p.W * p.H

	for i := 0; i < n; i++ {
		index := i / 32
		pv := ((p.Pix[index] & uint32(1<<uint32(i%32))) > 0)
		if pv {
			result.Pix[offset+0] = uint8(trueR >> 8)
			result.Pix[offset+1] = uint8(trueG >> 8)
			result.Pix[offset+2] = uint8(trueB >> 8)
			result.Pix[offset+3] = uint8(trueA >> 8)
		} else {
			result.Pix[offset+0] = uint8(falseR >> 8)
			result.Pix[offset+1] = uint8(falseG >> 8)
			result.Pix[offset+2] = uint8(falseB >> 8)
			result.Pix[offset+3] = uint8(falseA >> 8)
		}
		offset += 4
	}
	return result
}

// Use each monochrome bitmap as bit in color palette index. https://en.wikipedia.org/wiki/Planar_(computer_graphics)
func CreatePlanarColorImage(planes []MonoBitmap, palette []color.Color) (*image.RGBA, error) {
	if len(palette) == 0 {
		return nil, fmt.Errorf("palette missing")
	}
	if len(planes) == 0 {
		return nil, fmt.Errorf("no input images")
	}

	result := image.NewRGBA(planes[0].Bounds())
	//Check dimensions
	for _, plane := range planes {
		if plane.W != planes[0].W || plane.H != planes[0].H {
			return nil, fmt.Errorf("all bitplanes must have same dimensions")
		}
	}

	if 1<<len(planes) != len(palette) {
		return nil, fmt.Errorf("have %v bitplanes but pallette have %v entries. Must have %v entries", len(planes), len(palette), 1<<len(planes))
	}

	for y := 0; y < planes[0].H; y++ {
		for x := 0; x < planes[0].W; x++ {
			paletteIndex := uint32(0)
			for bit, plane := range planes {
				if plane.GetPixNoCheck(x, y) {
					paletteIndex |= 1 << bit
				}
			}
			result.Set(x, y, palette[paletteIndex])
		}
	}

	return result, nil
}

func (p *MonoBitmap) GetFgBgImage(fgPic image.Image, bgPic image.Image) (image.Image, error) {
	bb := bgPic.Bounds()
	fb := fgPic.Bounds()
	if fb.Dx() < p.W || bb.Dx() < p.W || fb.Dy() < p.H || bb.Dy() < p.H {
		return nil, fmt.Errorf("bg %vX%v and fg %vX%v must be equal or larger than mono bitmap %v,%v", bb.Dx(), bb.Dy(), fb.Dx(), fb.Dy(), p.W, p.H)
	}

	result := image.NewRGBA(p.Bounds())

	//TODO optimize byte -> 8 pixels
	for y := 0; y < p.H; y++ {
		for x := 0; x < p.W; x++ {
			if p.GetPixNoCheck(x, y) {
				result.Set(x, y, fgPic.At(x, y))
			} else {
				result.Set(x, y, bgPic.At(x, y))
			}
		}
	}
	return result, nil
}

/*
Generates image that is rendered like it was LCD. Space in between segments is transparent
upper vs lower color allows to render two color LCD's  (like cyan and yellow strip)
*/
func (p *MonoBitmap) GetDisplayImage(trueColorUpper color.Color, trueColorDowner color.Color, upperRows int, falseColor color.Color, pixelW int, pixelH int, gapW int, gapH int) *image.RGBA {
	totW := p.W*(pixelW+gapW) - gapW
	totH := p.H*(pixelH+gapH) - gapH
	result := image.NewRGBA(image.Rect(0, 0, totW, totH))

	for x := 0; x < p.W; x++ {
		xp := x * (pixelW + gapW)
		for y := 0; y < p.H; y++ {
			yp := y * (pixelW + gapW)
			colo := falseColor
			if p.GetPixNoCheck(x, y) {
				colo = trueColorDowner
				if upperRows > y {
					colo = trueColorUpper
				}
			}
			for i := 0; i < pixelW; i++ {
				for j := 0; j < pixelH; j++ {
					result.Set(xp+i, yp+j, colo)
				}
			}
		}
	}
	return result
}

// Get view (size w,h) for display. Starting from corner p0. Result is centered. If p0 goes outside, function clamps view
// This is meant only for producing scrollable output picture for display. Better scaling functions elsewhere
// pxStep=0, autoscale, so bitmap will fit
// pxStep=1 is 1:1
// pxStep=2 is 2:1 (50% scale)
// pxStep=3 is 3:1 (25% scale)
// pxStep is limited to point where whole bitmap is visible
// Returns: image, actual cornerpoint and zoom used. Useful if UI includes
func (p *MonoBitmap) GetView(w int, h int, p0 image.Point, pxStep int, edges bool) MonoBitmap {
	result := NewMonoBitmap(w, h, false)
	maxStep := math.Max(float64(p.W)/float64(w), float64(p.H)/float64(h)) //In decimal
	corner := image.Point{X: max(p0.X, 0), Y: max(p0.Y, 0)}               //Limit point inside

	var step float64

	step = math.Min(float64(pxStep), math.Ceil(maxStep)) //Limits zooming out too much
	if pxStep == 0 {                                     //Autoscale
		step = maxStep
		corner = image.Point{X: 0, Y: 0}
		if maxStep <= 0.5 { //Scale bigger
			//TODO: this is only reason why decimal step is now needed. Todo later integer step
		} else {
			step = math.Ceil(step)
		}
	}

	//Limit corner
	corner.X = min(corner.X, int(float64(p.W)-step*float64(w)))
	corner.Y = min(corner.Y, int(float64(p.H)-step*float64(h)))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			a := int(float64(x)*step) + corner.X
			b := int(float64(y)*step) + corner.Y
			if (a < 0) || (b < 0) || (p.W <= a) || (p.H <= b) {
				result.SetPix(x, y, edges)
			} else {
				result.SetPix(x, y, p.GetPix(a, b))
			}
		}
	}
	return result
}

// Fills rectangle area from map. Used for clearing image
func (p *MonoBitmap) Fill(area image.Rectangle, fillValue bool) {
	//Naive solution. TODO later faster solution
	for y := area.Min.Y; y <= area.Max.Y; y++ {
		p.Hline(area.Min.X, area.Max.X, y, fillValue)
	}
}

func (p *MonoBitmap) Rectangle(area image.Rectangle) {
	p.Hline(area.Min.X, area.Max.X, area.Min.Y, true)
	p.Hline(area.Min.X, area.Max.X, area.Max.Y, true)

	p.Vline(area.Min.X, area.Min.Y, area.Max.Y, true)
	p.Vline(area.Max.X, area.Min.Y, area.Max.Y, true)
}

// Inverts pixel values
func (p *MonoBitmap) Invert(area image.Rectangle) {
	//Naive solution. TODO later faster solution

	//TODO check limits
	x0 := min(max(0, area.Min.X), p.W-1)
	x1 := min(max(0, area.Max.X), p.W-1)

	y0 := min(max(0, area.Min.Y), p.H-1)
	y1 := min(max(0, area.Max.Y), p.H-1)

	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			p.SetPixNoCheck(x, y, !p.GetPixNoCheck(x, y))
		}
	}
}

// Flip with axis in vertical
func (p *MonoBitmap) FlipV() {
	var v bool
	var i int
	for x := 0; x < p.W/2; x++ {
		for y := 0; y < p.H; y++ {
			v = p.GetPixNoCheck(x, y)
			i = p.W - x - 1
			p.SetPix(x, y, p.GetPixNoCheck(i, y))
			p.SetPix(i, y, v)
		}
	}
}

func (p *MonoBitmap) FlipH() {
	var v bool
	var i int
	for x := 0; x < p.W; x++ {
		for y := 0; y < p.H/2; y++ {
			v = p.GetPixNoCheck(x, y)
			i = p.H - y - 1
			p.SetPix(x, y, p.GetPixNoCheck(x, i))
			p.SetPix(x, i, v)
		}
	}

}

// Rotates in 90 decree steps
// +1=90 clockwise
// -1=90 anticlockwise
// +2=180 clockwise etc...
func (p *MonoBitmap) Rotate90(turn90 int) {
	angle := turn90 % 4
	result := NewMonoBitmap(p.W, p.H, false)
	switch angle {
	case 0:
		return //NOP
	case 1, -3:
		result.W = p.H
		result.H = p.W
		for x := 0; x < p.W; x++ {
			for y := 0; y < p.H; y++ {
				result.SetPix(p.H-y-1, x, p.GetPixNoCheck(x, y))
			}
		}
	case 2, -2:
		for x := 0; x < p.W; x++ {
			for y := 0; y < p.H; y++ {
				result.SetPix(p.W-x-1, p.H-y-1, p.GetPixNoCheck(x, y))
			}
		}
	case 3, -1:
		result.W = p.H
		result.H = p.W
		for x := 0; x < p.W; x++ {
			for y := 0; y < p.H; y++ {
				result.SetPix(y, p.W-x-1, p.GetPixNoCheck(x, y))
			}
		}
	}
	p.W = result.W
	p.H = result.H
	p.Pix = result.Pix
}

// Bresenham's line, copied from http://41j.com/blog/2012/09/bresenhams-line-drawing-algorithm-implemetations-in-go-and-c/
func (p *MonoBitmap) Line(p0In image.Point, p1In image.Point, value bool) {

	bou := p.Bounds()
	bou.Max.X--
	bou.Max.Y--
	p0, p1 := ClipLine(p0In, p1In, bou)
	//fmt.Printf("p0 %#v  ->  %#v\np1 %#v  ->  %#v\n", p0In, p0, p1In, p1)

	//TODO CLIP INTO VIEWPORT!
	var cx int32 = int32(p0.X)
	var cy int32 = int32(p0.Y)

	var dx int32 = int32(p1.X) - cx
	var dy int32 = int32(p1.Y) - cy
	if dx < 0 {
		dx = 0 - dx
	}
	if dy < 0 {
		dy = 0 - dy
	}

	var sx int32
	var sy int32
	if cx < int32(p1.X) {
		sx = 1
	} else {
		sx = -1
	}
	if cy < int32(p1.Y) {
		sy = 1
	} else {
		sy = -1
	}
	var err int32 = dx - dy

	for {
		p.SetPixNoCheck(int(cx), int(cy), value)
		if (cx == int32(p1.X)) && (cy == int32(p1.Y)) {
			return
		}
		var e2 int32 = 2 * err
		if e2 > (0 - dy) {
			err = err - dy
			cx = cx + sx
		}
		if e2 < dx {
			err = err + dx
			cy = cy + sy
		}
	}

}

// Horizontal line for filling

func (p *MonoBitmap) Hline_(x0 int, x1 int, y int, value bool) {
	if y < 0 || p.H <= y {
		return
	}
	start := min(p.W-1, max(0, x0))
	end := min(p.W-1, max(0, x1))

	for i := start; i <= end; i++ {
		p.SetPixNoCheck(i, y, value)
	}
}

func (p *MonoBitmap) Hline(x0 int, x1 int, y int, value bool) {
	if y < 0 || p.H <= y {
		return
	}
	start := min(p.W, max(0, x0))
	end := min(p.W, max(0, x1+1))

	i0 := start + p.W*y
	i1 := end + p.W*y

	index0 := i0 >> 5
	index1 := i1 >> 5

	shift0 := uint32(i0 % 32)
	shift1 := uint32(i1 % 32)

	bm0 := uint32(0xFFFFFFFF) << shift0
	bm1 := ^(uint32(0xFFFFFFFF) << shift1)

	/*fmt.Printf("index0=%v index1=%v\nbm0=%s\nbm1=%s",
	index0, index1,
	strconv.FormatInt(int64(bm0), 2), strconv.FormatInt(int64(bm1), 2))*/

	if index0 == index1 { //short case
		bm := bm0 & bm1
		if value {
			p.Pix[index0] |= bm
		} else {
			p.Pix[index0] &= bm ^ uint32(0xFFFFFFFF)
		}
		return
	}

	//Full sets
	for i := index0 + 1; i < index1; i++ {
		if value {
			p.Pix[i] = uint32(0xFFFFFFFF)
		} else {
			p.Pix[i] = 0
		}
	}
	//First and last bitmask
	if value {
		p.Pix[index0] |= bm0
		p.Pix[index1] |= bm1
	} else {
		p.Pix[index0] &= bm0 ^ uint32(0xFFFFFFFF)
		p.Pix[index1] &= bm1 ^ uint32(0xFFFFFFFF)
	}
}

func (p *MonoBitmap) Vline(x int, y0 int, y1 int, value bool) {
	for i := y0; i <= y1; i++ {
		p.SetPixNoCheck(x, i, value)
	}
}

// Modified from C++ source https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func (p *MonoBitmap) CircleFill(p0 image.Point, r int, value bool) {
	x := r
	y := 0
	err := 0

	x0 := p0.X
	y0 := p0.Y

	for x >= y {
		p.Hline(x0-x, x0+x, y0+y, value)
		p.Hline(x0-x, x0+x, y0-y, value)

		p.Hline(x0-y, x0+y, y0+x, value)
		p.Hline(x0-y, x0+y, y0-x, value)
		y += 1
		err += 1 + 2*y
		if 2*(err-x)+1 > 0 {
			x -= 1
			err += 1 - 2*x
		}
	}
}

// Modified from C++ source https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func (p *MonoBitmap) Circle(p0 image.Point, r int, value bool) {
	x := r
	y := 0
	err := 0

	x0 := p0.X
	y0 := p0.Y

	for x >= y {
		p.SetPix(x0+x, y0+y, value)
		p.SetPix(x0+y, y0+x, value)
		p.SetPix(x0-y, y0+x, value)
		p.SetPix(x0-x, y0+y, value)
		p.SetPix(x0-x, y0-y, value)
		p.SetPix(x0-y, y0-x, value)
		p.SetPix(x0+y, y0-x, value)
		p.SetPix(x0+x, y0-y, value)
		y += 1
		err += 1 + 2*y
		if 2*(err-x)+1 > 0 {
			x -= 1
			err += 1 - 2*x
		}
	}
}

// Gets pixel. Returns false if out of range
func (p *MonoBitmap) GetPix(x int, y int) bool {
	index := (x + p.W*y) / 32
	alabitit := uint32((x + p.W*y) % 32)
	//alabitit:=byte(x)&7
	bittimaski := uint32(1 << alabitit)
	if index < len(p.Pix) {
		return ((p.Pix[index] & bittimaski) > 0)
	}
	return false
}

func (p *MonoBitmap) GetPixNoCheck(x int, y int) bool {
	index := (x + p.W*y) / 32
	alabitit := uint32((x + p.W*y) % 32)
	bittimaski := uint32(1 << alabitit)

	return ((p.Pix[index] & bittimaski) > 0)

}

// TODO BUG: does not work if not div by 8
func (p *MonoBitmap) SetPix(x int, y int, value bool) {
	i := x + p.W*y
	index := i / 32
	bittimaski := uint32(1 << uint32(i%32))

	if (0 <= x) && (0 <= y) && (x < p.W) && (y < p.H) {
		if value {
			p.Pix[index] |= bittimaski
		} else {
			p.Pix[index] &= (bittimaski ^ uint32(0xFFFFFFFF))
		}
	}
}

func (p *MonoBitmap) SetPixNoCheck(x int, y int, value bool) {
	i := x + p.W*y
	index := i / 32
	bittimaski := uint32(1 << uint32(i%32))

	if value {
		p.Pix[index] |= bittimaski
	} else {
		p.Pix[index] &= (bittimaski ^ uint32(0xFFFFFFFF))
	}

}

// Draws source bitmap on bitmap
// drawTrue, draw when point value is true
// drawFalse,  draw when point value is true
func (p *MonoBitmap) DrawBitmap(source MonoBitmap, sourceArea image.Rectangle, targetCorner image.Point, drawTrue bool, drawFalse bool, invert bool) {
	if !drawTrue && !drawFalse {
		return //NOP operation
	}
	//TODO naive solution, make optimized later
	dx := sourceArea.Dx()
	dy := sourceArea.Dy()

	targetEnd := image.Point{X: min(p.W, targetCorner.X+dx), Y: min(p.H, targetCorner.Y+dy)}

	x0 := min(max(0, targetCorner.X), p.W-1)
	y0 := min(max(0, targetCorner.Y), p.H-1)

	if drawTrue && drawFalse {
		if invert {
			for x := x0; x < targetEnd.X; x++ {
				for y := y0; y < targetEnd.Y; y++ {
					v := source.GetPixNoCheck(x-x0+sourceArea.Min.X, y-y0+sourceArea.Min.Y)
					p.SetPixNoCheck(x, y, !v) //TODO copy byte by byte
				}
			}
		} else {
			for x := x0; x < targetEnd.X; x++ {
				for y := y0; y < targetEnd.Y; y++ {
					v := source.GetPixNoCheck(x-x0+sourceArea.Min.X, y-y0+sourceArea.Min.Y)
					p.SetPixNoCheck(x, y, v) //TODO copy byte by byte
				}
			}
		}
		return
	}
	//Now only one
	if drawTrue {
		if invert {
			for x := x0; x < targetEnd.X; x++ {
				for y := y0; y < targetEnd.Y; y++ {
					v := source.GetPixNoCheck(x-x0+sourceArea.Min.X, y-y0+sourceArea.Min.Y)
					if v {
						p.SetPixNoCheck(x, y, true) //TODO OPTIMIZE with special true and false operations? or does compiler that?
					}
				}
			}
		} else {
			for x := x0; x < targetEnd.X; x++ {
				for y := y0; y < targetEnd.Y; y++ {
					v := source.GetPixNoCheck(x-x0+sourceArea.Min.X, y-y0+sourceArea.Min.Y)
					if v {
						p.SetPixNoCheck(x, y, false)
					}
				}
			}
		}
		return
	}
	//Now draw on false
	if drawTrue {
		if invert {
			for x := x0; x < targetEnd.X; x++ {
				for y := y0; y < targetEnd.Y; y++ {
					v := source.GetPixNoCheck(x-x0+sourceArea.Min.X, y-y0+sourceArea.Min.Y)
					if !v {
						p.SetPixNoCheck(x, y, true)
					}
				}
			}
		} else {
			for x := x0; x < targetEnd.X; x++ {
				for y := y0; y < targetEnd.Y; y++ {
					v := source.GetPixNoCheck(x-x0+sourceArea.Min.X, y-y0+sourceArea.Min.Y)
					if !v {
						p.SetPixNoCheck(x, y, false)
					}
				}
			}
		}
		return
	}

	/*
		for x := x0; x < targetEnd.X; x++ {
			for y := targetCorner.Y; y < targetEnd.Y; y++ {
				v := source.GetPix(x-targetCorner.X+sourceArea.Min.X, y-targetCorner.Y+sourceArea.Min.Y)
				if (v) && (drawTrue) {
					p.SetPix(x, y, !invert) //TODO copy byte by byte
				}
				if (!v) && (drawFalse) {
					p.SetPix(x, y, invert)
				}
			}
		}
	*/
}

// Prints message on screen.Creates new lines on \n
// Returns rectangle where text was printed
func (p *MonoBitmap) Print(text string, font MonoFont, lineSpacing int, gap int, area image.Rectangle, drawTrue bool, drawFalse bool, invert bool, wrap bool) image.Rectangle {
	result := image.Rectangle{Min: area.Min, Max: area.Min}

	x := area.Min.X
	y := area.Min.Y
	//dim:=target.Bounds().Max
	for _, c := range text {
		if c == '\n' {
			x = area.Min.X
			y += lineSpacing
			if y > area.Max.Y {
				break
			}
			continue
		}
		f, ok := font[c]
		if !ok {
			f = font['?'] //Not found in font set
		}
		if wrap {
			if x+f.W > area.Max.X {
				x = area.Min.X
				y += lineSpacing
				if y > area.Max.Y {
					break
				}
			}
		}
		if (!wrap) || (x+f.W <= area.Max.X) {
			p.DrawBitmap(f, f.Bounds(), image.Point{X: x, Y: y}, drawTrue, drawFalse, invert)
			result.Max.X = max(result.Max.X, x+f.W)
			result.Max.Y = max(result.Max.Y, y+f.H)
			x += f.W + gap
		}

	}
	return result
}
