package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main(){
	log.Println("Starting Server on 0.0.0.0:9000")
	http.HandleFunc("/", func (w http.ResponseWriter, req *http.Request ){
		log.Println(req.Header["Content-type"])
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(string(body))
		_, _ = w.Write([]byte("OK"))
	})
	panic(http.ListenAndServe(":9000", nil))
}