package main

import (
	aosong "github.com/d2r2/go-aosong"
	i2c "github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)

func main() {
	defer logger.FinalizeLogger()
	// Create new connection to i2c-bus on 1 line with address 0x5C.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x5c, 1)
	if err != nil {
		lg.Fatal(err)
	}
	defer i2c.Close()

	// Uncomment/comment next lines to suppress/increase verbosity of output
	// logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	// logger.ChangePackageLogLevel("aosong", logger.InfoLevel)

	// sensor := aosong.NewSensor(aosong.AM2320)
	sensor := aosong.NewSensor(aosong.DHT12)
	lg.Infof("Sensor type = %v", sensor.GetSensorType())

	rh, t, err := sensor.ReadRelativeHumidityAndTemperature(i2c)
	if err != nil {
		lg.Fatal(err)
	}

	lg.Infof("Relative humidity = %v%%", rh)
	lg.Infof("Temprature in celsius = %v*C", t)
}
