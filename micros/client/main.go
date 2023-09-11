package main

import (
	"client"
	"constants"
	"crypto/tls"
	"encoding/json"
	"flag"
	"interfaces"
	"log"
	"net/http"
	"time"
	"utils"
)

func main() {
	start := time.Now()
	var login string
	var uuid string
	var psk string
	var jwt string
	var privateKeySign string

	clienteHTTP := &http.Client{}
	clienteHTTP.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	utils.InitFrontCrypto()
	operation := flag.String("operation", "health", "a string")

	flag.Parse()
	var clearServerJwt string
	var clearServerJPk string

	if *operation == "onboarding" {
		log.Printf("Call publicKey")
		clearServerPublicKey :=
			client.DesencrypPublicServerKey(client.GetPublicKey(clienteHTTP))

		textBytes := []byte(clearServerPublicKey)

		resGetPublicKeyWithUuidAndSign := interfaces.ResGetPublicKeyWithUuidAndSign{}
		err := json.Unmarshal(textBytes, &resGetPublicKeyWithUuidAndSign)
		if err != nil {
			log.Println(err)
		}

		clearServerJwt = client.DesencrypJwtWithSign(resGetPublicKeyWithUuidAndSign)
		clearServerJPk = client.DesencrypPkWithSign(resGetPublicKeyWithUuidAndSign)

		//First onboarding call
		utils.InitFrontCrypto()
		response := client.PostReqOnboardingSign(
			clienteHTTP,
			string(clearServerJPk),
			string(clearServerJwt),
			resGetPublicKeyWithUuidAndSign.Uuid,
			constants.Step01,
			resGetPublicKeyWithUuidAndSign.PrivateSignKey)
		desencryptedResponse := utils.DesencryptedOneOnboarding(response.Payload)

		//Second onboarding call
		//time.Sleep(1 * time.Second)
		utils.InitFrontCrypto()

		response = client.PostReqOnboardingSign(
			clienteHTTP,
			desencryptedResponse.NextServerPublicKey,
			string(clearServerJwt),
			desencryptedResponse.Uuid,
			desencryptedResponse.Step,
			desencryptedResponse.PrivateSignKey)

		desencryptedResponse = utils.DesencryptedOneOnboarding(response.Payload)

		elapsed := time.Since(start)
		log.Printf("Onboarding syncr took %f", elapsed.Seconds())
	}

	if *operation == "login" || *operation == "test" {
		log.Printf("Call publicKey")
		clearServerPublicKey :=
			client.DesencrypPublicServerKey(client.GetPublicKey(clienteHTTP))

		textBytes := []byte(clearServerPublicKey)

		resGetPublicKeyWithUuidAndSign := interfaces.ResGetPublicKeyWithUuidAndSign{}
		err := json.Unmarshal(textBytes, &resGetPublicKeyWithUuidAndSign)
		if err != nil {
			log.Println(err)
		}

		clearServerJwt = client.DesencrypJwtWithSign(resGetPublicKeyWithUuidAndSign)
		clearServerJPk = client.DesencrypPkWithSign(resGetPublicKeyWithUuidAndSign)

		//Login call
		utils.InitFrontCrypto()
		response := client.PostReqLoginSign(clienteHTTP,
			string(clearServerJPk),
			string(clearServerJwt),
			resGetPublicKeyWithUuidAndSign.Uuid,
			constants.Step01,
			resGetPublicKeyWithUuidAndSign.PrivateSignKey)

		desencryptedResponse := utils.DesencryptedOneOnboarding(response.Payload)

		login = desencryptedResponse.Login
		uuid = desencryptedResponse.Uuid
		psk = desencryptedResponse.NextServerPublicKey
		jwt = desencryptedResponse.Jwt

		privateKeySign = desencryptedResponse.PrivateSignKey
		elapsed := time.Since(start)
		log.Printf("Login syncr took %f", elapsed.Seconds())

	}
	if *operation == "test" {
		//Consumption call
		utils.InitFrontCrypto()

		response := client.PostReqTestSign(
			clienteHTTP,
			psk,
			jwt,
			uuid,
			login,
			privateKeySign)
		desencryptedResponse := utils.DesencryptedOneOnboarding(response.Payload)
		log.Println(desencryptedResponse.PrivateSignKey)
		elapsed := time.Since(start)
		log.Printf("Consumption syncr took %f", elapsed.Seconds())
	}
}
