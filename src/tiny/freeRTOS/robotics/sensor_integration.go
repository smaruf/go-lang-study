package main

import (
	"machine"
	"time"
)

// Sensor integration example for robotics on FreeRTOS
// Supports IMU (accelerometer/gyroscope), GPS, and ultrasonic sensors

const (
	// Ultrasonic sensor pins
	TRIG_PIN = machine.GPIO10
	ECHO_PIN = machine.GPIO11
	
	// IMU I2C pins (MPU6050)
	I2C_SDA = machine.GPIO12
	I2C_SCL = machine.GPIO13
	
	// GPS UART pins
	GPS_TX = machine.GPIO14
	GPS_RX = machine.GPIO15
)

// UltrasonicSensor measures distance using ultrasonic waves
type UltrasonicSensor struct {
	trigPin machine.Pin
	echoPin machine.Pin
}

// NewUltrasonicSensor creates a new ultrasonic sensor
func NewUltrasonicSensor(trigPin, echoPin machine.Pin) *UltrasonicSensor {
	trigPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	echoPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	
	return &UltrasonicSensor{
		trigPin: trigPin,
		echoPin: echoPin,
	}
}

// GetDistance returns distance in centimeters
func (u *UltrasonicSensor) GetDistance() float32 {
	// Send 10us pulse to trigger
	u.trigPin.Low()
	time.Sleep(2 * time.Microsecond)
	u.trigPin.High()
	time.Sleep(10 * time.Microsecond)
	u.trigPin.Low()
	
	// Measure echo pulse duration
	timeout := time.Now().Add(30 * time.Millisecond)
	for u.echoPin.Get() == false && time.Now().Before(timeout) {
		// Wait for echo start
	}
	start := time.Now()
	
	for u.echoPin.Get() == true && time.Now().Before(timeout) {
		// Wait for echo end
	}
	duration := time.Since(start)
	
	// Calculate distance (speed of sound = 343 m/s)
	// distance = (duration * 343) / 2
	distance := float32(duration.Microseconds()) * 0.0343 / 2.0
	return distance
}

// IMUData represents IMU sensor readings
type IMUData struct {
	AccelX, AccelY, AccelZ float32
	GyroX, GyroY, GyroZ    float32
	Temperature            float32
}

// IMU represents an IMU sensor (MPU6050)
type IMU struct {
	i2c     machine.I2C
	address uint8
}

// NewIMU creates a new IMU sensor interface
func NewIMU(sda, scl machine.Pin) *IMU {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		SDA:       sda,
		SCL:       scl,
		Frequency: 400000, // 400kHz
	})
	
	imu := &IMU{
		i2c:     i2c,
		address: 0x68, // MPU6050 default address
	}
	
	// Wake up MPU6050
	imu.writeRegister(0x6B, 0x00)
	
	return imu
}

// writeRegister writes a value to an IMU register
func (imu *IMU) writeRegister(reg, value uint8) error {
	return imu.i2c.Tx(uint16(imu.address), []byte{reg, value}, nil)
}

// readRegister reads a value from an IMU register
func (imu *IMU) readRegister(reg uint8, data []byte) error {
	return imu.i2c.Tx(uint16(imu.address), []byte{reg}, data)
}

// Read reads all IMU sensor data
func (imu *IMU) Read() IMUData {
	data := make([]byte, 14)
	imu.readRegister(0x3B, data) // Start reading from ACCEL_XOUT_H
	
	// Convert raw data to signed 16-bit values
	accelX := int16(data[0])<<8 | int16(data[1])
	accelY := int16(data[2])<<8 | int16(data[3])
	accelZ := int16(data[4])<<8 | int16(data[5])
	temp := int16(data[6])<<8 | int16(data[7])
	gyroX := int16(data[8])<<8 | int16(data[9])
	gyroY := int16(data[10])<<8 | int16(data[11])
	gyroZ := int16(data[12])<<8 | int16(data[13])
	
	return IMUData{
		AccelX:      float32(accelX) / 16384.0,  // ±2g range
		AccelY:      float32(accelY) / 16384.0,
		AccelZ:      float32(accelZ) / 16384.0,
		GyroX:       float32(gyroX) / 131.0,     // ±250°/s range
		GyroY:       float32(gyroY) / 131.0,
		GyroZ:       float32(gyroZ) / 131.0,
		Temperature: float32(temp)/340.0 + 36.53, // Temperature in °C
	}
}

// GPSData represents GPS coordinates
type GPSData struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
	Speed     float32
	Satellites int
}

// GPS represents a GPS module
type GPS struct {
	uart machine.UART
}

// NewGPS creates a new GPS module interface
func NewGPS(tx, rx machine.Pin) *GPS {
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{
		TX:       tx,
		RX:       rx,
		BaudRate: 9600,
	})
	
	return &GPS{uart: uart}
}

// Read reads GPS data (simplified NMEA parsing)
func (g *GPS) Read() GPSData {
	// In a real implementation, this would parse NMEA sentences
	// This is a placeholder for demonstration
	return GPSData{
		Latitude:   37.7749,
		Longitude:  -122.4194,
		Altitude:   10.0,
		Speed:      0.0,
		Satellites: 8,
	}
}

// SensorReadTask is a FreeRTOS-style task for reading sensors
func SensorReadTask() {
	ultrasonic := NewUltrasonicSensor(TRIG_PIN, ECHO_PIN)
	imu := NewIMU(I2C_SDA, I2C_SCL)
	gps := NewGPS(GPS_TX, GPS_RX)
	
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	for {
		// Read ultrasonic distance
		distance := ultrasonic.GetDistance()
		
		// Read IMU data
		imuData := imu.Read()
		
		// Read GPS data
		gpsData := gps.Read()
		
		// Blink LED to indicate sensor reading
		led.High()
		time.Sleep(100 * time.Millisecond)
		led.Low()
		
		// Check for obstacles
		if distance < 20.0 {
			// Obstacle detected - take action
			led.High() // Keep LED on as warning
		}
		
		// Log sensor data (in real application)
		_ = imuData
		_ = gpsData
		
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// Run sensor reading task
	go SensorReadTask()
	
	// Keep main running
	select {}
}
