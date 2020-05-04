# superman_detections-v1 - Mike Finley

A service for identifying logins by a user that occur from
locations that are farther apart than a normal person can reasonably travel

## Run

### Build Image 

From the project root run:

```bash
docker build . -t superman_detector-v1
```

### Start API 

After building, run with
```bash 
docker run -p 8080:8080 superman_detector-v1:latest
```

Alternatively, skip building and from the project root run:

```bash
go run cmd/api/main.go
```

#### Notes
- When running in docker the database will not persist between restarts. I currently have it set up where it will only write to a file on the image itself and not a mounted volume. If we wanted something more persistent then a full RDS would be more ideal or we could write on a volume. 
- Running with go run outside of docker will persist it until the db file is deleted inside of the resources dir.

## Usage

### Add some logins
```bash 
curl localhost:8080/superman_detections/v1/logins -d '{"username": "mike", "unix_timestamp":1514764800,"event_uuid":"e6f40db2-7820-4030-989e-9aa46fef182d", "ip_address":"98.126.248.120"}'
```

```bash
curl localhost:8080/superman_detections/v1/logins -d '{"username": "mike", "unix_timestamp":1514864680,"event_uuid":"c8c47fe3-d57a-48ea-a910-5a820d706437", "ip_address":"183.59.81.103"}'
```

```bash
curl localhost:8080/superman_detections/v1/logins -d '{"username": "mike", "unix_timestamp":1514774800,"event_uuid":"29e0e668-5824-40b0-b308-2c31a990d2b3", "ip_address":"133.175.227.220"}'
```
