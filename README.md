Aosong Electronics DHT12, AM2320 humidity and temperature sensors
=================================================================

[![Build Status](https://travis-ci.org/d2r2/go-aosong.svg?branch=master)](https://travis-ci.org/d2r2/go-aosong)
[![Go Report Card](https://goreportcard.com/badge/github.com/d2r2/go-aosong)](https://goreportcard.com/report/github.com/d2r2/go-aosong)
[![GoDoc](https://godoc.org/github.com/d2r2/go-aosong?status.svg)](https://godoc.org/github.com/d2r2/go-aosong)
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

DHT12 ([pdf reference](https://raw.github.com/d2r2/go-aosong/master/docs/DHT12.pdf)) and AM2320 ([pdf reference](https://raw.github.com/d2r2/go-aosong/master/docs/AM2320.pdf)) are relatively cheap and popular among Arduino and Raspberry PI developers.
Both sensors may operate via i2c-bus interface:
![image](https://raw.github.com/d2r2/go-aosong/master/docs/am2320_dht12.jpg)

Here is a library written in [Go programming language](https://golang.org/) for Raspberry PI and clones, which gives you in the output relative humidity and temperature values (making all necessary i2c-bus interacting and values computing).

Pay attention that this library only employ i2c-bus interaction approach. Other option to work with the sensors - specific "single bus communication" protocol is not implemented here.
 
Golang usage
------------

```go
func main() {
	// Create new connection to i2c-bus on 1 line with address 0x5C.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x5C, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	sensor := aosong.NewSensor(aosong.DHT12)

	log.Printf("Sensor type = %v\n", sensor.GetSensorType())
	rh, t, err := sensor.ReadRelativeHumidityAndTemperature(i2c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Relative humidity = %v%%\n", rh)
	log.Printf("Temperature in celsius = %v*C\n", t)
}
```


Getting help
------------

GoDoc [documentation](http://godoc.org/github.com/d2r2/go-aosong)

Installation
------------

```bash
$ go get -u github.com/d2r2/go-aosong
```

Troubleshooting
--------------

- *How to obtain fresh Golang installation to RPi device (either any RPi clone):*
If your RaspberryPI golang installation taken by default from repository is outdated, you may consider
to install actual golang manually from official Golang [site](https://golang.org/dl/). Download
tar.gz file containing armv6l in the name. Follow installation instructions.

- *How to enable I2C bus on RPi device:*
If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS level.
Go to "Interfacing Options" menu, to active I2C bus.
Probably you will need to reboot to load i2c kernel module.
Finally you should have device like /dev/i2c-1 present in the system.

- *How to find I2C bus allocation and device address:*
Use i2cdetect utility in format "i2cdetect -y X", where X may vary from 0 to 5 or more,
to discover address occupied by peripheral device. To install utility you should run
`apt install i2c-tools` on debian-kind system. `i2cdetect -y 1` sample output:
	```
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- 76 --    
	```

Contribute authors
------------------

* [Efimov Ioan-Alexandru](https://github.com/efimovalex): report an [issue](https://github.com/d2r2/go-aosong/issues/2) and suggested a solution.


Contact
-------

Please use [Github issue tracker](https://github.com/d2r2/go-aosong/issues) for filing bugs or feature requests.


License
-------

Go-aosong is licensed under MIT License.

