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
	initialGraphValue = 0.0
	segmentLength     = 1.0
)

func main() {
	args := os.Args[1:]

	// Define whether each segment is a parabola (true) or a line (false)
	piecewiseFunctions := []bool{true, false, true, false, true}
	lineCount := 0
	parabolaCount := 0

	for _, isParabola := range piecewiseFunctions {
		if isParabola {
			parabolaCount++
		} else {
			lineCount++
		}
	}

	expectedArgCount := 1 + lineCount + parabolaCount*2 // 1 for initial value, lineCount for lines, parabolaCount*2 for parabolas

	if len(args) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		os.Exit(1)
	}

	var (
		y   = initialGraphValue
		xys plotter.XYs
	)

	argIndex := 0
	xOffset := 0.0
	for _, isParabola := range piecewiseFunctions {
		if isParabola {
			a, _ := strconv.ParseFloat(args[argIndex], 64)
			argIndex++
			b, _ := strconv.ParseFloat(args[argIndex], 64)
			argIndex++
			c := y - a*math.Pow(segmentLength, 2) - b*segmentLength

			for j := 0.0; j <= segmentLength; j += 0.01 {
				x := xOffset + j
				y = a*math.Pow(j, 2) + b*j + c
				xys = append(xys, plotter.XY{x, y})
			}
		} else {
			m, _ := strconv.ParseFloat(args[argIndex], 64)
			argIndex++
			b := y - m*segmentLength

			for j := 0.0; j <= segmentLength; j += 0.01 {
				x := xOffset + j
				y = m*j + b
				xys = append(xys, plotter.XY{x, y})
			}
		}
		xOffset += segmentLength
	}

	p := plot.New()

	lines, err := plotter.NewLine(xys)
	if err != nil {
		panic(err)
	}

	p.Add(lines)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "graph.png"); err != nil {
		panic(err)
	}
}
