package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {

	var graphsFile0 = dataLoading("strecken/s_1000_10.dat")
	var filteredGraphs0 = filterGraphs(graphsFile0)
	start0 := time.Now()
	intersectionCounter0 := lineSweep(filteredGraphs0)
	duration0 := time.Since(start0)
	fmt.Println("Total Intersections in the file s_1000_10:", intersectionCounter0)
	fmt.Println("Time taken for new calculation:", duration0)

	var graphsFile1 = dataLoading("strecken/s_1000_1.dat")
	var filteredGraphs1 = filterGraphs(graphsFile1)
	start1 := time.Now()
	intersectionCounter1 := lineSweep(filteredGraphs1)
	duration1 := time.Since(start1)
	fmt.Println("Total Intersections in the file s_1000_1:", intersectionCounter1)
	fmt.Println("Time taken for new calculation:", duration1)

	var graphsFile2 = dataLoading("strecken/s_10000_1.dat")
	var filteredGraphs2 = filterGraphs(graphsFile2)
	start2 := time.Now()
	intersectionCounter2 := lineSweep(filteredGraphs2)
	duration2 := time.Since(start2)
	fmt.Println("Total Intersections in the file s_10000_1:", intersectionCounter2)
	fmt.Println("Time taken for new calculation:", duration2)

	var graphsFile3 = dataLoading("strecken/s_100000_1.dat")
	var filteredGraphs3 = filterGraphs(graphsFile3)
	start3 := time.Now()
	intersectionCounter3 := lineSweep(filteredGraphs3)
	duration3 := time.Since(start3)
	fmt.Println("Total Intersections in the file s_100000_1:", intersectionCounter3)
	fmt.Println("Time taken for new calculation:", duration3)

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
		id := len(graphs)

		graphs = append(graphs, Graph{ID: id, Start: Point{X: values[0], Y: values[1]}, End: Point{X: values[2], Y: values[3]}})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return graphs
}

