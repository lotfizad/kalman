package kalman_test

import (
	"encoding/csv"
	"fmt"
	"github.com/mtlotfizad/kalman"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

func generatePoints(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))

	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	return pts
}

func TestKalmanGps_SinglePointProcess(t *testing.T) {

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

		klm.SinglePointProcess(x, y, 1.0, uint(i))
		xaryFiltered = append(xaryFiltered, klm.GetLatitude())
		yaryFiltered = append(yaryFiltered, klm.GetLongitude())
	}

	plotOriginFiltered(xary, yary, xaryFiltered, yaryFiltered)

}

func plotOriginFiltered(xary []float64, yary []float64, xaryFiltered []float64, yaryFiltered []float64) {
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

func TestKalmanGps_BatchProcess(t *testing.T) {

	latAry := []float64{35.6999312, 35.6999179, 35.6998761, 35.6999098, 35.699922, 35.6997343,
		35.6999536, 35.6999546, 35.7000039, 35.6999266}
	lngAry := []float64{51.3470473, 51.3471651, 51.3479653, 51.3483119, 51.3489923, 51.3492465, 51.3493252, 51.3496698,
		51.3499808, 51.3501218}
	accuracyAry := []float64{10, 10, 3, 10, 10, 12, 10, 10, 17, 10}
	timesString := []string{"2018-10-17 07:15:10", "2018-10-17 07:15:12", "2018-10-17 07:15:19", "2018-10-17 07:15:28",
		"2018-10-17 07:15:42", "2018-10-17 07:15:48", "2018-10-17 07:15:53", "2018-10-17 07:15:59", "2018-10-17 07:16:03",
		"2018-10-17 07:16:14"}

	timeAry := make([]uint, len(timesString))
	for i := 0; i < len(timesString); i++ {
		gpsTime, _ := time.Parse("2006-01-02 15:04:05", timesString[i])
		timeAry[i] = uint(gpsTime.Unix())
	}

	klm := kalman.New(3.0)

	latitudeAryFiltered, longitudeAryFiltered := klm.BatchProcess(latAry, lngAry, accuracyAry, timeAry)

	file, err := os.Create("/tmp/result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < len(timesString); i++ {
		data := []string{fmt.Sprintf("%.6f", latAry[i]), fmt.Sprintf("%.6f", lngAry[i]),
			fmt.Sprintf("%.6d", timeAry[i]), fmt.Sprintf("%.6f", latitudeAryFiltered[i]),
			fmt.Sprintf("%.6f", longitudeAryFiltered[i])}
		err := writer.Write(data)
		checkError("Cannot write to file", err)
	}

	plotOriginFiltered(latAry, lngAry, latitudeAryFiltered, longitudeAryFiltered)

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
