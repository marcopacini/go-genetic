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

import "math/rand"

type Mutator interface {
	Mutate(*Gene)
}

func bound(lower, value, upper float64) float64 {
	if value < lower {
		return lower
	}

	if value > upper {
		return upper
	}

	return value
}

type Uniform struct {
	Probability float64
}

// Apply uniform mutation on gene
func (u Uniform) Mutate(gene *Gene) {
	for i := range gene.Sequence {
		if rand.Float64() < u.Probability {
			gene.Sequence[i] = rand.Float64()
			// gene.Sequence[i] = bound(0., gene.Sequence[i]+(2.*rand.Float64()-1.)*u.Magnitude, 1.)
		}
	}
}

type Gaussian struct {
	Probability, Std, Mean float64
}

func (g Gaussian) Mutate(gene *Gene) {
	for i := range gene.Sequence {
		if rand.Float64() < g.Probability {
			gene.Sequence[i] = bound(0., gene.Sequence[i]+(rand.NormFloat64()*g.Std+g.Mean), 1.)
		}
	}
}
