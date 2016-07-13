/*
Package for testing/demonstrating ui rendering capabilities
*/
package uiRender_test

import (
	"fmt"
    "testing"
    "github.com/hjkoskel/gomonochromebitmap"
    "github.com/hjkoskel/gomonochromebitmap/uiRender"
    //"image"
    "image/color"
    "image/png"
    "os"
)
    
func TestSimple(t *testing.T){
    fmt.Printf("--- Testing uiRender ---\n")
        
        colTrue:=color.RGBA{R:255,G:255,B:255,A:255}    
    colFalse:=color.RGBA{R:0,G:0,B:0,A:255}
    //test1:=gomonochromebitmap.NewMonoBitmap(128,64,false)
    //test1.Fill(image.Rect(40,20,60,40),true)
    
    testfont1:=gomonochromebitmap.GetFont_8x8()    
   
    textArr:=[]string{ "Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrott", "Golf", "Hotel", "India", "Juliet" ,"Kilo" ,"Lima" ,"Mike" ,"November" ,"Oscar" ,"Papa" ,"Quebec" ,"Romeo" ,"Sierra" ,"Tango"}    
    //GetStringBitmaps(arr []string,font map[rune]gomonochromebitmap.MonoBitmap,w int,h int,lineSpacing,gap int)    
    menu1:=uiRender.ScrollVerticalSelectMenu{
        Bitmaps:uiRender.GetStringBitmaps(textArr,testfont1,127,32,8,1),
        SelectedIndex:2,
        Scroll:0,
        InvertSelection:true,
        Arrow:nil,
        ScrollBar:true,
    }
    
    menu2:=uiRender.ScrollVerticalSelectMenu{
        Bitmaps:uiRender.GetStringBitmaps(textArr,testfont1,127,32,8,1),
        SelectedIndex:5,
        Scroll:0,
        InvertSelection:true,
        Arrow:nil,
        ScrollBar:true,
    }
    
    test1:=menu1.Render(128,64)
    out1,_ := os.Create("test1.png")
    png.Encode(out1,test1.GetImage(colTrue,colFalse))
    out1.Close()
    
    test2:=menu2.Render(128,64)
    out2,_ := os.Create("test2.png")
    png.Encode(out2,test2.GetImage(colTrue,colFalse))
    out2.Close()
}