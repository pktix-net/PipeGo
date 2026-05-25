# PipeGo

PipeGo is a lightweight edge gateway and load balancer built with Go.


# pipego.routes

Example:

```text
tcp 0.0.0.0:80 -> 192.168.8.8:80
tcp 0.0.0.0:443 -> 192.168.8.8:443
```

# Build

```bash
go build -o pipego ./cmd/pipego
```

# Run

```bash
sudo ./pipego
```
