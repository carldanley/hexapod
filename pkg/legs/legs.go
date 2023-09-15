package legs

import (
	"time"

	"github.com/carldanley/hexapod/pkg/servos"
)

type Leg struct {
	coxa  *servos.Servo
	femur *servos.Servo
	tibia *servos.Servo
}

func New(coxa, femur, tibia *servos.Servo) Leg {
	return Leg{
		coxa,
		femur,
		tibia,
	}
}

func (l *Leg) MoveToAngles(coxaAngle, femurAngle, tibiaAngle float32, duration time.Duration) {
	l.MoveCoxaToAngle(coxaAngle, duration)
	l.MoveFemurToAngle(femurAngle, duration)
	l.MoveTibiaToAngle(tibiaAngle, duration)
}

func (l *Leg) MoveCoxaToAngle(angle float32, duration time.Duration) {
	l.coxa.MoveToAngle(angle, duration)
}

func (l *Leg) MoveFemurToAngle(angle float32, duration time.Duration) {
	l.femur.MoveToAngle(angle, duration)
}

func (l *Leg) MoveTibiaToAngle(angle float32, duration time.Duration) {
	l.tibia.MoveToAngle(angle, duration)
}

func (l *Leg) Start() {
	go l.coxa.Start()
	go l.femur.Start()
	go l.tibia.Start()
}

func (l *Leg) Stop() {
	l.coxa.Stop()
	l.femur.Stop()
	l.tibia.Stop()
}
