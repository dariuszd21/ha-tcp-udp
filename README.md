# ha-tcp-udp

## Running development environment

### Nix package manager / NixOS
```
nix develop
```

## Build component
```
cd server/
go build/
```

## Run server
```
./ha-tcp-udp
```

### Server Arguments

* -log_level=2 "Enable debug logs"
* -print_logs=true "Print logs on console"
* -tcp_port=12222 "Change default TCP port"
* -tcp_connections_limit=200 "Change number of TCP connections that may be open at once"
* -udp_port=12222 "Change default UDP port"
* -udp_connections_limit=200 "Change number of UDP connections that may be open at once"

### Clients Arguments

*  --stable 5 - number of stable connections (that continously receive and does not drop)
*  --port 13000 - which port to connect
*  --reconnecting 10 - number of connections that establish session after randomly selected number of packets
*  --dropping 15 - number of connections that are dropped and established again with a new session

## Containerized deployment

Every package has container image that may be built and used separately.

### Server

```
cd server/
podman build -f Dockerfile --tag ha-tcp-udp
```

### TCP Client
```
cd client_tcp/
podman build -f Dockerfile --tag client-tcp
```

### UDP Client
```
cd client_udp/
podman build -f Dockerfile --tag client-udp
```

