package main

import (
	"log"

	"github.com/carldanley/hexy/pkg/i2c"
	"github.com/carldanley/hexy/pkg/pca9685"
)

func main() {
	pca9865Board1, err := i2c.New(pca9685.DefaultAddress, "/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}

	defer pca9865Board1.Close()

	driver1, err := pca9685.New(pca9865Board1, &pca9685.Options{
		Frequency:  50,
		ClockSpeed: 26430000,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer driver1.Reset()
}
