package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	// go run . 9 0 -1 2 -1 0.55 -3.5 -1.4
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// CONSTANT DEFINITIONS:
var segmentLengths = []float64{200, 100, 200, 100}

// elevation/curvature data, evenly sampled over entire track
var inclineSlopeSampling = []float64{0.05, -0.05, 0.05, -0.05}
var curvatureSampling = []float64{1000, 31.83, 1000, 31.83}

// number of points in the graph to compute:
const numTicks = 1000

// input arguments:

// initial velocity, initial acceleration, then accel curve params
// 0: initial velocity
// 1: initial acceleration
// 2-4: parabola params
// next 3: parabola params

// this should sum to 0 at the end
var sum_mgh = 0.0

// this is to test if elevation even matters
var net_hill_speedup_work = 0.0

// returns (work done, centripetal force)
func CalculateWorkDone(velocity float64, curvature float64, step_distance float64, sinAngle float64) (float64, float64) {
	const carMassKg = 298.0
	const dragCoefficient = 0.1275
	const wheelCircumference = 1.875216

	airResistance := dragCoefficient * math.Pow(velocity, 2)

	var centripetalForce float64
	if curvature == 0 {
		centripetalForce = 0
	} else {
		centripetalForce = (math.Pow(velocity, 2) / math.Abs(curvature)) * carMassKg
	}
	//mgsin(theta)
	slope_force := carMassKg * 9.81 * sinAngle
	//force * distance
	sum_mgh += slope_force * step_distance

	//TOTAL FORCE
	force := airResistance + slope_force
	work := force * step_distance

	//Calculating motor efficiency (function of velocity)
	motorRpm := 60 * (velocity / wheelCircumference)
	motorCurrent := (-3*motorRpm - 2700) / 13
	var motorEfficiency float64
	if motorCurrent >= 14 {
		motorEfficiency = 0.9264 + 0.0015*(motorCurrent-14)
	} else {
		motorEfficiency = ((motorCurrent - 1.1) / (motorCurrent - .37)) - .02
	}

	//Calculating the work done because the slope changes the speed and thus the motor efficiency (not mgh) (Brachistochrone)
	net_hill_speedup_work += (1.0 / motorEfficiency) * slope_force * step_distance

	return motorEfficiency * work, centripetalForce
}

func integrand(x float64, q float64, w float64, e float64, r float64) float64 {
	return 1.0 / (q*math.Pow(x, 3) + w*math.Pow(x, 2) + e*x + r)
}

