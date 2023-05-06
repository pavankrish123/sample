package main

import (
	"fmt"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", welcomeHandleFunc).Methods("GET")

	// TODO: start reading these values from a configuration file??
	handler := errorSimulator{
		status:   http.StatusTooManyRequests,
		errorMsg: "way too many requests",
		m:        &protobufMarshaller{},
	}

	r.Handle("/v1/traces", handler).Methods("POST")
	r.Handle("/v1/metrics", handler).Methods("POST")
	r.Handle("/v1/logs", handler).Methods("POST")

	panic(http.ListenAndServe(":8080", r))
}

func welcomeHandleFunc(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, Simulated World!"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type marshaller interface {
	marshalStatus(resp *spb.Status) ([]byte, error)
	contentType() string
}

type protobufMarshaller struct{}

const (
	pbContentType = "application/x-protobuf"
)

func (protobufMarshaller) marshalStatus(resp *spb.Status) ([]byte, error) {
	return proto.Marshal(resp)
}

func (protobufMarshaller) contentType() string {
	return pbContentType
}

// Simulated Handler
type errorSimulator struct {
	status   int
	errorMsg string
	m        marshaller
}

func (s errorSimulator) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	// write headers
	w.Header().Set("Content-Type", s.m.contentType())
	w.WriteHeader(s.status)

	// write the status as proto
	msg, err := s.m.marshalStatus(errorMsgToStatus(s.errorMsg, s.status).Proto())
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
	_, _ = w.Write(msg)
}

func errorMsgToStatus(errMsg string, statusCode int) *status.Status {
	switch statusCode {
	case http.StatusBadRequest:
		return status.New(codes.InvalidArgument, errMsg)
	case http.StatusTooManyRequests:
		return status.New(codes.Unauthenticated, errMsg)
	case http.StatusInternalServerError:
		return status.New(codes.Internal, errMsg)
	case http.StatusNotFound:
		return status.New(codes.NotFound, errMsg)
	}
	return status.New(codes.Unknown, errMsg)
}
