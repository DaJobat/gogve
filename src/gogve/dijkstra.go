package gogve

import (
	"container/heap"
	"log"
	"math"
	"sort"
)

type WeightedDigraph interface {
	DirectedGraph
	Weights() map[Edge]float32
}

type RelaxableAttribute interface {
	Attribute
	ShortestEstimateFromSource() float32
	SetShortestEstimateFromSource(float32)
}

type RelaxableAttributes map[Vertex]RelaxableAttribute

type DijkstraAttribute struct {
	ShortestEstimate float32
	predecessor      Vertex
}

func (a *DijkstraAttribute) ShortestEstimateFromSource() float32 {
	return a.ShortestEstimate
}

func (a *DijkstraAttribute) SetShortestEstimateFromSource(est float32) {
	a.ShortestEstimate = est
}

func (a *DijkstraAttribute) Distance() int {
	return int(a.ShortestEstimate)
}

func (a *DijkstraAttribute) Predecessor() Vertex {
	return a.predecessor
}

func (a *DijkstraAttribute) SetPredecessor(v Vertex) {
	a.predecessor = v
}

func (da RelaxableAttributes) ToAttributeMap() AttributeMap {
	out := make(AttributeMap)
	for v, ba := range da {
		out[v] = ba
	}
	return out
}

func initSingleSource(graph DirectedGraph, source Vertex) RelaxableAttributes {
	dt := make(RelaxableAttributes)
	for _, v := range graph.Vertices() {
		dt[v] = &DijkstraAttribute{
			ShortestEstimate: float32(math.Inf(1)),
		}
	}

	dt[source].SetShortestEstimateFromSource(0)
	return dt
}

//relax finds the current shortest estimate for an edge. if an edge is changed, returns true
func relax(graph WeightedDigraph, edge Edge, attrs RelaxableAttributes) bool {
	fromAttr := attrs[edge.From()]
	toAttr := attrs[edge.To()]
	// check if this edge gives us a shorter path than the previous path to
	// the vertex we're going to
	changed := false
	fromSEFS := fromAttr.ShortestEstimateFromSource() + graph.Weights()[edge]
	if toAttr.ShortestEstimateFromSource() > fromSEFS {
		toAttr.SetShortestEstimateFromSource(fromSEFS)
		toAttr.SetPredecessor(edge.From())
		changed = true
	}

	return changed
}

func Dijkstra(graph WeightedDigraph, source Vertex) RelaxableAttributes {
	log.Print("dijkstra")
	attrs := initSingleSource(graph, source)

	return dijkstraLoop(graph, attrs, nil)
}

func dijkstraLoop(graph WeightedDigraph, attributes RelaxableAttributes, breakCondition func(Vertex) bool) RelaxableAttributes {
	queue := initDijkstraQueue(graph, attributes)
	outs := make(RelaxableAttributes)

	for queue.Len() > 0 {
		nextShortestVertex := heap.Pop(queue).(*VertexPriorityItem).vertex
		relaxed := false
		for _, edge := range graph.Edges()[nextShortestVertex] {
			relaxed = relax(graph, edge, attributes) || relaxed
		}
		if relaxed {
			sort.Sort(queue) //once we've relaxed an edge, we need to resort the heap. I want a better way of doing this
		}

		outs[nextShortestVertex] = attributes[nextShortestVertex]
	}
	return outs
}

func initDijkstraQueue(graph WeightedDigraph, attrs RelaxableAttributes) *MinPriorityQueue {
	queue := make(MinPriorityQueue, len(graph.Vertices()))
	for i, v := range graph.Vertices() {
		attr := attrs[v].(*DijkstraAttribute)
		queue[i] = &VertexPriorityItem{
			vertex:   v,
			index:    i,
			priority: &attr.ShortestEstimate,
		}
	}

	heap.Init(&queue)

	return &queue
}
