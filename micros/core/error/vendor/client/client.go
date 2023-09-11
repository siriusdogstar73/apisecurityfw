package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"constants"
	"utils"
)

func GetPublicKey(clienteHTTP *http.Client) utils.ResGetPublicKeyPayload {

	bresponse := GetGenericBody(clienteHTTP, constants.HostLocal,
		constants.TOKEN_URI_INTERNAL, "",
		"GET")

	responseString := string(bresponse)

	textBytes := []byte(responseString)
	resGetPublicKey := utils.ResGetPublicKeyPayload{}
	err := json.Unmarshal(textBytes, &resGetPublicKey)
	if err != nil {
		log.Println(err)
	}

	return resGetPublicKey

}

func GetGenericBody(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	request, err := http.NewRequest(method, host+uri, nil)
	if err != nil {
		log.Fatalf("Error creanting request getGenericBody: %v", err)
	}

	if authHeader != "" {
		request.Header.Add(constants.AUTH, authHeader)
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
func DesencrypPublicServerKey(resGetPublicKey utils.ResGetPublicKeyPayload) string {
	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKey.Payload)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)

	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypJwt(resGetPublicKey utils.ResGetPublicKey) string {
	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKey.Jwt)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypPk(resGetPublicKey utils.ResGetPublicKey) string {
	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKey.PublicKey)
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

	jwtResponse := utils.JwtResponse{}

	err := json.Unmarshal(textBytes, &jwtResponse)
	if err != nil {
		log.Println(err.Error())
	}
	return jwtResponse.Access_token

}

func PostGenericBodySimple(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	body := new(interface{})
	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}
	request, err := http.NewRequest("POST", host+uri, bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}
	/*
		NOT application/json!
		request.Header.Add("Content-Type", "application/json")
	*/
	request.Header.Add(constants.AUTH, authHeader)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error recibiendo respuesta postApisName: %v", err)
	}
	bresponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response: %v", err)
	}

	return bresponse

}

func PostRegisterWso2App(clienteHTTP *http.Client) utils.RegisterPayloadRes {

	bresponse := utils.PostRegister(clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.CLIENT_REGISTER_URI_WSO2,
		constants.BASIC_GENERIC_ADMIN,
		constants.POST)
	responseString := string(bresponse)

	textBytes := []byte(responseString)

	registerPayloadResponse := utils.RegisterPayloadRes{}

	err := json.Unmarshal(textBytes, &registerPayloadResponse)
	if err != nil {
		log.Println(err.Error())
	}

	return registerPayloadResponse
}

func GetAccessToken(clienteHTTP *http.Client, registerResponse utils.RegisterPayloadRes) utils.AccessTokenPayloadRes {

	base64EncodedCredentials := GetBase64EncodedCredentials(registerResponse)

	bresponse := utils.PostAccessToken(clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_ADMIN_URI_WSO2,
		base64EncodedCredentials,
		constants.POST)
	responseString := string(bresponse)

	textBytes := []byte(responseString)

	accessTokenPayloadRes := utils.AccessTokenPayloadRes{}

	err := json.Unmarshal(textBytes, &accessTokenPayloadRes)
	if err != nil {
		log.Println(err.Error())
	}
	return accessTokenPayloadRes

}

func GetBase64EncodedCredentials(registerResponse utils.RegisterPayloadRes) string {
	clientId := registerResponse.ClientId
	clientSecret := registerResponse.ClientSecret

	return b64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
}

func GenerateRandonApp(clienteHTTP *http.Client, accessToken string) utils.GenerateRandonAppRes {

	bresponse := utils.PostGenerateRandonApp(clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.URI_GENERATE_RANDOM_APP,
		accessToken,
		constants.POST)

	responseString := string(bresponse)

	textBytes := []byte(responseString)

	generateRandonAppRes := utils.GenerateRandonAppRes{}

	err := json.Unmarshal(textBytes, &generateRandonAppRes)
	if err != nil {
		log.Println(err.Error())
	}
	return generateRandonAppRes

}

func GenerateKeysApp(clienteHTTP *http.Client, accessToken string, applicationId string) utils.GenerateKeysAppRes {

	bresponse := utils.PostGenerateKeysApp(clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.URI_GENERATE_KEYS_APP,
		accessToken,
		constants.POST,
		applicationId)

	responseString := string(bresponse)

	textBytes := []byte(responseString)

	generateKeysAppRes := utils.GenerateKeysAppRes{}

	err := json.Unmarshal(textBytes, &generateKeysAppRes)
	if err != nil {
		log.Println(err.Error())
	}
	return generateKeysAppRes

}

func AddSuscription(clienteHTTP *http.Client, accessToken string, applicationId string) utils.AddSuscriptionAppRes {

	bresponse := utils.PostAddSuscriptionApp(clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.URI_ADD_SUSCRIPTION_APP,
		accessToken,
		constants.POST,
		applicationId)

	responseString := string(bresponse)

	textBytes := []byte(responseString)

	addSuscriptionAppRes := utils.AddSuscriptionAppRes{}

	err := json.Unmarshal(textBytes, &addSuscriptionAppRes)
	if err != nil {
		log.Println(err.Error())
	}
	return addSuscriptionAppRes

}

func TestUserGetAccessToken(clienteHTTP *http.Client, keys utils.GenerateKeysAppRes) utils.AccessTokenPayloadRes {
	base64EncodedCredentials := GetUserBase64EncodedCredentials(keys)

	bresponse := utils.PostAccessToken(clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_ADMIN_URI_WSO2,
		base64EncodedCredentials,
		constants.POST)
	responseString := string(bresponse)

	textBytes := []byte(responseString)

	accessTokenPayloadRes := utils.AccessTokenPayloadRes{}

	err := json.Unmarshal(textBytes, &accessTokenPayloadRes)
	if err != nil {
		log.Println(err.Error())
	}
	return accessTokenPayloadRes

}

func GetUserBase64EncodedCredentials(keys utils.GenerateKeysAppRes) string {
	clientId := keys.ConsumerKey
	clientSecret := keys.ConsumerSecret

	return b64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
}
