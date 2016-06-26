package main

import (
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	// DefaultLineStyle is the default style for vertical and horizontal lines.
	DefaultLineStyle = draw.LineStyle{
		Color: color.Black,
		Width: vg.Points(0.5),
	}
)

// Line implements the plot.Plotter interface, drawing
// a line at the specified value (vertial or horizontal).
type Line struct {
	// Vertical is the style of the vertical lines.
	Style draw.LineStyle

	Value float64

	horizontal bool
}

// NewVLine returns a new grid with both vertical and
// horizontal lines using the default grid line style.
func NewVLine(v float64) *Line {
	return &Line{
		Style:      DefaultLineStyle,
		horizontal: false,
		Value:      v,
	}
}

// NewHLine returns a new grid with both vertical and
// horizontal lines using the default grid line style.
func NewHLine(v float64) *Line {
	return &Line{
		Style:      DefaultLineStyle,
		horizontal: true,
		Value:      v,
	}
}

// Plot implements the plot.Plotter interface.
func (g *Line) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	if g.horizontal {
		x := trX(g.Value)
		c.StrokeLine2(g.Style, x, c.Min.Y, x, c.Max.Y)
	} else {
		y := trY(g.Value)
		c.StrokeLine2(g.Style, c.Min.X, y, c.Max.X, y)
	}
}
