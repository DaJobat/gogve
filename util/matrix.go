package util

import (
	"fmt"
	"strings"
)

var (
	ErrCannotMul = fmt.Errorf("cannot multiply matrices, invalid column x row count")
	ErrNoEntry   = fmt.Errorf("no entry in matrix")
	ErrNoRow     = fmt.Errorf("no row in matrix")
	ErrNoColumn  = fmt.Errorf("no column in matrix")
)

type Matrix interface {
	Dimensions() (m, n int)
	Entries() [][]float64
	Row(m int) ([]float64, error)
	Column(n int) ([]float64, error)
	Entry(m, n int) (float64, error)
	SetEntry(m, n int, value float64) error
	Transpose() Matrix
}

// NewMatrix returns a new matrix with m ROWS and n COLUMNS
// e.g.
//		NewMatrix(3, 2) = [ 1 2
//											  3 4
//												5 6 ]
//
func NewMatrix(m, n int, entries ...[]float64) Matrix {
	mat := matrix{
		entries:  make([][]float64, m),
		nRows:    m,
		nColumns: n,
	}

	for j := range mat.entries {
		mat.entries[j] = make([]float64, n)
		for i := range mat.entries[j] {
			if entries == nil || entries[j] == nil || i >= len(entries[j]) {
				break
			}
			mat.entries[j][i] = entries[j][i]
		}

	}

	return &mat
}

type matrix struct {
	nRows    int
	nColumns int
	entries  [][]float64
}

func (mat *matrix) String() string {
	b := &strings.Builder{}

	b.WriteString("[")
	for j := range mat.entries {
		b.WriteString("\n\t")
		for i := range mat.entries[j] {
			b.WriteString(fmt.Sprintf("%.2f, ", mat.entries[j][i]))
		}
	}
	b.WriteString("\n]")

	return b.String()
}

//Dimensions returns the m x n (rows x columns)
// values for the matrix
func (mat *matrix) Dimensions() (m, n int) {
	return mat.nRows, mat.nColumns
}

func (mat *matrix) Transpose() Matrix {
	transpose := NewMatrix(mat.nColumns, mat.nRows)
	for j := range mat.entries {
		for i := range mat.entries[j] {
			e, _ := mat.Entry(j+1, i+1)
			transpose.SetEntry(i+1, j+1, e)
		}
	}

	return transpose
}

func (mat *matrix) Entries() [][]float64 {
	return mat.entries
}

//SetEntry sets the matrix value at row i, column j
func (mat *matrix) SetEntry(i, j int, value float64) error {
	if i-1 < 0 || j-1 < 0 || i > len(mat.entries) || j > len(mat.entries[i-1]) {
		return ErrNoEntry
	}

	mat.entries[i-1][j-1] = value
	return nil
}

//Entry returns the matrix value at row j, column i
func (mat *matrix) Entry(i, j int) (float64, error) {
	if i-1 < 0 || j-1 < 0 || i > len(mat.entries) || j > len(mat.entries[i-1]) {
		return 0, ErrNoEntry
	}
	return mat.entries[i-1][j-1], nil
}

func (mat *matrix) Row(i int) ([]float64, error) {
	if i-1 < 0 || i > mat.nRows {
		return nil, ErrNoRow
	}

	return mat.entries[i-1], nil
}

func (mat *matrix) Column(j int) ([]float64, error) {
	if j-1 < 0 || j > mat.nColumns {
		return nil, ErrNoColumn
	}
	col := make([]float64, mat.nRows)

	for i, row := range mat.entries {
		col[i] = row[j-1]
	}

	return col, nil
}

//TODO: this shouldn't mutate the matrix
func (mat *matrix) ScalarMul(val float64) {
	for j := range mat.entries {
		for i := range mat.entries[j] {
			mat.entries[j][i] = mat.entries[j][i] * val
		}
	}
}

func (mat *matrix) Mul(m1 Matrix) error {
	return nil
}

func MatrixMul(m0, m1 Matrix) (Matrix, error) {
	m0m, m0n := m0.Dimensions()
	m1m, m1n := m1.Dimensions()
	if m0n != m1m {
		return nil, ErrCannotMul
	}

	out := make([][]float64, m0m)
	for i := range out {
		out[i] = make([]float64, m1n)
		for j := range out[i] {
			k := 1
			for k <= m0n {
				out[i][j] += MatrixEntryNoErr(m0, i+1, k) * MatrixEntryNoErr(m1, k, j+1)
				k++
			}
		}
	}

	return NewMatrix(m0m, m1n, out...), nil
}

func MatrixScalarMul(m0 Matrix, val float64) Matrix {
	m0m, m0n := m0.Dimensions()
	me := make([][]float64, m0m)

	for j := range me {
		me[j] = make([]float64, m0n)
		for i := range me[j] {
			e, _ := m0.Entry(j+1, i+1)
			me[j][i] = e * val
		}
	}

	return NewMatrix(m0m, m0n, me...)
}

func MatrixEntryNoErr(m0 Matrix, i, j int) float64 {
	v, err := m0.Entry(i, j)
	if err != nil {
		panic("MatrixEntryNoErr errored")
	}
	return v
}
