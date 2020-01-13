package main

import (
	"context"
	"encoding/json"
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
	"os/exec"
	"path"
	"runtime"
	"sync"
	"time"
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

type Evolution struct {
	sample image.Image
	running bool
	mutex sync.Mutex
	engine genetic.Engine
	result genetic.Phenotype
}

func newEvolution(img image.Image) *Evolution {
	ev := Evolution{
		sample: img,
		running: false,
		engine: genetic.Engine{},
	}

	init := func(e *genetic.Engine) {
		for i := range e.Population {
			for j := range e.Population[i].Chromosome.Genes {
				e.Population[i].Chromosome.Genes[j].Randomize()
				e.Population[i].Chromosome.Genes[j].Sequence[0] = .3 // hide
			}
		}

		for i := range e.Population {
			j := rand.Intn(len(e.Population[i].Chromosome.Genes))
			e.Population[i].Chromosome.Genes[j].Sequence[0] = .6
		}
	}

	eval := func(c genetic.Chromosome) float64 {
		difference := compareImage(Picture{Chromosome: c}.Draw(img.Bounds().Size().X, img.Bounds().Size().Y, color.Black), img)
		return 100 - (difference*100)/(float64(img.Bounds().Size().X)*float64(img.Bounds().Size().Y)*255*255*4)
	}

	observer := func(i int, e *genetic.Engine) {
		best := e.Best()

		if best.Fitness > ev.result.Fitness {
			ev.result = best
		}
	}

	ev.engine.Configuration = genetic.Configuration{
		GeneLength:       int(circle),
		ChromosomeLength: 150,
		PopulationSize:   75,
		MaxAge:           5,
		Selection:        genetic.TournamentSelection{Size: 10},
		Crossover:        genetic.UniformCrossover{},
		Mutation:         genetic.Gaussian{Probability: .001, Std: .1, Mean: 0.},
		Elitism:          .1,
		Iterations:       int(^uint(0) >> 1), // max int
		Init:             init,
		Evaluator:        eval,
		Observer:         observer,
	}

	return &ev
}

func (ev *Evolution) start() error {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()

	if ev.running {
		return fmt.Errorf("evolution already started")
	}

	go func() {
		_, _ = ev.engine.Start()
	}()
	ev.running = true

	return nil
}

func (ev *Evolution) stop() error {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()

	if !ev.running {
		return fmt.Errorf("evolution not yet started")
	}

	ev.engine.Stop()
	return nil
}

func (ev *Evolution) best() (*genetic.Phenotype, error) {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()

	if !ev.running {
		return nil, fmt.Errorf("evolution not yet started")
	}

	return &ev.result, nil
}

type Sample string

const (
	Small Sample = "sample-sm.png"
	Medium =  "sample-md.png"
	Large =  "sample-lg.png"
)

func getSample(s Sample) (image.Image, error) {
	file, err := os.Open(path.Join("resource", string(s)))
	if err != nil {
		return nil, err
	}

	sample, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func openGUI() error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, "http://localhost:3000")
	return exec.Command(cmd, args...).Start()
}

func main() {
	sample, err := getSample(Small)
	if err != nil {
		panic(err)
	}

	server := &http.Server{Addr: ":3001", Handler: nil}
	evolution := newEvolution(sample)

	if err := openGUI(); err != nil {
		panic(err)
	}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err := evolution.start(); err != nil {
			w.WriteHeader(http.StatusConflict)
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err := evolution.stop(); err != nil {
			log.Println(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/best", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		p, err := evolution.best()
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		img := Picture{Chromosome: p.Chromosome}.Draw(250, 250, color.Black)

		if err := png.Encode(w, img); err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stats := struct {
			IsRunning bool `json:"isRunning"`
		} {
			IsRunning: evolution.running,
		}

		if err := json.NewEncoder(w).Encode(stats); err != nil {
			fmt.Println(err)
		}
	})

	//http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
	//	if err := json.NewEncoder(w).Encode(stats); err != nil {
	//		fmt.Println(err)
	//	}
	//})

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}