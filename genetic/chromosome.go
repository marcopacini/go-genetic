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

import "fmt"

// A chromosome collects more genes
type Chromosome struct {
	Genes []Gene
}

// NewChromosome returns a new randomly chromosome
func NewChromosome(length int, geneLength int) Chromosome {
	if length < 1 {
		panic(fmt.Sprintf("invalid argument: length = %d (length must be greater than zero", length))
	}

	chromosome := Chromosome{make([]Gene, length)}

	for i := range chromosome.Genes {
		chromosome.Genes[i] = NewGene(geneLength)
	}

	return chromosome
}

// Clone the receiver chromosome
func (c Chromosome) Clone() Chromosome {
	chromosome := Chromosome{make([]Gene, len(c.Genes))}

	for i := range c.Genes {
		chromosome.Genes[i] = c.Genes[i].Clone()
	}

	return chromosome
}

// Execute mutation on receiver using mutator
func (c *Chromosome) Mutate(mutator Mutator) {
	for i := range c.Genes {
		mutator.Mutate(&c.Genes[i])
	}
}

// Phenotype encapsulates a chromosome with its relative fitness score and age
type Phenotype struct {
	Chromosome
	Fitness float64
	Age     int
}

// sort.Interface implementation for Phenotype slice
type decreasing []Phenotype

func (p decreasing) Len() int {
	return len(p)
}

func (p decreasing) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p decreasing) Less(i, j int) bool {
	return p[i].Fitness > p[j].Fitness
}
