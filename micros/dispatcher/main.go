package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/http"

	"client"
	"constants"
	"cored"
	"cryptod"
	"utils"
)

var clienteHTTP = &http.Client{}

func getPublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Called ")
	uuid, publicKHex, privateKHex := cryptod.CreateCryptoOnboarding()
	privateKSign, publicKSign := cryptod.CreateCryptoSign()

	log.Println(
		cored.SaveServerKeysSign(
			uuid.String(),
			publicKHex,
			privateKHex,
			privateKSign,
			publicKSign))

	w.Header().Set(constants.CONTENT_TYPE,
		constants.APPLICATION_JSON)
	fmt.Fprintf(
		w, `{"payload": "`+
			doubleSimetricEncrypt(
				publicKHex,
				privateKSign,
				uuid.String())+
			`"}`)

	return
}

func doubleSimetricEncrypt(
	publicKHex string,
	privateKSign string,
	uuid string) string {

	return b64.StdEncoding.EncodeToString(
		[]byte(cryptod.EncryptPublicKey(
			[]byte(
				`{"publicKey": "` +
					b64.StdEncoding.EncodeToString(
						[]byte(
							cryptod.EncryptPublicKey(
								[]byte(publicKHex)))) + `", ` +
					`"privateSignKey": "` +
					privateKSign + `", ` +
					`"uuid": "` +
					uuid + `", ` +
					`"jwt": "` +
					b64.StdEncoding.EncodeToString(
						[]byte(
							cryptod.EncryptPublicKey(
								[]byte(client.GetJWT(
									clienteHTTP))))) +
					`"}`))))
}

func onboardingHandler(
	w http.ResponseWriter,
	r *http.Request) {

	var codeHttpError = 200
	w.Header().Set(
		constants.CONTENT_TYPE,
		constants.APPLICATION_JSON)
	ret := utils.ReceiveReques(
		r,
		r.URL.Query().Get("uuid"),
		r.Header.Get("X-Sirius-Signature"))

	errFunctional := utils.FindError(r.URL.Query().Get("uuid"))

	if errFunctional != constants.EMPTY {
		ret = errFunctional
		codeHttpError = utils.MapErrorMessage(errFunctional)
		w.WriteHeader(codeHttpError)
	}
	fmt.Fprintf(w, ret)

	return
}

func testHandler(
	w http.ResponseWriter,
	r *http.Request) {

	var codeHttpError = 200
	w.Header().Set(
		constants.CONTENT_TYPE,
		constants.APPLICATION_JSON)

	ret := utils.HandleTest(
		r,
		r.URL.Query().Get("uuid"),
		r.Header.Get("X-Sirius-Signature"))

	errFunctional := utils.FindError(r.URL.Query().Get("uuid"))

	if errFunctional != constants.EMPTY {
		ret = errFunctional
		codeHttpError = utils.MapErrorMessage(errFunctional)
		w.WriteHeader(codeHttpError)
	}

	fmt.Fprintf(w, ret)

	return
}

func loginHandler(
	w http.ResponseWriter,
	r *http.Request) {

	var codeHttpError = 200
	w.Header().Set(
		constants.CONTENT_TYPE,
		constants.APPLICATION_JSON)

	ret := utils.ReceiveLogin(
		r,
		r.URL.Query().Get("uuid"),
		r.Header.Get("X-Sirius-Signature"))

	errFunctional := utils.FindError(r.URL.Query().Get("uuid"))

	if errFunctional != constants.EMPTY {
		ret = errFunctional
		codeHttpError = utils.MapErrorMessage(errFunctional)
		w.WriteHeader(codeHttpError)
	}
	fmt.Fprintf(w, ret)

	return
}
func main() {

	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cryptod.InitCrypto()
	// Index Handler
	http.HandleFunc(constants.SLASH, utils.Health)

	// JSON testHandler
	http.HandleFunc(constants.TEST_URI, testHandler)

	// JSON onboardingHandler
	http.HandleFunc(constants.ONBOARDING_URI, onboardingHandler)

	// JSON onboardingHandler
	http.HandleFunc(constants.LOGIN_URI, loginHandler)

	// JSON getPublicKeyHandler
	http.HandleFunc(constants.TOKEN_URI_INTERNAL, getPublicKeyHandler)

	fmt.Printf("Starting server for testing HTTP DISPATCHER...\n")
	if err := http.ListenAndServe(
		constants.DISPATCHER_PORT, nil); err != nil {
		log.Println(err)
	}

}
