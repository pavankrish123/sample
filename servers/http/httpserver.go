package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cfg, err := readConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error parsing configuration %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", welcomeHandleFunc).Methods("GET")

	handler := &errorSimulator{
		status:   cfg.Response.Code,
		errorMsg: cfg.Response.Message,
		m:        &protobufMarshaller{},
	}

	r.Handle("/v1/traces", handler).Methods("POST")
	r.Handle("/v1/metrics", handler).Methods("POST")
	r.Handle("/v1/logs", handler).Methods("POST")

	log.Printf("Starting OTLP HTTP Server :%v", cfg.Port)
	panic(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), r))
}

func welcomeHandleFunc(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, Simulated World!"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Simulated Handler
type errorSimulator struct {
	status   int
	errorMsg string
	m        marshaller
}

func (s *errorSimulator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request %v, %v", r.Method, r.RequestURI)

	// write headers
	w.Header().Set("Content-Type", s.m.contentType())
	w.WriteHeader(s.status)
	w.Header().Set("Retry-After", "15")

	// write the status as proto
	msg, err := s.m.marshalStatus(errorMsgToStatus(s.errorMsg, s.status).Proto())
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}

	// write message
	_, _ = w.Write(msg)
}

func errorMsgToStatus(errMsg string, statusCode int) *status.Status {
	switch statusCode {
	case http.StatusBadRequest:
		return status.New(codes.InvalidArgument, errMsg)
	case http.StatusTooManyRequests:
		return status.New(codes.ResourceExhausted, errMsg)
	case http.StatusInternalServerError:
		return status.New(codes.Internal, errMsg)
	case http.StatusNotFound:
		return status.New(codes.NotFound, errMsg)
	}
	return status.New(codes.Unknown, errMsg)
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

type Resp struct {
	Code    int    `mapstructure:"code"`
	Message string `mapstructure:"message"`
}

type ServerConfig struct {
	Port     int  `mapstructure:"port"`
	Response Resp `mapstructure:"response"`
}

func readConfig(cfgFile string) (*ServerConfig, error) {
	viper.SetDefault("port", 8080)
	viper.SetDefault("response", Resp{Code: 429, Message: "too many messages"})

	// Read the configuration file
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	// Unmarshal the configuration into a ServerConfig object
	var cfg ServerConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
