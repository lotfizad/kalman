package kalman

import (
	"math"
	"math/rand"

	kalman "github.com/ryskiwt/go-kalman"

	"github.com/gonum/matrix/mat64"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {

	//
	// kalman filter
	//

	sstd := 0.000001
	ostd := 0.1

	// trend model
	config := kalman.Config{
		F: mat64.NewDense(2, 2, []float64{2, -1, 1, 0}),
		G: mat64.NewDense(2, 1, []float64{1, 0}),
		Q: mat64.NewDense(1, 1, []float64{sstd}),
		H: mat64.NewDense(1, 2, []float64{1, 0}),
		R: mat64.NewDense(1, 1, []float64{ostd}),
	}
	filter, err := kalman.New(&config)
	if err != nil {
		panic(err)
	}

	n := 10000
	s := mat64.NewDense(1, n, nil)
	x, dx := 0.0, 0.01
	xary := make([]float64, 0, n)
	yaryOrig := make([]float64, 0, n)

	for i := 0; i < n; i++ {
		y := math.Sin(x) + 0.1*(rand.NormFloat64()-0.5)
		s.Set(0, i, y)
		x += dx

		xary = append(xary, x)
		yaryOrig = append(yaryOrig, y)
	}

	filtered := filter.Filter(s)
	yaryFilt := mat64.Row(nil, 0, filtered)

	//
	// plot
	//

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(p,
		"Original", generatePoints(xary, yaryOrig),
		"Filtered", generatePoints(xary, yaryFilt),
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(16*vg.Inch, 4*vg.Inch, "sample.png"); err != nil {
		panic(err)
	}
}

func generatePoints(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))

	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	return pts
}
