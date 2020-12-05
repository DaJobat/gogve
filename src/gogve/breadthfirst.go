package gogve

import (
	"container/list"
	"fmt"
)

// BFS uses colors to determine the visitation state of each
// vertex in a digraph
// white is unvisited, gray is visited via an edge during the search,
// black is visited and all edges are checked
// this prevents issues with graphs with cycles in
type BFSColor uint8

func (bc BFSColor) String() string {
	switch bc {
	case BFSWhite:
		return "white"
	case BFSGray:
		return "gray"
	case BFSBlack:
		return "black"
	default:
		return ""
	}
}

const (
	BFSWhite BFSColor = iota
	BFSGray
	BFSBlack
)

type BFSAttribute struct {
	*baseAttribute
	Color BFSColor
}

func (ba *BFSAttribute) String() string {
	return fmt.Sprintf("c: %s, d: %d, pre: %v\n", ba.Color, ba.distance, ba.predecessor)
}

type BreadthFirstTree map[Vertex]*BFSAttribute

func (bft BreadthFirstTree) ToAttributeMap() AttributeMap {
	out := make(AttributeMap)
	for v, ba := range bft {
		out[v] = ba
	}
	return out
}

type BFSCallback func(Vertex, BreadthFirstTree)

func initBFSTree(graph DirectedGraph, source Vertex) BreadthFirstTree {
	bfsTree := make(BreadthFirstTree)
	//Initialise base attributes of all nodes except source in graph to
	//white, with infinite distance (indicated by -1)
	for _, u := range graph.Vertices() {
		bfsTree[u] = &BFSAttribute{
			baseAttribute: &baseAttribute{
				distance: -1,
			},
			Color: BFSWhite,
		}
	}

	sAttrs, ok := bfsTree[source]
	if !ok {
		panic("source not in graph")
	}

	sAttrs.Color = BFSGray
	sAttrs.distance = 0
	return bfsTree
}

func breadthFirstSearch(graph DirectedGraph, source Vertex, callback BFSCallback) BreadthFirstTree {
	bfsTree := initBFSTree(graph, source)
	queue := list.New() // We are using a linked list as the queue, where PushBack is used as Enqueue and Remove is Dequeue
	queue.PushBack(source)

	for queue.Len() > 0 {
		fromVertex := queue.Front().Value.(Vertex)
		queue.Remove(queue.Front())                      // remove u from the list
		for _, edge := range graph.Edges()[fromVertex] { //Loop over all the edges of u
			toVertex := edge.To()
			if bfsTree[toVertex].Color == BFSWhite { //if we haven't visited this vert before
				bfsTree[toVertex].Color = BFSGray                             //make it so we have
				bfsTree[toVertex].distance = bfsTree[fromVertex].distance + 1 //set its distance
				bfsTree[toVertex].predecessor = fromVertex
				queue.PushBack(toVertex)
			}
		}
		bfsTree[fromVertex].Color = BFSBlack

		if callback != nil {
			callback(fromVertex, bfsTree)
		}
	}

	return bfsTree
}

func BreadthFirstSearch(graph DirectedGraph, source Vertex) BreadthFirstTree {
	return breadthFirstSearch(graph, source, nil)
}

func BreadthFirstSearchCallback(graph DirectedGraph, source Vertex, cb BFSCallback) BreadthFirstTree {
	return breadthFirstSearch(graph, source, cb)
}
