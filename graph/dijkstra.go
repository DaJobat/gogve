package graph

import (
	"container/heap"
	"math"
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

//Relax finds the current shortest estimate for an edge. if an edge is changed, returns true
func Relax(graph WeightedDigraph, edge Edge, attrs RelaxableAttributes) bool {
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
	attrs := initSingleSource(graph, source)

	return dijkstraLoop(graph, attrs, nil)
}

func dijkstraLoop(graph WeightedDigraph, attributes RelaxableAttributes, breakCondition func(Vertex) bool) RelaxableAttributes {
	queue, vvpm := initDijkstraQueue(graph, attributes)
	outs := make(RelaxableAttributes)

	for queue.Len() > 0 {
		nextShortest := heap.Pop(queue).(*VertexPriorityItem)
		relaxed := false
		for _, edge := range graph.Edges()[nextShortest.vertex] {
			relaxed = Relax(graph, edge, attributes) || relaxed
			heap.Fix(queue, vvpm[edge.To()].index)
		}

		outs[nextShortest.vertex] = attributes[nextShortest.vertex]
	}
	return outs
}

type vertexVPItemMap map[Vertex]*VertexPriorityItem

func initDijkstraQueue(graph WeightedDigraph, attrs RelaxableAttributes) (*MinPriorityQueue, vertexVPItemMap) {
	queue := make(MinPriorityQueue, len(graph.Vertices()))
	vvpm := make(vertexVPItemMap)
	for i, v := range graph.Vertices() {
		attr := attrs[v].(*DijkstraAttribute)
		queue[i] = &VertexPriorityItem{
			vertex:   v,
			index:    i,
			priority: &attr.ShortestEstimate,
		}

		vvpm[v] = queue[i]
	}

	heap.Init(&queue)

	return &queue, vvpm
}
