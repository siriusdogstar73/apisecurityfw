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

func login(input string) {

	utils.SaveLogin(input)

}

func initRSocketServer() {
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// bind responder
			return rsocket.NewAbstractSocket(
				rsocket.FireAndForget(func(msg payload.Payload) {
					login(string(msg.DataUTF8()))
				}),
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					//fmt.Print(string(setup.DataUTF8()))
					log.Println("login: ", msg.DataUTF8())
					var result = payload.NewString(searchLoginCore(msg.DataUTF8()), "Metadato")

					return mono.Just(result)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7881").Build()).
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
	log.Printf("Starting server for testing HTTP LOGIN...\n")
	if err := http.ListenAndServe(constants.LOGIN_PORT, nil); err != nil {
		log.Println(err)
	}
}

func searchLoginCore(input string) string {
	return utils.SearchLoginCore(input)
}
