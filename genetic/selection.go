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
	"sort"
)

type Selection interface {
	Select([]Phenotype, int) ([]Chromosome, error)
}

type RandomSelection struct{}

// Returns n Chromosomes selected randomly from population
func (r RandomSelection) Select(population []Phenotype, n int) ([]Chromosome, error) {
	if n > len(population) {
		return nil, fmt.Errorf("invalid selection size: %v > %v (population size)", n, len(population))
	}

	selection := make([]Chromosome, n)

	for i := 0; i < n; i++ {
		selection[i] = population[rand.Intn(len(population))].Chromosome
	}

	return selection, nil
}

type ElitismSelection struct {
	Size float64
}

// Returns n Chromosomes selected from the best individuals (elite). If the elite group size is lesser than
// the n, the elite group size ia automatically increased to n
func (e ElitismSelection) Select(population []Phenotype, n int) ([]Chromosome, error) {
	if n > len(population) {
		return nil, fmt.Errorf("invalid selection size: %v > %v (population size)", n, len(population))
	}

	sort.Sort(decreasing(population))

	size := int(bound(float64(n), e.Size*float64(len(population)), float64(len(population))))
	selection := make([]Chromosome, size)

	for i := range population[:size] {
		selection[i] = population[rand.Intn(len(population))].Chromosome
	}

	func(slice []Chromosome) {
		for i := range slice {
			j := rand.Intn(len(slice))
			slice[i], slice[j] = slice[j], slice[i]
		}
	}(selection)

	return selection[:n], nil
}
