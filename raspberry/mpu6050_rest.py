import math
import mpu6050
from flask import Flask, jsonify

app = Flask(__name__)
# Create a new Mpu6050 object
mpu6050 = mpu6050.mpu6050(0x68)

# Endpoint för att hämta mpu 6050 data
@app.route('/mpu-data', methods=['GET'])
def get_items():
  acc = mpu6050.get_accel_data()
  pitch = math.atan(acc["y"]/acc["z"])*180/math.pi;
  roll = math.atan(acc["x"]/acc["z"])*180/math.pi;
  return jsonify({"pitch":pitch,"roll":roll}), 200

if __name__ == '__main__':
  print("Starting MPU6050 REST server")
  app.run(debug=True, host='0.0.0.0', port=8080)

