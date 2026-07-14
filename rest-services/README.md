# PiNAS Backend Rest Services

## Useful Commands

### 1. Get the status and logs of the rest service <br>
`sudo systemctl status pinas-rest-services.service`

### 2. Restart the rest service <br>
`sudo systemctl restart pinas-rest-services.service`

### 3. Stop the rest service <br>
`sudo systemctl stop pinas-rest-services.service`

### 4.Disable the rest service <br>
`sudo systemctl disable pinas-rest-services.service`

### 5. Remove from system <br>
`sudo apt remove pinas-rest-services`

## Build process.
You will need to have golang and make installed <br>
`sudo apt install golang-go`

To Build from this directory, run <br>
`make build`

To create a deb package, run. This will put the deb package in the /tmp directory <br>
`make package-deb`
 
To install the deb package, run <br>
`sudo apt install /tmp/PINAS-REST-Services-0.0.1_arm64.deb`

