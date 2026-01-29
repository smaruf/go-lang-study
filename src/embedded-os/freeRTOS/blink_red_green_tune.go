package main

import (
    "time"
    "machine"
)

const (
    LED_PIN  = machine.LED
    RED_LED  = machine.GPIO0      // Assuming GPIO 0 for red LED
    BUZZER_PIN = machine.GPIO1    // Assuming GPIO 1 for buzzer
)

func GreenLEDTask() {
    for {
        LED_PIN.High()
        time.Sleep(1 * time.Second)
        LED_PIN.Low()
        time.Sleep(1 * time.Second)
    }
}

func RedLEDTask() {
    for {
        RED_LED.High()
        time.Sleep(100 * time.Millisecond)
        RED_LED.Low()
        time.Sleep(100 * time.Millisecond)
    }
}

func BuzzerTask() {
    // Frequencies and durations for the tune
    frequencies := []uint32{261, 293, 329, 349, 261, 293, 329, 349} // C, D, E, F in Hz
    durations := []time.Duration{500, 500, 500, 500, 500, 500, 500, 500} // Duration for each note in milliseconds
    lyrics := "I love my people, I love my earth, I love my nature, I love my love doll"

    for {
        for i, freq := range frequencies {
            playTone(BUZZER_PIN, freq, durations[i])
            syncLEDWithLyrics(LED_PIN, durations[i], lyrics)
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func playTone(pin machine.Pin, frequency uint32, duration time.Duration) {
    pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
    period := 1e6 / frequency // Period in microseconds

    for t := time.Duration(0); t < duration*time.Millisecond; t += time.Duration(period) * time.Microsecond {
        pin.High()
        time.Sleep(time.Duration(period/2) * time.Microsecond)
        pin.Low()
        time.Sleep(time.Duration(period/2) * time.Microsecond)
    }
}

func syncLEDWithLyrics(pin machine.Pin, duration time.Duration, lyrics string) {
    // Split the duration by the number of characters in the lyrics
    charDuration := duration / time.Duration(len(lyrics))
    
    for _, char := range lyrics {
        if char != ' ' {
            pin.High()
        } else {
            pin.Low()
        }
        time.Sleep(charDuration * time.Millisecond)
    }
    pin.Low() // Ensure the LED is turned off at the end
}

func main() {
    LED_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
    RED_LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
    BUZZER_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})

    go GreenLEDTask()
    go RedLEDTask()
    go BuzzerTask()

    // Keep the main function running
    select {}
}
