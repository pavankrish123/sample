module github.com/pavankrish123/sample-otlp-server

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/exporters/otlp v0.19.0
	go.opentelemetry.io/otel/metric v0.19.0
	go.opentelemetry.io/otel/sdk v0.19.0
	go.opentelemetry.io/otel/sdk/metric v0.19.0
	go.opentelemetry.io/proto/otlp v0.7.0
	google.golang.org/grpc v1.36.0
)
