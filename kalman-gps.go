package kalman

const minAccuracy = 1

type KalmanGps struct {
	metresPerSecond float64
	timeMilliSecond uint
	latitude        float64
	longitude       float64

	variance                   float64
	averageSpeedMeterPerSecond uint
}

func (kalmanGps *KalmanGps) Init(latitudeMeasured, longitudeMeasured, accuracyMeasured float64, timeMillisecond uint,
	averageSpeedMeterPerSecond uint) {
	kalmanGps.timeMilliSecond = timeMillisecond
	kalmanGps.latitude = latitudeMeasured
	kalmanGps.longitude = longitudeMeasured
	kalmanGps.variance = accuracyMeasured * accuracyMeasured
	kalmanGps.averageSpeedMeterPerSecond = averageSpeedMeterPerSecond
}

func (kalmanGps *KalmanGps) Process(latitudeMeasured, longitudeMeasured, accuracyMeasured float64, timeMillisecond uint) {
	if accuracyMeasured < minAccuracy {
		accuracyMeasured = minAccuracy
	}

	timeMillisecondIncremental := timeMillisecond - kalmanGps.timeMilliSecond
	if timeMillisecondIncremental > 0 {
		kalmanGps.variance += float64(timeMillisecondIncremental * kalmanGps.averageSpeedMeterPerSecond *
			kalmanGps.averageSpeedMeterPerSecond / 1000)
		kalmanGps.timeMilliSecond = timeMillisecond
	}

	var kalmanGain float64 = kalmanGps.variance / (kalmanGps.variance + accuracyMeasured*accuracyMeasured)

	kalmanGps.latitude += kalmanGain * (latitudeMeasured - kalmanGps.latitude)
	kalmanGps.longitude += kalmanGain * (longitudeMeasured - kalmanGps.longitude)
}
