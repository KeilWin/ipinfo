# IpInfo - selfhosted ip information service

Service can:
- Collect and update ip address information from RIRs

Service provide:
- Endpoint for get information by ip

## Make scripts
```
// Init project for dev
make init

// Build ipinfo app
make build-ipinfo

// Build ipinfo_updater app
make build-ipinfo-updater

// Run ipinfo app
run-ipinfo

// Run ipinfo_updater app
run-ipinfo-updater
```

## Migrations
```
// Up database migration
migrate -path migrations -database <database>://<user>:<password>@<host>:<port>/<database_name> up

// Down database migration
migrate -path migrations -database <database>://<user>:<password>@<host>:<port>/<database_name> down
```

## Using
```
// Ip v4
GET host/api/ipv4/127.0.0.1

// Ip v6
GET host/api/ipv6/::1

// Health
GET host/api/health
```
## How it works

1. Get info from all 5 top-level RIR(Regional Internet Registries)
    1. ARIN(American Registry for Internet Numbers)
    2. RIPE NCC(Reseaux IP Europeens)
    3. APNIC(Asia-Pacific Network Information Centre)
    4. LACNIC(Latin America and Caribbean Network Information Centre)
    5. AFRINIC(African Network Information Centre)
2. Merge with previous data in own database

P.S.
Map of RIRs areas

![Map](./docs/rir-map.svg)

## Preferences

Databases:
- PostgreSQL

Cache:
- Valkey