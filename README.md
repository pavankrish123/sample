# Sample OTLP

## Start the mock ingest server
```bash
go run server.go
```

## Start Prometheus
[Download](https://prometheus.io/download/)
```bash
./prometheus --config.file prometheus.yml
```


## Start metrics generator
```bash
go run main.go
```

## Start the OT Collector 
[Download](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.25.0)
```bash
./otelcol_darwin_amd64 --config=./otel-config.yml
```

## Deployment
```bash
metric gen app -> ot-collector -> [mock ingest server, prometheus]
```

## Prometheus 
- Doesn't show OT Distribution metrics 
- `flow_metric` is the metric 