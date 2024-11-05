# import mpu6050
import asyncio
import datetime
from xmlrpc.client import DateTime

import websockets

# Create a new Mpu6050 object
#mpu6050 = mpu6050.mpu6050(0x68)

async def pump(websocket):
  print("En klient har anslutit.")
  try:
    while True:
      await websocket.send("event:"+datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
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
