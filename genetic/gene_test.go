package genetic

import (
	"math/rand"
	"testing"
)

func TestNewGene(t *testing.T) {
	var table = []struct {
		Input int
		Want  Gene
		Panic bool
	}{
		{-rand.Int(), Gene{make([]float64, 1)}, true},
		{0, Gene{make([]float64, 1)}, true},
		{1, Gene{make([]float64, 1)}, false},
	}

	for _, test := range table {
		defer func() {
			if test.Panic && recover() != nil {
				t.Errorf("NewGene(%d) recover != nil, want panic!", test.Input)
			}
		}()

		gene := NewGene(test.Input)

		if len(gene.Sequence) != len(test.Want.Sequence) {
			t.Errorf("NewGene(%d) length = %d, want %d", test.Input, len(gene.Sequence), len(test.Want.Sequence))
		}

		for i := range gene.Sequence {
			if gene.Sequence[i] != test.Want.Sequence[i] {
				t.Errorf("NewGene(%d).Sequence[%d] = %f, want %d", test.Input, i, test.Want.Sequence[i])
			}
		}
	}
}

func TestGene_Clone(t *testing.T) {
	gene := NewGene(rand.Int())
	clone := gene.Clone()

	if &gene == &clone {
		t.Errorf("&gene == &gene.Clone(), want differents addresses")
	}

	for i := range gene.Sequence {
		if &gene.Sequence[i] == &clone.Sequence[i] {
			t.Errorf("&gene.Sequence[%d] = &gene.Clone().Sequence[%d] = %v, want differents addresses", i, i, &gene.Sequence[i])
		}

		if gene.Sequence[i] == clone.Sequence[i] {
			t.Errorf("gene.Clone().Sequence[%d] = %f, want %f", i, clone.Sequence[i], gene.Sequence[i])
		}
	}
}
