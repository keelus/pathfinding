package pair

import (
	"math"
	"strconv"
)

// A Pair is an I, J coordinate pair.
type Pair struct {
	I, J int
}

// String returns a string representation of p like "(3,4)".
func (p Pair) String() string {
	return "(" + strconv.Itoa(p.I) + "," + strconv.Itoa(p.J) + ")"
}

// Add returns the Pair p+q.
func (p Pair) Add(q Pair) Pair {
	return Pair{p.I + q.I, p.J + q.J}
}

// Sub returns the Pair p-q.
func (p Pair) Sub(q Pair) Pair {
	return Pair{p.I - q.I, p.J - q.J}
}

// Dot returns the dot product of p and q.
func (p Pair) Dot(q Pair) int {
	return p.I*q.I + p.J*q.J
}

// Mul returns the Pair p*k.
func (p Pair) Mul(k int) Pair {
	return Pair{p.I * k, p.J * k}
}

// Div returns the Pair p/k.
func (p Pair) Div(k int) Pair {
	return Pair{p.I / k, p.J / k}
}

// Eq reports whether p and q are equal.
func (p Pair) Eq(q Pair) bool {
	return p.I == q.I && p.J == q.J
}

// Zero reports whether the Pair p is (0, 0).
func (p Pair) Zero() bool {
	return p.I == 0 && p.J == 0
}

// Clone returns a copy of the Pair p.
func (p Pair) Copy() Pair {
	return Pair{p.I, p.J}
}

// New is shorthand for Pair{I, J}.
func New(I, J int) Pair {
	return Pair{I, J}
}

// Zero is shorthand for Pair{0, 0}.
func Zero() Pair {
	return Pair{}
}

// Dist returns the euclidean distance p to q.
func (p Pair) Dist(q Pair) float64 {
	return math.Sqrt(math.Pow(float64(q.I-p.I), 2) + math.Pow(float64(q.J-p.J), 2))
}

// MDist returns the Manhattan distance p to q.
func (p Pair) MDist(q Pair) int {
	return int(math.Abs(float64(q.J-p.J)) + math.Abs(float64(q.I-p.I)))
}

// Perp reports whether p and q are perpendicular.
func (p Pair) Perp(q Pair) bool {
	return p.Dot(q) == 0
}

// Neg returns the p vector but negated (-I, -J).
func (p Pair) Neg() Pair {
	return Pair{-p.I, -p.J}
}

// InBounds reports whether p is inside the bounds.
// minI and minJ inclusive. maxI and maxJ exclusive.
func (p Pair) InBounds(minI, minJ, maxI, maxJ int) bool {
	return p.I >= minI && p.I < maxI && p.J >= minI && p.J < maxJ
}

// The following functions use the direction logic based on the directions:
// UP, DOWN, LEFT, RIGHT = (-1, 0), (1, 0), (0, -1), (0, 1)

// Up returns a Pair in the direction UP (-1, 0)
func Up() Pair {
	return Pair{-1, 0}
}

// Down returns a Pair in the direction DOWN (1, 0)
func Down() Pair {
	return Pair{1, 0}
}

// Left returns a Pair in the direction LEFT (0, -1)
func Left() Pair {
	return Pair{0, -1}
}

// Right returns a Pair in the direction RIGHT (0, 1)
func Right() Pair {
	return Pair{0, 1}
}

// TurnL returns the p vector rotated to the left.
func (p Pair) TurnL() Pair {
	return Pair{-p.J, p.I}
}

// TurnL returns the p vector rotated to the right.
func (p Pair) TurnR() Pair {
	return Pair{p.J, -p.I}
}

// Opp gives the opposite direction.
func (p Pair) Opp() Pair {
	return Pair{-p.I, -p.J}
}
