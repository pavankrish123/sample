package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	mpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

type MetricsExportService struct {
	mpb.MetricsServiceServer
}

func (a *MetricsExportService) Export(ctx context.Context, request *mpb.ExportMetricsServiceRequest) (*mpb.ExportMetricsServiceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println(md)
	}

         marshaller := &jsonpb.Marshaler{Indent: "\t"}
	s, err := marshaller.MarshalToString(request)
        if err != nil {
	 	log.Print(err.Error())
         }
	log.Println(s)
	return &mpb.ExportMetricsServiceResponse{}, nil
}


func main() {
	flag.Parse()

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/tmp/certs/cert.pem", "/tmp/certs/cert-key.pem")
	if err != nil {
		panic(err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	creds := credentials.NewTLS(config)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server Listening on :%v", 5000)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	mpb.RegisterMetricsServiceServer(grpcServer, &MetricsExportService{})
	panic(grpcServer.Serve(lis))
}