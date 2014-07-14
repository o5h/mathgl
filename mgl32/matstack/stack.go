package matstack

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// A matrix stack is a linear fully-persistent data structure of matrix multiplications
// Each push to a MatStack multiplies the current top of the stack with thew new matrix
// and appends it to the top. Each pop undoes the previous multiplication.
//
// This is extremely useful for scenegraphs, where you can push the transformation of the current
// object for children to use, and pop the transformation before returning to the parent.
type MatStack []mgl32.Mat4

// Returns a matrix stack where the top element is the identity.
func NewMatStack() *MatStack {
	ms := make(MatStack, 1)
	ms[0] = mgl32.Ident4()

	return &ms
}

// Multiplies the current top matrix by m, and pushes the result
// on the stack.
func (ms *MatStack) Push(m mgl32.Mat4) {
	prev := (*ms)[len(*ms)-1]
	(*ms) = append(*ms, prev.Mul4(m))
}

// Pops the current matrix off the top of the stack and returns it.
// If the matrix stack only has one element left, this will return an error.
func (ms *MatStack) Pop() (mgl32.Mat4, error) {
	if len(*ms) == 1 {
		return mgl32.Mat4{}, errors.New("attempt to pop last element of the stack; Matrix Stack must have at least one element")
	}

	retVal := (*ms)[len(*ms)-1]

	(*ms) = (*ms)[:len(*ms)-1]

	return retVal, nil
}

// Returns the value of the current top element of the stack, without
// removing it.
func (ms *MatStack) Peek() mgl32.Mat4 {
	return (*ms)[len(*ms)-1]
}

// Returns the size of the matrix stack. This value will never be less
// than 1.
func (ms *MatStack) Len() int {
	return len(*ms)
}

// This cuts down the matrix as if Pop had been called n times. If n would
// bring the matrix down below 1 element, this does nothing and returns an error.
func (ms *MatStack) Unwind(n int) error {
	if n > len(*ms)-1 {
		return errors.New("Cannot unwind a matrix to below 1 value")
	}

	(*ms) = (*ms)[:len(*ms)-n]
	return nil
}

// Copy will create a new "branch" of the current matrix stack,
// the copy will contain all elements of the current stack in a new stack. Changes to
// one will never affect the other.
func (ms *MatStack) Copy() *MatStack {
	v := append(MatStack{}, (*ms)...)
	return &v
}

// Rebase is tricky. It attempts to seed an arbitrary point in the matrix and replay all transformations
// as if that point in the push at that point had been change instead of the original value.
// The matrix stack does NOT keep track of arguments so this is done via consecutive inverses.
// If the inverse of element i can be found, we can calculate the transformation that was given at point i+1.
// This transformation can then be multiplied by the NEW matrix at point i to complete the "what if".
// If no such inverse can be found at any given point along the rebase, it will be aborted, and the original
// stack will NOT be visibly affected. The error returned will be of type NoInverseError.
//
// If n is out of bounds (n <= 0 || n >= len(*ms)), a generic error from the errors package will be returned.
//
// If you have the old transformations retained, it is recommended
// that you use Unwind followed by Push(change) and then further calling Push for each transformation. Rebase is
// imprecise by nature, and sometimes impossible. It's also expensive due to the inverse calculation at each point.
func (ms *MatStack) Rebase(n int, change mgl32.Mat4) error {
	if n >= len(*ms) || n <= 0 {
		return errors.New("Cannot rebase at the given point on the stack, it is out of bounds.")
	}

	backup := []mgl32.Mat4((*ms)[n:])
	backup = append([]mgl32.Mat4{}, backup...) // copy into new slice

	curr := (*ms)[n]
	(*ms)[n] = (*ms)[n-1].Mul4(change)

	for i := n + 1; i < len(*ms); i++ {
		inv := curr.Inv()

		blank := mgl32.Mat4{}
		if inv == blank {
			ms.undoRebase(n, backup)
			return NoInverseError{Loc: i - 1, Mat: curr}
		}

		ghost := inv.Mul4((*ms)[i])

		curr = (*ms)[i]
		(*ms)[i] = ghost.Mul4((*ms)[i-1])
	}

	return nil
}

func (ms *MatStack) undoRebase(n int, prev []mgl32.Mat4) {
	for i := n; i < len(*ms); i++ {
		(*ms)[i] = prev[i-n]
	}
}

// A NoInverseError is returned on rebase when an inverse cannot be found along the chain,
// due to a transformation projecting the matrix into a singularity. The values include the matrix
// no inverse can be found for, and the location of that matrix.
type NoInverseError struct {
	Mat mgl32.Mat4
	Loc int
}

func (nie NoInverseError) Error() string {
	return fmt.Sprintf("cannot find inverse of matrix %v at location %d in matrix stack, aborting rebase")
}
