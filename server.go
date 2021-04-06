package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	mpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

type MetricsExportService struct {
	mpb.MetricsServiceServer
}

func (a *MetricsExportService) Export(_ context.Context, request *mpb.ExportMetricsServiceRequest) (*mpb.ExportMetricsServiceResponse, error) {
	marshaller := &jsonpb.Marshaler{Indent: "\t"}
	s, err := marshaller.MarshalToString(request)
	if err != nil {
		log.Print(err.Error())
	}
	log.Println(s)
	return &mpb.ExportMetricsServiceResponse{}, err
}


func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server Listening on :%v", 5000)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	mpb.RegisterMetricsServiceServer(grpcServer, &MetricsExportService{})
	panic(grpcServer.Serve(lis))
}
