package main

import (
	"log"

	"net/http"

	"crypto/tls"

	"constants"

	"client"
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
func registerProcess(input string) {

	if !utils.ExistsOnboardingKey(input) {
		//1. Dynamic client registration
		registerResponse := client.PostRegisterWso2App(
			clienteHTTP)

		//2. Get openid access token
		accessToken := client.GetAccessToken(
			clienteHTTP,
			registerResponse)

		//3. Create random application
		application := client.GenerateRandonApp(
			clienteHTTP,
			accessToken.
				Access_token)

		//4. Generate keys for application
		keys := client.GenerateKeysApp(
			clienteHTTP,
			accessToken.Access_token,
			application.ApplicationId)
		log.Println(keys.ConsumerKey)

		//5. Add subscriptions for application
		subscription := client.AddSuscription(
			clienteHTTP,
			accessToken.Access_token,
			application.ApplicationId)
		log.Println(subscription.SubscriptionId)

		//6. Get token simple user app
		userAccessToken := client.TestUserGetAccessToken(
			clienteHTTP,
			keys)
		log.Println(userAccessToken.Access_token)

		utils.SaveForSecondCall(
			input,
			application,
			keys,
			subscription,
			userAccessToken.Access_token)
	}

}

func SearchProcess(sHash string) string {

	return utils.SearchProcess(sHash)
}

func initRSocketServer() {

	err := rsocket.Receive().
		Acceptor(func(
			ctx context.Context,
			setup payload.SetupPayload,
			sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// bind responder
			return rsocket.NewAbstractSocket(
				rsocket.FireAndForget(func(msg payload.Payload) {
					registerProcess(string(msg.DataUTF8()))
				}),
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					//fmt.Print(string(setup.DataUTF8()))
					log.Println("sHash: ", msg.DataUTF8())
					var result = payload.NewString(saveServerKeys(msg.DataUTF8()), "Metadato")

					return mono.Just(result)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
		Serve(context.Background())
	log.Println(err)

	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: http client
	registerProcess("")
}

func main() {
	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	initRSocketServer()

	// JSON register
	http.HandleFunc(
		constants.ONBOARDING_URI,
		registerHandler)
	log.Printf("Starting server for testing HTTP REGISTER...\n")
	if err := http.ListenAndServe(constants.REGISTER_PORT, nil); err != nil {
		log.Println(err)
	}
}
