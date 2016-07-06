/*
Testing
*/

package gomonochromebitmap_test

import (
	"fmt"
    "testing"
    "github.com/hjkoskel/gomonochromebitmap"
    "os"
    "image"
    "image/color"
    "image/png"
    "time"
)

    
func TestSimple(t *testing.T){
    fmt.Printf("Preparing test data...\n");
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig) //TODO TARVIIKO
    imgfile,err:= os.Open("./testdata/dog.png")
    defer imgfile.Close()
    if(err!=nil){
        fmt.Printf("File error %v\n",err)
        return
    }
    
    
    testfont1:=gomonochromebitmap.GetFont_8x8()
    tStart:=time.Now() //Actual rendering starts here
    
    
    test1:=gomonochromebitmap.NewMonoBitmap(300,600,false)
    test1.Fill(image.Rect(40,20,60,40),true)
    test1.Fill(image.Rect(50,30,80,60),false)
    test1.Invert(image.Rect(55,35,70,50))
    
    //test1.Invert(test1.Bounds())
    //testfont1:=gomonochromebitmap.GetFont_5x7()
    
    
    test1.Print("Ok Text works\nTesting letters ABCDEFGHIJKLMNOPQRSTUVXYZÄÖ abcdefghijklmnopqrstuvxyzäö !\"'*/ []{} (). ~ &$",testfont1,8,2,test1.Bounds(),true,true,false,true)
    //test1.Print("oooooo\noooooooooo",testfont1,8,2,test1.Bounds(),true,true,false,true)    
        
    //test1.Print("!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!!",testfont1,8,test1.Bounds(),true,true,false,true)
    //test1.Rotate90(3)
        
        
    test1.Line(image.Point{X:30,Y:40},image.Point{X:100,Y:200},true)
        
    for a:=0;a<100;a+=7{
        test1.Line(image.Point{X:120+a,Y:60},image.Point{X:220-a,Y:160},true)
    }
        
    test1.Circle(image.Point{X:140,Y:370},100,true)
    test1.CircleFill(image.Point{X:140,Y:370},90,true)
    
    //test1.FlipH()  
        
    test2:=test1.GetView(128,64,image.Point{X:38,Y:260},0,true)
        
    //Small image, chip8 example
    chip8pic:=gomonochromebitmap.NewMonoBitmap(64,32,false)
    chip8pic.SetPix(3,5,true)
    chip8pic.CircleFill(image.Point{X:64,Y:32},32,true)
    test3:=chip8pic.GetView(128,64,image.Point{X:0,Y:0},0,true)
    
    
    pngimg, _, _ := image.Decode(imgfile)

    test4:=gomonochromebitmap.NewMonoBitmapFromImage(pngimg,pngimg.Bounds(),130,false)    
    test4.Rotate90(1)    
        
        
    test5:=gomonochromebitmap.NewMonoBitmap(800,600,false)
    for i:=10;i<400;i+=4{
        test5.Circle(image.Point{X:400,Y:300},i,true)    
    }
    
    fmt.Printf("Actual rendering took %v sec (ok, it is slow, optimizations coming soon)\n",float64(time.Since(tStart))/float64(time.Second))
    
    fmt.Printf("Printing images\n")
    
    colTrue:=color.RGBA{R:255,G:255,B:255,A:255}    
    colFalse:=color.RGBA{R:0,G:0,B:0,A:255}
        
    out1,_ := os.Create("test1.png")
    png.Encode(out1,test1.GetImage(colTrue,colFalse))
    out1.Close()
        
    out2,_ := os.Create("test2.png")
    png.Encode(out2,test2.GetImage(colTrue,colFalse))
    out2.Close()

    out3,_ := os.Create("test3.png")
    png.Encode(out3,test3.GetImage(colTrue,colFalse))
    out3.Close()

    out4,_ := os.Create("test4.png")
    png.Encode(out4,test4.GetImage(colTrue,colFalse))
    out4.Close()
        
    out5,_ := os.Create("test5.png")
    png.Encode(out5,test5.GetImage(colTrue,colFalse))
    out5.Close()

}