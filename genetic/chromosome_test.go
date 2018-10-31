package genetic

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestChromosome_Clone(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	chromosome := NewChromosome(rand.Intn(int(math.MaxUint8)), rand.Intn(int(math.MaxUint8)))
	clone := chromosome.Clone()

	if &chromosome == &clone {
		t.Errorf("&chromosome == &chromosome.Clone(), expected differents addresses")
	}

	for i := range chromosome.Genes {
		if &chromosome.Genes[i] == &clone.Genes[i] {
			t.Errorf("&chromosome.Genes[%d] == &chromosome.Clone().Genes[%d], expected differents addresses", i, i)
		}
	}
}
