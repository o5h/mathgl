package mathgl_test

import (
	"github.com/Jragonmiris/mathgl"
	"math"
	"testing"
)

func TestVecSetGet(t *testing.T) {
	v1 := mathgl.NewVector(mathgl.FLOAT32)
	if v1 == nil {
		t.Fatalf("Failed to create new vector")
	}

	input := []mathgl.Scalar{mathgl.ScalarInt32(1), mathgl.ScalarFloat32(2), mathgl.ScalarFloat32(3), mathgl.ScalarFloat32(4), mathgl.ScalarFloat32(5)}

	if err := v1.AddElements(input); err == nil {
		t.Errorf("Added list with bad element to vector")
	}

	input[0] = mathgl.ScalarFloat32(1)

	if err := v1.AddElements(input); err != nil {
		t.Fatalf("Failed to set vector with correct list")
	}

	if a := v1.GetElement(0); a == nil || float32(a.(mathgl.ScalarFloat32))-float32(1.) > 0.000001 {
		t.Errorf("Didn't get/set correct element of vector")
	}

	v2, e := mathgl.VectorOf(mathgl.FLOAT32, input)
	if v2 == nil || e != nil {
		t.Fatalf("VectorOf failed on good input")
	}

	if a, b := v1.GetElement(2), v2.GetElement(2); a == nil || b == nil || float32(a.(mathgl.ScalarFloat32))-float32(b.(mathgl.ScalarFloat32)) > .000000001 {
		t.Errorf("Two vectors not the same despite being made from same list")
	}

	if !v1.Equal(*v2) { // We should have checked if this was equal in the last step. So if this fails equal is PROBABLY bad
		t.Errorf("Vectors are not equal or equal function failed, v1: %v v2: %v", v1, v2)
	}

	if err := v1.SetElement(25, mathgl.ScalarFloat32(1)); err == nil {
		t.Errorf("Set out of bounds vector element")
	}

	if err := v1.SetElement(-3, mathgl.ScalarFloat32(1)); err == nil {
		t.Errorf("Set out of bounds vector element")
	}

	if err := v1.SetElement(4, mathgl.ScalarFloat32(42)); err != nil {
		t.Fatalf("Didn't set in-bounds vector element")
	}

	if a := v1.GetElement(4); math.Abs(float64(float32(a.(mathgl.ScalarFloat32))-float32(42))) > .0000001 {
		t.Errorf("Did not correctly set single-in bounds vector element")
	}

	if a := v2.GetElement(4); math.Abs(float64(float32(a.(mathgl.ScalarFloat32))-float32(42))) < .0000001 {
		t.Errorf("Changing one vector changed another")
	}

	if v1.Equal(*v2) {
		t.Errorf("Vectors are equal despite changing v1, or equal is wrong, v1: %v v2: %v", v1, v2)
	}

	//v3,_ := mathgl.VectorOf(mathgl.INT32, []mathgl.Scalar{mathgl.VecInt32(1),mathgl.VecInt32(2)})

}

func TestMatrixCreation(t *testing.T) {
	iden2 := mathgl.Identity(2, mathgl.FLOAT64)
	for i, el := range iden2.AsSlice() {
		if (i == 0 || i == 3) && math.Abs(float64(el.(mathgl.ScalarFloat64))-float64(1)) > .000001 {
			t.Errorf("Diagonals not 1 in 2x2 identity el: %v", el)
		} else if (i == 1 || i == 2) && math.Abs(float64(el.(mathgl.ScalarFloat64))-float64(0)) > .000001 {
			t.Errorf("Off-diagonals not 0 in 2x2 identity el: %v", el)
		}
	}
}
