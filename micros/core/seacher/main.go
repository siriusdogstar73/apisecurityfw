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
	"github.com/rsocket/rsocket-go/rx/mono"
)

var clienteHTTP = &http.Client{}

func SearchProcess(sHash string) string {

	sOnboarding := utils.SearchProcess(sHash)

	return utils.RenewalJwt(clienteHTTP, sOnboarding)
}

func initRSocketServer() {
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// bind responder
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					//fmt.Print(string(setup.DataUTF8()))
					var result = payload.NewString(SearchProcess(msg.DataUTF8()), "Metadato")
					return mono.Just(result)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7879").Build()).
		Serve(context.Background())
	log.Println(err)

	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: http client
	SearchProcess("")
}

func main() {
	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	initRSocketServer()
	//receiveRequestResponseRSocketServer()
	// JSON register
	http.HandleFunc(constants.ONBOARDING_URI, registerHandler)
	log.Printf("Starting server for testing HTTP REGISTER...\n")
	if err := http.ListenAndServe(constants.SEACHER_PORT, nil); err != nil {
		log.Println(err)
	}
}
