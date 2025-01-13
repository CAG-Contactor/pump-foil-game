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
4. Testa att bygg webserver för access till MPU6050

## 2025-01-13
Vi provade att få igång Websocket-baserad lösning som samplar MPU6050 och skickar data.

1. Logga in på raspberry 
```bash
ssh pumpfoil@<ip adress>
```
2. Kopiera ner `mpu6050_websocket.py` till `/home/pumpfoil/pumpfoiler`.
3. Starta applikationen (öppnasr websocket på port 8765)
```bash
python mpu6050_websocket.py
```
4. Installera [Websocket Test Client](https://chromewebstore.google.com/detail/websocket-test-client/fgponpodhbmadfljofbimhhlengambbn?utm_source=ext_app_menu)-pluginen 
   till Chrome.
5. Öppna _Websocket Test Client_ och koppla upp mot `ws://<ip adress>:8765` och kontrollera att det kommer events med data:
```json
{"x": 0.82360537109375, "y": 0.06703764648437499, "z": -8.197746484375}
{"x": 0.8020575561523438, "y": 0.04548983154296875, "z": -8.240842114257813}
{"x": 0.8188169677734375, "y": 0.05746083984375, "z": -8.118737829589843}
{"x": 0.833182177734375, "y": 0.00239420166015625, "z": -8.1594392578125}
```


