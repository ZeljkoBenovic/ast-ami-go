# CloudPBX to Asterisk AMI API
Program that connects to the Asterisk AMI interface and sends the AMI events to remote URL in json format (webhook).

## How to use
* Download appropriate release (i368 or amd64)
* Extract binary
* Run `cpbx-api -export`
* Edit config.yaml file to fit your needs
* Create `cpbx-api` service so that `cpbx-api` binary can run as daemon
* Configure service to use config file (example: `/usr/local/bin/cpbx-api -config /etc/cpbx-api/config.yaml`)
* Enable start on boot (example: `systemctl enable cpbx-api`)
* Start the service (example: `systemctl start cpbx-api`)

## Debugging
When `log-level` flag is set to `debug`, all debug logs will be printed out to the log file.   
It is recommended to use it only for debugging, and set `log-level` to `info` for production 