package servos

import (
	"fmt"
	"math"
)

type ServoType struct {
	// supplied
	minHardwarePWM   int
	maxHardwarePWM   int
	maxHardwareAngle float32
	minLimitAngle    float32
	maxLimitAngle    float32

	// needs to be calculated
	midHardwarePWM float32
	minLimitPWM    float32
	centerLimitPWM float32
	maxLimitPWM    float32
}

func NewServoType(minPWM, maxPWM, centerPWMOffset int, maxHardwareAngle, minLimitAngle, maxLimitAngle float32) ServoType {
	// cache all of the user values first
	st := ServoType{
		minHardwarePWM:   minPWM,
		maxHardwarePWM:   maxPWM,
		maxHardwareAngle: maxHardwareAngle,
		minLimitAngle:    minLimitAngle,
		maxLimitAngle:    maxLimitAngle,
	}

	// calculate all of the extra data points
	st.midHardwarePWM = float32(minPWM) + (float32(maxPWM-minPWM) / 2)
	st.centerLimitPWM = st.midHardwarePWM + float32(centerPWMOffset)
	st.minLimitPWM = st.calculateMinLimitPWM(centerPWMOffset)
	st.maxLimitPWM = st.calculateMaxLimitPWM(centerPWMOffset)

	// if there was a center pwm offset, adjust min/max limit angle values
	if centerPWMOffset != 0 {
		st.minLimitAngle = ((st.centerLimitPWM - st.minLimitPWM) / (0 - st.getPWMPerDegree()))
		st.maxLimitAngle = (st.maxLimitPWM - st.centerLimitPWM) / st.getPWMPerDegree()
	}

	// finally, return the object
	return st
}

func (st *ServoType) calculateMinLimitPWM(centerPWMOffset int) float32 {
	// make sure the min angle limit is not less than what the hardware can support (split down the middle)
	if st.minLimitAngle < (0 - (st.maxHardwareAngle / 2)) {
		return float32(st.minHardwarePWM)
	}

	// calculate the min limit pwm based on a negative min angle limit
	// NOTE: this code will break if the min angle limit is a positive number ie - (min: 30 max: 40)
	minLimitPWM := (st.midHardwarePWM - ((0 - st.getPWMPerDegree()) * st.minLimitAngle))

	// if possible, adjust the minimum limit pwm by the center offset pwm (to maintain proportions)
	minLimitPWM += float32(centerPWMOffset)

	// however, we still need to make sure the min limit pwm hasn't gone out of range of what the
	// hardware can handle
	// todo: maybe show some sort of warning this adjustment happened because of the offset
	if minLimitPWM < float32(st.minHardwarePWM) {
		minLimitPWM = float32(st.minHardwarePWM)
	}

	// finally, return the min limit pwm
	return minLimitPWM
}

func (st *ServoType) calculateMaxLimitPWM(centerPWMOffset int) float32 {
	// make sure the max angle limit is not more than what the hardware can support (split down the middle)
	if st.maxLimitAngle > (st.maxHardwareAngle / 2) {
		return float32(st.maxHardwarePWM)
	}

	// calculate the max limit pwm based on a positive angle limit
	// NOTE: this code will break if the max angle limit is a negative number -ie (min: -50 max: -30)
	maxLimitPWM := (st.midHardwarePWM + (st.getPWMPerDegree() * st.maxLimitAngle))

	// if possible, adjust the maximum limit pwm by the center offset pwm (to maintain proportions)
	maxLimitPWM += float32(centerPWMOffset)

	// however, we still need to make sure the max limit pwm hasn't gone out of range of what the
	// hardware can handle
	// todo: maybe show some sort of warning this adjustment happened because of the offset
	if maxLimitPWM > float32(st.maxHardwarePWM) {
		maxLimitPWM = float32(st.maxHardwarePWM)
	}

	// finally, return the max limit pwm
	return maxLimitPWM
}

