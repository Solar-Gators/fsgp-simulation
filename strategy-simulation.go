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

/*
  so essentially, if, on a given segment of track,
  our acceleration (i.e., x”) wrt x = some quadratic function of x (i.e. of the form ax^2 + bx + c)
  then x' = dx/dt = (a/3)x^3 + (b/2)x^2 + cx + initial_velocity
  thus, we derive dt = dx/[(a/3)x^3 + (b/2)x^2 + cx + initial_velocity]
  to find the time elapsed on the segment, we integrate the right side.
  this is a pretty complicated integral, so i've elected to use an approximation called simpson's method
*/

// for the following two functions, i use q,w,e,r as arbitrary coefficients. "n" denotes the increments of the approximation
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

	// run command to test:
	// go run . 9 0 1 1 1 1 -10 1

	rawArgs := os.Args[1:]

	args := make([]float64, len(rawArgs))
	for i, arg := range rawArgs {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		args[i] = val
	}

	// Define whether each segment is a straightaway (true) or a turn (false)
	piecewiseFunctions := []bool{true, false, true, false}
	segmentLengths := []float64{2, 1, 2, 1}

	turnCount := 0
	straightCount := 0

	for _, segmentType := range piecewiseFunctions {
		if segmentType {
			straightCount++
		} else {
			turnCount++
		}
	}

	// initial velocity, initial acceleration, then accel curve params
	expectedArgCount := 2 + turnCount + straightCount*2

	if len(rawArgs) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		os.Exit(1)
	}

	initialVelocity := args[0]
	currentTickAccel := args[1]
	var accelPlot plotter.XYs

	fmt.Println("initialVelocity: ", initialVelocity)

	argIndex := 2
	xOffset := 0.0
	tiempo := 0.00
	velo := initialVelocity
	a := 0.00
	b := 0.00
	c := 0.00
	m := 0.00

	for i, segmentIsStraight := range piecewiseFunctions {
		if segmentIsStraight {
			a := args[argIndex]
			argIndex++
			b := args[argIndex]
			argIndex++
			c := currentTickAccel - a*math.Pow(xOffset, 2) - b*xOffset

			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				currentTickAccel = a*math.Pow(x, 2) + b*x + c
				accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
			}
		} else {
			m := args[argIndex]
			argIndex++
			b := currentTickAccel - m*xOffset

			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				currentTickAccel = m*x + b
				accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
			}
		}

		if segmentIsStraight == true {
			velo += (a/3)*(math.Pow(segmentLengths[i], 3)) + (b/2)*(math.Pow(segmentLengths[i], 2)) + c*(segmentLengths[i])
			tiempo += simpson(0.00, float64(segmentLengths[i]), a, b, c, velo, 50)
		}
		if segmentIsStraight == false {
			velo += (m/2)*(math.Pow(segmentLengths[i], 2)) + b*(segmentLengths[i])
			tiempo += simpson(0.00, float64(segmentLengths[i]), 0.00, m, b, velo, 50)
		}

		xOffset += segmentLengths[i]
	}

	p := plot.New()

	lines, err := plotter.NewLine(accelPlot)
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