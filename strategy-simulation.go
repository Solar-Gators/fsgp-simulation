package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	graphResolution = 0.001
)

func main() {
	args := os.Args[1:]

	// Define whether each segment is a parabola (true) or a line (false)
	// piecewiseFunctions := []bool{true, false, true, false}
	piecewiseFunctions := []bool{true, false, true, false}
  segmentLengths := []float64{2, 1, 2, 1}

	lineCount := 0
	parabolaCount := 0

	for _, isParabola := range piecewiseFunctions {
		if isParabola {
			parabolaCount++
		} else {
			lineCount++
		}
	}

	expectedArgCount := 1 + lineCount + parabolaCount*2

	if len(args) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		os.Exit(1)
	}

	var (
		y, _ = strconv.ParseFloat(args[0], 64)
		xys  plotter.XYs
	)

	argIndex := 1
	xOffset := 0.0
	for i, isParabola := range piecewiseFunctions {
		if isParabola {
			a, _ := strconv.ParseFloat(args[argIndex], 64)
			argIndex++
			b, _ := strconv.ParseFloat(args[argIndex], 64)
			argIndex++
			c := y - a*math.Pow(xOffset, 2) - b*xOffset

			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				y = a*math.Pow(x, 2) + b*x + c
				xys = append(xys, plotter.XY{X: x, Y: y})
			}
		} else {
			m, _ := strconv.ParseFloat(args[argIndex], 64)
			argIndex++
		b := y - m*xOffset

			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				y = m*x + b
				xys = append(xys, plotter.XY{X: x, Y: y})
			}
		}
		xOffset += segmentLengths[i]
	}

	p := plot.New()

	lines, err := plotter.NewLine(xys)
	if err != nil {
		panic(err)
	}

	p.Add(lines)

	p.X.Min = 0
	p.Y.Min = 0

	p.X.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var ticks []plot.Tick
		for i := math.Ceil(min); i <= max+1; i++ {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		}
		return ticks
	})

	p.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var ticks []plot.Tick
		for i := math.Ceil(min); i <= max+1; i++ {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		}
		return ticks
	})

	p.Add(plotter.NewGrid())

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "graph.png"); err != nil {
		panic(err)
	}
}
