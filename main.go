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

	var graphsFirstFile = dataLoading("strecken/s_100000_1.dat")
	fmt.Println(len(graphsFirstFile))
	var filteredGraphs = filterGraphs(graphsFirstFile)
	fmt.Println(len(filteredGraphs))

	intersectionCounter := lineSweep(filteredGraphs)

	fmt.Println("Total Intersections:", intersectionCounter)
}

func dataLoading(filename string) []Graph {

	// Open the .dat file
	file, err := os.Open(filename)
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

type Graph struct {
	Start, End Point
}
type Point struct {
	X, Y float64
}
type EventType int

type Event struct {
	X    float64
	Y    float64
	Type EventType
	Seg1 int
	Seg2 int
}
type EventQueue struct {
	Events []Event
}

func (eq *EventQueue) Len() int {
	return len(eq.Events)
}
func (eq *EventQueue) Less(i, j int) bool {
	if eq.Events[i].X == eq.Events[j].X {
		return eq.Events[i].Type < eq.Events[j].Type
	}
	return eq.Events[i].X < eq.Events[j].X
}
func (eq *EventQueue) Swap(i, j int) {
	eq.Events[i], eq.Events[j] = eq.Events[j], eq.Events[i]
}
func (eq *EventQueue) Push(x interface{}) {
	eq.Events = append(eq.Events, x.(Event))
}
func (eq *EventQueue) Pop() interface{} {
	old := eq.Events
	n := len(old)
	event := old[n-1]
	eq.Events = old[:n-1]
	return event
}

type BalancedTree struct {
	SegmentMap map[int]int // Segment index -> position in balanced tree
}

func (bt *BalancedTree) Insert(segment int) {
	bt.SegmentMap[segment] = 0 // Insert at position 0 (dummy value)
}
func (bt *BalancedTree) Delete(segment int) {
	delete(bt.SegmentMap, segment)
}
func GetSuccessor(balancedTree *BalancedTree, segment int) int {
	successor := -1
	for s := range balancedTree.SegmentMap {
		if s > segment && (successor == -1 || s < successor) {
			successor = s
		}
	}
	return successor
}

func GetPredecessor(balancedTree *BalancedTree, segment int) int {
	predecessor := -1
	for s := range balancedTree.SegmentMap {
		if s < segment && s > predecessor {
			predecessor = s
		}
	}
	return predecessor
}
func GetNeighbor(balancedTree *BalancedTree, segment int) int {
	if _, ok := balancedTree.SegmentMap[segment]; ok {
		return segment
	}
	return -1
}
func lineSweep(filteredGraphs []Graph) int {
	intersectionCounter := 0
	segments := []Graph(filteredGraphs) // Example list of segments

	eq := &EventQueue{}   // Create an empty event queue
	bt := &BalancedTree{} // Create an empty balanced tree
	bt.SegmentMap = make(map[int]int)

	// Insert start and end events for each segment into the event queue
	for i, segment := range segments {
		startEvent := Event{Type: 0, X: segment.Start.X, Seg1: i}
		endEvent := Event{Type: 1, X: segment.End.X, Seg1: i}
		heap.Push(eq, startEvent)
		heap.Push(eq, endEvent)
	}

	for eq.Len() > 0 {
		event := heap.Pop(eq).(Event)

		if event.Type == 0 { // Start-Event
			bt.Insert(event.Seg1)
			CheckForIntersect(eq, segments, event.Seg1, GetSuccessor(bt, event.Seg1)) // if intersection, it will be added to eq
			CheckForIntersect(eq, segments, event.Seg1, GetPredecessor(bt, event.Seg1))
		} else if event.Type == 1 { // End-Event
			CheckForIntersect(eq, segments, GetPredecessor(bt, event.Seg1), GetSuccessor(bt, event.Seg1))
			bt.Delete(event.Seg1)
		} else if event.Type == 2 { // Intersection-Event
			intersectionCounter++
			fmt.Printf("Segments %d and %d intersect\n", event.Seg1, event.Seg2)
			bt.Delete(event.Seg1)
			bt.Delete(event.Seg2)
			// reinsert both so they are in the correct order in the tree
			bt.Insert(event.Seg1)
			bt.Insert(event.Seg2)
			CheckForIntersect(eq, segments, event.Seg1, GetNeighbor(bt, event.Seg1))
			CheckForIntersect(eq, segments, event.Seg2, GetNeighbor(bt, event.Seg2))
		}
	}
	return intersectionCounter
}

func CheckForIntersect(eq *EventQueue, segments []Graph, s1, s2 int) {
	if s1 != -1 && s2 != -1 {
		segment1 := segments[s1]
		segment2 := segments[s2]
		test := areIntercepting(segment1, segment2)
		if test == true {
			intersectionPoint, found := findIntersectionPoint(segment1, segment2) // You need to implement this function
			if found {
				fmt.Println(segment1, segment2)
				crossEvent := Event{Type: 2, X: intersectionPoint.X, Y: intersectionPoint.Y, Seg1: s1, Seg2: s2}
				heap.Push(eq, crossEvent)
			}
		}
	}
}

func areIntercepting(graph1 Graph, graph2 Graph) bool {
	// Bestimme die Punkte P und Q von Graph 1
	p1 := Point{graph1.Start.X, graph1.Start.Y}
	p2 := Point{graph1.End.X, graph1.End.Y}
	// Bestimme die Punkte R1 und R2 von Graph 2
	q1 := Point{graph2.Start.X, graph2.Start.Y}
	q2 := Point{graph2.End.X, graph2.End.Y}

	if ccw(p1, p2, q1) == 0 && ccw(p1, p2, q2) == 0 {
		return isPointOnLine2(p1, p2, q1) || isPointOnLine2(p1, p2, q2)
	} else if ccw(p1, p2, q1)*ccw(p1, p2, q2) <= 0 && ccw(q1, q2, p1)*ccw(q1, q2, p2) <= 0 {
		return true
	}
	return false
}

func isPointOnLine2(p1, p2, q Point) bool {
	// Überprüfung, ob q auf der Strecke p1-p2 liegt
	if (q.X >= p1.X && q.X <= p2.X) || (q.X >= p2.X && q.X <= p1.X) {
		return true
	}
	return false
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
func crossProduct(p1, p2, p3 Point) float64 {
	return (p2.X-p1.X)*(p3.Y-p1.Y) - (p2.Y-p1.Y)*(p3.X-p1.X)
}
