package main

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()
	driver := otlpgrpc.NewDriver(otlpgrpc.WithInsecure())
	exp, err := otlp.NewExporter(ctx, driver) //
	if err != nil {
		log.Fatalf("Failed to create the collector exporter: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exp,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("test-tracer")


	ctx, span := tracer.Start(ctx, "request")
	defer span.End() //

	span.SetName("metrics push")
	span.SetAttributes(
		attribute.String("metrics.name", "k8.metric"),
		attribute.String("appd.entity.type", "k8s"),
	)

	cl := http.Client{}
	cl.Do("POST", {})) -> data endpoint ("metrics")



	span.AddEvent("some event")

	log.Printf("trace id %v", span.SpanContext().TraceID().String())
}
