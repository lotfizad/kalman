package kalman_test

import (
	"github.com/mtlotfizad/kalman"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"testing"
)

func generatePoints(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))

	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	return pts
}

func TestKalmanGps_Process(t *testing.T) {

	x, dx := 0.0, 0.01
	n := 10000
	klm := kalman.New(3.0)

	klm.InitState(x, math.Sin(x), 1.0, 1)

	xary := make([]float64, 0, n)
	yary := make([]float64, 0, n)
	xaryFiltered := make([]float64, 0, n)
	yaryFiltered := make([]float64, 0, n)

	for i := 1; i < n; i++ {
		y := math.Sin(x) + 0.1*(rand.NormFloat64()-0.5)
		x += dx
		xary = append(xary, x)
		yary = append(yary, y)

		klm.ProcessSinglePoint(x, y, 1.0, uint(i))
		xaryFiltered = append(xaryFiltered, klm.GetLatitude())
		yaryFiltered = append(yaryFiltered, klm.GetLongitude())
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(p,
		"Original", generatePoints(xary, yary),
		"Filtered", generatePoints(xaryFiltered, yaryFiltered),
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(16*vg.Inch, 4*vg.Inch, "/tmp/sample.png"); err != nil {
		panic(err)
	}

}
