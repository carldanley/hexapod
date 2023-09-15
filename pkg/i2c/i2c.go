package i2c

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Options struct {
	addr uint8
	dev  string
	rc   *os.File
	Log  *logrus.Logger
}

func New(addr uint8, dev string) (*Options, error) {
	i2c := &Options{
		addr: addr,
		dev:  "/dev/i2c-0",
		Log:  logrus.New(),
	}

	if dev != "" {
		i2c.dev = dev
	}

	f, err := os.OpenFile(dev, os.O_RDWR, 0600)
	if err != nil {
		return i2c, err
	}

	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return i2c, err
	}

	i2c.rc = f

	return i2c, nil
}

func (o *Options) GetAddr() uint8 {
	return o.addr
}

func (o *Options) GetDev() string {
	return o.dev
}

func (o *Options) Close() error {
	return o.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0); err != 0 {
		return err
	}
	return nil
}
