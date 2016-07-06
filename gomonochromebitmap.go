//Package provides functions for operating monochrome images
package gomonochromebitmap

import (
    "image"
    "image/color"
    "math"
)

type MonoBitmap struct{
    Pix []uint32  //using byte vs uint16 vs uint32 vs uint64...  32bit shoud suit well for raspi1/2
    W int
    H int
}

//Initializes empty bitmap
//fill is default value
func NewMonoBitmap(w int,h int,fill bool) MonoBitmap{
    result:=MonoBitmap{W:w,H:h,Pix:make([]uint32,w*h/32+1)}
    if(fill){
        for index,_:=range result.Pix {
            result.Pix[index]=0xFFFFFFFF
        }
    }
    return result
}

//Initializes bitmap from image
//Color conversion: if any Red,Green or Blue value is over threshold then pixel is true
func NewMonoBitmapFromImage(img image.Image,area image.Rectangle,threshold byte,invert bool) MonoBitmap{
    b:=img.Bounds()
    w:=b.Max.X
    h:=b.Max.Y
    result:=NewMonoBitmap(w,h,false)
    for x:=0;x<=w;x++{
        for y:=0;y<h;y++{
            vr,vg,vb,_:=img.At(x,y).RGBA()
            v:=byte((intMax(int(vr),intMax(int(vg),int(vb))))>>8)
            if(v>threshold){
                result.SetPix(x,y,!invert);
            }else{
                result.SetPix(x,y,invert);
            }
        }
    }
    return result
}

func(p *MonoBitmap)Bounds() image.Rectangle{
    return image.Rect(0,0,p.W,p.H)
}

//Creates RGBA image from bitmap
func(p *MonoBitmap)GetImage(trueColor color.Color,falseColor color.Color) image.Image{
    result:=image.NewRGBA(image.Rect(0,0,p.W,p.H));
    for x:=0;x<p.W;x++{
        for y:=0;y<p.H;y++{
            if(p.GetPix(x,y)){
                result.Set(x,y,trueColor);
            }else{
                result.Set(x,y,falseColor);
            }
        }
    }
    return result
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
    result:=NewMonoBitmap(w,h,false)
    maxStep:=math.Max(float64(p.W)/float64(w),float64(p.H)/float64(h))  //In decimal
    corner:=image.Point{X:intMax(p0.X,0),Y:intMax(p0.Y,0)} //Limit point inside
            
    var step float64
    
    step=math.Min(float64(pxStep),math.Ceil(maxStep)) //Limits zooming out too much
    if(pxStep==0){//Autoscale
        step=maxStep
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
            a:=int(float64(x)*step)+corner.X
            b:=int(float64(y)*step)+corner.Y
            if((a<0)||(b<0)||(p.W<=a)||(p.H<=b)){
                result.SetPix(x,y,edges)
            }else{
                result.SetPix(x,y,p.GetPix(a,b))
            }
        }
    }
    return result
}

//Fills rectangle area from map. Used for clearing image                           
func(p *MonoBitmap)Fill(area image.Rectangle,fillValue bool){
    //Naive solution. TODO later faster solution
    for x:=area.Min.X;x<=area.Max.X;x++{
        for y:=area.Min.Y;y<=area.Max.Y;y++{
            p.SetPix(x,y,fillValue)
        }
    }
}

//Inverts pixel values
func(p *MonoBitmap)Invert(area image.Rectangle){
    //Naive solution. TODO later faster solution
    for x:=area.Min.X;x<=area.Max.X;x++{
        for y:=area.Min.Y;y<=area.Max.Y;y++{
            p.SetPix(x,y,!p.GetPix(x,y))
        }
    }
}

//Flip with axis in vertical
func(p *MonoBitmap)FlipV(){
    var v bool
    var i int
    for x:=0;x<p.W/2;x++{
        for y:=0;y<p.H;y++{
            v=p.GetPix(x,y)
            i=p.W-x-1
            p.SetPix(x,y,p.GetPix(i,y))
            p.SetPix(i,y,v)
        }
    }
}

func(p *MonoBitmap)FlipH(){
    var v bool
    var i int
    for x:=0;x<p.W;x++{
        for y:=0;y<p.H/2;y++{
            v=p.GetPix(x,y)
            i=p.H-y-1
            p.SetPix(x,y,p.GetPix(x,i))
            p.SetPix(x,i,v)
        }
    }

}


//Rotates in 90 decree steps
//+1=90 clockwise
//-1=90 anticlockwise
//+2=180 clockwise etc...
func(p *MonoBitmap)Rotate90(turn90 int){
    angle:=turn90%4
    result:=NewMonoBitmap(p.W,p.H,false)
    switch(angle){
        case 0:
            return //NOP
        case 1,-3:
            result.W=p.H
            result.H=p.W
            for x:=0;x<p.W;x++{
                for y:=0;y<p.H;y++{
                    result.SetPix(p.H-y-1,x,p.GetPix(x,y))
                }
            }
        case 2,-2:
            for x:=0;x<p.W;x++{
                for y:=0;y<p.H;y++{
                    result.SetPix(p.W-x-1,p.H-y-1,p.GetPix(x,y))
                }
            }
        case 3,-1:
            result.W=p.H
            result.H=p.W
            for x:=0;x<p.W;x++{
                for y:=0;y<p.H;y++{
                    result.SetPix(y,p.W-x-1,p.GetPix(x,y))
                }
            }
    }
    p.W=result.W
    p.H=result.H
    p.Pix=result.Pix
}


