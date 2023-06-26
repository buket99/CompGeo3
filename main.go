package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var graphsFirstFile = dataLoading("s_1000_1.dat")
	start := time.Now()
	// Call your function
	var amount1 = amountOfInterceptingGraphs(graphsFirstFile)
	// Get the time again and calculate the duration
	duration := time.Since(start)
	fmt.Println("In the first data set the amount of crossing graphs is ", amount1)
	// Print the duration
	fmt.Println("Time taken for first calculation:", duration)

	var graphsSecondFile = dataLoading("s_1000_10.dat")
	start = time.Now()
	// Call your function
	var amount2 = amountOfInterceptingGraphs(graphsSecondFile)
	// Get the time again and calculate the duration
	duration = time.Since(start)
	fmt.Println("In the second data set the amount of crossing graphs is ", amount2)
	// Print the duration
	fmt.Println("Time taken for second calculation:", duration)

	var graphsThirdFile = dataLoading("s_10000_1.dat")
	start = time.Now()
	// Call your function
	var amount3 = amountOfInterceptingGraphs(graphsThirdFile)
	// Get the time again and calculate the duration
	duration = time.Since(start)
	fmt.Println("In the third data set the amount of crossing graphs is ", amount3)
	// Print the duration
	fmt.Println("Time taken for third calculation:", duration)

	var graphsFourthFile = dataLoading("s_100000_1.dat")
	start = time.Now()
	// Call your function
	var amount4 = amountOfInterceptingGraphs(graphsFourthFile)
	// Get the time again and calculate the duration
	duration = time.Since(start)
	fmt.Println("In the fourth data set the amount of crossing graphs is ", amount4)
	// Print the duration
	fmt.Println("Time taken for fourth calculation:", duration)
}

func dataLoading(filename string) []Graph {

	// Open the .dat file
	file, err := os.Open("strecken/" + filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	var graphs []Graph

	// Loop through each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into fields
		fields := strings.Fields(line)

		// Convert the fields to float64 values
		var values []float64
		for _, field := range fields {
			value, err := strconv.ParseFloat(field, 64)
			if err != nil {
				fmt.Println("Error parsing float value:", err)
			}
			values = append(values, value)
		}

		graphs = append(graphs, Graph{p1X: values[0], p1Y: values[1], p2X: values[2], p2Y: values[3]})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return graphs
}

func amountOfInterceptingGraphs(graphs []Graph) int {
	amount := 0
	// add compgeo1
	// add line sweep algorithm
	return amount
}

type Point struct {
	x, y float64
}

type Graph struct {
	p1X float64
	p1Y float64
	p2X float64
	p2Y float64
}
