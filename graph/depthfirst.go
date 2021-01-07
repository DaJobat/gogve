package graph

import (
	"fmt"
)

type DFSAttribute struct {
	color             BFSColor
	predecessor       Vertex
	discoverTime      int // time vert was discovered
	finishTime        int // time all child verts were searched
	children          int // number of children of this vert
	lowestReachable   int // lowest reachable vert from this vert, not including predecessor
	articulationPoint bool
}

func (d *DFSAttribute) IsArticulationPoint() bool {
	return d.articulationPoint
}

func (d *DFSAttribute) String() string {
	return fmt.Sprintf("d: %d, f: %d, children: %d, lr: %d, ap: %t,\n(pre: %v)\n",
		d.discoverTime, d.finishTime, d.children, d.lowestReachable, d.articulationPoint, d.predecessor)
}

func (d *DFSAttribute) Distance() int {
	return 0
}

func (d *DFSAttribute) Predecessor() Vertex {
	return d.predecessor
}

func (d *DFSAttribute) SetPredecessor(pre Vertex) {
	d.predecessor = pre
}

type DFSTree map[Vertex]*DFSAttribute

var time int

func DepthFirstSearch(graph DirectedGraph, source Vertex) DFSTree {
	time = 0
	attrs := make(DFSTree)
	for _, u := range graph.Vertices() {
		attrs[u] = &DFSAttribute{
			color: BFSWhite,
		}
	}

	//start by visiting the source
	dfsVisit(graph, attrs, source)

	for _, u := range graph.Vertices() {
		if attrs[u].color == BFSWhite {
			dfsVisit(graph, attrs, u)
		}
	}

	time = 0
	return attrs
}

func dfsVisit(graph DirectedGraph, attrs DFSTree, u Vertex) {
	time = time + 1
	attrs[u].discoverTime = time
	attrs[u].lowestReachable = time
	attrs[u].color = BFSGray

	for _, edge := range graph.Edges()[u] {
		if edge.To() == attrs[u].predecessor {
			continue
		}
		switch attrs[edge.To()].color {
		case BFSWhite:
			attrs[u].children++
			attrs[edge.To()].predecessor = u
			dfsVisit(graph, attrs, edge.To())
			attrs[u].lowestReachable = min(attrs[u].lowestReachable, attrs[edge.To()].lowestReachable)
			if attrs[u].predecessor == nil && attrs[u].children > 1 {
				// if this is the root of the tree and it has >1 child, it's an articulation point
				attrs[u].articulationPoint = true
			}
			if attrs[u].predecessor != nil && attrs[edge.To()].lowestReachable >= attrs[u].discoverTime {
				// if this vertex has a child which cannot reach a vertex with a lower discoverability (i.e.
				// a vert that was discovered before this one), removing this vert would bisect the tree
				// so this is an articulation point
				attrs[u].articulationPoint = true
			}
		case BFSGray:
			attrs[u].lowestReachable = min(attrs[u].lowestReachable, attrs[edge.To()].discoverTime)
		}
	}

	attrs[u].color = BFSBlack
	time = time + 1
	attrs[u].finishTime = time
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
