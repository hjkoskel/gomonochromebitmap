//Package provides functions for operating monochrome images
package gomonochromebitmap

import (
	"fmt"
    "image"
    "image/color"
    "math"
)

type MonoBitmap struct{
    Pix []byte
    W int
    H int
}


//Initializes empty bitmap
//fill is default value
func NewMonoBitmap(w int,h int,fill bool) MonoBitmap{
    result:=MonoBitmap{W:w,H:h,Pix:make([]byte,w*h/8)}
    if(fill){
        for index,_:=range result.Pix {
            result.Pix[index]=255
        }
    }
    return result
}

//Initializes bitmap from image
//Color conversion: if any Red,Green or Blue value is over threshold then pixel is true
//
func NewMonoBitmapFromImage(img image.Image,area image.Rectangle,threshold byte,invert byte) MonoBitmap{
    b:=img.Bounds()
    fmt.Printf("Paivitetaan kuvasta kokoa %#v\n",b)
    w:=b.Max.X
    h:=b.Max.Y
    result:=MonoBitmap{W:w,H:h,Pix:make([]byte,w*h/8)}
    for x:=0;x<=w;x++{
        for y:=0;y<h;y++{
            vr,vg,vb,_:=img.At(x,y).RGBA()
            //fmt.Printf("R=%v G=%v B=%v\n",vr,vg,vb)
            v:=byte((intMax(int(vr),intMax(int(vg),int(vb))))>>8)
            if(v>threshold){
                index:=x/8+w*y/8
                alabitit:=byte(x)&7
                bittimaski:=byte(1<<alabitit)
                result.Pix[index]|=bittimaski
            }
        }
    }
    return result
}


func(p *MonoBitmap)Bounds() image.Rectangle{
    return image.Rect(0,0,p.W,p.H)
}

//Creates image from bitmap
func(p *MonoBitmap)GetImage(trueColor color.Color,falseColor color.Color) image.Image{
    tulos:=image.NewRGBA(image.Rect(0,0,p.W,p.H));
    //colTrue:=color.RGBA{R:255,G:255,B:255,A:255}    
    //colFalse:=color.RGBA{R:0,G:0,B:0,A:255}
    
    //TODO
    //BUG: if pixel count is not div by 8
    
    /*
    for i:=0;i<len(p.Pix);i++{
        for j:=0;j<8;j++{
            x:=(i%p.W)
            y:=((i/p.W)&0x7)*8+j
            //fmt.Printf("TAKAS: %v->[%v,%v]\n",i,x,y);                
            if((p.Pix[i]&(1<<byte(j)))>0){
                tulos.Set(x,y,trueColor);
            }else{
                tulos.Set(x,y,falseColor);
            }
        }
    }*/
    
    //Naive solution
    for x:=0;x<p.W;x++{
        for y:=0;y<p.H;y++{
            if(p.GetPix(x,y)){
                tulos.Set(x,y,trueColor);
            }else{
                tulos.Set(x,y,falseColor);
            }
        }
    }
    
    
    return tulos
}

//Get view (size w,h) for display. Starting from corner p0. Result is centered. If p0 goes outside, function clamps view
//This is meant only for producing scrollable output picture for display. Better scaling functions elsewhere
//pxStep=0, autoscale, so bitmap will fit
//pxStep=1 is 1:1
//pxStep=2 is 2:1 (50% scale)
//pxStep=3 is 3:1 (25% scale)
//pxStep is limited to point where whole bitmap is visible
//Returns: image, actual cornerpoint and zoom used. Useful if UI includes
func(p *MonoBitmap)GetView(w int, h int,p0 image.Point,pxStep int,edges bool) MonoBitmap{
    result:=MonoBitmap{W:w,H:h,Pix:make([]byte,w*h/8)}
    
    maxStep:=math.Max(float64(w)/float64(p.W),float64(w)/float64(p.W))  //In decimal    

    corner:=image.Point{X:intMax(p0.X,0),Y:intMax(p0.Y,0)} //Limit point inside
        
        
    var step float64
    
    step=math.Min(float64(pxStep),math.Ceil(maxStep)) //Limits zooming out too much
    if(pxStep==0){//Autoscale
        corner=image.Point{X:0,Y:0}
        if (maxStep<=0.5){//Scale bigger
            //TODO: this is only reason why decimal step is now needed. Todo later integer step
        }else{
            step=math.Ceil(step)
        }
    }
    //Limit corner
    corner.X=intMin(corner.X,int(float64(p.W)-step*float64(w)))
    corner.Y=intMin(corner.Y,int(float64(p.H)-step*float64(h)))
    
    for x:=0;x<w;x++{
        for y:=0;y<h;y++{
            result.SetPix(x,y,p.GetPix(int(float64(x)*step)+corner.X,int(float64(y)*step)+corner.Y))
        }
    }
    return result
}