type Graph struct {
	ID         int
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

func (eq *EventQueue) SortByX() {
	sort.SliceStable(eq.Events, func(i, j int) bool {
		// If the x values are equal, sort by the y value
		if eq.Events[i].X == eq.Events[j].X {
			return eq.Events[i].Y < eq.Events[j].Y
		}
		// Otherwise sort by the x value
		return eq.Events[i].X < eq.Events[j].X
	})
}

func lineSweep(filteredGraphs []Graph) int {
	intersectionCounter := 0
	segments := []Graph(filteredGraphs)
	defaultGraph := Graph{
		ID:    -1,
		Start: Point{},
		End:   Point{},
	}
	processedEvents := make(map[Event]bool) // Map to track processed events

	eq := &EventQueue{} // Create an empty event queue
	root := &Node{}     // Create an empty AVL tree root
	root = nil          // Initialize as nil since the tree is initially empty

	// Insert start and end events for each segment into the event queue
	for i, segment := range segments {
		startEvent := Event{Type: 0, X: segment.Start.X, Y: segment.Start.Y, Seg1: i}
		endEvent := Event{Type: 1, X: segment.End.X, Y: segment.End.Y, Seg1: i}
		heap.Push(eq, endEvent)   // Push end event first
		heap.Push(eq, startEvent) // Push start event second
	}

	eq.SortByX()

	for eq.Len() > 0 {
		event := heap.Pop(eq).(Event)
		if processedEvents[event] {
			continue // Skip the event if already processed
		}
		processedEvents[event] = true // Mark the event as processed

		if event.Type == 0 { // Start-Event
			root = insertNode(root, segments[event.Seg1], event.X)
			CheckForIntersect(eq, segments[event.Seg1], getSuccessor(root, segments[event.Seg1], event.X), getPredecessor(root, segments[event.Seg1], event.X))
			CheckForIntersect(eq, segments[event.Seg1], getPredecessor(root, segments[event.Seg1], event.X), defaultGraph)
		} else if event.Type == 1 { // End-Event
			root = deleteNode(root, segments[event.Seg1], event.X)
			CheckForIntersect(eq, segments[event.Seg1], getSuccessor(root, segments[event.Seg1], event.X), getPredecessor(root, segments[event.Seg1], event.X))
		} else if event.Type == 2 {
			if !processedEvents[Event{Type: 2, X: event.X, Seg1: event.Seg2, Seg2: event.Seg1}] {
				intersectionCounter++
				/*
					seg1 := segments[event.Seg1]
					seg2 := segments[event.Seg2]
					fmt.Printf("Segments %d and %d intersect - %d (%.5f %.5f %.5f %.5f) and %d (%.5f %.5f %.5f %.5f)\n", seg1.ID, seg2.ID, seg1.ID, seg1.Start.X, seg1.Start.Y, seg1.End.X, seg1.End.Y, seg2.ID, seg2.Start.X, seg2.Start.Y, seg2.End.X, seg2.End.Y)
				*/
				root = deleteNode(root, segments[event.Seg1], event.X)
				root = deleteNode(root, segments[event.Seg2], event.X)
				// reinsert both so they are in the correct order in the tree
				root = insertNode(root, segments[event.Seg1], event.X)
				root = insertNode(root, segments[event.Seg2], event.X)
				preSeg1 := getPredecessor(root, segments[event.Seg1], event.X)
				sucSeg1 := getSuccessor(root, segments[event.Seg1], event.X)
				preSeg2 := getPredecessor(root, segments[event.Seg2], event.X)
				sucSeg2 := getSuccessor(root, segments[event.Seg2], event.X)
				CheckForIntersect(eq, segments[event.Seg1], preSeg1, sucSeg1)
				CheckForIntersect(eq, segments[event.Seg2], preSeg2, sucSeg2)
			}
		}

	}
	return intersectionCounter
}
func getPredecessor(node *Node, key Graph, sweepX float64) Graph {
	var predecessor Graph

	for node != nil {
		if key.getYatX(sweepX) < node.key.getYatX(sweepX) {
			node = node.left
		} else if key.getYatX(sweepX) > node.key.getYatX(sweepX) {
			predecessor = node.key
			node = node.right
		} else {
			if node.left != nil {
				predecessor = nodeWithMaximumValue(node.left).key
			}
			break
		}
	}

	return predecessor
}

func getSuccessor(node *Node, key Graph, sweepX float64) Graph {
	var successor Graph

	for node != nil {
		if key.getYatX(sweepX) < node.key.getYatX(sweepX) {
			successor = node.key
			node = node.left
		} else if key.getYatX(sweepX) > node.key.getYatX(sweepX) {
			node = node.right
		} else {
			if node.right != nil {
				successor = nodeWithMinimumValue(node.right).key
			}
			break
		}
	}

	return successor
}

func CheckForIntersect(eq *EventQueue, seg1 Graph, preSeg Graph, sucSeg Graph) {
	if preSeg.ID != -1 && preSeg.ID != seg1.ID && areIntercepting(seg1, preSeg) {
		intersectionPoint, _ := findIntersectionPoint(seg1, preSeg)
		crossEvent := Event{Type: 2, X: intersectionPoint.X, Y: intersectionPoint.Y, Seg1: seg1.ID, Seg2: preSeg.ID}
		heap.Push(eq, crossEvent)
	}

	if sucSeg.ID != -1 && sucSeg.ID != seg1.ID && areIntercepting(seg1, sucSeg) {
		intersectionPoint, _ := findIntersectionPoint(seg1, sucSeg)
		crossEvent := Event{Type: 2, X: intersectionPoint.X, Y: intersectionPoint.Y, Seg1: seg1.ID, Seg2: sucSeg.ID}
		heap.Push(eq, crossEvent)
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
func (graph *Graph) getYatX(x float64) float64 {
	if graph.Start.X == graph.End.X {
		return graph.Start.Y
	}
	if math.Abs(graph.Start.X-graph.End.X) < epsilon {
		return math.MaxFloat64
	}
	slope := (graph.End.Y - graph.Start.Y) / (graph.End.X - graph.Start.X)
	return slope*(x-graph.Start.X) + graph.Start.Y
}

const epsilon = 1e-9

func isPointOnLine2(p1, p2, q Point) bool {
	// Überprüfung, ob q auf der Strecke p1-p2 liegt
	if ((q.X >= p1.X && q.X <= p2.X) || (q.X >= p2.X && q.X <= p1.X)) && math.Abs(q.X-p1.X) > epsilon && math.Abs(q.X-p2.X) > epsilon {
		return true
	}
	return false
}

func ccw(p1, p2, p3 Point) float64 {
	return (p2.Y-p1.Y)*(p3.X-p2.X) - (p2.X-p1.X)*(p3.Y-p2.Y)
}

type Node struct {
	key    Graph
	left   *Node
	right  *Node
	height int
	index  int // index of the graph in the segments array
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Calculates the height of the node
func height(N *Node) int {
	if N == nil {
		return 0
	}
	return N.height
}

// Performs a right rotation on the node
func rightRotate(y *Node) *Node {
	x := y.left
	T2 := x.right
	x.right = y
	y.left = T2
	y.height = max(height(y.left), height(y.right)) + 1
	x.height = max(height(x.left), height(x.right)) + 1
	return x
}

// Performs a left rotation on the node
func leftRotate(x *Node) *Node {
	y := x.right
	T2 := y.left
	y.left = x
	x.right = T2
	x.height = max(height(x.left), height(x.right)) + 1
	y.height = max(height(y.left), height(y.right)) + 1
	return y
}

// Calculates the balance factor
// of the node
func getBalanceFactor(N *Node) int {
	if N == nil {
		return 0
	}
	return height(N.left) - height(N.right)
}

func newNode(graph Graph) *Node {
	node := &Node{key: graph}
	node.left = nil
	node.right = nil
	node.height = 1
	return node
}

func insertNode(node *Node, graph Graph, sweepX float64) *Node {
	if node == nil {
		return newNode(graph)
	}
	if graph.getYatX(sweepX) < node.key.getYatX(sweepX) {
		node.left = insertNode(node.left, graph, sweepX)
	} else if graph.getYatX(sweepX) > node.key.getYatX(sweepX) {
		node.right = insertNode(node.right, graph, sweepX)
	} else {
		return node
	}

	node.height = 1 + max(height(node.left), height(node.right))
	balanceFactor := getBalanceFactor(node)

	if balanceFactor > 1 {
		if node.left != nil && graph.getYatX(sweepX) < node.left.key.getYatX(sweepX) {
			return rightRotate(node)
		} else if node.left != nil && graph.getYatX(sweepX) > node.left.key.getYatX(sweepX) {
			node.left = leftRotate(node.left)
			return rightRotate(node)
		}
	}

	if balanceFactor < -1 {
		if node.right != nil && graph.getYatX(sweepX) > node.right.key.getYatX(sweepX) {
			return leftRotate(node)
		} else if node.right != nil && graph.getYatX(sweepX) < node.right.key.getYatX(sweepX) {
			node.right = rightRotate(node.right)
			return leftRotate(node)
		}
	}

	return node
}
func deleteNode(root *Node, graph Graph, sweepX float64) *Node {
	if root == nil {
		return root
	}
	if graph.getYatX(sweepX) < root.key.getYatX(sweepX) {
		root.left = deleteNode(root.left, graph, sweepX)
	} else if graph.getYatX(sweepX) > root.key.getYatX(sweepX) {
		root.right = deleteNode(root.right, graph, sweepX)
	} else {
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		} else {
			temp := nodeWithMinimumValue(root.right)
			root.key = temp.key
			root.right = deleteNode(root.right, temp.key, sweepX)
		}
	}
	if root == nil {
		return root
	}
	root.height = 1 + max(height(root.left), height(root.right))
	balanceFactor := getBalanceFactor(root)

	if balanceFactor > 1 {
		if getBalanceFactor(root.left) >= 0 {
			return rightRotate(root)
		} else {
			root.left = leftRotate(root.left)
			return rightRotate(root)
		}
	}
	if balanceFactor < -1 {
		if getBalanceFactor(root.right) <= 0 {
			return leftRotate(root)
		} else {
			root.right = rightRotate(root.right)
			return leftRotate(root)
		}
	}
	return root
}

// Fetches the Node with maximum
// value from the AVL tree
func nodeWithMaximumValue(node *Node) *Node {
	current := node
	for current.right != nil {
		current = current.right
	}
	return current
}
func nodeWithMinimumValue(node *Node) *Node {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

func filterGraphs(graphs []Graph) []Graph {
	newGraphs := make([]Graph, 0)
	filteredGraphs := make([]Graph, 0)

	for _, graph := range graphs {
		// Aussortieren, wenn x-Werte oder y-Werte der Start- und Endpunkte gleich sind
		if graph.Start.X == graph.End.X || graph.Start.Y == graph.End.Y {
			// fmt.Println("X und Y nicht paarweise verschieden:", graph)
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
				// fmt.Println("Berührt:", graph, otherGraph)
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
				// fmt.Println("Der folgende Graph hat mehr als zwei anderen Graphen im selben Schnittpunkt: ", graph)
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
