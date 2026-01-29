import time

class Pin:
    """A class to represent a GPIO pin."""
    
    OUT = 0

    def __init__(self, pin_number):
        self.pin_number = pin_number
        self.state = False

    def configure(self, mode):
        """Configure the pin mode."""
        pass

    def high(self):
        """Set the pin high."""
        self.state = True
        print(f"Pin {self.pin_number} is HIGH")

    def low(self):
        """Set the pin low."""
        self.state = False
        print(f"Pin {self.pin_number} is LOW")

LED_PIN = Pin(13)  # Assuming GPIO 13 for green LED
RED_LED = Pin(0)   # Assuming GPIO 0 for red LED
BUZZER_PIN = Pin(1)  # Assuming GPIO 1 for buzzer

def green_led_task():
    """Task to blink the green LED."""
    while True:
        LED_PIN.high()
        time.sleep(1)
        LED_PIN.low()
        time.sleep(1)

def red_led_task():
    """Task to blink the red LED."""
    while True:
        RED_LED.high()
        time.sleep(0.1)
        RED_LED.low()
        time.sleep(0.1)

def buzzer_task():
    """Task to play a tune with the buzzer."""
    frequencies = [261, 293, 329, 349, 261, 293, 329, 349]  # C, D, E, F in Hz
    durations = [0.5] * 8  # Duration for each note in seconds
    lyrics = "I love my people, I love my earth, I love my nature, I love my love doll"
    
    while True:
        for i, freq in enumerate(frequencies):
            play_tone(BUZZER_PIN, freq, durations[i])
            sync_led_with_lyrics(LED_PIN, durations[i], lyrics)
            time.sleep(0.1)

def play_tone(pin, frequency, duration):
    """Play a tone on the buzzer."""
    period = 1.0 / frequency  # Period in seconds
    end_time = time.time() + duration
    
    while time.time() < end_time:
        pin.high()
        time.sleep(period / 2)
        pin.low()
        time.sleep(period / 2)

def sync_led_with_lyrics(pin, duration, lyrics):
    """Sync LED blinking with the lyrics."""
    char_duration = duration / len(lyrics)
    
    for char in lyrics:
        if char != ' ':
            pin.high()
        else:
            pin.low()
        time.sleep(char_duration)
    pin.low()  # Ensure the LED is turned off at the end

if __name__ == "__main__":
    LED_PIN.configure(Pin.OUT)
    RED_LED.configure(Pin.OUT)
    BUZZER_PIN.configure(Pin.OUT)
    
    import threading
    threading.Thread(target=green_led_task).start()
    threading.Thread(target=red_led_task).start()
    threading.Thread(target=buzzer_task).start()
    
    while True:
        time.sleep(1)
