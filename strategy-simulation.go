package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	// "gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	piecewiseCount    = 10
	segmentLength     = 1.0
	initialGraphValue = 0.0
)

func main() {
	args := os.Args[1:]

	if len(args) != piecewiseCount {
		fmt.Println("Invalid number of arguments")
		os.Exit(1)
	}

	var (
		parabola = true
		y        = initialGraphValue
		xys      plotter.XYs
	)

  for i := 0; i < piecewiseCount; i++ {
    if parabola {
      a, _ := strconv.ParseFloat(args[i], 64)
      i++
      if i < len(args) { // Add this check
        b, _ := strconv.ParseFloat(args[i], 64)
        c := y - a*math.Pow(segmentLength, 2) - b*segmentLength

        for j := 0.0; j <= segmentLength; j += 0.01 {
          y = a*math.Pow(j, 2) + b*j + c
          xys = append(xys, plotter.XY{j, y})
        }
      }
    } else {
      if i < len(args) { // Add this check
        m, _ := strconv.ParseFloat(args[i], 64)
        b := y - m*segmentLength

        for j := 0.0; j <= segmentLength; j += 0.01 {
          y = m*j + b
          xys = append(xys, plotter.XY{j, y})
        }
      }
    }

    parabola = !parabola
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
