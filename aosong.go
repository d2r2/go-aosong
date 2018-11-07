package aosong

import (
	i2c "github.com/d2r2/go-i2c"
)

// SensorType identify which Aosong Electronics
// humidity and temperature sensor is used.
// DHT12, AM2320 are supported.
type SensorType int

// Implement Stringer interface.
func (v SensorType) String() string {
	if v == DHT12 {
		return "DHT12"
	} else if v == AM2320 {
		return "AM2320"
	} else {
		return "!!! unknown !!!"
	}
}

const (
	// Aosong Electronics humidity and temperature sensor model DHT12.
	DHT12 SensorType = iota
	// Aosong Electronics humidity and temperature sensor model AM2320.
	AM2320
)

// Abstract Aosong Electronics sensor interface
// to control and gather data via I2C-bus.
type SensorInterface interface {
	ReadRelativeHumidityAndTemperatureMult10(i2c *i2c.I2C) (humidity int16, temperature int16, err error)
	//ReadTemperatureMult10(i2c *i2c.I2C) (int32, error)
	//ReadRelativeHumidityMult10(i2c *i2c.I2C) (int32, error)
}

type Sensor struct {
	sensorType SensorType
	sensor     SensorInterface
}

func NewSensor(sensorType SensorType) *Sensor {
	v := &Sensor{sensorType: sensorType}
	switch sensorType {
	case AM2320:
		v.sensor = &SensorAM2320{}
	case DHT12:
		v.sensor = &SensorDHT12{}
	}

	return v
}

func (v *Sensor) GetSensorType() SensorType {
	return v.sensorType
}

func (v *Sensor) ReadRelativeHumidityAndTemperature(i2c *i2c.I2C) (humidity float32,
	temperature float32, err error) {
	rh, temp, err := v.sensor.ReadRelativeHumidityAndTemperatureMult10(i2c)
	if err != nil {
		return 0, 0, err
	}
	rhf, tempf := float32(rh)/10, float32(temp)/10
	return rhf, tempf, nil
}
