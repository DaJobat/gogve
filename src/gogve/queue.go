package gogve

type VertexPriorityItem struct {
	vertex   Vertex
	priority *float32
	index    int
}

type MinPriorityQueue []*VertexPriorityItem

func (pq MinPriorityQueue) Len() int {
	return len(pq)
}

// I think this is the only thing that changes between a min and a max priority queue
func (pq MinPriorityQueue) Less(i, j int) bool {
	return *pq[i].priority < *pq[j].priority
}

func (pq MinPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *MinPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*VertexPriorityItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *MinPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
