package main

import (
	"log"

	"net/http"

	"crypto/tls"

	"constants"
	"fmt"
)

var clienteHTTP = &http.Client{}

func dockerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(
		constants.CONTENT_TYPE,
		constants.APPLICATION_JSON)

	ret := "Hello Docker"
	log.Printf("Receiving Hello Docker...\n")
	fmt.Fprintf(w, ret)

	return

}

func main() {
	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// JSON sample
	http.HandleFunc(
		constants.DOCKER_URI,
		dockerHandler)
	log.Printf("Starting server for testing HTTP DOCKER...\n")
	if err := http.ListenAndServe(constants.DOCKER_PORT, nil); err != nil {
		log.Println(err)
	}
}
