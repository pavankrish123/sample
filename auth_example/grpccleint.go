package main

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcOAuth "google.golang.org/grpc/credentials/oauth"
	"net/http"
	"time"
)


type ClientCredentialsAuthenticator struct {
	// OAuth2 Client Credentials
	clientCredentials *clientcredentials.Config

	// client used by auth framework to fetch and manage tokens
	client            *http.Client
}


type Config struct {
	ClientID string
	ClientSecret string
	TokenURL string
	Scopes []string
	Timeout time.Duration `mapstructure:"timeout,omitempty"`
}

func newClientCredentialsAuthenticator(cfg *Config) (*ClientCredentialsAuthenticator, error) {
	return &ClientCredentialsAuthenticator{
		clientCredentials: &clientcredentials.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			TokenURL:     cfg.TokenURL,
			Scopes:       cfg.Scopes,
		},

		// TLS stuff can also be added to Transport of this metricsclient.
		client: &http.Client{
			Timeout:   cfg.Timeout,
		},
	}, nil
}


func (o *ClientCredentialsAuthenticator) PerRPCCredentials() (credentials.PerRPCCredentials, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, o.client)
	return grpcOAuth.TokenSource{
		TokenSource: o.clientCredentials.TokenSource(ctx),
	}, nil
}


func main(){

	cc, err := newClientCredentialsAuthenticator(&Config{
		ClientID: "someclientid",
		ClientSecret: "comeclientsecret",
		TokenURL: "https://some.url/v1/token",
		Scopes: []string{"api.write"},
		Timeout: time.Second * 10,
	})
	if err != nil {
		panic(err)
	}

	perRPC, err := cc.PerRPCCredentials()
	if err != nil {
		panic(err)
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
	}

	opts = append(opts, grpc.WithBlock())
	// conn, err := grpc.Dial(*addr, opts...)
	// ....
}