//Fills rectangle area from map. Used for clearing image                           
func(p *MonoBitmap)Fill(area image.Rectangle,fillValue bool){
    //Naive solution. TODO later faster solution
    fmt.Printf("DEBUG: Filling area %#v\n",area)
    for x:=area.Min.X;x<=area.Max.X;x++{
        for y:=area.Min.Y;y<=area.Max.Y;y++{
            p.SetPix(x,y,fillValue)
        }
    }
}

//Inverts pixel values
func(p *MonoBitmap)Invert(area image.Rectangle){
    //Naive solution. TODO later faster solution
    fmt.Printf("DEBUG: Inverting area %#v\n",area)
    for x:=area.Min.X;x<=area.Max.X;x++{
        for y:=area.Min.Y;y<=area.Max.Y;y++{
            p.SetPix(x,y,!p.GetPix(x,y))
        }
    }
}


func(p *MonoBitmap)Line(p0 image.Point, p1 image.Point){    
}

//func(p *MonoBitmap)Hline(


func(p *MonoBitmap)Circle(p0 image.Point, r int){
}

//Gets pixel. Returns false if out of range
func(p *MonoBitmap)GetPix(x int,y int) bool{
    index:=x/8+p.W*y/8
    alabitit:=byte(x)&7
    bittimaski:=byte(1<<alabitit)
    if index<len(p.Pix){    
        return ((p.Pix[index]&bittimaski)>0)
    }
    return false
}

func(p *MonoBitmap)SetPix(x int,y int,value bool){
    index:=x/8+p.W*y/8
    alabitit:=byte(x)&7
    bittimaski:=byte(1<<alabitit)
    
    
    if (0<=x)&&(0<=y)&&(x<p.W)&&(y<p.H){
        if(value){
            p.Pix[index]|=bittimaski
        }else{
            p.Pix[index]&=(bittimaski^byte(255))
        }
    }
}

//Draws source bitmap on bitmap
//drawTrue, draw when point value is true
//drawFalse,  draw when point value is true
func(p *MonoBitmap)DrawBitmap(source MonoBitmap,sourceArea image.Rectangle,targetCorner image.Point,drawTrue bool, drawFalse bool,invert bool){
    //TODO naive solution, make optimized later
    dx:=sourceArea.Dx()
    dy:=sourceArea.Dy()
    
    targetEnd:=image.Point{X:intMin(p.W,targetCorner.X+dx),Y:intMin(p.H,targetCorner.Y+dy)}
    
    for x:=targetCorner.X;x<=targetEnd.X;x++{
        for y:=targetCorner.Y;y<=targetEnd.Y;y++{
            v:=source.GetPix(x-targetCorner.X+sourceArea.Min.X,y-targetCorner.Y+sourceArea.Min.Y)
            if (v)&&(drawTrue){
                p.SetPix(x,y,!invert)
            }
            if (!v)&&(drawFalse){
                p.SetPix(x,y,invert)
            }
        }
    }
}

//Utils
func intMax(a int,b int) int{
    if(a>b){return a;}
    return b;
}
func intMin(a int,b int) int{
    if(a<b){return a;}
    return b;
}