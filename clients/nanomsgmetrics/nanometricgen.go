package main

import (
	"fmt"
	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus"
	"os"

	"go.opentelemetry.io/collector/model/otlp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.nanomsg.org/mangos/v3/protocol/push"
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

// setLabels sets the labels for the metric
func (b builder) setLabels(labels []string) builder {
	labelKeys := make([]*metricspb.LabelKey, len(labels))
	for i, l := range labels {
		labelKeys[i] = &metricspb.LabelKey{
			Key: l,
		}
	}
	b.metric.MetricDescriptor.LabelKeys = labelKeys
	return b
}

// addTimeseries adds new timeseries with the labelValuesVal and startTimestamp
func (b builder) addTimeseries(startTimestampSeconds int64, labelValuesVal []string) builder {
	labelValues := make([]*metricspb.LabelValue, len(labelValuesVal))
	for i, v := range labelValuesVal {
		labelValues[i] = &metricspb.LabelValue{
			Value:    v,
			HasValue: true,
		}
	}

	var startTimestamp *timestamppb.Timestamp
	if startTimestampSeconds != 0 {
		startTimestamp = &timestamppb.Timestamp{Seconds: startTimestampSeconds}
	}

	timeseries := &metricspb.TimeSeries{
		StartTimestamp: startTimestamp,
		LabelValues:    labelValues,
		Points:         nil,
	}
	b.metric.Timeseries = append(b.metric.Timeseries, timeseries)
	return b
}

// setDataType sets the data type of this metric
func (b builder) setDataType(dataType metricspb.MetricDescriptor_Type) builder {
	b.metric.MetricDescriptor.Type = dataType
	return b
}

// setUnit sets the unit of this metric
func (b builder) setUnit(unit string) builder {
	b.metric.MetricDescriptor.Unit = unit
	return b
}

// addInt64Point adds a int64 point to the tidx-th timseries
func (b builder) addInt64Point(tidx int, val int64, timestampVal int64) builder {
	point := &metricspb.Point{
		Timestamp: &timestamppb.Timestamp{
			Seconds: timestampVal,
			Nanos:   0,
		},
		Value: &metricspb.Point_Int64Value{
			Int64Value: val,
		},
	}
	points := b.metric.Timeseries[tidx].Points
	b.metric.Timeseries[tidx].Points = append(points, point)
	return b
}

// addDoublePoint adds a double point to the tidx-th timseries
func (b builder) addDoublePoint(tidx int, val float64, timestampVal int64) builder {
	point := &metricspb.Point{
		Timestamp: &timestamppb.Timestamp{
			Seconds: timestampVal,
			Nanos:   0,
		},
		Value: &metricspb.Point_DoubleValue{
			DoubleValue: val,
		},
	}
	points := b.metric.Timeseries[tidx].Points
	b.metric.Timeseries[tidx].Points = append(points, point)
	return b
}

// addDistributionPoints adds a distribution point to the tidx-th timseries
func (b builder) addDistributionPoints(tidx int, count int64, sum float64, bounds []float64, bucketsVal []int64) builder {
	buckets := make([]*metricspb.DistributionValue_Bucket, len(bucketsVal))
	for buIdx, bucket := range bucketsVal {
		buckets[buIdx] = &metricspb.DistributionValue_Bucket{
			Count: bucket,
		}
	}
	point := &metricspb.Point{
		Timestamp: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   0,
		},
		Value: &metricspb.Point_DistributionValue{
			DistributionValue: &metricspb.DistributionValue{
				BucketOptions: &metricspb.DistributionValue_BucketOptions{
					Type: &metricspb.DistributionValue_BucketOptions_Explicit_{
						Explicit: &metricspb.DistributionValue_BucketOptions_Explicit{
							Bounds: bounds,
						},
					},
				},
				Count:   count,
				Sum:     sum,
				Buckets: buckets,
			},
		},
	}
	points := b.metric.Timeseries[tidx].Points
	b.metric.Timeseries[tidx].Points = append(points, point)
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