// simpson calculates the definite integral of a function using Simpson's rule.
// a: the lower limit of integration.
// b: the upper limit of integration.
// q, w, e, r: parameters of the function to be integrated.
// n: the number of subintervals to use in the approximation; should be even.
func simpson(a float64, b float64, q float64, w float64, e float64, r float64, n int) float64 {
	h := (b - a) / float64(n)
	sum := integrand(a, q, w, e, r) + integrand(b, q, w, e, r)

	for i := 1; i < n; i += 2 {
		sum += 4 * integrand(a+float64(i)*h, q, w, e, r)
	}

	for i := 2; i < n-1; i += 2 {
		sum += 2 * integrand(a+float64(i)*h, q, w, e, r)
	}

	return (h / 3) * sum
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

	graphOutput := false
	if !hasEndArg || rawArgs[len(rawArgs)-1] != "none" {
		graphOutput = true
	}

	expectedArgCount := 2 + len(segmentLengths)*2
	if hasEndArg {
		expectedArgCount++
	}
	if len(rawArgs) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		fmt.Println("Recieved: ", len(rawArgs))
		os.Exit(1)
	}

	var totalLength float64 = 0
	for i := range segmentLengths {
		totalLength += segmentLengths[i]
	}

	var graphResolution float64 = 1 / float64(numTicks)
	graphResolution *= totalLength

	currentTickVelo := math.Abs(args[0]) + 1
	currentTickAccel := args[1]

	var accelPlot plotter.XYs
	var veloPlot plotter.XYs
	var forcePlot plotter.XYs
	var energyPlot plotter.XYs
	var curvaturePlot plotter.XYs
	var elevPlot plotter.XYs

	argIndex := 2
	segmentStart := 0.0
	tiempo := 0.00
	velo := currentTickVelo
	var trackDrawingVelocities = ""
	var totalEnergyLost = 0.0
	var maxAccel, minAccel, maxVelo, minVelo, maxCentripetal float64 = math.Inf(-1), math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)
	var colorOffsetVar = 0.0
	for _, segmentLength := range segmentLengths {
		//checking different accel and velocity for different curves
		a := args[argIndex] / segmentLength
		argIndex++
		b := args[argIndex] / segmentLength
		argIndex++
		c := currentTickAccel - a*math.Pow(segmentStart, 2) - b*segmentStart
		var red = 255
		var green = 0
		var blue = 0
		for x := segmentStart; x <= segmentStart+segmentLength; x += graphResolution {
			timeToTravel := graphResolution / currentTickVelo
			currentTickAccel = a*math.Pow(x, 2) + b*x + c
			currentTickVelo += currentTickAccel * timeToTravel

			if currentTickAccel > maxAccel {
				maxAccel = currentTickAccel
			}
			if currentTickAccel < minAccel {
				minAccel = currentTickAccel
			}

			// Update max and min velocity
			if currentTickVelo > maxVelo {
				maxVelo = currentTickVelo
			}
			if currentTickVelo < minVelo {
				minVelo = currentTickVelo
			}

			currentCurvature := curvatureSampling[int(float64(x)/totalLength*float64(len(curvatureSampling)))]
			currentElevation := inclineSlopeSampling[int(float64(x)/totalLength*float64(len(inclineSlopeSampling)))]

			if currentCurvature > 500 {
				currentCurvature = 0
			}

			var currentTickEnergy, currentTickCentripetal = CalculateWorkDone(currentTickVelo, currentCurvature, graphResolution, currentElevation)

			if currentTickCentripetal > maxCentripetal {
				maxCentripital = currentTickCentripetal
			}

			totalEnergyLost += currentTickEnergy
			accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
			veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})
			energyPlot = append(energyPlot, plotter.XY{X: x, Y: totalEnergyLost})
			curvaturePlot = append(curvaturePlot, plotter.XY{X: x, Y: currentCurvature})
			elevPlot = append(elevPlot, plotter.XY{X: x, Y: currentElevation})

			//converts and makes velocity string
			colorOffsetStr := strconv.FormatFloat(colorOffsetVar/totalLength, 'f', 4, 64)

			// if statment only needed to prevent printing final point
			if graphOutput && colorOffsetVar/totalLength <= 1.0 {
				trackDrawingVelocities += "<stop offset=\"" + colorOffsetStr + "\" style=\"stop-color:rgb(" + strconv.Itoa(red) + "," + strconv.Itoa(green) + "," + strconv.Itoa(blue) + ");stop-opacity:1\"/>\n"

				colorOffsetVar += graphResolution

				//max of 16 units of speed... can change scale later by putting in for denominator
				red = int(math.Round(255 * currentTickVelo / 16))
				blue = 0
				green = 0
			}
		}

		velo += (a/3)*(math.Pow(segmentLength, 3)) + (b/2)*(math.Pow(segmentLength, 2)) + c*(segmentLength)
		tiempo += simpson(0.00, float64(segmentLength), a, b, c, velo, 50)

		segmentStart += segmentLength
	}

	if graphOutput {
		os.MkdirAll("./plots", 0755)

		outputGraph(accelPlot, "./plots/acceleration.png")
		outputGraph(veloPlot, "./plots/velocity.png")
		outputGraph(forcePlot, "./plots/force.png")
		outputGraph(energyPlot, "./plots/energy.png")
		outputGraph(curvaturePlot, "./plots/curvature.png")
		outputGraph(elevPlot, "./plots/elevation.png")

		// fmt.Println(trackDrawingVelocities)
	}

	fmt.Println("Initial Velocity (m/s):", veloPlot[0].Y)
	fmt.Println("Final Velocity (m/s):", veloPlot[len(veloPlot)-1].Y)
	fmt.Println("Max Velocity (m/s):", maxVelo)
	fmt.Println("Min Velocity (m/s):", minVelo)
	fmt.Println("Max Acceleration (m/s^2):", maxAccel)
	fmt.Println("Min Acceleration (m/s^2):", minAccel)
	fmt.Println("Max Centripetal Force (N): ", maxCentripetal)
	fmt.Println("Final Velocity (m/s):", veloPlot[len(veloPlot)-1].Y)
	fmt.Println("Time Elapsed (s): ", tiempo)
	fmt.Println("Energy Consumed (J): ", energyPlot[len(energyPlot)-1].Y)
	fmt.Println("Energy Consumption (W): ", energyPlot[len(energyPlot)-1].Y/tiempo)

	fmt.Println(sum_mgh)
	fmt.Println(net_hill_speedup_work)
}
