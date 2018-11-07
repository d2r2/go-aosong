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
	"time"

	i2c "github.com/d2r2/go-i2c"
	"github.com/davecgh/go-spew/spew"
)

// Command byte's sequences
const (
	CMD_AM2320_READ_REGISTERS  byte = 0x03 // Reading register data
	CMD_AM2320_WRITE_REGISTERS byte = 0x10 // Write  multiple registers
)

// SensorAM2320 specific type
type SensorAM2320 struct {
}

// Static cast to verify at compile time
// that type implement interface.
var _ SensorInterface = &SensorAM2320{}

func (v *SensorAM2320) ReadRelativeHumidityAndTemperatureMult10(i2c *i2c.I2C) (humidity int16,
	temperature int16, err error) {
	// Ping sensor: try to read 1 byte to wake up sensor.
	// Never check up error here, since one will be ever
	buf1 := make([]byte, 1)
	_, _ = i2c.ReadBytes(buf1)
	// Send command to read registers
	const startRegAddr = 0
	const dataBytesCount = 4 // Maximum 32 bytes of registers
	_, err = i2c.WriteBytes([]byte{CMD_AM2320_READ_REGISTERS,
		startRegAddr, dataBytesCount})
	if err != nil {
		return 0, 0, err
	}
	// Wait 3 millisecond according to specification
	time.Sleep(time.Millisecond * 3)
	// Read register's results
	const responsePrefixBytesCount = 2
	const crcBytesCount = 2
	buf2 := make([]byte, responsePrefixBytesCount+
		dataBytesCount+crcBytesCount)
	_, err = i2c.ReadBytes(buf2)
	if err != nil {
		return 0, 0, err
	}
	// Construct AM2320 read response
	data := &struct {
		FunctionCode byte
		BytesCount   byte
		Data         [dataBytesCount]byte
		CRC1         byte
		CRC2         byte
	}{}
	err = readDataToStruct(i2c, responsePrefixBytesCount+
		dataBytesCount+crcBytesCount, binary.BigEndian, data)
	if err != nil {
		return 0, 0, err
	}

	rh := getS16BE(data.Data[0:2])
	if rh > 100*10 {
		return -1, -1, spew.Errorf("humidity value exceed 100%%: %v", humidity)
	}
	temp := getS16BE(data.Data[2:4])
	var crc uint16 = getU16LE([]byte{data.CRC1, data.CRC2})
	crcBuf := append([]byte{data.FunctionCode, data.BytesCount},
		data.Data[0:dataBytesCount]...)
	calcCrc := calcCRC_AM2320(crcBuf)
	if crc != calcCrc {
		err := errors.New(spew.Sprintf(
			"CRCs doesn't match: CRC from sensor(%v) != calculated CRC(%v)",
			crc, calcCrc))
		return 0, 0, err
	} else {
		lg.Debugf("CRCs verified: CRC from sensor(%v) = calculated CRC(%v)",
			crc, calcCrc)
	}
	if rh > 100*10 {
		return -1, -1, spew.Errorf("humidity value exceed 100%%: %v", float32(humidity)/10)
	}

	return rh, temp, nil
}
