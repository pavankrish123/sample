package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/v1/metrics", func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.Header)
	})
	http.HandleFunc("/v1/traces", func(writer http.ResponseWriter, r *http.Request){
		log.Println(r.Header)
		defer r.Body.Close()
		all, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		log.Println(string(all))
	})

	log.Println("Listening on :9000 server")
	panic(http.ListenAndServe(":9000", nil))
}
