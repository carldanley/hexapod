package servos

import (
	"math"
	"time"

	"github.com/carldanley/hexy/pkg/easings"
	"github.com/carldanley/hexy/pkg/pca9685"
)

const ServoMovementSpeedMS = 18

type Servo struct {
	channel         int
	controller      *pca9685.PCA9685
	servoType       ServoType
	stopWorkChannel chan bool

	currentPWM      float32
	beginningPWM    float32
	endingPWM       float32
	easingStartTime time.Time
	easingDuration  time.Duration
}

func New(channel int, controller *pca9685.PCA9685, servoType ServoType, defaultAngle float32) (*Servo, error) {
	servo := &Servo{
		channel:         channel,
		controller:      controller,
		servoType:       servoType,
		stopWorkChannel: make(chan bool),
		beginningPWM:    servoType.ConvertAngleToPWM(defaultAngle),
		endingPWM:       servoType.ConvertAngleToPWM(defaultAngle),
		currentPWM:      servoType.ConvertAngleToPWM(defaultAngle),
		easingStartTime: time.Now(),
		easingDuration:  time.Duration(0),
	}

	// we have to set the servo's position right off the bat (in order
	// to accurately do the math for easing)
	controller.SetPWM(channel, 0, int(servo.currentPWM))

	return servo, nil
}

func (s *Servo) MoveToAngle(angle float32, duration time.Duration) {
	s.MoveToPWM(s.servoType.ConvertAngleToPWM(angle), duration)
}

func (s *Servo) MoveToPWM(pwm float32, duration time.Duration) {
	if pwm < s.servoType.GetMinLimitPWM() {
		pwm = s.servoType.GetMinLimitPWM()
	} else if pwm > s.servoType.GetMaxLimitPWM() {
		pwm = s.servoType.GetMaxLimitPWM()
	}

	// setup a few of the easing variables
	s.beginningPWM = s.currentPWM
	s.endingPWM = pwm
	s.easingStartTime = time.Now()
	s.easingDuration = duration

	// handle cases where the servo needs to move directly to the pwm
	if s.easingDuration.Milliseconds() < ServoMovementSpeedMS {
		s.easingDuration = ServoMovementSpeedMS
	}
}

func (s *Servo) Stop() {
	if s.stopWorkChannel != nil {
		s.stopWorkChannel <- true
		close(s.stopWorkChannel)
	}
}

func (s *Servo) Start() {
	for {
		select {
		case <-s.stopWorkChannel:
			return
		case <-time.After(time.Duration(ServoMovementSpeedMS) * time.Millisecond):
			s.performStep()
		}
	}
}

func (s *Servo) performStep() {
	elapsedTime := float32(time.Since(s.easingStartTime).Milliseconds())
	changeInPWM := float32(s.endingPWM - s.beginningPWM)

	newPWM := easings.LinearNone(elapsedTime, s.beginningPWM, changeInPWM, float32(s.easingDuration.Milliseconds()))

	if math.IsNaN(float64(newPWM)) {
		return
	}

	if (changeInPWM > 0) && (newPWM > s.endingPWM) {
		newPWM = s.endingPWM
	} else if (changeInPWM < 0) && (newPWM < s.endingPWM) {
		newPWM = s.endingPWM
	}

	if int(newPWM) != int(s.currentPWM) {
		s.controller.SetPWM(s.channel, 0, int(newPWM))
		s.currentPWM = newPWM
	}
}
