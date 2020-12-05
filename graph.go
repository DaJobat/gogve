package gogve

type DirectedGraph interface {
	Edges() map[Vertex][]Edge
	AddEdge(from, to Vertex)
	RemoveEdge(Edge)
	Vertices() []Vertex
}

type Edge interface {
	From() Vertex
	To() Vertex
}

type Vertex interface{}

type Attribute interface {
	Distance() int
	Predecessor() Vertex
	SetPredecessor(Vertex)
}

type baseAttribute struct {
	predecessor Vertex
	distance    int
}

func (b *baseAttribute) Distance() int {
	return b.distance
}

func (b *baseAttribute) Predecessor() Vertex {
	return b.predecessor
}

func (b *baseAttribute) SetPredecessor(pre Vertex) {
	b.predecessor = pre
}

type AttributeMap map[Vertex]Attribute

type baseEdge struct {
	from Vertex
	to   Vertex
}

func NewEdge(from, to Vertex) Edge {
	be := baseEdge{
		from: from,
		to:   to,
	}

	return &be
}

func (be *baseEdge) From() Vertex {
	return be.from
}

func (be *baseEdge) To() Vertex {
	return be.to
}

func DeadEnds(g DirectedGraph) []Vertex {
	verts := make([]Vertex, 0)
	for v, e := range g.Edges() {
		if len(e) == 1 {
			verts = append(verts, v)
		}
	}
	return verts
}
