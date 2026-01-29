package main

import (
	"machine"
	"time"
)

// Motor control example for FreeRTOS on Raspberry Pi Pico and Arduino
// Supports DC motors, servo motors, and stepper motors

const (
	// DC Motor pins
	MOTOR_PWM_PIN = machine.GPIO2
	MOTOR_DIR_PIN = machine.GPIO3
	
	// Servo motor pin
	SERVO_PIN = machine.GPIO4
	
	// Stepper motor pins (4-wire)
	STEPPER_PIN1 = machine.GPIO5
	STEPPER_PIN2 = machine.GPIO6
	STEPPER_PIN3 = machine.GPIO7
	STEPPER_PIN4 = machine.GPIO8
)

// DCMotor represents a DC motor with PWM speed control
type DCMotor struct {
	pwmPin machine.Pin
	dirPin machine.Pin
	pwm    machine.PWM
}

// NewDCMotor creates a new DC motor controller
func NewDCMotor(pwmPin, dirPin machine.Pin) *DCMotor {
	pwmPin.Configure(machine.PinConfig{Mode: machine.PinPWM})
	dirPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	pwm := machine.PWM0
	pwm.Configure(machine.PWMConfig{Period: 16384})
	
	return &DCMotor{
		pwmPin: pwmPin,
		dirPin: dirPin,
		pwm:    pwm,
	}
}

// SetSpeed sets motor speed (-100 to 100, negative for reverse)
func (m *DCMotor) SetSpeed(speed int) {
	if speed < -100 {
		speed = -100
	}
	if speed > 100 {
		speed = 100
	}
	
	// Set direction
	if speed >= 0 {
		m.dirPin.High()
	} else {
		m.dirPin.Low()
		speed = -speed
	}
	
	// Set PWM duty cycle (0-65535)
	duty := uint32(speed * 655)
	channel, err := m.pwm.Channel(m.pwmPin)
	if err == nil {
		m.pwm.Set(channel, duty)
	}
}

// Stop stops the motor
func (m *DCMotor) Stop() {
	m.SetSpeed(0)
}

// ServoMotor represents a servo motor controller
type ServoMotor struct {
	pin machine.Pin
	pwm machine.PWM
}

// NewServoMotor creates a new servo motor controller
func NewServoMotor(pin machine.Pin) *ServoMotor {
	pin.Configure(machine.PinConfig{Mode: machine.PinPWM})
	
	pwm := machine.PWM1
	pwm.Configure(machine.PWMConfig{Period: 20000000}) // 20ms period for servo
	
	return &ServoMotor{
		pin: pin,
		pwm: pwm,
	}
}

// SetAngle sets servo angle (0-180 degrees)
func (s *ServoMotor) SetAngle(angle int) {
	if angle < 0 {
		angle = 0
	}
	if angle > 180 {
		angle = 180
	}
	
	// Servo pulse width: 500us (0°) to 2500us (180°)
	pulseWidth := 500 + (angle * 2000 / 180)
	channel, err := s.pwm.Channel(s.pin)
	if err == nil {
		s.pwm.Set(channel, uint32(pulseWidth))
	}
}

// StepperMotor represents a 4-wire stepper motor
type StepperMotor struct {
	pins     [4]machine.Pin
	position int
}

// NewStepperMotor creates a new stepper motor controller
func NewStepperMotor(pin1, pin2, pin3, pin4 machine.Pin) *StepperMotor {
	pins := [4]machine.Pin{pin1, pin2, pin3, pin4}
	for _, pin := range pins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	
	return &StepperMotor{
		pins:     pins,
		position: 0,
	}
}

// Step moves the stepper motor one step
func (s *StepperMotor) Step(direction int) {
	// Half-step sequence for smoother operation
	sequence := [][]bool{
		{true, false, false, false},
		{true, true, false, false},
		{false, true, false, false},
		{false, true, true, false},
		{false, false, true, false},
		{false, false, true, true},
		{false, false, false, true},
		{true, false, false, true},
	}
	
	if direction > 0 {
		s.position = (s.position + 1) % len(sequence)
	} else {
		s.position = (s.position - 1 + len(sequence)) % len(sequence)
	}
	
	for i, state := range sequence[s.position] {
		if state {
			s.pins[i].High()
		} else {
			s.pins[i].Low()
		}
	}
}

// Rotate rotates the stepper motor by specified steps
func (s *StepperMotor) Rotate(steps int, delayMs time.Duration) {
	direction := 1
	if steps < 0 {
		direction = -1
		steps = -steps
	}
	
	for i := 0; i < steps; i++ {
		s.Step(direction)
		time.Sleep(delayMs * time.Millisecond)
	}
}

// MotorControllerTask is a FreeRTOS-style task for motor control
func MotorControllerTask() {
	// Initialize motors
	dcMotor := NewDCMotor(MOTOR_PWM_PIN, MOTOR_DIR_PIN)
	servo := NewServoMotor(SERVO_PIN)
	stepper := NewStepperMotor(STEPPER_PIN1, STEPPER_PIN2, STEPPER_PIN3, STEPPER_PIN4)
	
	for {
		// DC Motor forward at 75% speed
		dcMotor.SetSpeed(75)
		servo.SetAngle(90)
		time.Sleep(2 * time.Second)
		
		// DC Motor reverse at 50% speed
		dcMotor.SetSpeed(-50)
		servo.SetAngle(0)
		stepper.Rotate(200, 5) // 200 steps clockwise
		time.Sleep(2 * time.Second)
		
		// Stop DC motor
		dcMotor.Stop()
		servo.SetAngle(180)
		stepper.Rotate(-200, 5) // 200 steps counter-clockwise
		time.Sleep(2 * time.Second)
	}
}

func main() {
	// Run motor controller as a task
	go MotorControllerTask()
	
	// Keep main running
	select {}
}
