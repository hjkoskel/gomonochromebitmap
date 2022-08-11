package gomonochromebitmap

import (
	"math"
	"strings"
)

type AnsiColorString string

const (
	CLEARDISPLAY = "\033c"
)

const (
	ANSI_NO    = ""
	ANSI_RESET = "\033[0m"
)
const (
	FGANSI_BLACK  = "\033[30m"
	FGANSI_RED    = "\033[31m"
	FGANSI_GREEN  = "\033[32m"
	FGANSI_YELLOW = "\033[33m"
	FGANSI_BLUE   = "\033[34m"
	FGANSI_PURPLE = "\033[35m"
	FGANSI_CYAN   = "\033[36m"
	FGANSI_WHITE  = "\033[37m"

	FGANSI_BRIGHT_BLACK  = "\033[30;1m"
	FGANSI_BRIGHT_RED    = "\033[31;1m"
	FGANSI_BRIGHT_GREEN  = "\033[32;1m"
	FGANSI_BRIGHT_YELLOW = "\033[33;1m"
	FGANSI_BRIGHT_BLUE   = "\033[34;1m"
	FGANSI_BRIGHT_PURPLE = "\033[35;1m"
	FGANSI_BRIGHT_CYAN   = "\033[36;1m"
	FGANSI_BRIGHT_WHITE  = "\033[37;1m"
)

const (
	BGANSI_BLACK  = "\033[40m"
	BGANSI_RED    = "\033[41m"
	BGANSI_GREEN  = "\033[42m"
	BGANSI_YELLOW = "\033[43m"
	BGANSI_BLUE   = "\033[44m"
	BGANSI_PURPLE = "\033[45m"
	BGANSI_CYAN   = "\033[46m"
	BGANSI_WHITE  = "\033[47m"

	BGANSI_BRIGHT_BLACK  = "\033[40;1m"
	BGANSI_BRIGHT_RED    = "\033[41;1m"
	BGANSI_BRIGHT_GREEN  = "\033[42;1m"
	BGANSI_BRIGHT_YELLOW = "\033[43;1m"
	BGANSI_BRIGHT_BLUE   = "\033[44;1m"
	BGANSI_BRIGHT_PURPLE = "\033[45;1m"
	BGANSI_BRIGHT_CYAN   = "\033[46;1m"
	BGANSI_BRIGHT_WHITE  = "\033[47;1m"
)

type BlockGraphics struct {
	Clear       bool
	HaveBorder  bool
	BorderColor AnsiColorString
	TextColor   AnsiColorString
}

func (p *BlockGraphics) titleRow(xdim int) string {
	var sb strings.Builder
	if p.Clear {
		sb.WriteString(CLEARDISPLAY)
	}

	if p.HaveBorder {
		sb.WriteString(string(p.BorderColor))
		sb.WriteString("╔" + strings.Repeat("═", xdim) + "╗\n")
		if 0 < len(p.BorderColor) {
			sb.WriteString(ANSI_RESET)
		}
	}
	return sb.String()
}

func (p *BlockGraphics) bottomRow(xdim int) string {
	var sb strings.Builder
	if p.HaveBorder {
		sb.WriteString(string(p.BorderColor))
		sb.WriteRune('╚')
		for x := 0; x < xdim; x++ {
			sb.WriteRune('═')
		}
		sb.WriteString("╝")
		if 0 < len(p.BorderColor) {
			sb.WriteString(ANSI_RESET)
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func (p *BlockGraphics) edgeLeft() string {
	var sb strings.Builder
	if p.HaveBorder {
		sb.WriteString(string(p.BorderColor))
		sb.WriteString("║")
	}
	if 0 < len(p.BorderColor) {
		sb.WriteString(ANSI_RESET)
	}
	sb.WriteString(string(p.TextColor))
	return sb.String()
}

func (p *BlockGraphics) edgeRight() string {
	var sb strings.Builder
	if p.HaveBorder {
		sb.WriteString(string(p.BorderColor))
		sb.WriteString("║")
		sb.WriteString(ANSI_RESET)
	}
	if 0 < len(p.TextColor) {
		sb.WriteString(ANSI_RESET)
	}

	sb.WriteString("\n")
	return sb.String()
}

//ToFullBlockChars creates console printable version of image ' ','█'
func (p *BlockGraphics) ToFullBlockChars(bitmap *MonoBitmap) string {
	var sb strings.Builder

	sb.WriteString(p.titleRow(bitmap.W))
	for y := 0; y < bitmap.H; y++ {
		sb.WriteString(p.edgeLeft())
		for x := 0; x < bitmap.W; x++ {
			if bitmap.GetPix(x, y) {
				sb.WriteRune('█')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteString(p.edgeRight())
	}
	sb.WriteString(p.bottomRow(bitmap.W))
	return sb.String()
}

//ToHalfBlockChars ' ', '▀', '▄', '█'
func (p *BlockGraphics) ToHalfBlockChars(bitmap *MonoBitmap) string {
	var sb strings.Builder
	sb.WriteString(p.titleRow(bitmap.W))
	for y := 0; y < bitmap.H; y += 2 {
		sb.WriteString(p.edgeLeft())
		for x := 0; x < bitmap.W; x++ {
			i := 0
			if bitmap.GetPix(x, y) { //upper
				i++
			}
			if bitmap.GetPix(x, y+1) { //lower, GetPix returns false if goes over in y direction
				i += 2
			}
			sb.WriteRune([]rune{' ', '▀', '▄', '█'}[i])
		}
		sb.WriteString(p.edgeRight())
	}
	sb.WriteString(p.bottomRow(bitmap.W))
	return sb.String()
}

//ToQuadBlockChars ' ', '▘', '▝', '▀','▖', '▌', '▞', '▛','▗', '▚', '▐', '▜','▄', '▙', '▟', '█'
func (p *BlockGraphics) ToQuadBlockChars(bitmap *MonoBitmap) string {
	var sb strings.Builder
	w := int(math.Ceil(float64(bitmap.W) / 2))
	sb.WriteString(p.titleRow(w))

	for y := 0; y < bitmap.H; y += 2 {
		sb.WriteString(p.edgeLeft())
		for x := 0; x < bitmap.W; x += 2 {
			i := 0
			if bitmap.GetPix(x, y) { //Q0▘
				i++
			}
			if bitmap.GetPix(x+1, y) { //Q1▝
				i += 2
			}
			if bitmap.GetPix(x, y+1) { //Q2▖
				i += 4
			}
			if bitmap.GetPix(x+1, y+1) { //Q3▗
				i += 8
			}
			sb.WriteRune([]rune{
				' ', '▘', '▝', '▀',
				'▖', '▌', '▞', '▛',
				'▗', '▚', '▐', '▜',
				'▄', '▙', '▟', '█'}[i])
		}
		sb.WriteString(p.edgeRight())
	}

	sb.WriteString(p.bottomRow(w))
	return sb.String()
}
