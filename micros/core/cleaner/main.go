package main

import (
	"log"

	"net/http"

	"crypto/tls"

	"constants"

	"utils"

	"context"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

var clienteHTTP = &http.Client{}

func cleaner(input string) {

	utils.CleanerLogin(input)

}

func initRSocketServer() {
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// bind responder
			return rsocket.NewAbstractSocket(
				rsocket.FireAndForget(func(msg payload.Payload) {
					cleaner(string(msg.DataUTF8()))
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7882").Build()).
		Serve(context.Background())
	log.Println(err)

	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: http client

}

func main() {
	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	initRSocketServer()

	// JSON register
	http.HandleFunc(constants.ONBOARDING_URI, registerHandler)
	log.Printf("Starting server for testing HTTP CLEANER...\n")
	if err := http.ListenAndServe(constants.CLEANER_PORT, nil); err != nil {
		log.Println(err)
	}
}
