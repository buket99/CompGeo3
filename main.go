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

	var graphsFirstFile = dataLoading("s_1000_10.dat")
	fmt.Println(len(graphsFirstFile))
	var filteredGraphs = filterGraphs(graphsFirstFile)
	fmt.Println(len(filteredGraphs))

	/*
		myGraphs := []Graph{
			// Punkt und nicht paarweise verschieden
			{Start: Point{20, 20}, End: Point{20, 20}},
			{Start: Point{10, 10}, End: Point{10, 30}},
			// Berühren
			{Start: Point{0, 2}, End: Point{0.5, 2.5}},
			{Start: Point{0, 2}, End: Point{-0.5, 2.5}},
			// Zweifach Schnittpunkt in 0,0
			{Start: Point{-1, -1}, End: Point{0.3, 0.3}},
			{Start: Point{0.5, 1}, End: Point{-0.5, -1}},
			{Start: Point{-1, 1}, End: Point{1, -1}},
		}
		var myFilteredgraphs = filterGraphs(myGraphs)
		print(len(myFilteredgraphs))

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
	newGraphs := make([]Graph, 0)
	filteredGraphs := make([]Graph, 0)

	for _, graph := range graphs {
		// Aussortieren, wenn x-Werte oder y-Werte der Start- und Endpunkte gleich sind
		if graph.Start.X == graph.End.X || graph.Start.Y == graph.End.Y {
			fmt.Println("X und Y nicht paarweise verschieden:", graph)
			continue
		}
		newGraphs = append(newGraphs, graph)
	}

	for i, graph := range newGraphs {

		// Aussortieren, wenn der Graph einen anderen Graphen nur berührt
		intersects := false
		for h, otherGraph := range newGraphs {
			if i != h && doGraphsTouch(graph, otherGraph) {
				intersects = true
				fmt.Println("Berührt:", graph, otherGraph)
				break
			}
		}
		if intersects {
			continue
		}

		intersects = false
		samePointIntersections := make(map[Point]int) // Track intersecting points and their counts

		for j, graph2 := range newGraphs {
			if i != j && doGraphsIntersect(graph, graph2) {
				if !doGraphsTouch(graph, graph2) {
					// Graphs intersect at a point
					intersectingPoint, found := findIntersectionPoint(graph, graph2) // Assuming a helper function to find the intersection point
					if found {
						samePointIntersections[intersectingPoint]++
						intersects = true
					}
				}
			}
		}

		areIntersectingInTheSamePoint := false
		// Check if there are more than 2 intersections at the same point
		for _, count := range samePointIntersections {
			if count >= 2 {
				areIntersectingInTheSamePoint = true
				fmt.Println("Der folgende Graph hat mehr als zwei anderen Graphen im selben Schnittpunkt: ", graph)
				continue
			}
		}
		if areIntersectingInTheSamePoint {
			continue
		}

		// Graph zu den gefilterten Graphen hinzufügen
		filteredGraphs = append(filteredGraphs, graph)
	}

	return filteredGraphs
}
func findIntersectionPoint(graph1, graph2 Graph) (Point, bool) {
	p1 := graph1.Start
	p2 := graph1.End
	q1 := graph2.Start
	q2 := graph2.End

	// Calculate the slopes of the line segments
	m1 := (p2.Y - p1.Y) / (p2.X - p1.X)
	m2 := (q2.Y - q1.Y) / (q2.X - q1.X)

	// Check if the line segments are parallel
	if m1 == m2 {
		return Point{}, false
	}

	// Calculate the intersection point coordinates
	intersectionX := (m1*p1.X - m2*q1.X + q1.Y - p1.Y) / (m1 - m2)
	intersectionY := m1*(intersectionX-p1.X) + p1.Y

	return Point{X: intersectionX, Y: intersectionY}, true
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
	p1 := graph1.Start
	p2 := graph1.End
	// Bestimme die Punkte R1 und R2 von Graph 2
	q1 := graph2.Start
	q2 := graph2.End

	if ccw(p1, p2, q1)*ccw(p1, p2, q2) <= 0 && ccw(q1, q2, p1)*ccw(q1, q2, p2) <= 0 {
		return true
	}
	return false
}

func doGraphsTouch(graph1, graph2 Graph) bool {
	// Extract the start and end points of Graph 1
	p1 := graph1.Start
	p2 := graph1.End
	// Extract the start and end points of Graph 2
	q1 := graph2.Start
	q2 := graph2.End

	// Check if the start or end point of Graph 1 lies on Graph 2
	if isPointOnLine(p1, q1, q2) || isPointOnLine(p2, q1, q2) {
		return true
	}

	// Check if the start or end point of Graph 2 lies on Graph 1
	if isPointOnLine(q1, p1, p2) || isPointOnLine(q2, p1, p2) {
		return true
	}

	return false
}

func isPointOnLine(p, q1, q2 Point) bool {
	// Check if point p lies on the line segment q1-q2
	return (ccw(p, q1, q2) == 0) && (q1.X <= p.X && p.X <= q2.X || q2.X <= p.X && p.X <= q1.X) &&
		(q1.Y <= p.Y && p.Y <= q2.Y || q2.Y <= p.Y && p.Y <= q1.Y)
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
