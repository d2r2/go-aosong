package aosong

// Utility functions

// getS16BE extract 2-byte integer as signed big-endian.
func getS16BE(buf []byte) int16 {
	v := int16(buf[0])<<8 + int16(buf[1])
	return v
}

// getS16LE extract 2-byte integer as signed little-endian.
func getS16LE(buf []byte) int16 {
	w := getS16BE(buf)
	// exchange bytes
	v := (w&0xFF)<<8 + w>>8
	return v
}

// getU16BE extract 2-byte integer as unsigned big-endian.
func getU16BE(buf []byte) uint16 {
	v := uint16(buf[0])<<8 + uint16(buf[1])
	return v
}

// getU16LE extract 2-byte integer as unsigned little-endian.
func getU16LE(buf []byte) uint16 {
	w := getU16BE(buf)
	// exchange bytes
	v := (w&0xFF)<<8 + w>>8
	return v
}

// Calc CRC according to AM2320 specification.
func calcCRC_AM2320(buf []byte) uint16 {
	var seed uint16 = 0xFFFF
	for i := 0; i < len(buf); i++ {
		seed ^= uint16(buf[i])
		for j := 0; j < 8; j++ {
			if seed&0x01 != 0 {
				seed >>= 1
				seed ^= 0xA001
			} else {
				seed >>= 1
			}
		}
	}
	return seed
}

func calcCRC1(seed byte, buf []byte) byte {
	for i := 0; i < len(buf); i++ {
		b := buf[ /*len(buf)-1-*/ i]
		for j := 0; j < 8; j++ {
			if (seed^b)&0x01 != 0 {
				seed ^= 0x18
				seed >>= 1
				seed |= 0x80
				// crc = crc ^ 0x8c
			} else {
				seed >>= 1
			}
			b >>= 1
		}
	}
	return seed
}
