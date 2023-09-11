package client

import (
	"bytes"
	"constants"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"interfaces"
	"io/ioutil"
	"log"
	"net/http"
	"utils"
)

func GetPublicKey(clienteHTTP *http.Client) interfaces.ResGetPublicKeyPayload {

	bresponse := GetGenericBody(clienteHTTP, constants.HostWso2Local,
		constants.TOKEN_URI_CRYPTO, constants.API_KEY_CRYPTO,
		"GET")

	responseString := string(bresponse)

	textBytes := []byte(responseString)
	resGetPublicKey := interfaces.ResGetPublicKeyPayload{}
	err := json.Unmarshal(textBytes, &resGetPublicKey)
	if err != nil {
		log.Println(err)
	}

	return resGetPublicKey

}

func GetGenericBody(clienteHTTP *http.Client,
	host string,
	uri string,
	authHeader string,
	method string) []byte {

	request, err := http.NewRequest(method, host+uri, nil)
	if err != nil {
		log.Fatalf("Error creanting request getGenericBody: %v", err)
	}

	if authHeader != "" {
		request.Header.Add(constants.API_KEY_HEADER, authHeader)
	}

	request.Header.Add(constants.CONTENT_TYPE, constants.APPLICATION_JSON)

	response, err := clienteHTTP.Do(request)

	if err != nil {
		log.Fatalf("Error receiving response getGenericBody: %v", err)
	}
	bresponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	return bresponse
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New(constants.CIPHERTEXT_TOO_SHORT)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
func DesencrypPublicServerKey(resGetPublicKey interfaces.ResGetPublicKeyPayload) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKey.Payload)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)

	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypJwt(resGetPublicKeyWithUuid interfaces.ResGetPublicKeyWithUuid) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKeyWithUuid.Jwt)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypPk(resGetPublicKeyWithUuid interfaces.ResGetPublicKeyWithUuid) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKeyWithUuid.PublicKey)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func GetJWT(clienteHTTP *http.Client) string {
	bresponse := PostGenericBodySimple(clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_URI_WSO2,
		constants.BASIC_GENERIC_CREDENTIALS,
		constants.POST)

	responseString := string(bresponse)

	textBytes := []byte(responseString)

	jwtResponse := interfaces.JwtResponse{}

	err := json.Unmarshal(textBytes, &jwtResponse)
	if err != nil {
		log.Println(err.Error())
	}
	return jwtResponse.Access_token

}

func PostGenericBodySimple(clienteHTTP *http.Client,
	host string,
	uri string,
	authHeader string,
	method string) []byte {

	body := new(interface{})
	req, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Error codificando usuario como JSON: %v", err)
	}
	request, err := http.NewRequest("POST", host+uri, bytes.NewBuffer(req))
	if err != nil {
		log.Fatalf("Error creando petición postApisName: %v", err)
	}
	/*
		NOT application/json!
		request.Header.Add("Content-Type", "application/json")
	*/
	request.Header.Add(constants.AUTH, authHeader)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Fatalf("Error recibiendo respuesta postApisName: %v", err)
	}
	bresponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	return bresponse

}

func PostReqOnboarding(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string) utils.ResGeneralPayload {

	response := utils.PostOnboarding(clienteHTTP, clearServerJPk, clearServerJwt, sUuid, step)

	return response
}

func PostReqTest(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	login string) utils.ResGeneralPayload {

	response := utils.PostTest(clienteHTTP, clearServerJPk, clearServerJwt, sUuid, login)

	return response
}

func PostReqTestSign(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	login string,
	privateSignKey string) utils.ResGeneralPayload {

	response := utils.PostTestSign(clienteHTTP, clearServerJPk, clearServerJwt, sUuid, login, privateSignKey)

	return response
}

func PostReqLogin(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string) utils.ResGeneralPayload {

	response := utils.PostLogin(clienteHTTP, clearServerJPk, clearServerJwt, sUuid, step)

	return response
}
func PostReqLoginSign(
	clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string,
	privateSignKey string) utils.ResGeneralPayload {

	response := utils.PostLoginSign(clienteHTTP, clearServerJPk, clearServerJwt, sUuid, step, privateSignKey)

	return response
}

func DesencrypJwtWithSign(resGetPublicKeyWithUuidAndSign interfaces.ResGetPublicKeyWithUuidAndSign) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKeyWithUuidAndSign.Jwt)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypPkWithSign(resGetPublicKeyWithUuidAndSign interfaces.ResGetPublicKeyWithUuidAndSign) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKeyWithUuidAndSign.PublicKey)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func PostReqOnboardingSign(
	clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string,
	privateSignKey string) utils.ResGeneralPayload {

	return utils.PostOnboardingSign(
		clienteHTTP,
		clearServerJPk,
		clearServerJwt,
		sUuid,
		step,
		privateSignKey)
}
