# Sample OTLP

## Start metrics generator
```bash
go run main.go
```

## Edit OT Collector 

Modify endpoint to point to OTLP HTTP Ingestion Point
```bash 
# cat otel-config.yml

  otlphttp: # for additional settings see https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlphttpexporter
    endpoint: <endpoint> # point this to oltp http based endpoint
    insecure: true
    
```


## Start the OT Collector 
[Download](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.25.0)
```bash
./otelcol_darwin_amd64 --config=./otel-config.yml
```

## Deployment
```bash
metric gen app -> ot-collector -> [OTLP HTTP Ingest]
```

