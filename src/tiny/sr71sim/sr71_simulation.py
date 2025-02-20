"""
SR-71 Blackbird Flight Simulation

This script provides an animated visualization of the SR-71 Blackbird's flight, demonstrating changes in altitude, speed, and engine mode over time. It features interactive controls and dynamic data visualization techniques, enhancing understanding of the aircraft's performance characteristics.

Requirements:
- Python 3.6 or higher.
- NumPy: For handling numerical operations.
- Matplotlib: For plotting and animating the data.

Libraries Installation:
Use pip to install the necessary Python packages:
$ pip install numpy matplotlib

How to Run:
1. Ensure Python 3.6+ is installed on your system.
2. Install the required Python libraries using pip if they are not already installed.
3. Save this script to a local file, for example, 'sr71_simulation.py'.
4. Run the script using Python:
$ python sr71_simulation.py

Functionality:
- The animation illustrates the SR-71 Blackbird's flight path, changing altitude, speed, and engine mode over a period of simulated time.
- The plot dynamically updates showing altitude and speed data.
- The user can interact with the animation using a slider to control the simulation time or scan through the flight data.
- Color coding and performance metrics highlight critical flight conditions and changes.

Features:
- Altitude and speed are plotted in real-time as the animation progresses.
- A slider allows users to control and navigate different points in the simulation dynamically.
- Dynamic axis scaling and conditional color coding are used to enhance data visualization.
- Annotations and a legend provide insights into flight conditions and data interpretation.
"""

import numpy as np
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation
from matplotlib.widgets import Slider

def simulate_flight():
    # Time range from 0 to 100 with 1000 points
    times = np.linspace(0, 100, num=1000)
    # Example sinusoidal altitude change
    altitudes = 25000 + 10000 * np.sin(0.04 * times)
    # Example cosinusoidal speed change
    speeds = 1800 + 800 * np.cos(0.05 * times)
    # Define Turbojet or Ramjet based on speed
    engine_modes = ["Turbojet" if speed < 2200 else "Ramjet" for speed in speeds]
    return times, altitudes, speeds, engine_modes

def update(frame):
    current_time = times[frame]
    current_alt = altitudes[frame]
    current_speed = speeds[frame]
    current_mode = engine_modes[frame]

    flight_path.set_data(times[:frame], altitudes[:frame])
    speed_graph.set_data(times[:frame], speeds[:frame])
    
    # Updating scatter points
    scat.set_offsets([[current_time, current_alt]])
    scat_speed.set_offsets([[current_time, current_speed]])

    # Updating text annotations
    speed_text.set_text(f"Speed: {current_speed:.1f} mph")
    altitude_text.set_text(f"Altitude: {current_alt:.1f} ft")
    engine_text.set_text(f"Engine Mode: {current_mode}")

    # Color coding for different conditions
    if current_speed > 2400:
        flight_path.set_color('red')  # Danger: high speed
        speed_graph.set_color('red')
    else:
        flight_path.set_color('green')  # Normal operation
        speed_graph.set_color('blue')
    
    # Display performance metrics (example: acceleration)
    if frame > 0:
        dt = times[frame] - times[frame - 1]
        d_speed = speeds[frame] - speeds[frame - 1]
        acceleration = d_speed / dt
        perf_metrics.set_text(f"Acceleration: {acceleration:.2f} mph/s")

    return flight_path, scat, speed_graph, scat_speed, speed_text, altitude_text, engine_text, perf_metrics

# Simulation data
times, altitudes, speeds, engine_modes = simulate_flight()

# Set up the plots (Figure and Axes)
fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 10))

# Axis for Altitude plot
ax1.set_title('SR-71 Blackbird Flight Simulation')
ax1.set_xlim(times.min(), times.max())
ax1.set_ylim(20000, 40000)
ax1.set_xlabel('Time (s)')
ax1.set_ylabel('Altitude (ft)')
flight_path, = ax1.plot([], [], 'g-')
scat = ax1.scatter([], [], color='red')

# Axis for Speed plot
ax2.set_xlim(times.min(), times.max())
ax2.set_ylim(1500, 3000)
ax2.set_xlabel('Time (s)')
ax2.set_ylabel('Speed (mph)')
speed_graph, = ax2.plot([], [], 'b-')
scat_speed = ax2.scatter([], [], color='orange')

# Text displays for information
speed_text = ax1.text(0.05, 0.95, '', transform=ax1.transAxes, color='blue')
altitude_text = ax1.text(0.05, 0.90, '', transform=ax1.transAxes, color='blue')
engine_text = ax1.text(0.05, 0.85, '', transform=ax1.transAxes, color='purple')
perf_metrics = ax2.text(0.05, 0.95, '', transform=ax2.transAxes, color='green')

# Time control slider
slider_ax = plt.axes([0.25, 0.01, 0.50, 0.02], facecolor='lightgoldenrodyellow')
time_slider = Slider(ax=slider_ax, label='Time', valmin=times.min(), valmax=times.max(), valinit=times[0])

def slider_update(val):
    frame = np.searchsorted(times, val)
    update(frame)
    fig.canvas.draw_idle()

time_slider.on_changed(slider_update)

# Create the Animation
ani = FuncAnimation(fig, update, frames=len(times), interval=50, blit=True)

plt.tight_layout()
plt.show()
