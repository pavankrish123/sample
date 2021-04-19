package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() func() {
	ctx := context.Background()

	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint("localhost:30080"),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)
	exp, err := otlp.NewExporter(ctx, driver)
	handleErr(err, "failed to create exporter")

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.CloudProviderAzure,
		),
	)
	handleErr(err, "failed to create resource")

	cont := controller.New(
		processor.New(
			// aggregation goes here
			simple.NewWithInexpensiveDistribution(),
			exp,
		),
		controller.WithExporter(exp),
		controller.WithCollectPeriod(time.Minute),
		controller.WithResource(res),
	)

	global.SetMeterProvider(cont.MeterProvider())

	handleErr(cont.Start(context.Background()), "failed to start controller")

	return func() {
		// Push any last metric events to the exporter.
		handleErr(cont.Stop(context.Background()), "failed to stop controller")
	}
}

func main() {
	log.Printf("Waiting for connection...")

	shutdown := initProvider()
	defer shutdown()

	meter := global.Meter("appdynamics-flow-meter")

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []attribute.KeyValue{
		attribute.String("flowId", "HAGSA_HAHSH_IN"),
		attribute.String("flowName", "/bt/checkout"),
	}

	// Recorder metric example
	valueRecorder := metric.Must(meter).NewFloat64ValueRecorder("flow_metric",
		metric.WithDescription("Measures flow metrics"),
		).Bind(commonLabels...)
	defer valueRecorder.Unbind()

	for  {
		r := rand.Float64() * 10.0
		log.Printf("Adding Measurement %5.2f \n", r)
		valueRecorder.Record(context.Background(), r)
		<-time.After(time.Second * 10)
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
