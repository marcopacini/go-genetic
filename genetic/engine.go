package genetic

import (
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Configuration struct {
	GeneLength       int
	ChromosomeLength int
	PopulationSize   int
	Selection
	Crossover
	Mutation   Mutator
	Elitism    float64
	Iterations int
	Evaluator  func(Chromosome) float64
	Observer   func(int, *Engine)
}

type Engine struct {
	Configuration
	Population []Phenotype
	running    bool
}

func (e *Engine) Start() (Phenotype, time.Duration) {
	rand.Seed(time.Now().UnixNano())

	e.Population = make([]Phenotype, e.PopulationSize)

	for i := range e.Population {
		chromosome := NewChromosome(e.ChromosomeLength, e.GeneLength)
		e.Population[i] = Phenotype{chromosome, e.Evaluator(chromosome)}
	}

	e.running = true
	start := time.Now()

	for i := 0; e.running && i < e.Iterations; i++ {
		sort.Sort(decreasing(e.Population))

		// Elitism
		survivors := int(bound(0., e.Elitism*float64(e.PopulationSize), float64(e.PopulationSize)))

		var wg sync.WaitGroup
		mutex := &sync.Mutex{}

		offspring := e.Population[:survivors]

		for survivors < e.PopulationSize {
			wg.Add(1)
			survivors += e.Crossover.Children()

			go func() {
				defer wg.Done()

				parents, err := e.Selection.Select(e.Population, e.Crossover.Children())
				if err != nil {
					panic(err)
				}

				children, err := e.Crossover.Cross(parents)
				if err != nil {
					panic(err)
				}

				for i := range children {
					children[i].Mutate(e.Mutation)

					phenotype := Phenotype{children[i], e.Evaluator(children[i])}

					mutex.Lock()
					offspring = append(offspring, phenotype)
					mutex.Unlock()
				}
			}()
		}

		wg.Wait()
		e.Population = offspring[:e.PopulationSize]

		if e.Observer != nil {
			e.Observer(i, e)
		}
	}

	return e.Best(), time.Since(start)
}

func (e *Engine) Stop() {
	e.running = false
}

func (e *Engine) Best() Phenotype {
	best := e.Population[0]

	for i := range e.Population[1:] {
		if e.Population[i].Fitness > best.Fitness {
			best = e.Population[i]
		}
	}

	return best
}

func (e *Engine) Worst() Phenotype {
	worst := e.Population[0]

	for i := range e.Population[1:] {
		if e.Population[i].Fitness < worst.Fitness {
			worst = e.Population[i]
		}
	}

	return worst
}
