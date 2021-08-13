package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)


func main(){
	// original client
	client := http.DefaultClient

	cfg := clientcredentials.Config{
		ClientID:     "0oapexn8bAgIB4rR05d6",
		ClientSecret: "_Yya0CZSBy_okGOuxdoAe0NeUG7E4aj__6dE0Ii8",
		Scopes:       []string{"api.metrics"},
		TokenURL:  "https://dev-37157674.okta.com/oauth2/default/v1/token",
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)

	fmt.Println("ere")


	client = cfg.Client(ctx)

	for {
                client.Get("http://localhost:9000/v1/metrics")
                time.Sleep(time.Minute)
        }	

}
