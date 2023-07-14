package tvd_test

import (
	"image/color"
	"math/rand"
	"testing"

	"github.com/soypat/tvd"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
)

func TestDenoisingWithGonumGraphVisualization(t *testing.T) {
	// Generate noisy data
	const (
		width, height = 19 * font.Centimeter, 10 * font.Centimeter
		xstart, xmax  = 0, 10
		div           = 10000
		lambda        = 32
		noiseAmp      = 4
	)
	randf := rand.New(rand.NewSource(1)).Float64

	fn := func(x float64) (y float64) {
		switch {
		case x < xmax/6:
			y = 4
		case x < 2*xmax/3:
			y = 2
		default:
			y = x
		}
		return y + noiseAmp*randf()/2
	}
	var noisy plotter.XYs
	noiseSlice := make([]float64, div)
	for i := 0; i < div; i++ {
		x := xstart + xmax*float64(i)/float64(div)
		noiseSlice[i] = fn(x)
		noisy = append(noisy, plotter.XY{X: x, Y: noiseSlice[i]})
	}
	denoiseSlice := tvd.Denoise1D(noiseSlice, lambda)

	var denoised plotter.XYs
	for i := 0; i < div; i++ {
		x := xstart + xmax*float64(i)/float64(div)
		denoised = append(denoised, plotter.XY{X: x, Y: denoiseSlice[i]})
	}

	p := plot.New()
	p.Title.Text = "Total variation denoising"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"
	plttr1, _ := plotter.NewLine(noisy)
	plttr2, _ := plotter.NewLine(denoised)
	plttr2.Color = color.RGBA{R: 255, A: 255}
	plttr1.Color = color.RGBA{B: 255, A: 255}
	p.Add(plttr1, plttr2)
	p.Legend.Add("Noisy", plttr1)
	p.Legend.Add("Denoised", plttr2)
	p.Save(width, height, "tvd1d_test.png")
}
