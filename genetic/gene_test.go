package genetic

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestGene_Clone(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	gene := NewGene(rand.Intn(int(math.MaxUint8)))
	clone := gene.Clone()

	if &gene == &clone {
		t.Errorf("&gene == &gene.Clone(), expected differents addresses")
	}

	for i := range gene.Sequence {
		if &gene.Sequence[i] == &clone.Sequence[i] {
			t.Errorf("&gene.Sequence[%d] == &gene.Clone().Sequence[%d], want differents addresses", i, i)
		}

		if gene.Sequence[i] != clone.Sequence[i] {
			t.Errorf("gene.Clone().Sequence[%d] = %f, want %f", i, clone.Sequence[i], gene.Sequence[i])
		}
	}
}
