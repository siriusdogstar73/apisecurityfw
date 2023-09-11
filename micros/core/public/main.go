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

func saveServerKeys(input string) string {
	return utils.SaveServerKeys(input)
}
func publicProcess(input string) {

	utils.SaveServerKeys(input)

}

func SearchProcess(sHash string) string {

	return utils.SearchProcess(sHash)
}

func initRSocketServer() {
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// bind responder
			return rsocket.NewAbstractSocket(
				rsocket.FireAndForget(func(msg payload.Payload) {
					publicProcess(string(msg.DataUTF8()))
				}),
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					//fmt.Print(string(setup.DataUTF8()))
					log.Println("uuid: ", msg.DataUTF8())
					var result = payload.NewString(searchServerKeys(msg.DataUTF8()), "Metadato")

					return mono.Just(result)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7880").Build()).
		Serve(context.Background())
	log.Println(err)

	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: http client
	publicProcess("")
}

func main() {
	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	initRSocketServer()

	// JSON register
	http.HandleFunc(constants.ONBOARDING_URI, registerHandler)
	log.Printf("Starting server for testing HTTP REGISTER...\n")
	if err := http.ListenAndServe(constants.REGISTER_PORT, nil); err != nil {
		log.Println(err)
	}
}

func searchServerKeys(input string) string {
	return utils.SearchServerKeys(input)
}
