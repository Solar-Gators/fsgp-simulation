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

// wind direction and speed
var windDirectionRadians = 0.0
var windSpeed = 10.0 //placeholder value

// elevation data, evenly sampled over entire track
var elevationSampling = []float64{400, 400, 400, 400}

// positive curvature is clockwise, negative is counterclockwise
var curvatureSampling = []float64{1000, 1000, 31.83, 1000, 1000, 31.83}

// number of points in the graph to compute:
const numTicks = 1000

// input arguments:
// Solver dictates the velo and accel. Sim dicatates energy required to execute and time elasped. Solver constrained by energy, optimized for time.
// initial velocity, initial acceleration, then accel curve params
// 0: initial velocity
// 1: initial acceleration
// 2-4: parabola params
// next 3: parabola params

func CalculateWorkDone(velocity float64, step_distance float64, slope float64, prev_velo float64, facing_direction float64) float64 {
	const carMassKg = 298.0
	const dragCoefficient = 0.1275
	const wheelCircumference = 1.875216

	//change the velocity to make relative to wind for airResistance only
	relativeVelocity := velocity - windSpeed*math.Cos(math.Abs(windDirectionRadians-facing_direction))
	airResistance := dragCoefficient * math.Pow(relativeVelocity, 2)

	//mgsin(theta)
	slope_force := carMassKg * 9.81 * math.Sin(math.Atan(slope)) // slope = tan(Theta)
	net_velo_energy := .5 * carMassKg * (velocity - prev_velo)
	total_force := airResistance + slope_force
	total_work := total_force*step_distance + net_velo_energy

	//Calculating motor efficiency (function of velocity)
	motorRpm := 60 * (velocity / wheelCircumference)
	motorCurrent := (-3*motorRpm - 2700) / 13
	var motorEfficiency float64
	if motorCurrent >= 14 {
		motorEfficiency = 0.9264 + 0.0015*(motorCurrent-14)
	} else {
		motorEfficiency = ((motorCurrent - 1.1) / (motorCurrent - .37)) - .02
	}

	// work motor does is "positive"
	return motorEfficiency * total_work
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

	var stepDistance float64 = 1 / float64(numTicks)
	stepDistance *= totalLength

	initialVelocity := math.Abs(args[0]) + 1

	currentTickVelo := initialVelocity
	currentTickAccel := args[1]

	var accelPlot plotter.XYs
	var veloPlot plotter.XYs
	var forcePlot plotter.XYs
	var energyPlot plotter.XYs
	var curvaturePlot plotter.XYs
	var elevPlot plotter.XYs
	var centripetalPlot plotter.XYs
	var facingDirectionPlot plotter.XYs

	facingDirectionRadians := 0.0
	slope := 0.0
	argIndex := 2
	segmentStart := 0.0
	tiempo := 0.00
	velo := currentTickVelo
	prevVelo := 0.0
	var trackDrawingVelocities = ""
	var totalEnergyUsed = 0.0
	var maxAccel, minAccel, maxVelo, minVelo, maxCentripetal float64 = math.Inf(-1), math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)
	var colorOffsetVar = 0.0
	for _, segmentLength := range segmentLengths {
		//checking different accel and velocity for different curves
		a := args[argIndex] / segmentLength
		argIndex++
		b := args[argIndex] / segmentLength
		argIndex++
		c := currentTickAccel
		for i := 0.0; i < segmentLength; i += stepDistance {
			track_pos := segmentStart + i
			timeToTravel := stepDistance / currentTickVelo
			prevVelo = currentTickVelo
			currentTickAccel = a*math.Pow(i, 2) + b*i + c
			currentTickVelo += currentTickAccel * timeToTravel

			maxAccel = max(maxAccel, currentTickAccel)
			minAccel = min(minAccel, currentTickAccel)
			maxVelo = max(maxVelo, currentTickVelo)
			minVelo = min(minVelo, currentTickVelo)

			currentCurvature := curvatureSampling[int(float64(track_pos)/totalLength*float64(len(curvatureSampling)))]
			currentElevation := elevationSampling[int(float64(track_pos)/totalLength*float64(len(elevationSampling)))]

			// slope calculated after start of movement
			if i > 0 {
				previousElevation := elevationSampling[int((float64(track_pos)-stepDistance)/totalLength*float64(len(elevationSampling)))]
				slope = (currentElevation - previousElevation) / stepDistance
			}

			var currentTickCentripetal = 0.0
			if currentCurvature > 500 {
				currentCurvature = 0
			} else {
				currentTickCentripetal = (math.Pow(currentTickVelo, 2) / math.Abs(currentCurvature))
				maxCentripetal = max(maxCentripetal, currentTickCentripetal)

				// update facing direction
				facingDirectionRadians += stepDistance / currentCurvature
				if facingDirectionRadians >= 2*math.Pi {
					facingDirectionRadians -= 2 * math.Pi
				} else if facingDirectionRadians < 0 {
					facingDirectionRadians += 2 * math.Pi
				}
			}

			var currentTickEnergy = CalculateWorkDone(currentTickVelo, stepDistance, slope, prevVelo, facingDirectionRadians)
			totalEnergyUsed += currentTickEnergy

			if graphOutput {
				accelPlot = append(accelPlot, plotter.XY{X: track_pos, Y: currentTickAccel})
				veloPlot = append(veloPlot, plotter.XY{X: track_pos, Y: currentTickVelo})
				energyPlot = append(energyPlot, plotter.XY{X: track_pos, Y: totalEnergyUsed})
				curvaturePlot = append(curvaturePlot, plotter.XY{X: track_pos, Y: currentCurvature})
				elevPlot = append(elevPlot, plotter.XY{X: track_pos, Y: currentElevation})
				centripetalPlot = append(centripetalPlot, plotter.XY{X: track_pos, Y: currentTickCentripetal})
				facingDirectionPlot = append(facingDirectionPlot, plotter.XY{X: track_pos, Y: facingDirectionRadians})
			}

			// if statment only needed to prevent printing final point
			if graphOutput && colorOffsetVar/totalLength <= 1.0 {
				//max of 16 units of speed... can change scale later by putting in for denominator
				const redDivisor = 16

				const blue = 0
				const green = 0
				//converts and makes velocity string
				colorOffsetStr := strconv.FormatFloat(colorOffsetVar/totalLength, 'f', 4, 64)

				trackDrawingVelocities += "<stop offset=\"" + colorOffsetStr + "\" style=\"stop-color:rgb(" + strconv.Itoa(int(math.Round(255*currentTickVelo/redDivisor))) + "," + strconv.Itoa(green) + "," + strconv.Itoa(blue) + ");stop-opacity:1\"/>\n"

				colorOffsetVar += stepDistance
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
		outputGraph(centripetalPlot, "./plots/centripetal.png")
		outputGraph(facingDirectionPlot, "./plots/facingDir.png")

		// fmt.Println(trackDrawingVelocities)
	}

	fmt.Println("Initial Velocity (m/s): ", initialVelocity)
	fmt.Println("Final Velocity (m/s): ", currentTickVelo)
	fmt.Println("Max Velocity (m/s): ", maxVelo)
	fmt.Println("Min Velocity (m/s): ", minVelo)
	fmt.Println("Max Acceleration (m/s^2): ", maxAccel)
	fmt.Println("Min Acceleration (m/s^2): ", minAccel)
	fmt.Println("Max Centripetal Acceleration (m/s^2): ", maxCentripetal)
	fmt.Println("Time Elapsed (s): ", tiempo)
	fmt.Println("Energy Consumed (J): ", totalEnergyUsed)
	fmt.Println("Energy Consumption (W): ", totalEnergyUsed/tiempo)
}
