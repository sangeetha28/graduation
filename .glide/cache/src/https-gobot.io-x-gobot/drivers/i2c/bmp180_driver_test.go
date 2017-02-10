package i2c

import (
	"bytes"
	"encoding/binary"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BMP180Driver)(nil)

// --------- HELPERS
func initTestBMP180Driver() (driver *BMP180Driver) {
	driver, _ = initTestBMP180DriverWithStubbedAdaptor()
	return
}

func initTestBMP180DriverWithStubbedAdaptor() (*BMP180Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP180Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBMP180Driver(t *testing.T) {
	// Does it return a pointer to an instance of BMP180Driver?
	var bmp180 interface{} = NewBMP180Driver(newI2cTestAdaptor())
	_, ok := bmp180.(*BMP180Driver)
	if !ok {
		t.Errorf("NewBMP180Driver() should have returned a *BMP180Driver")
	}
}

func TestBMP180Driver(t *testing.T) {
	bmp180 := initTestBMP180Driver()
	gobottest.Refute(t, bmp180.Connection(), nil)
}

func TestBMP180DriverStart(t *testing.T) {
	bmp180, _ := initTestBMP180DriverWithStubbedAdaptor()
	gobottest.Assert(t, bmp180.Start(), nil)
}

func TestBMP180DriverHalt(t *testing.T) {
	bmp180 := initTestBMP180Driver()

	gobottest.Assert(t, bmp180.Halt(), nil)
}

func TestBMP180DriverMeasurements(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		if adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB {
			binary.Write(buf, binary.BigEndian, int16(408))
			binary.Write(buf, binary.BigEndian, int16(-72))
			binary.Write(buf, binary.BigEndian, int16(-14383))
			binary.Write(buf, binary.BigEndian, uint16(32741))
			binary.Write(buf, binary.BigEndian, uint16(32757))
			binary.Write(buf, binary.BigEndian, uint16(23153))
			binary.Write(buf, binary.BigEndian, int16(6190))
			binary.Write(buf, binary.BigEndian, int16(4))
			binary.Write(buf, binary.BigEndian, int16(-32768))
			binary.Write(buf, binary.BigEndian, int16(-8711))
			binary.Write(buf, binary.BigEndian, int16(2868))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterTempMSB {
			binary.Write(buf, binary.BigEndian, int16(27898))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterPressureMSB {
			binary.Write(buf, binary.BigEndian, int16(23843))
			// XLSB, not used in this test.
			buf.WriteByte(0)
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp180.Start()
	temp, err := bmp180.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(15.0))
	pressure, err := bmp180.Pressure(BMP180UltraLowPower)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(69964))
}
