package hexapod

import (
	"log"
	"time"

	"github.com/carldanley/hexapod/pkg/i2c"
	"github.com/carldanley/hexapod/pkg/legs"
	"github.com/carldanley/hexapod/pkg/pca9685"
	"github.com/carldanley/hexapod/pkg/servos"
)

type Hexapod struct {
	legs         []legs.Leg
	i2cSlaves    []*i2c.Options
	servoDrivers []*pca9685.PCA9685
}

func New() (*Hexapod, error) {
	hexapod := Hexapod{
		legs:         []legs.Leg{},
		i2cSlaves:    []*i2c.Options{},
		servoDrivers: []*pca9685.PCA9685{},
	}

	// open i2c connections to the 2x servo controller slaves
	slave1, _ := hexapod.addI2CSlave(0x40, "/dev/i2c-1")
	// slave2, _ := hexapod.addI2CSlave(0x41, "/dev/i2c-1")

	// initialize all of the servo drivers
	servoDriver1, err := hexapod.addServoDriver(slave1)
	if err != nil {
		log.Fatal(err)
	}
	// servoDriver2, _ := hexapod.addServoDriver(slave2)

	// initialize all of the legs
	hexapod.addLeg(0, servoDriver1)
	hexapod.addLeg(3, servoDriver1)
	hexapod.addLeg(6, servoDriver1)
	// hexapod.addLeg(9, servoDriver1)
	// hexapod.addLeg(12, servoDriver1)
	// hexapod.addLeg(0, servoDriver2)

	return &hexapod, nil
}

func (hp *Hexapod) addI2CSlave(address uint8, dev string) (*i2c.Options, error) {
	slave, err := i2c.New(address, dev)
	if err != nil {
		return nil, err
	}

	hp.i2cSlaves = append(hp.i2cSlaves, slave)
	return slave, nil
}

func (hp *Hexapod) addServoDriver(slave *i2c.Options) (*pca9685.PCA9685, error) {
	servoDriver, err := pca9685.New(slave, &pca9685.Options{
		Frequency:  50,
		ClockSpeed: 26430000,
	})

	if err != nil {
		return nil, err
	}

	hp.servoDrivers = append(hp.servoDrivers, servoDriver)

	return servoDriver, nil
}

func (hp *Hexapod) addLeg(channelOffset int, servoDriver *pca9685.PCA9685) (legs.Leg, error) {
	coxa, err := servos.New(channelOffset, servoDriver, servos.ServoType_DS3225_90, 0)
	if err != nil {
		return legs.Leg{}, err
	}

	femur, err := servos.New(channelOffset+1, servoDriver, servos.ServoType_DS3225_135, 0)
	if err != nil {
		return legs.Leg{}, err
	}

	tibia, err := servos.New(channelOffset+2, servoDriver, servos.ServoType_DS3225_135, 0)
	if err != nil {
		return legs.Leg{}, err
	}

	leg := legs.New(coxa, femur, tibia)
	hp.legs = append(hp.legs, leg)
	leg.Start()

	return leg, nil
}

func (hp *Hexapod) Shutdown() {
	// iterate through the legs, stopping each one
	for _, leg := range hp.legs {
		leg.Stop()
	}

	// iterate through the boards and reset each one
	for _, driver := range hp.servoDrivers {
		driver.Reset()
	}

	// iterate through the slaves and closeout communication over i2c
	for _, slave := range hp.i2cSlaves {
		slave.Close()
	}
}

func (hp *Hexapod) MoveAllLegsToAngles(coxaAngle, femurAngle, tibiaAngle float32, duration time.Duration) {
	for _, leg := range hp.legs {
		leg.MoveToAngles(coxaAngle, femurAngle, tibiaAngle, duration)
	}
}

func (hp *Hexapod) GetLeg(index int) legs.Leg {
	return hp.legs[index]
}
