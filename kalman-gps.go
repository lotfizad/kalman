package kalman

import (
	"errors"
)

const minAccuracy = 1

type kalmanGps struct {
	timeEpoch                  uint
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

func (gps *kalmanGps) InitState(latitudeMeasured, longitudeMeasured, accuracyMeasured float64, timeEpoch uint) {
	gps.timeEpoch = timeEpoch
	gps.latitude = latitudeMeasured
	gps.longitude = longitudeMeasured
	gps.variance = accuracyMeasured * accuracyMeasured
}

func (kalmanGps *kalmanGps) SinglePointProcess(latitudeMeasured, longitudeMeasured, accuracyMeasured float64, timeEpoch uint) {
	if accuracyMeasured < minAccuracy {
		accuracyMeasured = minAccuracy
	}
	if kalmanGps.variance < 0 {
		panic(errors.New("InitState should be called first"))
	}

	timeEpochIncremental := timeEpoch - kalmanGps.timeEpoch
	if timeEpochIncremental > 0 {
		kalmanGps.variance += float64(timeEpochIncremental) * kalmanGps.averageSpeedMeterPerSecond *
			kalmanGps.averageSpeedMeterPerSecond / 1000.0
		kalmanGps.timeEpoch = timeEpoch
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
	timeEpochs []uint) (latitudeAryFiltered, longitudeAryFiltered []float64) {

	kalmanGps.InitState(latitudeAry[0], longitudeAry[0], accuracyArray[0], timeEpochs[0])

	inputPointsLength := len(latitudeAry)
	if inputPointsLength != len(longitudeAry) || inputPointsLength != len(accuracyArray) || inputPointsLength != len(timeEpochs) {
		panic(errors.New("Length of input arrays should be equal"))
	}

	latitudeAryFiltered = make([]float64, 0, inputPointsLength)
	longitudeAryFiltered = make([]float64, 0, inputPointsLength)

	for i := 1; i < inputPointsLength; i++ {

		kalmanGps.SinglePointProcess(latitudeAry[i], longitudeAry[i], accuracyArray[1], uint(i))
		latitudeAryFiltered = append(latitudeAry, kalmanGps.GetLatitude())
		longitudeAryFiltered = append(longitudeAry, kalmanGps.GetLongitude())
	}

	return
}
