/*
This sub library is for simulating gadgets on desktop
Supports keys and monochrome display.

Rough but can help when developing ui and button layout for gadget.

*/
package gadgetSimUi

import (
	"image"
	"image/color"
	"strconv"

	"github.com/hjkoskel/gomonochromebitmap"
	"github.com/nfnt/resize"
	"github.com/veandco/go-sdl2/sdl"
)

type XyIntPair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MonochromeDisplay struct {
	Bitmap      gomonochromebitmap.MonoBitmap
	ID          string     `json:"id"`
	Corner      XyIntPair  `json:"corner"` //Coordinate on bitmap
	PixelSize   XyIntPair  `json:"pixelSize"`
	PixelGap    XyIntPair  `json:"pixelGap"`
	UpperRows   int        `json:"upperRows"`   //Allows to simulate two color OLED displays with 8 rows up or down with yellow and others are blue
	OnColor     color.RGBA `json:"onColor"`     //Hex format with transparent, upper rows
	OnColorDown color.RGBA `json:"onColorDown"` //Hex format with transparent
	OffColor    color.RGBA `json:"offColor"`    //Hex format with transparent
}

func parseStringToColor(s string) (color.RGBA, error) {
	value, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: byte(value & 0xFF), G: byte((value >> 8) & 0xFF), B: byte((value >> 16) & 0xFF), A: byte((value >> 24) & 0xFF)}, nil
}

type ButtonSettings struct {
	ID         string    `json:"id"` //Can have duplicates like control and control :D
	Corner     XyIntPair `json:"corner"`
	Dimensions XyIntPair `json:"dimensions,omitempty"`

	DebugEdges bool //should print debug edges for checking placement
	DebugColor color.RGBA
}

func (p *ButtonSettings) Hits(x int, y int) bool {
	return ((p.Corner.X <= x) && (x <= (p.Corner.X + p.Dimensions.X))) && ((p.Corner.Y <= y) && (y <= (p.Corner.Y + p.Dimensions.Y)))
}

type GadgetWindow struct {
	Title string
	//ImageFile    string
	MonoDisplays []MonochromeDisplay
	Buttons      []ButtonSettings

	BgImage image.Image //image.RGBA

	FromKeys  chan KeyboardStatus
	ToDisplay chan DisplayUpdate

	//Internal
	usedWindowArea sdl.Rect //simplifies some coding and button mapping etc...
	window         *sdl.Window
}

func (p *GadgetWindow) ScaleWinOnPicCoord(x int, y int) (int, int) {
	b := p.BgImage.Bounds()
	a := p.usedWindowArea
	//fmt.Printf("Used window corner=(%v,%v) dim=(%v,%v)\n", a.X, a.Y, a.W, a.H)
	return (int(b.Dx()) * (x - int(a.X))) / int(a.W), (int(b.Dy()) * (y - int(a.Y))) / int(a.H)
}

