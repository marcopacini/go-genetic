package main

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/marcopacini/go-genetic/genetic"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

type Shape int

const (
	circle  Shape = 8
)

type Picture struct {
	genetic.Chromosome
}

func (p Picture) Draw(width int, height int, background color.Gray16) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, image.Rect(0, 0, width, height), &image.Uniform{background}, image.ZP, draw.Src)

	gc := draw2dimg.NewGraphicContext(img)

	for _, gene := range p.Chromosome.Genes {
		drawCircle(gc, gene, width, height)
	}

	return img
}

func drawCircle(gc *draw2dimg.GraphicContext, gene genetic.Gene, width int, height int) {
	if len(gene.Sequence) != int(circle) {
		panic(fmt.Sprintf("Input length is not valid: %d != %d", len(gene.Sequence), circle))
	}

	if gene.Sequence[0] < .5 {
		return
	}

	normalizeUint8 := func(value float64) uint8 {
		return uint8(value * float64(math.MaxUint8))
	}

	nrgba := color.NRGBA{
		normalizeUint8(gene.Sequence[1]),
		normalizeUint8(gene.Sequence[2]),
		normalizeUint8(gene.Sequence[3]),
		normalizeUint8(gene.Sequence[4]),
	}

	x, y := gene.Sequence[5]*float64(width), gene.Sequence[6]*float64(height)
	radius := gene.Sequence[7] * math.Min(float64(width), float64(height)) / 4

	draw2dkit.Circle(gc, x, y, radius)

	gc.SetFillColor(nrgba)
	gc.Fill()
}

func compareImage(img1 image.Image, img2 image.Image) float64 {
	if img1.Bounds().Size() != img2.Bounds().Size() {
		panic(fmt.Sprintf("Images have to be the same size: %v != %v", img1.Bounds().Size(), img2.Bounds().Size()))
	}

	compare := func(c1 color.Color, c2 color.Color) float64 {
		r1, g1, b1, a1 := c1.RGBA()
		r2, g2, b2, a2 := c2.RGBA()

		difference := .0

		squareDiff := func(a, b uint32) float64 {
			math.Abs(float64(a/257) - float64(b/257))
			return math.Pow(float64(a/257)-float64(b/257), 2)
		}

		difference += squareDiff(r1, r2)
		difference += squareDiff(g1, g2)
		difference += squareDiff(b1, b2)
		difference += squareDiff(a1, a2)

		return difference
	}

	var wg sync.WaitGroup

	mutex := &sync.Mutex{}
	result := 0.

	width, height := img1.Bounds().Size().X, img1.Bounds().Size().Y

	for i := 0; i < width; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			difference := 0.
			for j := 0; j < height; j++ {
				difference += compare(img1.At(i, j), img2.At(i, j))
			}

			mutex.Lock()

			result += difference

			mutex.Unlock()
		}(i)
	}

	wg.Wait()

	return result
}

func main() {
	file, err := os.Open("resource/sample-md.png")
	if err != nil {
		panic(err)
	}

	sample, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	width, height := sample.Bounds().Size().X, sample.Bounds().Size().Y

	init := func(e *genetic.Engine) {
		for i := range e.Population {
			for j := range e.Population[i].Chromosome.Genes {
				e.Population[i].Chromosome.Genes[j].Randomize()
				e.Population[i].Chromosome.Genes[j].Sequence[0] = .3 // hide circle
			}
		}

		for i := range e.Population {
			j := rand.Intn(len(e.Population[i].Chromosome.Genes))
			e.Population[i].Chromosome.Genes[j].Sequence[0] = .6
		}
	}

	eval := func(chromosome genetic.Chromosome) float64 {
		difference := compareImage(Picture{chromosome}.Draw(width, height, color.Black), sample)
		return 100 - (difference*100)/(float64(width)*float64(height)*255*255*4)
	}

	var best genetic.Phenotype

	configuration := genetic.Configuration{
		GeneLength:       int(circle),
		ChromosomeLength: 300,
		PopulationSize:   100,
		MaxAge:           5,
		Selection:        genetic.TournamentSelection{10},
		Crossover:        genetic.UniformCrossover{},
		Mutation:         genetic.Gaussian{.001, .1, 0.},
		Elitism:          .1,
		Iterations:       100000000,
		Init:             init,
		Evaluator:        eval,
		Observer: 		  func(i int, e *genetic.Engine) { best = e.Best() },
	}

	engine := genetic.Engine{Configuration: configuration}

	isStarted := false

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if (isStarted) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		go func() {
			_, _ = engine.Start()
		}()
		isStarted = true

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		engine.Stop()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/best", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var img image.Image
		if (isStarted) {
			img = Picture{best.Chromosome}.Draw(width, height, color.Black)
		} else {
			img = Picture{}.Draw(width, height, color.Gray16{})
		}

		if err := png.Encode(w, img); err != nil {
			fmt.Println(err)
		}
	})

	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatal(err)
	}
}