package util

import (
	"testing"
)

func TestMatrixCreation(t *testing.T) {
	mat := NewMatrix(3, 2)
	acMat := mat.(*matrix)
	matM, matN := mat.Dimensions()
	if matM != 3 || matN != 2 {
		t.Error("incorrect dimensions of matrix")
	}

	if len(acMat.entries) != 3 {
		t.Error("backing row array not created correctly")
	}
	if len(acMat.entries[2]) != 2 {
		t.Error("backing column array not created correctly")
	}
}

func TestMatrixEntries(t *testing.T) {
	mat := NewMatrix(4, 4, [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}...)

	val, err := mat.Entry(3, 4)
	if err != nil {
		t.Error("Entry() errors on valid address")
	} else if val != 12 {
		t.Error("Entry() returns incorrect value")
	}

	row, err := mat.Row(1)
	if err != nil {
		t.Error("Row() errors on valid row")
	} else {
		out := []float64{1, 2, 3, 4}
		for i := range out {
			if row[i] != out[i] {
				t.Error("Row() returns incorrect row")
				break
			}
		}
	}

	col, err := mat.Column(3)
	if err != nil {
		t.Error("Column() errors on valid column")
	} else {
		out := []float64{3, 7, 11, 15}
		for i := range out {
			if col[i] != out[i] {
				t.Error("Column() returns incorrect column")
				break
			}
		}
	}

}

func TestMatrixTranspose(t *testing.T) {
	mat := NewMatrix(1, 2, []float64{3, 4})
	tMat := mat.Transpose()
	tMatM, tMatN := tMat.Dimensions()
	if tMatM != 2 || tMatN != 1 {
		t.Error("transposed matrix has incorrect dimensions")
	}

	_, err := tMat.Entry(0, 1)
	if err == nil {
		t.Error("transposed matrix allows access of invalid address")
	}

	val, err := tMat.Entry(2, 1)
	if err != nil {
		t.Error("transposed matrix not allowing access of valid address")
	}

	if val != 4 {
		t.Error("incorrect result at transposed matrix (1,0)")
	}
}

func TestMatrixScalarMul(t *testing.T) {
	mat := NewMatrix(3, 2, [][]float64{
		{3, 4},
		{10, -4},
		{0, -1},
	}...)

	mulMat := MatrixScalarMul(mat, 2)

	if v, _ := mulMat.Entry(2, 2); v != -8 {
		t.Errorf("incorrect scalar multiplication of matrix (cloned): %.2f", v)
		t.Log(mulMat)
	}

	mat = MatrixScalarMul(mat, 3)
	if v, _ := mat.Entry(2, 2); v != -12 {
		t.Error("incorrect scalar multiplication of matrix (in place)")
	}

}

func TestMatrixMul(t *testing.T) {
	mat := NewMatrix(2, 3, []float64{2, 1, 4}, []float64{9, -10})
	mulMat := NewMatrix(3, 2, []float64{-1, 30}, []float64{2, 3}, []float64{-4, -0})
	invalidMulMat := NewMatrix(4, 5)

	if _, err := MatrixMul(mat, invalidMulMat); err == nil {
		t.Error("invalid multiplication does not error")
	}

	tMat, err := MatrixMul(mat, mulMat)
	if err != nil {
		t.Error("valid multiplication errors")
	}

	matM, _ := mat.Dimensions()
	_, mulMatN := mulMat.Dimensions()

	if tMatM, tMatN := tMat.Dimensions(); tMatM != matM || tMatN != mulMatN {
		t.Error("matrix multiplication result has invalid dimensions")
		t.Log(tMatM, tMatN)
		return
	}

	if val, err := tMat.Entry(2, 2); err != nil {
		t.Errorf("error for multipled matrix getting entry 2,2\n%v", tMat)
	} else if val != 240 {
		t.Errorf("incorrect value for entry 2,2 in multiplied matrix\n%v", tMat)
	}
}
