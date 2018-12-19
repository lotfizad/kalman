package kalman

import (
	"errors"
)

const minAccuracy = 1

type kalmanGps struct {
	timeMilliSecond            uint
	latitude                   float64
	longitude                  float64
	variance                   float64
	averageSpeedMeterPerSecond float64
}

func New(averageSpeedMeterPerSecond float64) *kalmanGps {
	kGps := kalmanGps{}
	kGps.averageSpeedMeterPerSecond = averageSpeedMeterPerSecond
	kGps.variance = -1
	return &kGps
}

func (gps *kalmanGps) InitState(latitudeMeasured, longitudeMeasured, accuracyMeasured float64, timeMillisecond uint) {
	gps.timeMilliSecond = timeMillisecond
	gps.latitude = latitudeMeasured
	gps.longitude = longitudeMeasured
	gps.variance = accuracyMeasured * accuracyMeasured
}

func (kalmanGps *kalmanGps) ProcessSinglePoint(latitudeMeasured, longitudeMeasured, accuracyMeasured float64, timeMillisecond uint) {
	if accuracyMeasured < minAccuracy {
		accuracyMeasured = minAccuracy
	}
	if kalmanGps.variance < 0 {
		panic(errors.New("InitState should be called first"))
	}

	timeMillisecondIncremental := timeMillisecond - kalmanGps.timeMilliSecond
	if timeMillisecondIncremental > 0 {
		kalmanGps.variance += float64(timeMillisecondIncremental) * kalmanGps.averageSpeedMeterPerSecond *
			kalmanGps.averageSpeedMeterPerSecond / 1000.0
		kalmanGps.timeMilliSecond = timeMillisecond
	}

	var kalmanGain float64 = kalmanGps.variance / (kalmanGps.variance + accuracyMeasured*accuracyMeasured)

	kalmanGps.latitude += kalmanGain * (latitudeMeasured - kalmanGps.latitude)
	kalmanGps.longitude += kalmanGain * (longitudeMeasured - kalmanGps.longitude)
	kalmanGps.variance = (1 - kalmanGain) * kalmanGps.variance
}

func (gps kalmanGps) GetLatitude() float64 {
	return gps.latitude
}

func (gps kalmanGps) GetLongitude() float64 {
	return gps.longitude
}

func (kalmanGps *kalmanGps) BatchProcess(latitudeAry, longitudeAry, accuracyArray []float64,
	timeMilliseconds []uint) (latitudeAryFiltered, longitudeAryFiltered []float64) {

	inputPointsLength := len(latitudeAry)
	if inputPointsLength != len(longitudeAry) || inputPointsLength != len(accuracyArray) || inputPointsLength != len(timeMilliseconds) {
		panic(errors.New("Length of input arrays should be equal"))
	}

	latitudeAryFiltered = make([]float64, 0, inputPointsLength)
	longitudeAryFiltered = make([]float64, 0, inputPointsLength)

	for i := 1; i < inputPointsLength; i++ {

		kalmanGps.ProcessSinglePoint(latitudeAry[i], longitudeAry[i], accuracyArray[1], uint(i))
		latitudeAryFiltered = append(latitudeAry, kalmanGps.GetLatitude())
		longitudeAryFiltered = append(longitudeAry, kalmanGps.GetLongitude())
	}

	return
}
