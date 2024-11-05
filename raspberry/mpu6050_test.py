import mpu6050
import time
import math

# Define a function to read the sensor data
def read_sensor_data():
    # Read the accelerometer values
    accelerometer_data = mpu6050.get_accel_data()

    # Read the gyroscope values
    gyroscope_data = mpu6050.get_gyro_data()

    # Read temp
    temperature = mpu6050.get_temp()

    return accelerometer_data, gyroscope_data, temperature

# Start a while loop to continuously read the sensor data
while True:

    # Read the sensor data
    accelerometer_data, gyroscope_data, temperature = read_sensor_data()

    # Print the sensor data
    print("Accelerometer data:", accelerometer_data)
    print("Gyroscope data:", gyroscope_data)
    print("Temp:", temperature)

    pitch = math.atan(accelerometer_data["y"]/accelerometer_data["z"])*180/math.pi;
    print("Pitch:", pitch)

    # Wait for 1 second
    time.sleep(1)
