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

const ()

func CalculateForce(velocity float64, currentCurvature float64) float64 {
	carMassKg := 298.0

	var dragCoefficient = 0.1275
	airResistance := dragCoefficient * math.Pow(velocity, 2)

	// todo use slope of elevation

	centripitalForce := (math.Pow(velocity, 2) / currentCurvature) * carMassKg

	if currentCurvature >= 1000 {
		centripitalForce = 0
	}

	return airResistance + centripitalForce
}

func CalculateWorkDone(force float64, distance float64) float64 {

	// todo use motor efficciency

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
	curvatureSampling := []float64{1000, 1000, 31.83, 1000, 1000, 31.83}
	segmentLengths := []float64{200, 100, 200, 100}
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

	var totalLength float64 = 0

	for i := range segmentLengths {
		totalLength += segmentLengths[i]
	}

	graphResolution *= totalLength

	// initial velocity, initial acceleration, then accel curve params
	expectedArgCount := 2 + len(segmentLengths)*3

	if hasEndArg {
		expectedArgCount++
	}

	if len(rawArgs) != expectedArgCount {
		fmt.Println("Expected argument count: ", expectedArgCount)
		fmt.Println("Recieved: ", len(rawArgs))
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

	var totalEnergyLost = 0.0
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
			accelPlot = append(accelPlot, plotter.XY{X: x, Y: currentTickAccel})
			veloPlot = append(veloPlot, plotter.XY{X: x, Y: currentTickVelo})

			var currentTickForce = CalculateForce(currentTickVelo, curvatureSampling[int(float64(xOffset)/totalLength*float64(len(curvatureSampling)))])
			var currentTickEnergy = CalculateWorkDone(currentTickForce, graphResolution)
			totalEnergyLost += currentTickEnergy
			forcePlot = append(forcePlot, plotter.XY{X: x, Y: currentTickForce})
			energyPlot = append(energyPlot, plotter.XY{X: x, Y: totalEnergyLost})
		}

		velo += (a/3)*(math.Pow(segmentLengths[i], 3)) + (b/2)*(math.Pow(segmentLengths[i], 2)) + c*(segmentLengths[i])
		tiempo += simpson(0.00, float64(segmentLengths[i]), a, b, c, velo, 50)

		xOffset += segmentLengths[i]
	}

	if !hasEndArg || rawArgs[len(rawArgs)-1] != "none" {
		os.MkdirAll("./plots", 0755)

		outputGraph(accelPlot, "./plots/acceleration.png")
		outputGraph(veloPlot, "./plots/velocity.png")
		outputGraph(forcePlot, "./plots/force.png")
		outputGraph(energyPlot, "./plots/energy.png")
	}

	fmt.Println("Initial Velocity (m/s):", veloPlot[0].Y)
	fmt.Println("Final Velocity (m/s):", veloPlot[len(veloPlot)-1].Y)
	fmt.Println("Time Elapsed (s): ", tiempo)
	fmt.Println("Energy Consumed (J): ", energyPlot[len(energyPlot)-1].Y)
	fmt.Println("Energy Consumption (W): ", energyPlot[len(energyPlot)-1].Y/tiempo)
}
