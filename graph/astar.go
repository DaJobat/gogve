package graph

import (
	"container/heap"
	"math"
)

type EstimatedVertex interface {
	Vertex
	EstimatedDistance(Vertex) float32
}

type DestinationEstimateAttribute interface {
	RelaxableAttribute
	ShortestEstimateToDestination() float32
	SetShortestEstimateToDestination(float32)
	TotalCostEstimate() float32
}

type AStarAttribute struct {
	*DijkstraAttribute
	estDest float32
}

func (a *AStarAttribute) ShortestEstimateToDestination() float32 {
	return a.estDest
}

func (a *AStarAttribute) SetShortestEstimateToDestination(est float32) {
	a.estDest = est
}

func (a *AStarAttribute) TotalCostEstimate() float32 {
	return a.estDest + a.ShortestEstimate
}

func (a *AStarAttribute) Distance() int {
	return int(a.estDest)
}

type AStarAttributes map[EstimatedVertex]*AStarAttribute

func (aa AStarAttributes) ToAttributeMap() AttributeMap {
	out := make(AttributeMap)
	for v, a := range aa {
		out[v] = a
	}
	return out
}

func initAStarSingleSource(wg WeightedDigraph, source, destination EstimatedVertex) AStarAttributes {
	at := make(AStarAttributes)
	for _, v := range wg.Vertices() {
		at[v.(EstimatedVertex)] = &AStarAttribute{
			DijkstraAttribute: &DijkstraAttribute{
				ShortestEstimate: float32(math.Inf(1)),
			},
			estDest: float32(v.(EstimatedVertex).EstimatedDistance(destination)),
		}
	}

	at[source].SetShortestEstimateFromSource(0)

	return at
}

func AStar(wg WeightedDigraph, source, destination EstimatedVertex) AStarAttributes {
	// A Star basically is a mix of dijkstra and BFS.
	// From BFS we use the concept of an expanding frontier of cells
	// that neighbour the source, rather than using the dijkstra style
	// of putting all cells in a queue and working through them by shortest
	// distance to the source.
	// The reason for this is that instead of setting all
	// initial costs to infinite, we set the initial costs to the
	// provided estimated distance from the vertex being pathed to.
	// if we then used this with Dijkstra's algorithm, we would start far from the source,
	// and the algorithm would have to path backwards to the source

	attrs := initAStarSingleSource(wg, source, destination)
	relaxableAttrs := make(RelaxableAttributes)
	for v, a := range attrs {
		relaxableAttrs[v] = a
	}
	// This is the only changed bit, we need to change all the shortest estimates
	// to mirror the estimatedVertex estimates

	outs := make(AStarAttributes)
	queue := make(MinPriorityQueue, 0)
	queue = append(queue, NewVertexPriorityItem(
		source,
		&attrs[source].DijkstraAttribute.ShortestEstimate,
	))
	heap.Init(&queue)

	for queue.Len() > 0 {
		current := heap.Pop(&queue).(*VertexPriorityItem).Vertex().(EstimatedVertex)
		if current == destination {
			break
		}

		for _, edge := range wg.Edges()[current] {
			if Relax(wg, edge, relaxableAttrs) {
				queue.Push(NewVertexPriorityItem(
					edge.To(),
					&attrs[edge.To().(EstimatedVertex)].DijkstraAttribute.ShortestEstimate,
				))
			}
		}
	}

	return outs
}
