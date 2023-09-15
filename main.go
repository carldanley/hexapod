package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carldanley/hexy/pkg/i2c"
	"github.com/carldanley/hexy/pkg/legs"
	"github.com/carldanley/hexy/pkg/pca9685"
	"github.com/carldanley/hexy/pkg/servos"
)

var signalChannel chan os.Signal

func init() {
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)
}

func main() {
	pca9865Board1, err := i2c.New(0x40, "/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}

	defer pca9865Board1.Close()

	pca9865Board2, err := i2c.New(0x41, "/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}

	defer pca9865Board2.Close()

	driver1, err := pca9685.New(pca9865Board1, &pca9685.Options{
		Frequency:  50,
		ClockSpeed: 26430000,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer driver1.Reset()

	driver2, err := pca9685.New(pca9865Board2, &pca9685.Options{
		Frequency:  50,
		ClockSpeed: 26430000,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer driver2.Reset()

	// startup pose
	// coxa = 0
	// femur = 120
	// tibia = 135

	// standing pose
	// coxa = 0
	// femur = 25
	// tibia = 110

	// mid height pose
	// coxa = 0
	// femur = -20
	// tibia = 40

	r1Coxa, _ := servos.New(0, driver1, servos.ServoType_RightLeg1_Coxa, 0)
	r1Femur, _ := servos.New(1, driver1, servos.ServoType_RightLeg1_Femur, 0)
	r1Tibia, _ := servos.New(2, driver1, servos.ServoType_RightLeg1_Tibia, 0)
	r1Leg := legs.New(r1Coxa, r1Femur, r1Tibia)

	r2Coxa, _ := servos.New(3, driver1, servos.ServoType_RightLeg2_Coxa, 0)
	r2Femur, _ := servos.New(4, driver1, servos.ServoType_RightLeg2_Femur, 0)
	r2Tibia, _ := servos.New(5, driver1, servos.ServoType_RightLeg2_Tibia, 0)
	r2Leg := legs.New(r2Coxa, r2Femur, r2Tibia)

	r3Coxa, _ := servos.New(6, driver1, servos.ServoType_RightLeg3_Coxa, 0)
	r3Femur, _ := servos.New(7, driver1, servos.ServoType_RightLeg3_Femur, 0)
	r3Tibia, _ := servos.New(8, driver1, servos.ServoType_RightLeg3_Tibia, 0)
	r3Leg := legs.New(r3Coxa, r3Femur, r3Tibia)

	l1Coxa, _ := servos.New(0, driver2, servos.ServoType_LeftLeg1_Coxa, 0)
	l1Femur, _ := servos.New(1, driver2, servos.ServoType_LeftLeg1_Femur, 0)
	l1Tibia, _ := servos.New(2, driver2, servos.ServoType_LeftLeg1_Tibia, 0)
	l1Leg := legs.New(l1Coxa, l1Femur, l1Tibia)

	l2Coxa, _ := servos.New(3, driver2, servos.ServoType_LeftLeg2_Coxa, 0)
	l2Femur, _ := servos.New(4, driver2, servos.ServoType_LeftLeg2_Femur, 0)
	l2Tibia, _ := servos.New(5, driver2, servos.ServoType_LeftLeg2_Tibia, 0)
	l2Leg := legs.New(l2Coxa, l2Femur, l2Tibia)

	l3Coxa, _ := servos.New(6, driver2, servos.ServoType_LeftLeg3_Coxa, 0)
	l3Femur, _ := servos.New(7, driver2, servos.ServoType_LeftLeg3_Femur, 0)
	l3Tibia, _ := servos.New(8, driver2, servos.ServoType_LeftLeg3_Tibia, 0)
	l3Leg := legs.New(l3Coxa, l3Femur, l3Tibia)

	r1Leg.Start()
	defer r1Leg.Stop()

	r2Leg.Start()
	defer r2Leg.Stop()

	r3Leg.Start()
	defer r3Leg.Stop()

	l1Leg.Start()
	defer l1Leg.Stop()

	l2Leg.Start()
	defer l2Leg.Stop()

	l3Leg.Start()
	defer l3Leg.Stop()

	fmt.Println("doing work")

	defer func() {
		r1Leg.MoveToAngles(0, 0, 0, time.Second*2)
		r2Leg.MoveToAngles(0, 0, 0, time.Second*2)
		r3Leg.MoveToAngles(0, 0, 0, time.Second*2)
		l1Leg.MoveToAngles(0, 0, 0, time.Second*2)
		l2Leg.MoveToAngles(0, 0, 0, time.Second*2)
		l3Leg.MoveToAngles(0, 0, 0, time.Second*2)
		time.Sleep(time.Second * 2)
	}()

	time.Sleep(time.Second)
	r1Leg.MoveToAngles(10, 120, 135, time.Second)
	r2Leg.MoveToAngles(10, 120, 135, time.Second)
	r3Leg.MoveToAngles(10, 120, 135, time.Second)
	l1Leg.MoveToAngles(10, 120, 135, time.Second)
	l2Leg.MoveToAngles(10, 120, 135, time.Second)
	l3Leg.MoveToAngles(10, 120, 135, time.Second)
	time.Sleep(time.Second * 2)
	r1Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	r2Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	r3Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	l1Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	l2Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	l3Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	time.Sleep(time.Second * 2)
	r1Leg.MoveToAngles(0, -20, 40, time.Millisecond*500)
	r2Leg.MoveToAngles(0, -20, 40, time.Millisecond*500)
	r3Leg.MoveToAngles(0, -20, 40, time.Millisecond*500)
	l1Leg.MoveToAngles(0, -20, 40, time.Millisecond*500)
	l2Leg.MoveToAngles(0, -20, 40, time.Millisecond*500)
	l3Leg.MoveToAngles(0, -20, 40, time.Millisecond*500)
	time.Sleep(time.Second * 2)
	r1Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	r2Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	r3Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	l1Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	l2Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	l3Leg.MoveToAngles(10, 25, 110, time.Millisecond*500)
	time.Sleep(time.Second * 2)

	// <-signalChannel
}
