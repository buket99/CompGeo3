package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	var graphsFirstFile = dataLoading("s_1000_1.dat")
	fmt.Println(len(graphsFirstFile))
	var filteredGraphs = filterGraphs(graphsFirstFile)
	fmt.Println(len(filteredGraphs))
	/*
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

	*/
}

func filterGraphs(graphs []Graph) []Graph {
	filteredGraphs := make([]Graph, 0)
	intersectingPoints := make(map[Point]int)

	for i, graph := range graphs {
		// Aussortieren, wenn x-Werte oder y-Werte der Start- und Endpunkte gleich sind
		if graph.Start.X == graph.End.X || graph.Start.Y == graph.End.Y {
			fmt.Println("ist ein punkt:", graph)
			continue
		}

		intersects := false
		for j, graph2 := range graphs {
			if i != j && doGraphsIntersect(graph, graph2) {
				// Graphs intersect at a point
				intersectingPoints[graph.Start]++
				intersectingPoints[graph.End]++
				intersects = true
				break
			}
		}
		// Aussortieren, wenn sich mehr als 2 Graphen im selben Punkt schneiden
		if !intersects && intersectingPoints[graph.Start] > 2 || intersectingPoints[graph.End] > 2 {
			fmt.Println("mehr als zwei schnitte im selben punkt:", intersectingPoints)
			continue
		}

		// Aussortieren, wenn der Graph einen anderen Graphen nur berührt
		intersects = false
		for _, otherGraph := range graphs {
			if graph != otherGraph && doGraphsIntersect(graph, otherGraph) {
				intersects = true
				fmt.Println("Berührt:", graph, otherGraph)
				break
			}
		}
		if intersects {
			continue
		}

		// Graph zu den gefilterten Graphen hinzufügen
		filteredGraphs = append(filteredGraphs, graph)
	}

	return filteredGraphs
}

// Berechnet das Kreuzprodukt zweier Vektoren.
func crossProduct(p1, p2, p3 Point) float64 {
	return (p2.X-p1.X)*(p3.Y-p1.Y) - (p2.Y-p1.Y)*(p3.X-p1.X)
}

// Überprüft die Orientierung von drei Punkten im Uhrzeigersinn, gegen den Uhrzeigersinn oder kollinear.
func ccw(p1, p2, p3 Point) int {
	result := crossProduct(p1, p2, p3)
	if result > 0 {
		return 1 // Gegen den Uhrzeigersinn (CCW)
	} else if result < 0 {
		return -1 // Im Uhrzeigersinn (CW)
	}
	return 0 // Kollinear
}

// Überprüft, ob sich zwei Graphen echt schneiden.
func doGraphsIntersect(graph1, graph2 Graph) bool {
	orientation1 := ccw(graph1.Start, graph1.End, graph2.Start)
	orientation2 := ccw(graph1.Start, graph1.End, graph2.End)
	orientation3 := ccw(graph2.Start, graph2.End, graph1.Start)
	orientation4 := ccw(graph2.Start, graph2.End, graph1.End)

	// Überprüfe, ob die Linien sich echt schneiden
	if (orientation1 != orientation2) && (orientation3 != orientation4) {
		return true
	}

	return false
}

func getIntersectingPoint(graph1 Graph, graph2 Graph) Point {
	x1, y1 := graph1.Start.X, graph1.Start.Y
	x2, y2 := graph1.End.X, graph1.End.Y
	x3, y3 := graph2.Start.X, graph2.Start.Y
	x4, y4 := graph2.End.X, graph2.End.Y

	denominator := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if denominator == 0 {
		// Linien sind parallel, kein Schnittpunkt vorhanden
		return Point{0, 0}
	}

	intersectionX := ((x1 * y2) - (y1*x2)*(x3-x4) - (x1-x2)*(x3*y4-y3*x4)) / denominator
	intersectionY := ((x1 * y2) - (y1*x2)*(y3-y4) - (y1-y2)*(x3*y4-y3*x4)) / denominator

	intersectionPoint := Point{X: intersectionX, Y: intersectionY}
	return intersectionPoint

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

		graphs = append(graphs, Graph{Start: Point{X: values[0], Y: values[1]}, End: Point{X: values[2], Y: values[3]}})
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

type Graph struct {
	Start Point
	End   Point
}

type Point struct {
	X float64
	Y float64
}

type EventType int

const (
	Start EventType = iota
	End
	Intersection
)

type Event struct {
	x, y      float64
	eventType EventType
}

type StatusNode struct {
	line  Graph
	left  *StatusNode
	right *StatusNode
}

type EventQueue []Event

func (eq EventQueue) Len() int           { return len(eq) }
func (eq EventQueue) Less(i, j int) bool { return eq[i].x < eq[j].x }
func (eq EventQueue) Swap(i, j int)      { eq[i], eq[j] = eq[j], eq[i] }

func (eq *EventQueue) Push(x interface{}) {
	*eq = append(*eq, x.(Event))
}

func (eq *EventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	x := old[n-1]
	*eq = old[0 : n-1]
	return x
}

func lineSweep(events EventQueue) {
	// Ereigniswarteschlange sortieren
	heap.Init(&events)

	// Initialisierung der Statusstruktur (z. B. leeren Baum)

	for events.Len() > 0 {
		// Nächstes Ereignis aus der Warteschlange abrufen
		event := heap.Pop(&events).(Event)

		switch event.eventType {
		case Start:
			// Linie zur Statusstruktur hinzufügen
		case End:
			// Linie aus der Statusstruktur entfernen
		case Intersection:
			// Behandlung von Schnittpunkten
		}
	}
}
