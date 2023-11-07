package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	//"gonum.org/v1/gonum/integrate"
	"github.com/fogleman/gg"
)

const (
	graphResolution = 0.001
)
const (
	pi = 3.1415926535979
)

//type Pattern interface {
//	ColorAt(x, y int) color.Color
//}
//
//type Gradient interface {
//	Pattern
//	AddColorStop(offset float64, color color.Color)
//}

func drawTrack(lengths []float64, angles []float64) {
	//angle is not absolute - is relative to current direction in radians
	//code can be simplified if turns are input as their true compass direction and not relative to the turn before
	var currentAngle = 0.0
	var currentX = 75.0
	var currentY = 75.0
	dc := gg.NewContext(250, 250)
	for i := 0; i < len(lengths); i++ {
		if angles[i] > 0.01 || angles[i] < -0.01 {
			//curve - .01 rad is tolerance
			var radius = lengths[i] / angles[i]
			dc.DrawArc(currentX + math.Cos(-(angles[i]-pi/2)), currentY + math.Sin(-(angles[i]-pi/2)), radius, currentAngle-pi/2, currentAngle+angles[i]-pi/2)
			currentX += 2 * radius * math.Cos(currentAngle + angles[i])
			currentY += 2 * radius * math.Sin(currentAngle + angles[i])
		} else {
			//straight
				dc.DrawLine(currentX, currentY, currentX + lengths[i] * math.Cos(currentAngle), currentY + lengths[i] * math.Sin(currentAngle))
				currentX += radius * math.Cos(angles[i])
				currentY += radius * math.Sin(angles[i])
			} 
		}
		dc.SetRGB(50, 100, 0)
		dc.Stroke()
		currentAngle += angles[i]
		if currentAngle >= 2*pi {
			currentAngle -= 2 * pi
		}
		if currentAngle <= -2*pi {
			currentAngle += 2 * pi
		}
	}
	dc.SavePNG("trackLayout.png")
}

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

	const W = 10
	const H = 10

	// run command to test:
	// go run . 9 0 1 1 1 1 -10 1
	//bug in track if angle is less than 1 rad it floor functions to 0 and draws as a straight + drawing circles 90 degrees off

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
	segmentLengths := []float64{100, 100, 100, 100}
	segmentRotation := []float64{0, pi, 0, pi}

	drawTrack(segmentLengths, segmentRotation)
	/*turnCount := 0
	straightCount := 0

	for _, segmentType := range segmentRotation {
		if segmentType == 0 {
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
	currentTickVelo := args[1]
	var accelPlot plotter.XYs
	var veloPlot plotter.XYs
	fmt.Println("initialVelocity: ", initialVelocity)
	argIndex := 2
	xOffset := 0.0
	tiempo := 0.00
	velo := initialVelocity
	a := 0.00
	b := 0.00
	c := 0.00
	m := 0.00
	for i, angle := range segmentRotation {
		if int(angle) == 0 {
			a := args[argIndex]
			argIndex++
			b := args[argIndex]
			argIndex++
			c := currentTickAccel - a*math.Pow(xOffset, 2) - b*xOffset
			d := currentTickVelo - (a/3)*math.Pow(xOffset, 3) - (b/2)*math.Pow(xOffset, 2) - c*xOffset
			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				currentTickAccel = a*math.Pow(x, 2) + b*x + c
				currentTickVelo = (a/3)*math.Pow(x, 3) + (b/2)*math.Pow(x, 2) + c*x + d
				accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
				veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})
			}
		} else {
			m := args[argIndex]
			argIndex++
			b := currentTickAccel - m*xOffset
			d := currentTickVelo - (m/2)*math.Pow(xOffset, 2) - b*xOffset
			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				currentTickAccel = m*x + b
				currentTickVelo = (m/2)*math.Pow(x, 2) + b*x + d
				accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
				veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})
			}
		}
		if int(angle) == 0 {
			velo += (a/3)*(math.Pow(segmentLengths[i], 3)) + (b/2)*(math.Pow(segmentLengths[i], 2)) + c*(segmentLengths[i])
			tiempo += simpson(0.00, float64(segmentLengths[i]), a, b, c, velo, 50)
		}
		if int(angle) == 0 {
			velo += (m/2)*(math.Pow(segmentLengths[i], 2)) + b*(segmentLengths[i])
			tiempo += simpson(0.00, float64(segmentLengths[i]), 0.00, m, b, velo, 50)
		}
		xOffset += segmentLengths[i]
	}
	accelpo := plot.New()
	velopo := plot.New()
	lines, err := plotter.NewLine(accelPlot)
	if err != nil {
		panic(err)
	}
	lines2, err2 := plotter.NewLine(veloPlot)
	if err2 != nil {
		panic(err)
	}
	accelpo.Add(lines)
	velopo.Add(lines2)
	accelpo.X.Min = 0
	accelpo.Y.Min = 0
	velopo.X.Min = 0
	velopo.Y.Min = 0
	accelpo.X.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var ticks []plot.Tick
		for i := math.Ceil(min); i <= max+1; i++ {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		}
		return ticks
	})
	velopo.X.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var ticks []plot.Tick
		for i := math.Ceil(min); i <= max+1; i++ {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		}
		return ticks
	})
	accelpo.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var ticks []plot.Tick
		for i := math.Ceil(min); i <= max+1; i++ {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		}
		return ticks
	})
	velopo.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		var ticks []plot.Tick
		for i := math.Ceil(min); i <= max+1; i++ {
			ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
		}
		return ticks
	})
	accelpo.Add(plotter.NewGrid())
	velopo.Add(plotter.NewGrid())
	//given an array of tuples w/ speeds and angles,
	//draw circle and straightaways with angles and speeds to form a circle given an array of turns [0,180,0,180]
	if err := accelpo.Save(4*vg.Inch, 4*vg.Inch, "graph.png"); err != nil {
		panic(err)
	}
	if err2 := accelpo.Save(4*vg.Inch, 4*vg.Inch, "velograph.png"); err != nil {
		panic(err2)
	}
	accelpo = plot.New()
	velopo = plot.New()
	lines, err = plotter.NewLine(accelPlot)
	if err != nil {
		panic(err)
	}
	lines2, err2 = plotter.NewLine(accelPlot)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("Total time:", tiempo)*/
}
