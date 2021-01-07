package util

import (
	"fmt"
	"strings"
)

type FVec interface {
	Degree() int
	X() float64
	Y() float64
	Z() float64
	W() float64
	Cross(FVec) float64
}

type fVec struct {
	degree  int
	members []float64
}

func NewFVec(degree int, members ...float64) FVec {
	if degree < 2 {
		panic("cannot have a 1 member vector")
	}

	v := fVec{
		degree:  degree,
		members: members,
	}
	return &v
}

func NewFVec2(x, y float64) FVec {
	v := NewFVec(2, x, y)
	return v
}

func (f *fVec) Degree() int {
	return f.degree
}

func (f *fVec) X() float64 {
	return f.members[0]
}

func (f *fVec) Y() float64 {
	return f.members[1]
}

func (f *fVec) Z() float64 {
	switch f.Degree() {
	case 2:
		return 0
	default:
		return f.members[2]
	}
}

func (f *fVec) W() float64 {
	switch f.Degree() {
	case 2:
		fallthrough
	case 3:
		return 0
	default:
		return f.members[2]
	}
}

func (f *fVec) String() string {
	s := &strings.Builder{}
	d := 0
	s.WriteString("[")
	for d < f.Degree()-1 {
		s.WriteString(fmt.Sprintf("%.2f, ", f.members[d]))
		d++
	}
	s.WriteString(fmt.Sprintf("%.2f]", f.members[d+1]))
	return s.String()
}

func (f *fVec) Cross(f1 FVec) float64 {
	switch f.Degree() {
	case 2:
		return fVec2CrossProduct(f, f1)
	default:
		return 0
	}
}

func fVec2CrossProduct(p0, p1 FVec) float64 {
	return (p0.X() * p1.Y()) - (p1.X() * p0.Y())
}

type ComparableResult int

const (
	ComparableLess    = -1
	ComparableEqual   = 0
	ComparableGreater = 1
)

type Comparable interface {
	Compare(Comparable) ComparableResult //Returns -1, 0 or 1 as is standard (-1 less, 0 equal, 1 greater)
}

type ComparableFloat float64

func (i ComparableFloat) Compare(n Comparable) ComparableResult {
	switch num := n.(type) {
	case ComparableFloat:
		if i < num {
			return ComparableLess
		} else if i == num {
			return ComparableEqual
		} else {
			return ComparableGreater
		}
	default:
		panic("invalid comparison betweenn types")
	}
}
