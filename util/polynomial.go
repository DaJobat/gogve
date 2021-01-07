package util

import (
	"math"
)

type Polynomial interface {
	Calculate(x float64) float64
	Coefficients() []float64
	Degree() int
}

func QuadraticIntersection(p Polynomial, r Polynomial) (p0, p1 FVec) {
	if p.Degree() > 3 || r.Degree() > 3 {
		panic("this only works for quadratic or lower polynomials")
	}

	var big, small Polynomial
	if p.Degree() < r.Degree() {
		small, big = p, r
	} else {
		small, big = r, p
	}

	diffCoeffs := SparseRPCoefficients{}

	for i := range big.Coefficients() {
		if i < small.Degree() {
			diffCoeffs[i] = big.Coefficients()[i] - small.Coefficients()[i]
		} else {
			diffCoeffs[i] = big.Coefficients()[i]
		}
	}

	np := NewRPolynomial(diffCoeffs, big.Degree())

	x1, x2 := Quadratic(np)
	y1, y2 := small.Calculate(x1), small.Calculate(x2)

	return &fvec2{x1, y1}, &fvec2{x2, y2}
}

func Quadratic(p Polynomial) (result1, result2 float64) {
	if p.Degree() > 3 {
		panic("this only works for quadratic or lower polynomials")
	}

	var tc []float64 //temp coefficient slice so we can pad it up
	if len(p.Coefficients()) < 3 {
		//pad the coefficents so we have enough
		tc = make([]float64, 3)
		for i, c := range p.Coefficients() {
			tc[i] = c
		}
	} else {
		tc = p.Coefficients()
	}

	divBy := tc[2] * 2
	rootPart := math.Sqrt((tc[1] * tc[1]) - (4 * tc[2] * tc[0]))

	return ((-tc[1] + rootPart) / divBy), ((-tc[1] - rootPart) / divBy)
}

//ZPolynomial represents a polynomial where all coefficients are in Z, the set of integers
type ZPolynomial struct {
	coeffs []int
	degree int
}

//NewZPolynomial creates a new ZPolynomial. The map is used to allow easy representation of polynomials like
// x^2 + x^5  without having to put a bunch of zeroes in for x^0, x^1, x^3 and x^4
func NewZPolynomial(sparseCoeffs SparseZPCoefficients, degree int) *ZPolynomial {

	zp := ZPolynomial{
		degree: degree,
		coeffs: coeffsFromSparseZCoeffs(sparseCoeffs, degree),
	}

	return &zp
}

func (zp *ZPolynomial) Calculate(x float64) (y float64) {
	var floati float64 = 1
	y = float64(zp.coeffs[0])
	for i := 1; i < zp.degree; i++ {
		y += float64(zp.coeffs[i]) * math.Pow(x, floati)
		floati += 1
	}
	return y
}

type RPolynomial struct {
	coeffs []float64
	degree int
}

func NewRPolynomial(sparseCoeffs SparseRPCoefficients, degree int) *RPolynomial {
	p := RPolynomial{
		degree: degree,
		coeffs: coeffsFromSparseRCoeffs(sparseCoeffs, degree),
	}

	return &p
}

func (rp *RPolynomial) Calculate(x float64) (y float64) {
	var floati float64 = 1
	y = rp.coeffs[0]
	for i := 1; i < rp.degree; i++ {
		y += float64(rp.coeffs[i]) * math.Pow(x, floati)
		floati += 1
	}
	return y
}

func (rp *RPolynomial) Coefficients() []float64 {
	return rp.coeffs
}

func (rp *RPolynomial) Degree() int {
	return rp.degree
}

func coeffsFromSparseZCoeffs(sparseCoeffs SparseZPCoefficients, degree int) (coeffs []int) {
	coeffs = make([]int, degree)
	for i := 0; i < degree; i++ {
		if d, ok := sparseCoeffs[i]; ok {
			//coefficient exists in sparseCoeffs
			coeffs[i] = d
		} else {
			//no coefficient, add 0 (the zero value of int, so do nothing)
		}
	}

	return coeffs
}

func coeffsFromSparseRCoeffs(sparseCoeffs SparseRPCoefficients, degree int) (coeffs []float64) {
	coeffs = make([]float64, degree)
	for i := 0; i < degree; i++ {
		if d, ok := sparseCoeffs[i]; ok {
			//coefficient exists in sparseCoeffs
			coeffs[i] = d
		} else {
			//no coefficient, add 0 (the zero value of int, so do nothing)
		}
	}

	return coeffs
}

type SparseZPCoefficients map[int]int
type SparseRPCoefficients map[int]float64
