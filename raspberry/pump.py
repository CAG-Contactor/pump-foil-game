import time
import mpu6050
import asyncio
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
        pumping = False
        turn = None
        update = False
        while True:
            accelerometer_data, _, _ = read_sensor_data()

            pumpValue = accelerometer_data["z"]
            turn_value = accelerometer_data["y"]
            if turn_value > 1:
                turn = "right"
                update = True
            elif turn_value < -1:
                turn = "left"
                update = True
            else:
                if turn is not None:
                    turn = None
                    update = True

            if pumpValue > -7:
                if not pumping:
                    pumping = True
                    print("Pumping")
                    update = True
            else:
                if pumping:
                    print("Not pumping")
                    pumping = False

            if update:
                message = json.dumps({"pumping":  pumping, "turn": turn, "timestamp": time.time()})
                await websocket.send(message)
                update = False
            # Clear the z_values list for the next round of sampling
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