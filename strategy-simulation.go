package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

<<<<<<< Updated upstream
	//go run . 9 0 -1 2 -1 0.55 -3.5 -1.4
=======
	// NEED A NEW TEST ARG FOR 14 ARGS
>>>>>>> Stashed changes
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const ()

<<<<<<< Updated upstream
// CalculateForce calculates the force based on the given velocity.
func CalculateForce(velocity float64) float64 {
	var dragCoefficient = 1.0
	return dragCoefficient * math.Pow(velocity, 2)
=======
func CalculateForce(velocity float64, curvature float64) float64 {
	carMassKg := 298.0

	var dragCoefficient = 0.1275
	airResistance := dragCoefficient * math.Pow(velocity, 2)

	// todo slope of elevation

	var centripetalForce float64

	if curvature == 0 {
		centripetalForce = 0
	} else {
		centripetalForce = (math.Pow(velocity, 2) / math.Abs(curvature)) * carMassKg
	}

	return airResistance + centripetalForce
>>>>>>> Stashed changes
}

// CalculateWorkDone calculates the work done by a force over a segment.
func CalculateWorkDone(force float64, distance float64) float64 {
<<<<<<< Updated upstream
=======

	// todo use motor efficiency

>>>>>>> Stashed changes
	return force * distance
}

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

func outputGraph(inputArr plotter.XYs, fileName string) {

	toPlot := plot.New()

	lines, err := plotter.NewLine(inputArr)
	if err != nil {
		panic(err)
	}

	toPlot.Add(lines)

	toPlot.X.Tick.Marker = plot.DefaultTicks{}
	toPlot.Y.Tick.Marker = plot.DefaultTicks{}
	toPlot.Add(plotter.NewGrid())
	if err := toPlot.Save(4*vg.Inch, 4*vg.Inch, fileName); err != nil {
		panic(err)
	}
}

func main() {

	// Define whether each segment is a straightaway (true) or a turn (false)
	piecewiseFunctions := []bool{true, false, true, false}
	segmentLengths := []float64{2, 1, 2, 1}
	var graphResolution float64 = 0.001

	rawArgs := os.Args[1:]
	hasEndArg := false

	args := make([]float64, len(rawArgs))
	for i, arg := range rawArgs {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			if i != len(rawArgs)-1 {
				fmt.Println(err)
				return
			} else {
				hasEndArg = true
			}
		}
		args[i] = val
	}

	turnCount := 0
	straightCount := 0

	var totalLength float64 = 0

	for i, segmentType := range piecewiseFunctions {
		if segmentType {
			straightCount++
		} else {
			turnCount++
		}
		totalLength += segmentLengths[i]
	}

	graphResolution *= totalLength

	// initial velocity, initial acceleration, then accel curve params
	expectedArgCount := 2 + turnCount + straightCount*2

	if hasEndArg {
		expectedArgCount++
	}

	if len(rawArgs) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		os.Exit(1)
	}

	currentTickVelo := args[0]
	currentTickAccel := args[1]

	var accelPlot plotter.XYs
	var veloPlot plotter.XYs
	var forcePlot plotter.XYs
	var energyPlot plotter.XYs

	argIndex := 2
	xOffset := 0.0
	tiempo := 0.00
	velo := currentTickVelo
	a := 0.00
	b := 0.00
	c := 0.00
	m := 0.00
	var trackDrawingVelocities = ""

	var totalEnergyLost = 0.0
