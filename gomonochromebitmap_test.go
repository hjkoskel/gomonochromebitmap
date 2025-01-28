/*
Testing
*/

package gomonochromebitmap_test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/hjkoskel/gomonochromebitmap"
)

func BenchmarkGetBlankImage(b *testing.B) {
	img := gomonochromebitmap.NewMonoBitmap(800, 601, false)
	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	b.ResetTimer()
	for range b.N {
		img.GetImage(colTrue, colFalse)
	}
}

// BenchmarkHlines-16    	    7452	    157976 ns/op	       0 B/op	       0 allocs/op
// BenchmarkHlines-16    	  273817	      4361 ns/op	       0 B/op	       0 allocs/op
// BenchmarkHlines-16    	  353409	      3398 ns/op	       0 B/op	       0 allocs/op
func BenchmarkHlines(bench *testing.B) {
	b := gomonochromebitmap.NewMonoBitmap(640, 480, false)
	bench.ResetTimer()
	for range bench.N {
		for y := 0; y < b.H/2; y++ {
			b.Hline(0, y, y, true)
			b.Hline(b.W/2+y, b.W, b.H/2+y, true)
		}
	}
}

// BenchmarkVlines-16    	    3828	    308432 ns/op	       0 B/op	       0 allocs/op
func BenchmarkVlines(bench *testing.B) {
	b := gomonochromebitmap.NewMonoBitmap(640, 480, false)
	bench.ResetTimer()
	for range bench.N {
		for x := 0; x < b.W/2; x++ {
			b.Vline(x, 0, x, true)
			b.Vline(x+b.W/2, x, b.W, true)
		}
	}
}

func isVline(bm gomonochromebitmap.MonoBitmap, x int, y0 int, y1 int, value bool) error {
	for i := y0; i <= y1; i++ {
		if bm.GetPix(x, i) != value && i < bm.W {
			return fmt.Errorf("v line is not draw enough y0=%v y1=%v i=%v", y0, y1, i)
		}
	}
	for i := 0; i < y0; i++ {
		if bm.GetPix(x, i) == value {
			return fmt.Errorf("v line draw over before start=%v i=%v", y0, i)
		}
	}

	for i := y1 + 1; i < bm.W; i++ {
		if bm.GetPix(x, i) == value {
			return fmt.Errorf("v line draw over y1=%v after i=%v", y1, i)
		}
	}
	return nil
}

func TestVlines(t *testing.T) {
	b := gomonochromebitmap.NewMonoBitmap(40, 40, false)
	for x := 0; x < b.W/2; x++ {
		b.Vline(x, 0, x, true)
		e1 := isVline(b, x, 0, x, true)
		if e1 != nil {
			t.Error(e1)
		}

		y0 := b.H / 2
		y1 := b.H/2 + x
		b.Vline(x+b.W/2, y0, y1, true)
		e2 := isVline(b, x+b.W/2, y0, y1, true)
		if e2 != nil {
			t.Error(e2)
		}
	}
	out1, _ := os.Create("testVline.png")
	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	png.Encode(out1, b.GetImage(colTrue, colFalse))
	out1.Close()
}

// Is hline what wanted and other pixels different
func isHline(bm gomonochromebitmap.MonoBitmap, x0 int, x1 int, y int, value bool) error {
	for i := x0; i <= x1; i++ {
		if bm.GetPix(i, y) != value && i < bm.W {
			return fmt.Errorf("h line is not draw enough x0=%v x1=%v i=%v", x0, x1, i)
		}
	}
	for i := 0; i < x0; i++ {
		if bm.GetPix(i, y) == value {
			return fmt.Errorf("h line draw over before start=%v i=%v", x0, i)
		}
	}

	for i := x1 + 1; i < bm.W; i++ {
		if bm.GetPix(i, y) == value {
			return fmt.Errorf("h line draw over x1=%v after i=%v", x1, i)
		}
	}
	return nil
}

func TestHlines(t *testing.T) {
	b := gomonochromebitmap.NewMonoBitmap(40, 40, false)
	for y := 0; y < b.H/2; y++ {
		b.Hline(0, y, y, true)
		e1 := isHline(b, 0, y, y, true)
		if e1 != nil {
			t.Error(e1)
		}
		b.Hline(b.W/2+y, b.W, b.H/2+y, true)
		e2 := isHline(b, b.W/2+y, b.W, b.H/2+y, true)
		if e2 != nil {
			t.Error(e2)
		}
	}
	out1, _ := os.Create("testHline.png")
	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	png.Encode(out1, b.GetImage(colTrue, colFalse))
	out1.Close()
}

/*
func TestHlines2(t *testing.T) {
	b := gomonochromebitmap.NewMonoBitmap(128, 64, true)
	//53->59  y=37
	b.Hline(53, 59, 37, false)

	out, _ := os.Create("testHline2.png")
	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	png.Encode(out, b.GetImage(colTrue, colFalse))
	out.Close()
	//t.FailNow()

}*/

