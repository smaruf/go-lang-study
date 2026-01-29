package main

import (
	"machine"
	"time"
)

// Robot movement patterns for autonomous navigation
// Supports differential drive robots, holonomic robots, and wheeled platforms

const (
	// Left motor pins
	LEFT_MOTOR_PWM = machine.GPIO2
	LEFT_MOTOR_DIR = machine.GPIO3
	
	// Right motor pins
	RIGHT_MOTOR_PWM = machine.GPIO4
	RIGHT_MOTOR_DIR = machine.GPIO5
	
	// Obstacle sensor pins
	FRONT_SENSOR_TRIG = machine.GPIO10
	FRONT_SENSOR_ECHO = machine.GPIO11
	LEFT_SENSOR_TRIG  = machine.GPIO12
	LEFT_SENSOR_ECHO  = machine.GPIO13
	RIGHT_SENSOR_TRIG = machine.GPIO14
	RIGHT_SENSOR_ECHO = machine.GPIO15
)

// DifferentialDrive represents a two-wheeled robot with differential drive
type DifferentialDrive struct {
	leftMotor  *Motor
	rightMotor *Motor
}

// Motor represents a single motor
type Motor struct {
	pwmPin machine.Pin
	dirPin machine.Pin
	speed  int
}

// NewMotor creates a new motor controller
func NewMotor(pwmPin, dirPin machine.Pin) *Motor {
	pwmPin.Configure(machine.PinConfig{Mode: machine.PinPWM})
	dirPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return &Motor{pwmPin: pwmPin, dirPin: dirPin, speed: 0}
}

// SetSpeed sets motor speed (-100 to 100)
func (m *Motor) SetSpeed(speed int) {
	m.speed = speed
	if speed >= 0 {
		m.dirPin.High()
	} else {
		m.dirPin.Low()
		speed = -speed
	}
	// Simplified PWM - in real implementation, use proper PWM API
	_ = speed
}

// NewDifferentialDrive creates a new differential drive robot
func NewDifferentialDrive(leftPWM, leftDir, rightPWM, rightDir machine.Pin) *DifferentialDrive {
	return &DifferentialDrive{
		leftMotor:  NewMotor(leftPWM, leftDir),
		rightMotor: NewMotor(rightPWM, rightDir),
	}
}

// Forward moves robot forward
func (d *DifferentialDrive) Forward(speed int) {
	d.leftMotor.SetSpeed(speed)
	d.rightMotor.SetSpeed(speed)
}

// Backward moves robot backward
func (d *DifferentialDrive) Backward(speed int) {
	d.leftMotor.SetSpeed(-speed)
	d.rightMotor.SetSpeed(-speed)
}

// TurnLeft turns robot left
func (d *DifferentialDrive) TurnLeft(speed int) {
	d.leftMotor.SetSpeed(-speed / 2)
	d.rightMotor.SetSpeed(speed)
}

// TurnRight turns robot right
func (d *DifferentialDrive) TurnRight(speed int) {
	d.leftMotor.SetSpeed(speed)
	d.rightMotor.SetSpeed(-speed / 2)
}

// Stop stops the robot
func (d *DifferentialDrive) Stop() {
	d.leftMotor.SetSpeed(0)
	d.rightMotor.SetSpeed(0)
}

// Rotate rotates robot in place
func (d *DifferentialDrive) Rotate(speed int, direction int) {
	if direction > 0 {
		// Rotate right
		d.leftMotor.SetSpeed(speed)
		d.rightMotor.SetSpeed(-speed)
	} else {
		// Rotate left
		d.leftMotor.SetSpeed(-speed)
		d.rightMotor.SetSpeed(speed)
	}
}

// NavigationController handles autonomous navigation
type NavigationController struct {
	robot         *DifferentialDrive
	frontSensor   *SimpleSensor
	leftSensor    *SimpleSensor
	rightSensor   *SimpleSensor
	obstacleThreshold float32
}

// SimpleSensor represents a distance sensor
type SimpleSensor struct {
	trigPin machine.Pin
	echoPin machine.Pin
}

// NewSimpleSensor creates a new sensor
func NewSimpleSensor(trig, echo machine.Pin) *SimpleSensor {
	trig.Configure(machine.PinConfig{Mode: machine.PinOutput})
	echo.Configure(machine.PinConfig{Mode: machine.PinInput})
	return &SimpleSensor{trigPin: trig, echoPin: echo}
}

// GetDistance returns distance in cm (simplified)
func (s *SimpleSensor) GetDistance() float32 {
	// Simplified distance measurement
	return 50.0 // Placeholder
}

// NewNavigationController creates a new navigation controller
func NewNavigationController(robot *DifferentialDrive) *NavigationController {
	return &NavigationController{
		robot:             robot,
		frontSensor:       NewSimpleSensor(FRONT_SENSOR_TRIG, FRONT_SENSOR_ECHO),
		leftSensor:        NewSimpleSensor(LEFT_SENSOR_TRIG, LEFT_SENSOR_ECHO),
		rightSensor:       NewSimpleSensor(RIGHT_SENSOR_TRIG, RIGHT_SENSOR_ECHO),
		obstacleThreshold: 20.0, // 20cm
	}
}

