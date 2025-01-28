// Package providing clipping on line
package gomonochromebitmap

import (
	"image"
)

// Define region codes for Cohen-Sutherland algorithm
const (
	INSIDE = 0 // 0000
	LEFT   = 1 // 0001
	RIGHT  = 2 // 0010
	BOTTOM = 4 // 0100
	TOP    = 8 // 1000
)

// computeRegionCode calculates the region code for a point relative to the rectangle
func computeRegionCode(p image.Point, area image.Rectangle) int {
	code := INSIDE

	if p.X < area.Min.X {
		code |= LEFT
	} else if p.X > area.Max.X {
		code |= RIGHT
	}

	if p.Y < area.Min.Y {
		code |= BOTTOM
	} else if p.Y > area.Max.Y {
		code |= TOP
	}

	return code
}

// ClipLine clips a line segment from a to b inside the rectangle area
func ClipLine(a, b image.Point, area image.Rectangle) (*image.Point, *image.Point) {
	codeA := computeRegionCode(a, area)
	codeB := computeRegionCode(b, area)

	for {
		// Both points are inside the rectangle
		if codeA == INSIDE && codeB == INSIDE {
			return &a, &b
		}

		// Both points are outside the rectangle and in the same region
		if (codeA & codeB) != 0 {
			return nil, nil
		}

		// At least one point is outside the rectangle, pick it
		var codeOut int
		var x, y float64
		if codeA != INSIDE {
			codeOut = codeA
		} else {
			codeOut = codeB
		}

		// Find intersection point
		if (codeOut & TOP) != 0 {
			x = float64(a.X) + (float64(b.X)-float64(a.X))*(float64(area.Max.Y)-float64(a.Y))/(float64(b.Y)-float64(a.Y))
			y = float64(area.Max.Y)
		} else if (codeOut & BOTTOM) != 0 {
			x = float64(a.X) + (float64(b.X)-float64(a.X))*(float64(area.Min.Y)-float64(a.Y))/(float64(b.Y)-float64(a.Y))
			y = float64(area.Min.Y)
		} else if (codeOut & RIGHT) != 0 {
			y = float64(a.Y) + (float64(b.Y)-float64(a.Y))*(float64(area.Max.X)-float64(a.X))/(float64(b.X)-float64(a.X))
			x = float64(area.Max.X)
		} else if (codeOut & LEFT) != 0 {
			y = float64(a.Y) + (float64(b.Y)-float64(a.Y))*(float64(area.Min.X)-float64(a.X))/(float64(b.X)-float64(a.X))
			x = float64(area.Min.X)
		}

		// Update the point outside the rectangle to the intersection point
		if codeOut == codeA {
			a.X = int(x)
			a.Y = int(y)
			codeA = computeRegionCode(a, area)
		} else {
			b.X = int(x)
			b.Y = int(y)
			codeB = computeRegionCode(b, area)
		}
	}
}
