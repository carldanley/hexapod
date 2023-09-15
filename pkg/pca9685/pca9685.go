package pca9685

import (
	"fmt"
	"time"

	"github.com/carldanley/hexapod/pkg/i2c"
)

const (
	DefaultAddress = 0x40

	// Taken from: https://github.com/adafruit/Adafruit-PWM-Servo-Driver-Library/blob/master/Adafruit_PWMServoDriver.h
	Mode1    byte = 0x00
	Mode2    byte = 0x01
	ModeTest byte = 0xFF

	Led0OnLow   byte = 0x06
	Led0OnHigh  byte = 0x07
	Led0OffLow  byte = 0x08
	Led0OffHigh byte = 0x09

	AllLedOnLow   byte = 0xFA
	AllLedOnHigh  byte = 0xFB
	AllLedOffLow  byte = 0xFC
	AllLedOffHigh byte = 0xFD

	Mode1Sleep         byte = 0x10
	Mode1AutoIncrement byte = 0x20
	Mode1Restart       byte = 0x80

	Prescale byte = 0xFE

	ReferenceClockSpeed float32 = 25000000.0 // 25MHz
	StepCount           float32 = 4096.0     // 12-bit
	DefaultPWMFrequency float32 = 50.0       // 50Hz
)

type PCA9685 struct {
	i2c     *i2c.Options
	options *Options
}

type Options struct {
	Frequency  float32
	ClockSpeed float32
}

func New(i2c *i2c.Options, options *Options) (*PCA9685, error) {
	address := i2c.GetAddr()
	if address == 0 {
		return nil, fmt.Errorf("I2C device is not initialized")
	}

	pca := &PCA9685{
		i2c: i2c,
		options: &Options{
			Frequency:  DefaultPWMFrequency,
			ClockSpeed: ReferenceClockSpeed,
		},
	}

	if options != nil {
		pca.options = options
	}

	// next, set the frequency for the board to communicate
	if err := pca.SetOscillatorFrequency(pca.options.Frequency); err != nil {
		return nil, err
	}

	// finally, return the pca
	return pca, nil
}

func (pca *PCA9685) SetOscillatorFrequency(frequency float32) error {
	prescaleVal := pca.options.ClockSpeed/StepCount/frequency + 0.5

	if prescaleVal < 3.0 {
		return fmt.Errorf("PCA9685 cannot output at the given frequency")
	}

	oldMode, err := pca.i2c.ReadRegU8(Mode1)
	if err != nil {
		return err
	}

	newMode := (oldMode &^ Mode1Restart) | Mode1Sleep
	if err := pca.i2c.WriteRegU8(Mode1, newMode); err != nil {
		return err
	}

	if err := pca.i2c.WriteRegU8(Prescale, byte(prescaleVal)); err != nil {
		return err
	}

	pca.options.Frequency = frequency
	if err := pca.i2c.WriteRegU8(Mode1, oldMode); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)
	if err := pca.i2c.WriteRegU8(Mode1, oldMode|Mode1Restart|Mode1AutoIncrement); err != nil {
		return err
	}

	return nil
}

func (pca *PCA9685) GetOscillatorFrequency() float32 {
	return pca.options.Frequency
}

func (pca *PCA9685) Reset() error {
	err := pca.i2c.WriteRegU8(Mode1, Mode1Restart)
	time.Sleep(time.Millisecond * 10)

	return err
}

func (pca *PCA9685) Sleep() error {
	awake, err := pca.i2c.ReadRegU8(Mode1)
	if err != nil {
		return err
	}

	sleep := awake | Mode1Sleep
	err = pca.i2c.WriteRegU8(Mode1, sleep)
	time.Sleep(time.Millisecond * 5)

	return err
}

func (pca *PCA9685) Wakeup() error {
	sleep, err := pca.i2c.ReadRegU8(Mode1)
	if err != nil {
		return err
	}

	wakeup := sleep &^ Mode1Sleep
	return pca.i2c.WriteRegU8(Mode1, wakeup)
}

func (pca *PCA9685) SetPWM(channel, on, off int) error {
	if (channel < 0) || (channel > 15) {
		return fmt.Errorf("invalid channel value")
	}

	if (on < 0) || (on > int(StepCount)) {
		return fmt.Errorf("invalid on value")
	}

	if (off < 0) || (off > int(StepCount)) {
		return fmt.Errorf("invalid off value")
	}

	buffer := []byte{
		Led0OnLow + byte(4*channel),
		byte(on),
		byte(on >> 8),
		byte(off),
		byte(off >> 8),
	}

	_, err := pca.i2c.WriteBytes(buffer)
	return err
}

func (pca *PCA9685) GetPWM(channel int, off bool) (int, error) {
	addressByte := byte(Led0OnLow + byte(4*channel))

	if off {
		addressByte += byte(2)
	}

	data, bytesRead, err := pca.i2c.ReadRegBytes(addressByte, 2)
	if err != nil {
		return 0, nil
	}

	if bytesRead != 2 {
		return 0, fmt.Errorf("invalid number of bytes read")
	}

	return int(uint16(data[0]) | uint16(data[1])<<8), nil
}