func imageToSdlSurf(i image.Image) (*sdl.Surface, error) {
	bou := i.Bounds()
	rgba := image.NewRGBA(bou)

	s, err := sdl.CreateRGBSurface(0, int32(bou.Max.X), int32(bou.Max.Y), 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
	if err != nil {
		return s, err
	}
	rgba.Pix = s.Pixels()

	for x := 0; x < bou.Max.X; x++ {
		for y := 0; y < bou.Max.Y; y++ {
			rgba.Set(x, y, i.At(x, y))
		}
	}
	return s, err
}

/*
Assuming that background image stays same and LCD/OLED display is static
*/
func (p *GadgetWindow) Initialize(backgroundPicture image.Image) error {

	p.FromKeys = make(chan KeyboardStatus, 10)
	p.ToDisplay = make(chan DisplayUpdate, 1)

	p.BgImage = backgroundPicture
	err := sdl.Init(sdl.INIT_EVERYTHING)

	if err != nil {
		return err
	}
	maxB := backgroundPicture.Bounds().Max

	/*
		TODO android support?
	*/
	p.window, err = sdl.CreateWindow(p.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(maxB.X), int32(maxB.Y), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}
	return nil
}

func (p *GadgetWindow) Quit() {
	p.window.Destroy()
	sdl.Quit()
}

func (p *GadgetWindow) Run() error {
	running := true
	p.window.UpdateSurface()
	for running {

		if 0 < len(p.ToDisplay) {
			newBitmap := <-p.ToDisplay
			for i := range p.MonoDisplays {
				if p.MonoDisplays[i].ID == newBitmap.ID {
					p.MonoDisplays[i].Bitmap = newBitmap.Bitmap
				}
			}
			p.Render()
			p.window.UpdateSurface()
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.MouseButtonEvent:
				//fmt.Printf("TODO MOUSE BUTTON %#v\n", event)
				//fmt.Printf("MOUSE=%#v\n", t)
				if 0 < t.State {
					//fmt.Printf("Nappi alas x:%v y:%v\n", t.X, t.Y)
					//Let's do scaling
					xp, yp := p.ScaleWinOnPicCoord(int(t.X), int(t.Y))
					hits := p.HitsToId(xp, yp)
					//hits := p.HitsToId(int(t.X), int(t.Y))
					if 0 < len(hits) {
						//fmt.Printf("HITTED to %v\n", hits)
						p.FromKeys <- KeyboardStatus{KeysDown: []string{hits}}
					}
				} else {
					//fmt.Printf("Nappi ylös x:%v y:%v\n", t.X, t.Y)
					p.FromKeys <- KeyboardStatus{KeysDown: []string{}} //Mark as cleared
				}
			case *sdl.WindowEvent:

				if t.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					//fmt.Printf("TODO WINDOW EVENT %#v\n", event)
					//fmt.Printf("Windowevent size changed\n")
					p.Render()
					p.window.UpdateSurface()
				}

			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
	return nil
}

/*
Render if content changes or window resizes
*/
func (p *GadgetWindow) Render() error {
	winW, winH := p.window.GetSize()
	//golang image have nice resize functions.
	b := p.BgImage.Bounds().Max

	//by ratio
	picW := int(0)
	picH := int(0)
	//window is skinnier than background  -> scale width to same
	if (float32(winW) / float32(winH)) < (float32(b.X) / float32(b.Y)) {
		picW = int(winW)
		picH = (b.Y * int(winW)) / b.X
	} else { //Height is constraint
		picW = (b.X * int(winH)) / b.Y
		picH = int(winH)
	}

	surface, err := p.window.GetSurface()
	if err != nil {
		return err
	}

	//Clear background
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}, 0xFF000000)

	//Draw background, resize then draw  TODO CACHE AND GAIN PERFORMANCE
	bgSurf, err := imageToSdlSurf(resize.Resize(uint(picW), uint(picH), p.BgImage, resize.Bilinear))
	if err != nil {
		return err
	}

	//Centered picture, fixed from Y=0
	Xoff := (int32(winW) - int32(bgSurf.W)) / 2
	//fmt.Printf("Xoff=%v\n", Xoff)

	p.usedWindowArea = sdl.Rect{X: Xoff, Y: 0, W: bgSurf.W, H: bgSurf.H}
	for _, but := range p.Buttons {
		if but.DebugEdges {
			r := sdl.Rect{
				X: int32(but.Corner.X*picW) / int32(b.X),
				Y: int32(but.Corner.Y*picH) / int32(b.Y),
				W: int32(but.Dimensions.X*picW) / int32(b.X),
				H: int32(but.Dimensions.Y*picH) / int32(b.Y),
			}
			c := but.DebugColor
			bgSurf.FillRect(&r, uint32(c.A)<<24|uint32(c.R)<<16|uint32(c.G)<<8|uint32(c.B))
		}
	}

	bgSurf.Blit(&sdl.Rect{X: 0, Y: 0, W: bgSurf.W, H: bgSurf.H}, surface, &p.usedWindowArea)
	//Draw all displays, scaled and placed
	scaleFactor := float32(picW) / float32(b.X)
	for _, dis := range p.MonoDisplays {
		//scale that image
		disImg := dis.Bitmap.GetDisplayImage(dis.OnColor, dis.OnColorDown, dis.UpperRows, dis.OffColor, dis.PixelSize.X, dis.PixelSize.Y, dis.PixelGap.X, dis.PixelGap.Y)
		dboundMax := disImg.Bounds().Max
		scaledDisImg := resize.Resize(uint(scaleFactor*float32(dboundMax.X)), uint(scaleFactor*float32(dboundMax.Y)), disImg, resize.Bilinear)
		displaySurf, err := imageToSdlSurf(scaledDisImg)
		if err != nil {
			return err
		}

		err = displaySurf.Blit(
			&sdl.Rect{X: 0, Y: 0, W: displaySurf.W, H: displaySurf.H}, surface,
			&sdl.Rect{X: int32(float32(dis.Corner.X)*scaleFactor) + Xoff, Y: int32(float32(dis.Corner.Y) * scaleFactor), W: displaySurf.W, H: displaySurf.H})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *GadgetWindow) HitsToId(x int, y int) string {
	for _, b := range p.Buttons {
		if b.Hits(x, y) {
			return b.ID
		}
	}
	return ""
}

/*
Message channels
* For getting keyboard events
* Updating displays
*/

type KeyboardStatus struct {
	KeysDown []string //array if multitouch support
}

type DisplayUpdate struct {
	ID     string
	Bitmap gomonochromebitmap.MonoBitmap
}
