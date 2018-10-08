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

type Crossover interface {
	Cross([]Chromosome) ([]Chromosome, error)
	Children() int
}

type None struct{}

func (n None) Cross(parents []Chromosome) ([]Chromosome, error) {
	var children []Chromosome

	for _, p := range parents {
		children = append(children, p.Clone())
	}

	return children, nil
}

func (n None) Children() int {
	return 1
}

const alphabet = "abcdefghijklmnopqrtuvwxyz"

type Text struct {
	Chromosome
}

func (t Text) String() string {
	var s string

	for _, gene := range t.Chromosome.Genes {
		s += string(alphabet[int(gene.Sequence[0]*float64(len(alphabet)-1))])
	}

	return s
}

type SinglePointCrossover struct{}

// Return two child applying single point crossover on parents
// https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)#Single-point_crossover
func (s SinglePointCrossover) Cross(parents []Chromosome) ([]Chromosome, error) {
	if len(parents) != 2 {
		return nil, fmt.Errorf("invalid parents number: %v != 2", len(parents))
	}

	mother, father := parents[0], parents[1]
	pivot := rand.Intn(len(mother.Genes))

	children := make([]Chromosome, 2)

	children[0] = NewChromosome(len(mother.Genes), len(mother.Genes[0].Sequence))
	children[1] = NewChromosome(len(mother.Genes), len(mother.Genes[0].Sequence))

	for i := 0; i < pivot; i++ {
		children[0].Genes[i] = mother.Genes[i].Clone()
		children[1].Genes[i] = father.Genes[i].Clone()
	}

	for i := pivot; i < len(mother.Genes); i++ {
		children[0].Genes[i] = father.Genes[i].Clone()
		children[1].Genes[i] = mother.Genes[i].Clone()
	}

	return children, nil
}

func (s SinglePointCrossover) Children() int {
	return 2
}

type UniformCrossover struct{}

func (u UniformCrossover) Cross(parents []Chromosome) ([]Chromosome, error) {
	if len(parents) != 2 {
		return nil, fmt.Errorf("invalid parents number: %v != 2", len(parents))
	}

	mother, father := parents[0], parents[1]

	children := make([]Chromosome, 2)

	children[0] = NewChromosome(len(mother.Genes), len(mother.Genes[0].Sequence))
	children[1] = NewChromosome(len(mother.Genes), len(mother.Genes[0].Sequence))

	for i := 0; i < len(mother.Genes); i++ {
		if rand.Float64() < .5 {
			children[0].Genes[i] = mother.Genes[i].Clone()
			children[1].Genes[i] = father.Genes[i].Clone()
		} else {
			children[0].Genes[i] = father.Genes[i].Clone()
			children[1].Genes[i] = mother.Genes[i].Clone()
		}
	}

	return children, nil
}

func (u UniformCrossover) Children() int {
	return 2
}
