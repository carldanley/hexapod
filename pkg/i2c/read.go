package i2c

import "encoding/hex"

// ReadBytes read bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (o *Options) ReadBytes(buf []byte) (int, error) {
	n, err := o.rc.Read(buf)

	if err != nil {
		return n, err
	}

	o.Log.Debugf("Read %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return n, nil
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
func (o *Options) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	o.Log.Debugf("Read %d bytes starting from reg 0x%0X...", n, reg)
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return nil, 0, err
	}

	buf := make([]byte, n)
	c, err := o.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}

	return buf, c, nil
}

// ReadRegU8 reads byte from I2C-device register specified in reg.
func (o *Options) ReadRegU8(reg byte) (byte, error) {
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 1)
	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	o.Log.Debugf("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// ReadRegU16BE reads unsigned big endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU16BE(reg byte) (uint16, error) {
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 2)
	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	w := uint16(buf[0])<<8 + uint16(buf[1])
	o.Log.Debugf("Read U16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegU16LE reads unsigned little endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := o.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}

	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// ReadRegS16BE reads signed big endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegS16BE(reg byte) (int16, error) {
	if _, err := o.WriteBytes([]byte{reg}); err != nil {
		return 0, err
	}

	buf := make([]byte, 2)
	if _, err := o.ReadBytes(buf); err != nil {
		return 0, err
	}

	w := int16(buf[0])<<8 + int16(buf[1])
	o.Log.Debugf("Read S16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegS16LE reads signed little endian word (16 bits)
// from I2C-device starting from address specified in reg.
func (o *Options) ReadRegS16LE(reg byte) (int16, error) {
	w, err := o.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}

	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}
