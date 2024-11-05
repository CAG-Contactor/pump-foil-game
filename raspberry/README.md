Vad har vi gjort?
=================
1. Installera Standard Debian Desktop på Raspberry Pi
2. Sätt upp en virtuell miljö för Python
   - `python -venv /home/pumpfoil/pythonenv`
   - sätt upp alias för `pip` och `python` i _.bashrc_
     - `alias python=/home/pumpfoil/pythonenv/bin/python`
     - `alias pip=/home/pumpfoil/pythonenv/bin/pip`
2. Följ https://www.instructables.com/How-to-Use-the-MPU6050-With-the-Raspberry-Pi-4/
  - men installera _smbus_ via `pip` istället för _python3-smbus_
3. Testa access till MPU6050 med [mpu6050_test.py](./mpu6050_test.py):
   - `python mpu6050_test.py`
4 Testa att bygg webserver för access till MPU6050
   - 
