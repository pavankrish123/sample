package main

import (
	"fmt"
	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus"
	"os"

	"go.nanomsg.org/mangos/v3/protocol/push"
	"go.opentelemetry.io/collector/model/otlp"
	// imported for side effects to register ipc/tcp transport
	_ "go.nanomsg.org/mangos/v3/transport/ipc"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"

	"time"
)

type builder struct {
	metric *metricspb.Metric
}

// metricBuilder is used to build metrics for testing
func metricBuilder() builder {
	return builder{
		metric: &metricspb.Metric{
			MetricDescriptor: &metricspb.MetricDescriptor{},
			Timeseries:       make([]*metricspb.TimeSeries, 0),
		},
	}
}

// setName sets the name of the metric
func (b builder) setName(name string) builder {
	b.metric.MetricDescriptor.Name = name
	return b
}

// setDataType sets the data type of this metric
func (b builder) setDataType(dataType metricspb.MetricDescriptor_Type) builder {
	b.metric.MetricDescriptor.Type = dataType
	return b
}


// Build builds from the builder to the final metric
func (b builder) build() *metricspb.Metric {
	return b.metric
}

const (
	metricsAddr = "tcp://0.0.0.0:45001"
)


func main(){
	mSock, err := push.NewSocket()
	if err != nil {
		panic(err)
	}

	err = mSock.Dial(metricsAddr)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i ++ {

		r := &resourcepb.Resource{
			Labels: map[string]string{
				"original": "label",
			},
		}

		for k := 0; k < len(os.Args); k++ {
			r.Labels[fmt.Sprintf("resource-%d", k)] = os.Args[k]
		}

		m := []*metricspb.Metric{
			metricBuilder().setName("foo/metric").
				setDataType(metricspb.MetricDescriptor_GAUGE_INT64).build(),
			metricBuilder().setName("bar/metric").
				setDataType(metricspb.MetricDescriptor_GAUGE_INT64).build(),
		}

		metrics, err := otlp.NewProtobufMetricsMarshaler().MarshalMetrics(opencensus.OCToMetrics(nil, r, m))
		if err != nil {
			panic(err)
		}

		err = mSock.Send(metrics)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 30)
	}
}
