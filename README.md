# The pump foil game

# Networking

We have a local accesspoint which will provide network to the devices.

Admin Access:
Password: Raspberry5

Wifi Access: 
SSID: rpi-net-5 or rpi-net-24 (depending on 5Ghz or 2.4Ghz)
Password: Jfokus2024

To administer things:

* Connect to the wifi above with your laptop (you will loose internet access...)
* Then goto http://192.168.0.1
* Login with password Raspberry5

In DHCP settings I have locked addresses based on MAC address so for example 
qrscanner raspberry has IP 192.168.0.51

# qrscanner raspberry

Login to qrscanner raspberry with 

```ssh ubuntu@192.168.0.51```

password is the same as the username.

qrscanner is installed as a service, so to deploy do

```
scp badgescanner/target/badgescanner-0.1.1-SNAPSHOT.jar ubuntu@192.168.1.51
ssh ubuntu@192.168.1.51
sudo systemctl restart qrscanner
# check logs with (-f means follow...)
journalctl -f -u qrscanner
```