//--------------------------------------------------------------------------------------------------
//
// Copyright (c) 2018 Denis Dyakov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
// BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
//--------------------------------------------------------------------------------------------------

package aosong

import (
	"encoding/binary"
	"errors"

	i2c "github.com/d2r2/go-i2c"
	"github.com/davecgh/go-spew/spew"
)

// DHT12 sensor memory map
const (
	DHT12_HUM_INT    = 0x00
	DHT12_HUM_SCALE  = 0x01
	DHT12_TEMP_INT   = 0x02
	DHT12_TEMP_SCALE = 0x03
	DHT12_CHECKSUM   = 0x04
	DHT12_DATA_BYTES = 5
	DHT12_DATA_START = DHT12_HUM_INT
)

// SensorDHT12 specific type
type SensorDHT12 struct {
}

// Static cast to verify at compile time
// that type implement interface.
var _ SensorInterface = &SensorDHT12{}

func (v *SensorDHT12) ReadRelativeHumidityAndTemperatureMult10(i2c *i2c.I2C) (humidity int16,
	temperature int16, err error) {

	_, err = i2c.WriteBytes([]byte{DHT12_DATA_START})
	if err != nil {
		return 0, 0, err
	}

	// Construct DHT12 read response
	data := &struct {
		Humidity         byte
		HumidityScale    byte
		Temperature      byte
		TemperatureScale byte
		Checksum         byte
	}{}
	err = readDataToStruct(i2c, DHT12_DATA_BYTES, binary.BigEndian, data)
	if err != nil {
		return 0, 0, err
	}

	calcCrc := byte(data.Humidity + data.HumidityScale +
		data.Temperature + data.TemperatureScale)
	if data.Checksum != calcCrc {
		return 0, 0, errors.New(spew.Sprintf(
			"Checksums doesn't match: CRC from sensor(%v) != calculated CRC(%v)",
			data.Checksum, calcCrc))
	} else {
		lg.Debugf("Checksums verified: CRC from sensor(%v) = calculated CRC(%v)",
			data.Checksum, calcCrc)
	}

	rh := int16(data.Humidity)*10 + int16(data.HumidityScale)
	if rh > 100*10 {
		return -1, -1, spew.Errorf("Humidity value exceed 100%%: %v", float32(humidity)/10)
	}
	temp := int16(data.Temperature)*10 + int16(data.TemperatureScale)

	return rh, temp, nil
}
