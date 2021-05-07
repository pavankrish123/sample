package main

import (
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/v1/metrics", func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.Header)
	})
	log.Println("Listening on :9000 server")
	panic(http.ListenAndServe(":9000", nil))
}