// Bresenham's line, copied from http://41j.com/blog/2012/09/bresenhams-line-drawing-algorithm-implemetations-in-go-and-c/
func(p *MonoBitmap)Line(p0 image.Point, p1 image.Point,value bool){
  var cx int32 = int32(p0.X);
  var cy int32 = int32(p0.Y);
 
  var dx int32 = int32(p1.X) - cx;
  var dy int32 = int32(p1.Y) - cy;
  if dx<0 { dx = 0-dx; }
  if dy<0 { dy = 0-dy; }
 
  var sx int32;
  var sy int32;
  if cx < int32(p1.X) { sx = 1; } else { sx = -1; }
  if cy < int32(p1.Y) { sy = 1; } else { sy = -1; }
  var err int32 = dx-dy;
 
  var n int32;
  for n=0;n<1000;n++ {
    p.SetPix(int(cx),int(cy),value);
    if((cx==int32(p1.X)) && (cy==int32(p1.Y))) {return;}
    var e2 int32 = 2*err;
    if e2 > (0-dy) { err = err - dy; cx = cx + sx; }
    if e2 < dx     { err = err + dx; cy = cy + sy; }
  }
    
}

//Horizontal line for filling
func(p *MonoBitmap)Hline(x0 int, x1 int, y int,value bool){
    for i:=x0;i<=x1;i++{
        p.SetPix(i,y,value)
    }
}

func(p *MonoBitmap)Vline(x int, y0 int, y1 int,value bool){
    for i:=y0;i<=y1;i++{
        p.SetPix(x,i,value)
    }
}

// Modified from C++ source https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func(p *MonoBitmap)CircleFill(p0 image.Point, r int,value bool){
    x:=r
    y:=0
    err:=0
        
    x0:=p0.X
    y0:=p0.Y

    for(x>=y){
        p.Hline(x0-x,x0+x,y0+y,value)
        p.Hline(x0-x,x0+x,y0-y,value)
        
        p.Hline(x0-y,x0+y,y0+x,value)
        p.Hline(x0-y,x0+y,y0-x,value)
        y += 1;
        err += 1 + 2*y;
        if (2*(err-x) + 1 > 0){
            x -= 1;
            err += 1 - 2*x;
        }
    }    
}

// Modified from C++ source https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func(p *MonoBitmap)Circle(p0 image.Point, r int,value bool){
    x:=r
    y:=0
    err:=0
        
    x0:=p0.X
    y0:=p0.Y

    for(x>=y){
        p.SetPix(x0 + x, y0 + y,value)
        p.SetPix(x0 + y, y0 + x,value)
        p.SetPix(x0 - y, y0 + x,value)
        p.SetPix(x0 - x, y0 + y,value)
        p.SetPix(x0 - x, y0 - y,value)
        p.SetPix(x0 - y, y0 - x,value)
        p.SetPix(x0 + y, y0 - x,value)
        p.SetPix(x0 + x, y0 - y,value)
        y += 1;
        err += 1 + 2*y;
        if (2*(err-x) + 1 > 0){
            x -= 1;
            err += 1 - 2*x;
        }
    }    
}

//Gets pixel. Returns false if out of range
func(p *MonoBitmap)GetPix(x int,y int) bool{
    index:=(x+p.W*y)/32
    alabitit:=uint32((x+p.W*y)%32)
    //alabitit:=byte(x)&7
    bittimaski:=uint32(1<<alabitit)
    if index<len(p.Pix){    
        return ((p.Pix[index]&bittimaski)>0)
    }
    return false
}

//TODO BUG: does not work if not div by 8
func(p *MonoBitmap)SetPix(x int,y int,value bool){
    index:=(x+p.W*y)/32
    //alabitit:=byte(x)&7
    alabitit:=uint32((x+p.W*y)%32)
    bittimaski:=uint32(1<<alabitit)
    
    if (0<=x)&&(0<=y)&&(x<p.W)&&(y<p.H){
        if(value){
            p.Pix[index]|=bittimaski
        }else{
            p.Pix[index]&=(bittimaski^uint32(0xFFFFFFFF))
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
    
    //fmt.Printf("Haluu piirtää bitmapin %#v ---> %v\n",targetCorner,targetEnd)
    
    for x:=targetCorner.X;x<targetEnd.X;x++{
        for y:=targetCorner.Y;y<targetEnd.Y;y++{
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

//Prints message on screen.Creates new lines on \n 
func(p *MonoBitmap)Print(text string,font map[rune]MonoBitmap,lineSpacing int,gap int,area image.Rectangle, drawTrue bool, drawFalse bool,invert bool,wrap bool){
    x:=area.Min.X
    y:=area.Min.Y
    //dim:=target.Bounds().Max
    for _, c := range text {
        if(c=='\n'){
            x=area.Min.X
            y+=lineSpacing
            if(y>area.Max.Y){
                break
            }        
        }else{
            f,ok:=font[c]
            if(!ok){
                f=font['?'] //Not found in font set
            }
            if(wrap){
                if(x+f.W>area.Max.X){
                    x=area.Min.X
                    y+=lineSpacing
                    if(y>area.Max.Y){
                        break
                    }
                }
            }
            if(!wrap)||(x+f.W<=area.Max.X){
                p.DrawBitmap(f,f.Bounds(),image.Point{X:x,Y:y},drawTrue,drawFalse,invert)
                x+=f.W+gap
            }
        }
    }
}

//Private Utils
func intMax(a int,b int) int{
    if(a>b){return a;}
    return b;
}
func intMin(a int,b int) int{
    if(a<b){return a;}
    return b;
}