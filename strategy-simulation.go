package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	//"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	graphResolution = 0.001
)

func integrand(x float64, q float64, w float64, e float64, r float64) float64 {
	return 1.0 / (q*math.Pow(x, 3) + w*math.Pow(x, 2) + e*x + r)
}
func simpson(a float64, b float64, q float64, w float64, e float64, r float64, n float64) float64 {
	h := (b - a) / n
	secondpart := 0.00
	for zz := 1.00; zz <= n/2; zz += 1.00 {
		secondpart += integrand((a + (2*zz-1)*h), q, w, e, r)
	}
	secondpart *= 2.00
	thirdpart := 0.00
	for zz := 1.00; zz <= (n/2)-1; zz += 1.00 {
		thirdpart += integrand((a + (2*zz)*h), q, w, e, r)
	}
	thirdpart *= 2.00
	return ((h / 3.00) * (integrand(a, q, w, e, r) + secondpart + thirdpart + integrand(b, q, w, e, r)))
}
func main() {
	args := os.Args[1:]

	// test command:
	// go run . 9 0 1 1 1 1 -10 1

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

	// initial velocity, initial acceleration, then accel curve params
	expectedArgCount := 2 + lineCount + parabolaCount*2

	if len(args) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		os.Exit(1)
	}

	var (
		initialVelocity, _ = strconv.ParseFloat(args[0], 64)
		y, _               = strconv.ParseFloat(args[1], 64)
		xys                plotter.XYs
	)

	fmt.Println("initialVelocity: ", initialVelocity)

	argIndex := 2
	xOffset := 0.0
	tiempo := 0.00
	velo := initialVelocity
	a := 0.00
	b := 0.00
	c := 0.00
	m := 0.00
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
		if isParabola == true {
			velo += (a/3)*(math.Pow(segmentLengths[i], 3)) + (b/2)*(math.Pow(segmentLengths[i], 2)) + c*(segmentLengths[i])
			tiempo += simpson(0.00, float64(segmentLengths[i]), a, b, c, velo, 50)
		}
		if isParabola == false {
			velo += (m/2)*(math.Pow(segmentLengths[i], 2)) + b*(segmentLengths[i])
			tiempo += simpson(0.00, float64(segmentLengths[i]), 0.00, m, b, velo, 50)
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
	fmt.Println("Total time:", tiempo)
}
