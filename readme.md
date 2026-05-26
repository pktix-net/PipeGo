# PipeGo

PipeGo is a lightweight edge gateway and load balancer built with Go.

- [x] Supports 4-layer inbound ASN filtering via MaxMind GeoLite2-ASN database.

# pipego.routes

Example:

```text
# No ASN filter
tcp 0.0.0.0:80 -> 192.168.8.8:80

# ASN whitelist
tcp 0.0.0.0:443 -> 192.168.8.8:443 allow_asn=4134,4837,9808

# ASN blacklist
tcp 0.0.0.0:8080 -> 192.168.8.8:8080 deny_asn=13335,209
```

## ASN Database

Download GeoLite2-ASN.mmdb (free account required):
https://dev.maxmind.com/geoip/geolite2-free-geolocation-data

Place the `.mmdb` file in the same directory as `pipego`.

# Build

```bash
go build -ldflags="-s -w" ./
```

# Run

> **-p:** Web management port

> **-a:** Web management password, default account is admin

> **--asn-db:** Path to GeoLite2-ASN.mmdb (default: `GeoLite2-ASN.mmdb`, set empty to disable)

```bash
sudo ./pipego -p 12138 -a 123456 --asn-db GeoLite2-ASN.mmdb
```

# UI

![image](./docs/image1.png)