<<<<<<< Updated upstream
	for i, segmentIsStraight := range piecewiseFunctions {
		//checking different accel and velocity for different curves
		if segmentIsStraight {
			a := args[argIndex]
			argIndex++
			b := args[argIndex]
			argIndex++
			c := currentTickAccel - a*math.Pow(xOffset, 2) - b*xOffset
			d := currentTickVelo - (a/3)*math.Pow(xOffset, 3) - (b/2)*math.Pow(xOffset, 2) - c*xOffset
=======
	var maxAccel, minAccel, maxVelo, minVelo float64 = math.Inf(-1), math.Inf(1), math.Inf(-1), math.Inf(1)
	for i := range segmentLengths {
		//checking different accel and velocity for different curves
		a := args[argIndex]
		argIndex++
		b := args[argIndex]
		argIndex++
		c := currentTickAccel - a*math.Pow(xOffset, 2) - b*xOffset
		d := currentTickVelo - (a/3)*math.Pow(xOffset, 3) - (b/2)*math.Pow(xOffset, 2) - c*xOffset
		for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
			currentTickAccel = a*math.Pow(x, 2) + b*x + c
			currentTickVelo = (a/3)*math.Pow(x, 3) + (b/2)*math.Pow(x, 2) + c*x + d
>>>>>>> Stashed changes

			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				currentTickAccel = a*math.Pow(x, 2) + b*x + c
				currentTickVelo = (a/3)*math.Pow(x, 3) + (b/2)*math.Pow(x, 2) + c*x + d
				accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
				veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})

				var currentTickForce = CalculateForce(currentTickVelo)
				var currentTickEnergy = CalculateWorkDone(currentTickForce, graphResolution)
				totalEnergyLost += currentTickEnergy
				forcePlot = append(forcePlot, plotter.XY{X: x, Y: currentTickForce})
				energyPlot = append(energyPlot, plotter.XY{X: x, Y: totalEnergyLost})
			}
		} else {
			m := args[argIndex]
			argIndex++
			b := currentTickAccel - m*xOffset
			d := currentTickVelo - (m/2)*math.Pow(xOffset, 2) - b*xOffset
			var colorOffsetVar = 0.0
			var red = 255
			var green = 0
			var blue = 0
			for x := xOffset; x <= xOffset+segmentLengths[i]; x += graphResolution {
				currentTickAccel = m*x + b
				currentTickVelo = (m/2)*math.Pow(x, 2) + b*x + d
				accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
				veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})

				//converts and makes velocity string
				colorOffsetStr := strconv.FormatFloat(colorOffsetVar, 'f', 2, 64)
				trackDrawingVelocities += "<stop offset=\"" + colorOffsetStr + "\" style=\"stop-color:rgb(" + strconv.Itoa(red) + "," + strconv.Itoa(green) + "," + strconv.Itoa(blue) + ");stop-opacity:1\"/>"
				var currentTickForce = CalculateForce(currentTickVelo)
				var currentTickEnergy = CalculateWorkDone(currentTickForce, graphResolution)
				totalEnergyLost += currentTickEnergy
				forcePlot = append(forcePlot, plotter.XY{X: x, Y: currentTickForce})
				energyPlot = append(energyPlot, plotter.XY{X: x, Y: totalEnergyLost})

				//166 velo points rn? might scale
				colorOffsetVar += segmentLengths[i] / 166

				//max of 16 units of speed... can change scale later by putting in for denominator
				red = int(math.Round(255 * currentTickVelo / 16))
				blue = 0
				green = 0
			}
<<<<<<< Updated upstream
=======

			// Update max and min velocity
			if currentTickVelo > maxVelo {
				maxVelo = currentTickVelo
			}
			if currentTickVelo < minVelo {
				minVelo = currentTickVelo
			}

			currentCurvature := curvatureSampling[int(float64(xOffset)/totalLength*float64(len(curvatureSampling)))]

			if currentCurvature > 500 {
				currentCurvature = 0
			}

			var currentTickForce = CalculateForce(currentTickVelo, currentCurvature)
			var currentTickEnergy = CalculateWorkDone(currentTickForce, graphResolution)
			totalEnergyLost += currentTickEnergy
			accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
			veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})
			forcePlot = append(forcePlot, plotter.XY{X: x, Y: currentTickForce})
			energyPlot = append(energyPlot, plotter.XY{X: x, Y: totalEnergyLost})
			curvaturePlot = append(curvaturePlot, plotter.XY{X: x, Y: currentCurvature})
>>>>>>> Stashed changes
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

	if !hasEndArg || rawArgs[len(rawArgs)-1] != "none" {
		os.MkdirAll("./plots", 0755)

		outputGraph(accelPlot, "./plots/acceleration.png")
		outputGraph(veloPlot, "./plots/velocity.png")
		outputGraph(forcePlot, "./plots/force.png")
		outputGraph(energyPlot, "./plots/energy.png")
	}
<<<<<<< Updated upstream
=======
	var colorOffsetVar = 0.0
	var red = 255
	var green = 0
	var blue = 0

	//Chromatic Plot
	for i, currentVelocity := range veloPlot {
		colorOffsetStr := strconv.FormatInt(int64(i)/int64(len(veloPlot)), 10)
		red = 255 - int(math.Round(125*((currentVelocity.Y-minVelo)/(maxVelo-minVelo))))
		green = int(math.Round(255 * ((currentVelocity.Y - minVelo) / (maxVelo - minVelo))))
		blue = 0
		//for not printing last point
		if colorOffsetVar/totalLength <= 1.0 {
			trackDrawingVelocities += "<stop offset=\"" + colorOffsetStr + "\" style=\"stop-color:rgb(" + strconv.Itoa(red) + "," + strconv.Itoa(green) + "," + strconv.Itoa(blue) + ");stop-opacity:1\"/>\n"
		}
		colorOffsetVar += graphResolution
	}
	fmt.Println(trackDrawingVelocities)

>>>>>>> Stashed changes
	fmt.Println("Initial Velocity (m/s):", veloPlot[0].Y)
	fmt.Println("Final Velocity (m/s):", veloPlot[len(veloPlot)-1].Y)
	fmt.Println("Time Elapsed (s): ", tiempo)
	fmt.Println("Energy Consumed (J): ", energyPlot[len(energyPlot)-1].Y)
	fmt.Println("Energy Consumption (W): ", energyPlot[len(energyPlot)-1].Y/tiempo)
}
