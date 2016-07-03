/*
Testing
*/

package gomonochromebitmap_test

import (
	"fmt"
    "testing"
    "github.com/hjkoskel/gomonochromebitmap"
    "os"
    //"image"
    "image/color"
    "image/png"
)

func TestSimple(t *testing.T){
    fmt.Printf("AJETAAN TESTIÄ\n")
    
    test1:=gomonochromebitmap.NewMonoBitmap(80,60,false)
    /*test1.Fill(image.Rect(40,20,60,40),true)
    test1.Fill(image.Rect(50,30,80,60),false)
    test1.Invert(image.Rect(55,35,70,50))
    */
    test1.Invert(test1.Bounds())
        
    fmt.Printf("Printing images\n")
    
    colTrue:=color.RGBA{R:255,G:255,B:255,A:255}    
    colFalse:=color.RGBA{R:0,G:0,B:0,A:255}
        
    out1,_ := os.Create("test1.png")
    png.Encode(out1,test1.GetImage(colTrue,colFalse))
    out1.Close()
        
    fmt.Printf("Vakio on %v\n",gomonochromebitmap.Vakio)
}