package util

import (
	"math"
)

type DirectrixDirection int

const (
	DirectrixX DirectrixDirection = iota
	DirectrixY
)

var directions = []DirectrixDirection{DirectrixX, DirectrixY}

func NextDirectrixDirection(currentDirection DirectrixDirection) DirectrixDirection {
	var i int
	for _, dir := range directions {
		if dir == currentDirection {
			break
		}
		i++
	}
	return directions[(i+1)%len(directions)]
}

type Parabola struct {
	focus              FVec
	directrix          float64
	directrixDirection DirectrixDirection
}

func NewParabola(focus FVec, directrix float64, direction DirectrixDirection) *Parabola {
	p := Parabola{
		focus:              focus,
		directrix:          directrix,
		directrixDirection: direction,
	}
	return &p
}

func (p *Parabola) Directrix() float64 {
	return p.directrix
}

func (p *Parabola) SetDirectrix(d float64) {
	p.directrix = d
}

func (p *Parabola) Focus() FVec {
	return p.focus
}

func (p *Parabola) SetFocus(f FVec) {
	p.focus = f
}

func (p *Parabola) DirectrixDirection() DirectrixDirection {
	return p.directrixDirection
}
func (p *Parabola) SetDirectrixDirection(dd DirectrixDirection) {
	p.directrixDirection = dd
}

func (p *Parabola) GetParabolaPoint(point float64) (v float64) {
	switch p.directrixDirection {
	case DirectrixX:
		v = p.getParabolaPointDX(point)
	case DirectrixY:
		v = p.getParabolaPointDY(point)
	default:
		panic("invalid directix direction set")
	}

	return v
}

func (p *Parabola) getParabolaPointDX(py float64) (x float64) {
	a := p.focus.X()
	b := p.focus.Y()
	k := p.directrix

	x = ((a + k) / 2) + (math.Pow((py-b), 2) / (2 * (a - k)))

	return x
}

func (p *Parabola) getParabolaPointDY(px float64) (y float64) {
	a := p.focus.X()
	b := p.focus.Y()
	k := p.directrix

	y = ((b + k) / 2) + (math.Pow((px-a), 2) / (2 * (b - k)))
	return y
}