func TestSimple(t *testing.T) {
	fmt.Printf("Preparing test data...\n")
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig) //TODO TARVIIKO
	imgfile, err := os.Open("./testdata/dog.png")
	//defer imgfile.Close()
	if err != nil {
		fmt.Printf("File error %v\n", err)
		return
	}

	testfont1 := gomonochromebitmap.GetFont_8x8()
	tStart := time.Now() //Actual rendering starts here

	test1 := gomonochromebitmap.NewMonoBitmap(300, 600, false)
	test1.Fill(image.Rect(40, 20, 60, 40), true)
	test1.Fill(image.Rect(50, 30, 80, 60), false)

	//test1.Invert(test1.Bounds())
	//testfont1:=gomonochromebitmap.GetFont_5x7()

	test1.Print("Ok Text works\nTesting letters ABCDEFGHIJKLMNOPQRSTUVXYZÄÖ abcdefghijklmnopqrstuvxyzäö !\"'*/ []{} (). ~ &$", testfont1, 8, 2, test1.Bounds(), true, true, false, true)
	//test1.Print("oooooo\noooooooooo",testfont1,8,2,test1.Bounds(),true,true,false,true)

	//test1.Print("!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!!",testfont1,8,test1.Bounds(),true,true,false,true)
	//test1.Rotate90(3)

	test1.Line(image.Point{X: -10, Y: 500}, image.Point{X: 200, Y: 800}, true) //Test over edge
	test1.Line(image.Point{X: 250, Y: 700}, image.Point{X: 400, Y: 50}, true)  //Test over edge

	test1.Line(image.Point{X: 30, Y: 40}, image.Point{X: 100, Y: 200}, true)

	for a := 0; a < 100; a += 7 {
		test1.Line(image.Point{X: 120 + a, Y: 60}, image.Point{X: 220 - a, Y: 160}, true)
	}

	test1.Circle(image.Point{X: 140, Y: 370}, 100, true)
	test1.CircleFill(image.Point{X: 140, Y: 370}, 90, true)

	//test1.FlipH()
	test1.Invert(image.Rect(55, 35, 270, 350))
	test2 := test1.GetView(128, 64, image.Point{X: 38, Y: 260}, 0, true)

	//Small image, chip8 example
	chip8pic := gomonochromebitmap.NewMonoBitmap(64, 32, false)
	chip8pic.SetPix(3, 5, true)
	chip8pic.CircleFill(image.Point{X: 64, Y: 32}, 32, true)
	test3 := chip8pic.GetView(128, 64, image.Point{X: 0, Y: 0}, 0, true)

	pngimg, _, _ := image.Decode(imgfile)

	bou := pngimg.Bounds()
	test4 := gomonochromebitmap.NewMonoBitmapFromImage(pngimg, bou, 130, false)
	test4.Rotate90(1)

	test5 := gomonochromebitmap.NewMonoBitmap(800, 600, false)
	for i := 10; i < 400; i += 4 {
		test5.Circle(image.Point{X: 400, Y: 300}, i, true)
	}

	fmt.Printf("Actual rendering took %v sec (ok, it is slow, optimizations coming soon)\n", float64(time.Since(tStart))/float64(time.Second))

	fmt.Printf("Printing images\n")

	colTrue := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colFalse := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	out1, _ := os.Create("test1.png")
	png.Encode(out1, test1.GetImage(colTrue, colFalse))
	out1.Close()

	out2, _ := os.Create("test2.png")
	png.Encode(out2, test2.GetImage(colTrue, colFalse))
	out2.Close()

	out3, _ := os.Create("test3.png")
	png.Encode(out3, test3.GetImage(colTrue, colFalse))
	out3.Close()

	out4, _ := os.Create("test4.png")
	png.Encode(out4, test4.GetImage(colTrue, colFalse))
	out4.Close()

	out5, _ := os.Create("test5.png")
	png.Encode(out5, test5.GetImage(colTrue, colFalse))
	out5.Close()

}

/*
func TestPlanarColors(t *testing.T) {
	bitmaps := []gomonochromebitmap.MonoBitmap{}
	for i := 0; i < 3; i++ {
		newBm := gomonochromebitmap.NewMonoBitmap(6, 4, false)
		for x:=
		newBm.SetPix(i, 1+i%2, true)
		newBm.SetPix(i, 3, true)
		bitmaps = append(bitmaps, newBm)
	}

	planar, errPlanar := gomonochromebitmap.CreatePlanarColorImage(bitmaps, []color.Color{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 255, A: 255},
		color.RGBA{R: 0, G: 255, B: 0, A: 255},
		color.RGBA{R: 0, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		color.RGBA{R: 255, G: 0, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 0, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	})

	if errPlanar != nil {
		t.Error(errPlanar)
		return
	}

	out, errCreateOut := os.Create("testplanar.png")
	if errCreateOut != nil {
		t.Error(errCreateOut)
		return
	}
	errEncode := png.Encode(out, planar)
	if errEncode != nil {
		t.Error(errEncode)
	}
	out.Close()
}
*/
