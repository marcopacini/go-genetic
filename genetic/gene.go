/*
 *  MIT License
 *
 *  Copyright (c) 2018 Marco Pacini
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package genetic

import (
	"fmt"
	"math/rand"
)

// A gene is a slice of float64 values in interval [0, 1]
type Gene struct {
	Sequence []float64
}

// NewGene returns a new 'zero' Gene
func NewGene(n int) Gene {
	return Gene{make([]float64, n)}
}

func (g *Gene) Randomize() {
	for i := range g.Sequence {
		g.Sequence[i] = rand.Float64()
	}
}

// Clone the receiver gene
func (g Gene) Clone() Gene {
	gene := Gene{make([]float64, len(g.Sequence))}

	for i := range g.Sequence {
		gene.Sequence[i] = g.Sequence[i]
	}

	return gene
}

func (g Gene) String() string { return fmt.Sprint(g.Sequence) }
