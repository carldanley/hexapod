package i2c

import "encoding/hex"

// WriteBytes send bytes to the remote I2C-device. The interpretation of
// the message is implementation-dependent.
func (o *Options) WriteBytes(buf []byte) (int, error) {
	o.Log.Debugf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return o.rc.Write(buf)
}

// WriteRegU8 writes byte to I2C-device register specified in reg.
func (o *Options) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}

	o.Log.Debugf("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU16BE writes unsigned big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}

	o.Log.Debugf("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU16LE writes unsigned little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8

	return o.WriteRegU16BE(reg, w)
}

// WriteRegS16BE writes signed big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	if _, err := o.WriteBytes(buf); err != nil {
		return err
	}

	o.Log.Debugf("Write S16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegS16LE writes signed little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
func (o *Options) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8

	return o.WriteRegS16BE(reg, w)
}

// WriteRegU24BE writes unsigned big endian word (24 bits)
// value to I2C-device starting from address specified in reg.
func (v *Options) WriteRegU24BE(reg byte, value uint32) error {
	buf := []byte{reg, byte(value >> 16 & 0xFF), byte(value >> 8 & 0xFF), byte(value & 0xFF)}
	if _, err := v.WriteBytes(buf); err != nil {
		return err
	}

	v.Log.Debugf("Write U24 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU32BE writes unsigned big endian word (32 bits)
// value to I2C-device starting from address specified in reg.
func (v *Options) WriteRegU32BE(reg byte, value uint32) error {
	buf := []byte{reg, byte(value >> 24 & 0xFF), byte(value >> 16 & 0xFF), byte(value >> 8 & 0xFF), byte(value & 0xFF)}
	if _, err := v.WriteBytes(buf); err != nil {
		return err
	}

	v.Log.Debugf("Write U32 %d to reg 0x%0X", value, reg)
	return nil
}
