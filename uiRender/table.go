/*
Table control

Like vertical scroll menu but have title and

*/

package uirender

import (
	"image"
	"strings"

	"github.com/hjkoskel/gomonochromebitmap"
)

type Allig byte

const (
	CENTER    Allig = 0
	WEST      Allig = 1
	EAST      Allig = 2
	NORTH     Allig = 3
	NORTHWEST Allig = 4
	NORTHEAST Allig = 5
	SOUTH     Allig = 6
	SOUTHWEST Allig = 7
	SOUTHEAST Allig = 8
)

func (p *Allig) Placement(in image.Rectangle, area image.Rectangle) image.Rectangle {
	xc := (area.Dx() - in.Dx()) / 2
	yc := (area.Dy() - in.Dy()) / 2

	switch *p {
	case CENTER:
		return image.Rect(xc, yc, area.Dx()-xc, area.Dy()-yc)
	case WEST:
		return image.Rect(0, yc, in.Dx(), yc+in.Dy())
	case EAST:
		return image.Rect(xc*2, yc, area.Dx(), yc+in.Dy())
	case NORTH:
		return image.Rect(xc, 0, area.Dx()-xc, in.Dy())
	case NORTHWEST:
		return image.Rect(0, 0, in.Dx(), in.Dy())
	case NORTHEAST:
		return image.Rect(xc*2, 0, area.Dx(), in.Dy())
	case SOUTH:
		return image.Rect(xc, yc*2, area.Dx()-xc, area.Dy())
	case SOUTHWEST:
		return image.Rect(0, yc*2, in.Dx(), area.Dy())
	case SOUTHEAST:
		return image.Rect(xc*2, yc, area.Dx(), yc+in.Dy())
	}
	return image.Rect(0, 0, 0, 0)
}

type TableColumnSetting struct {
	Title       string
	TitleAllign Allig
	Rows        []string
	RowsAllign  Allig

	MinWidth int
	MaxWidth int
}

func numberOfTextRows(txt string) int {
	return len(strings.Split(strings.TrimSpace(txt), "\n"))
}

func (p *TableColumnSetting) MaxStrLen() int {
	result := 0
	for _, r := range p.Rows {
		result = max(result, len(r))
	}
	return result
}

type TableSetting struct {
	TitleFont      gomonochromebitmap.MonoFont
	Font           gomonochromebitmap.MonoFont
	Columns        []TableColumnSetting
	Outline        int
	HLines         int //How many pixels
	VLines         int
	TitleSeparator int

	//Printin options
	TitleLineSpacing int
	TitleGap         int
	LineSpacing      int
	Gap              int
}

func (p *TableSetting) RowCount() int {
	result := 0
	for _, c := range p.Columns {
		result = max(result, len(c.Rows))
	}
	return result
}

/*
func (p *TableSetting) TitleHeight() int {
	titleTextHeight := 0
	for _, c := range p.Columns {
		a := p.TitleFont.AreaEstimated(c.Title, p.TitleLineSpacing, p.TitleGap)
		titleTextHeight = max(titleTextHeight, a.Dy())
	}
	return titleTextHeight
}
*/
/*
func (p *TableSetting) Render(target gomonochromebitmap.MonoBitmap) error {

	//calc column widths
	colw := make([]int, len(p.Columns))
	rowh := make([]int, p.RowCount())
	textH, textW := p.Font.GetWH()
	//_, titleFontHeight := p.TitleFont.GetWH()
	titleh := p.TitleHeight()

	totalw := p.Outline*2 + (len(p.Columns)-1)*p.VLines
	for i, c := range p.Columns {
		colw[i] = c.MaxStrLen() * textW
		totalw += colw[i]
	}
	totalh := p.Outline*2 + titleh //top and bottom
	for rownumber, _ := range rowh {
		h := 0
		for _, c := range p.Columns {
			if rownumber < len(c.Rows) {
				h = max(h, numberOfTextRows(c.Rows[rownumber])*textH)
			}
		}
		rowh[rownumber] = h
		totalh += h + p.HLines
	}
	//ok, dimensions are known. Lets draw
	//Clear bitmap?
	target.Fill(target.Bounds(), false)

	target.Fill(image.Rect(0, 0, 0, p.VLines), true) //left side
	target.Fill(image.Rect(totalw, 0, totalw, p.VLines), true)
	target.Fill(image.Rect(0, 0, target.W, p.Outline), true)                         //Top line
	target.Fill(image.Rect(0, p.Outline+titleh, target.W, p.Outline*2+titleh), true) //Under title line

}
*/

//--------------------------------

type TableColum struct {
	Title       gomonochromebitmap.MonoBitmap
	TitleAlling Allig

	Cells      []gomonochromebitmap.MonoBitmap
	CellAlling Allig
}

type Table struct {
	Columns       []TableColum
	PositionIndex int
}

func (p *Table) TitleHeight() int {
	result := 0
	for _, c := range p.Columns {
		result = max(result, c.Title.H)
	}
	return result
}

func (p *TableColum) Width() int {
	result := p.Title.W
	for _, cel := range p.Cells {
		result = max(result, cel.W)
	}
	return result
}

/*
func (p *Table) Height(rownumber int32) int32 {
	if rownumber < 0 {
		return 0
	}
	result := 0
	for _, c := range p.Columns {
		if rownumber < len(c.Cells) {
			result = c.Cells[rownumber].H
		}
	}
	return result
}*/

/*
func (p *Table) Render(height int, scrollPosition float64) gomonochromebitmap.MonoBitmap {
	r := gomonochromebitmap.NewMonoBitmap(p.Width(), height, false)

	r.Hline(0, r.W, 0, true)
	for _, c := range p.Columns {
		c.Title
	}

	return r
}
*/
