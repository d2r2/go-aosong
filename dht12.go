package aosong

import (
	"bytes"
	"encoding/binary"
	"errors"

	i2c "github.com/d2r2/go-i2c"
	"github.com/davecgh/go-spew/spew"
)

// SensorDHT12 specific type
type SensorDHT12 struct {
}

// Static cast to verify at compile time
// that type implement interface.
var _ SensorInterface = &SensorDHT12{}

// DHT12 sensor read responce according to specification.
type rawDHT12Responce struct {
	Humidity         byte
	HumidityScale    byte
	Temperature      byte
	TemperatureScale byte
	Checksum         byte
}

func (v *SensorDHT12) ReadRelativeHumidityAndTemperatureMult10(i2c *i2c.I2C) (humidity int16,
	temperature int16, err error) {
	// read 1 byte to wake up sensor
	// never check error, since one will be every time
	const bytesExpected = 5
	buf1, _, err := i2c.ReadRegBytes(0, bytesExpected)
	if err != nil {
		return 0, 0, err
	}

	resp := &rawDHT12Responce{}
	err = binary.Read(bytes.NewBuffer(buf1), binary.BigEndian, resp)
	if err != nil {
		return 0, 0, err
	}

	calcCrc := byte(resp.Humidity + resp.HumidityScale +
		resp.Temperature + resp.TemperatureScale)
	if resp.Checksum != calcCrc {
		return 0, 0, errors.New(spew.Sprintf(
			"Checksums doesn't match: CRC from sensor(%v) != calculated CRC(%v)",
			resp.Checksum, calcCrc))
	} else {
		lg.Debugf("Checksums verified: CRC from sensor(%v) = calculated CRC(%v)",
			resp.Checksum, calcCrc)
	}

	rh := int16(resp.Humidity)*10 + int16(resp.HumidityScale)
	if rh > 100*10 {
		return -1, -1, spew.Errorf("Humidity value exceed 100%%: %v", float32(humidity)/10)
	}
	temp := int16(resp.Temperature)*10 + int16(resp.TemperatureScale)
	return rh, temp, nil
}
