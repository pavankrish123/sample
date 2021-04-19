# Sample OTLP

## Start the mock ingest server
```bash
go run server.go
```

## Start metrics generator
```bash
go run main.go
```

## Start the OT Collector 
```bash
./otelcol_darwin_amd64 --config=./otel-config.yml
```

## Deployment
```bash
metric gen -> ot-collector -> mock ingest server
```