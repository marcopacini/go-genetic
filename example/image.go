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
	"flag"
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"go-genetic/genetic"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	circle = 7
)

type Picture struct {
	genetic.Chromosome
}

func (p Picture) Draw(width int, height int, background color.Gray16) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, image.Rect(0, 0, width, height), &image.Uniform{background}, image.ZP, draw.Src)

	gc := draw2dimg.NewGraphicContext(img)

	for _, gene := range p.Chromosome.Genes {
		DrawCircle(gc, gene, width, height)
	}

	return img
}

func DrawCircle(gc *draw2dimg.GraphicContext, gene genetic.Gene, width int, height int) {
	if len(gene.Sequence) != int(circle) {
		panic(fmt.Sprintf("Input length is not valid: %d != %d", len(gene.Sequence), circle))
	}

	normalizeUint8 := func(value float64) uint8 {
		return uint8(value * float64(math.MaxUint8))
	}

	nrgba := color.NRGBA{
		normalizeUint8(gene.Sequence[0]),
		normalizeUint8(gene.Sequence[1]),
		normalizeUint8(gene.Sequence[2]),
		normalizeUint8(gene.Sequence[3]),
	}

	x, y := gene.Sequence[4]*float64(width), gene.Sequence[5]*float64(height)
	radius := gene.Sequence[6] * math.Min(float64(width), float64(height)) / 4

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
	var path string
	var iterations int
	var verbose bool

	flag.IntVar(&iterations, "n", int(^uint(0)>>1), "number of iterations")
	flag.BoolVar(&verbose, "v", false, "verbose")

	flag.Usage = func() {
		fmt.Println("Usage: image options target_image.png")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	} else {
		path = flag.Args()[0]
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	sample, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	width, height := sample.Bounds().Size().X, sample.Bounds().Size().Y

	eval := func(chromosome genetic.Chromosome) float64 {
		difference := compareImage(Picture{chromosome}.Draw(width, height, color.White), sample)
		return 100 - (difference*100)/(float64(width)*float64(height)*256*256*4)
	}

	observer := func(i int, e *genetic.Engine) {
		if verbose {
			fmt.Printf("%d - best: %f - worst: %f\n", i, e.Best().Fitness, e.Worst().Fitness)
		}
	}

	configuration := genetic.Configuration{
		GeneLength:       circle,
		ChromosomeLength: 100,
		PopulationSize:   50,
		Selection:        genetic.ElitismSelection{.1},
		Crossover:        genetic.UniformCrossover{},
		Mutation:         genetic.Gaussian{.001, .1, 0.},
		Elitism:          .1,
		Iterations:       iterations,
		Evaluator:        eval,
		Observer:         observer,
	}

	if verbose {
		fmt.Println(configuration)
	}

	engine := genetic.Engine{Configuration: configuration}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		engine.Stop()
	}()

	best, elapsed := engine.Start()

	if verbose {
		fmt.Printf("Completed in %v\n", elapsed)
	}

	img := Picture{best.Chromosome}.Draw(4*width, 4*height, color.White)
	draw2dimg.SaveToPngFile("monna-lisa.png", img)
}
