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

package main

import (
	"fmt"
	"go-genetic/genetic"
)

const alphabet = " abcdefghijklmnopqrtuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,.;:!?"

type Text struct {
	genetic.Chromosome
	s string
}

func (t Text) String() string {
	var s string

	for _, gene := range t.Chromosome.Genes {
		s += string(alphabet[int(gene.Sequence[0]*float64(len(alphabet)))])
	}

	return s
}

func distance(s1, s2 string) (int, error) {
	if len(s1) != len(s2) {
		return 0, fmt.Errorf("invalid strings length: %v != %v", len(s1), len(s2))
	}

	fitness := 0

	for i := range s1 {
		if s1[i] == s2[i] {
			fitness++
		}
	}

	return fitness, nil
}

func main() {
	const text = "Hello, World! cit. MP"

	eval := func(c genetic.Chromosome) float64 {
		if fitness, err := distance(Text{c, ""}.String(), text); err == nil {
			return float64(fitness)
		}

		return 0.
	}

	observer := func(i int, e *genetic.Engine) {
		if int(e.Best().Fitness) == len(text) {
			e.Stop()
		}
	}

	configuration := genetic.Configuration{
		GeneLength:       1,
		ChromosomeLength: len(text),
		PopulationSize:   100,
		Selection:        genetic.ElitismSelection{.1},
		Crossover:        genetic.SinglePointCrossover{},
		Mutation:         genetic.Uniform{.1},
		Elitism:          .1,
		Iterations:       100000,
		Evaluator:        eval,
		Observer:         observer,
	}

	engine := genetic.Engine{Configuration: configuration}
	best, elapsed := engine.Start()

	fmt.Printf("Completed in %v\n", elapsed)
	fmt.Printf("%v\t{%v/%v}\n", Text{best.Chromosome, ""}, best.Fitness, len(text))
}