// AvoidObstacles implements obstacle avoidance behavior
func (n *NavigationController) AvoidObstacles() {
	frontDist := n.frontSensor.GetDistance()
	leftDist := n.leftSensor.GetDistance()
	rightDist := n.rightSensor.GetDistance()
	
	if frontDist < n.obstacleThreshold {
		// Obstacle in front - decide which way to turn
		if leftDist > rightDist {
			n.robot.TurnLeft(50)
			time.Sleep(500 * time.Millisecond)
		} else {
			n.robot.TurnRight(50)
			time.Sleep(500 * time.Millisecond)
		}
	} else if leftDist < n.obstacleThreshold {
		// Obstacle on left - turn right
		n.robot.TurnRight(40)
		time.Sleep(300 * time.Millisecond)
	} else if rightDist < n.obstacleThreshold {
		// Obstacle on right - turn left
		n.robot.TurnLeft(40)
		time.Sleep(300 * time.Millisecond)
	} else {
		// No obstacles - move forward
		n.robot.Forward(70)
	}
}

// FollowWall implements wall-following behavior
func (n *NavigationController) FollowWall(preferredSide int) {
	var sensorDist float32
	targetDist := float32(15.0) // Target distance from wall
	
	if preferredSide > 0 {
		// Follow right wall
		sensorDist = n.rightSensor.GetDistance()
	} else {
		// Follow left wall
		sensorDist = n.leftSensor.GetDistance()
	}
	
	if sensorDist < targetDist-5 {
		// Too close to wall - turn away
		if preferredSide > 0 {
			n.robot.TurnLeft(30)
		} else {
			n.robot.TurnRight(30)
		}
	} else if sensorDist > targetDist+5 {
		// Too far from wall - turn towards
		if preferredSide > 0 {
			n.robot.TurnRight(30)
		} else {
			n.robot.TurnLeft(30)
		}
	} else {
		// Good distance - move forward
		n.robot.Forward(60)
	}
}

// PatternSquare moves robot in a square pattern
func (d *DifferentialDrive) PatternSquare(sideLength int, speed int) {
	for i := 0; i < 4; i++ {
		// Move forward
		d.Forward(speed)
		time.Sleep(time.Duration(sideLength*10) * time.Millisecond)
		
		// Stop
		d.Stop()
		time.Sleep(200 * time.Millisecond)
		
		// Turn 90 degrees right
		d.Rotate(speed, 1)
		time.Sleep(500 * time.Millisecond)
		
		// Stop
		d.Stop()
		time.Sleep(200 * time.Millisecond)
	}
}

// PatternCircle moves robot in a circular pattern
func (d *DifferentialDrive) PatternCircle(speed int, radius int) {
	// Adjust speeds for circular motion
	outerSpeed := speed
	innerSpeed := speed * radius / (radius + 10) // Adjust for wheel base
	
	d.leftMotor.SetSpeed(outerSpeed)
	d.rightMotor.SetSpeed(innerSpeed)
}

// PatternFigureEight moves robot in a figure-8 pattern
func (d *DifferentialDrive) PatternFigureEight(speed int) {
	// First circle (clockwise)
	d.PatternCircle(speed, 20)
	time.Sleep(3 * time.Second)
	
	d.Stop()
	time.Sleep(200 * time.Millisecond)
	
	// Second circle (counter-clockwise)
	d.leftMotor.SetSpeed(speed * 2 / 3)
	d.rightMotor.SetSpeed(speed)
	time.Sleep(3 * time.Second)
	
	d.Stop()
}

// RobotNavigationTask is the main FreeRTOS task for robot navigation
func RobotNavigationTask() {
	// Initialize robot
	robot := NewDifferentialDrive(
		LEFT_MOTOR_PWM, LEFT_MOTOR_DIR,
		RIGHT_MOTOR_PWM, RIGHT_MOTOR_DIR,
	)
	
	nav := NewNavigationController(robot)
	
	mode := 0
	for {
		switch mode {
		case 0: // Obstacle avoidance mode
			nav.AvoidObstacles()
			time.Sleep(100 * time.Millisecond)
			
		case 1: // Wall following mode
			nav.FollowWall(1) // Follow right wall
			time.Sleep(100 * time.Millisecond)
			
		case 2: // Pattern mode - square
			robot.PatternSquare(100, 60)
			mode = (mode + 1) % 4
			
		case 3: // Pattern mode - figure eight
			robot.PatternFigureEight(60)
			mode = 0
		}
		
		// Switch modes periodically (in real app, use sensor input or commands)
		// mode = (mode + 1) % 4
	}
}

func main() {
	// Run robot navigation task
	go RobotNavigationTask()
	
	// Keep main running
	select {}
}
