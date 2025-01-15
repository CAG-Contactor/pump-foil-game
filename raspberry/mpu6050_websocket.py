import mpu6050
import asyncio
import math
import websockets
import json

# Create a new Mpu6050 object
mpu6050 = mpu6050.mpu6050(0x68)

# Define a function to read the sensor data
def read_sensor_data():
  # Read the accelerometer values
  accelerometer_data = mpu6050.get_accel_data()

  # Read the gyroscope values
  gyroscope_data = mpu6050.get_gyro_data()

  # Read temp
  temperature = mpu6050.get_temp()

  return accelerometer_data, gyroscope_data, temperature

async def pump(websocket):
  print("En klient har anslutit.")
  try:
    while True:
      # Read the sensor data
      accelerometer_data, gyroscope_data, temperature = read_sensor_data()

      # Print the sensor data
      print("Accelerometer data:", accelerometer_data)
      print("Gyroscope data:", gyroscope_data)
      print("Temp:", temperature)

      pitch = math.atan(accelerometer_data["y"]/accelerometer_data["z"])*180/math.pi;
      print("Pitch:", pitch)

      message = json.dumps(accelerometer_data)
      await websocket.send(message)
      await asyncio.sleep(10.0)
  except websockets.ConnectionClosed:
    print("Klienten kopplade från.")

async def handler(websocket, path):
  print("En klient har anslutit.")
  try:
    async for message in websocket:
      print(f"Meddelande från klient: {message}")
      # Svara tillbaka till klienten
      response = f"Servern mottog ditt meddelande: {message}"
      await websocket.send(response)
  except websockets.ConnectionClosed:
    print("Klienten kopplade från.")

# Starta servern på localhost:8765
async def main():
  print("WS main")
  async with websockets.serve(pump, "0.0.0.0", 8765):
#  async with websockets.serve(handler, "0.0.0.0", 8765):
    print("WebSocket-servern körs på ws://0.0.0.0:8765")
    await asyncio.Future()  # Håll servern igång

if __name__ == '__main__':
  print("Starting MPU6050 Websocket server")
  asyncio.run(main())