func (st *ServoType) getPWMPerDegree() float32 {
	return float32(math.Abs(float64(float32(st.maxHardwarePWM-st.minHardwarePWM) / st.maxHardwareAngle)))
}

func (st *ServoType) GetMinHardwarePWM() int {
	return st.minHardwarePWM
}

func (st *ServoType) GetMidHardwarePWM() float32 {
	return st.midHardwarePWM
}

func (st *ServoType) GetMaxHardwarePWM() int {
	return st.maxHardwarePWM
}

func (st *ServoType) GetMaxHardwareAngle() float32 {
	return st.maxHardwareAngle
}

func (st *ServoType) GetMinLimitAngle() float32 {
	return st.minLimitAngle
}

func (st *ServoType) GetMaxLimitAngle() float32 {
	return st.maxLimitAngle
}

func (st *ServoType) GetMinLimitPWM() float32 {
	return st.minLimitPWM
}

func (st *ServoType) GetCenterLimitPWM() float32 {
	return st.centerLimitPWM
}

func (st *ServoType) GetMaxLimitPWM() float32 {
	return st.maxLimitPWM
}

func (st *ServoType) ConvertAngleToPWM(angle float32) float32 {
	// keep the angle within bounds
	if angle < st.minLimitAngle {
		return st.minLimitPWM
	} else if angle > st.maxLimitAngle {
		return st.maxLimitPWM
	}

	return st.centerLimitPWM + (angle * st.getPWMPerDegree())
}

func (st *ServoType) ConvertPWMToAngle(pwm int) float32 {
	pwmF := float32(pwm)

	if pwmF < st.centerLimitPWM {
		return (st.centerLimitPWM - pwmF) / (0 - st.getPWMPerDegree())
	} else if pwmF > st.centerLimitPWM {
		return (pwmF - st.centerLimitPWM) / st.getPWMPerDegree()
	}

	return st.centerLimitPWM
}

func (st *ServoType) DebugOutput() {
	fmt.Printf("Min Hardware PWM: %d\n", st.GetMinHardwarePWM())
	fmt.Printf("Mid Hardware PWM: %f\n", st.GetMidHardwarePWM())
	fmt.Printf("Max Hardware PWM: %d\n", st.GetMaxHardwarePWM())
	fmt.Printf("Min Hardware Angle: %f\n", st.GetMaxHardwareAngle())
	fmt.Printf("Min Limit Angle: %f\n", st.GetMinLimitAngle())
	fmt.Printf("Max Limit Angle: %f\n", st.GetMaxLimitAngle())
	fmt.Printf("Min Limit PWM: %f\n", st.GetMinLimitPWM())
	fmt.Printf("Center Limit PWM: %f\n", st.GetCenterLimitPWM())
	fmt.Printf("Max Limit PWM: %f\n", st.GetMaxLimitPWM())
	fmt.Printf("PWM Per Degree: %f\n", st.getPWMPerDegree())
	fmt.Printf("Converting %f PWM to Angle: %f\n", st.GetMinLimitPWM(), st.ConvertPWMToAngle(int(st.GetMinLimitPWM())))
	fmt.Printf("Converting %f PWM to Angle: %f\n", st.GetMaxLimitPWM(), st.ConvertPWMToAngle(int(st.GetMaxLimitPWM())))
	fmt.Printf("Converting %f Angle to PWM: %f\n", st.GetMinLimitAngle(), st.ConvertAngleToPWM(st.GetMinLimitAngle()))
	fmt.Printf("Converting %f Angle to PWM: %f\n", st.GetMaxLimitAngle(), st.ConvertAngleToPWM(st.GetMaxLimitAngle()))
	fmt.Printf("Converting %f Angle to PWM: %f\n", 45.0, st.ConvertAngleToPWM(45))
	fmt.Printf("Converting %f Angle to PWM: %f\n", -45.0, st.ConvertAngleToPWM(-45))
}
