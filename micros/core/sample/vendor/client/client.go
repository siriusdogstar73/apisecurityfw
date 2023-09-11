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
	"interfaces"
	"utils"
)

func GetPublicKey(clienteHTTP *http.Client) interfaces.ResGetPublicKeyPayload {

	bresponse := GetGenericBody(
		clienteHTTP,
		constants.HostLocal,
		constants.TOKEN_URI_INTERNAL, "",
		"GET")

	responseString := string(bresponse)

	textBytes := []byte(responseString)
	resGetPublicKey := interfaces.ResGetPublicKeyPayload{}
	err := json.Unmarshal(
		textBytes,
		&resGetPublicKey)
	if err != nil {
		log.Println(err)
	}

	return resGetPublicKey

}

func GetGenericBody(
	clienteHTTP *http.Client,
	host string,
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

func Decrypt(
	ciphertext []byte,
	key []byte) ([]byte, error) {

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
func DesencrypPublicServerKey(
	resGetPublicKey interfaces.ResGetPublicKeyPayload) string {
	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKey.Payload)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)

	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypJwt(resGetPublicKey interfaces.ResGetPublicKey) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGetPublicKey.Jwt)
	plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(plaintext)
}

func DesencrypPk(resGetPublicKey interfaces.ResGetPublicKey) string {

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

	jwtResponse := interfaces.JwtResponse{}

	err := json.Unmarshal(bresponse, &jwtResponse)
	if err != nil {
		log.Println(err.Error())
	}
	return jwtResponse.Access_token

}

func PostGenericBodySimple(
	clienteHTTP *http.Client,
	host string,
	uri string,
	authHeader string,
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

func PostRegisterWso2App(clienteHTTP *http.Client) interfaces.RegisterPayloadRes {

	bresponse := utils.PostRegister(
		clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.CLIENT_REGISTER_URI_WSO2,
		constants.BASIC_GENERIC_ADMIN,
		constants.POST)

	registerPayloadResponse := interfaces.RegisterPayloadRes{}

	err := json.Unmarshal(bresponse, &registerPayloadResponse)
	if err != nil {
		log.Println(err.Error())
	}

	return registerPayloadResponse
}

func GetAccessToken(
	clienteHTTP *http.Client,
	registerResponse interfaces.RegisterPayloadRes) interfaces.AccessTokenPayloadRes {

	base64EncodedCredentials := GetBase64EncodedCredentials(registerResponse)

	bresponse := utils.PostAccessToken(clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_ADMIN_URI_WSO2,
		base64EncodedCredentials,
		constants.POST)

	accessTokenPayloadRes := interfaces.AccessTokenPayloadRes{}

	err := json.Unmarshal(bresponse, &accessTokenPayloadRes)
	if err != nil {
		log.Println(err.Error())
	}
	return accessTokenPayloadRes

}

func GetBase64EncodedCredentials(registerResponse interfaces.RegisterPayloadRes) string {

	return b64.StdEncoding.EncodeToString(
		[]byte(registerResponse.ClientId +
			":" + registerResponse.ClientSecret))
}

func GenerateRandonApp(
	clienteHTTP *http.Client,
	accessToken string) interfaces.GenerateRandonAppRes {

	bresponse := utils.PostGenerateRandonApp(
		clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.URI_GENERATE_RANDOM_APP,
		accessToken,
		constants.POST)

	generateRandonAppRes := interfaces.GenerateRandonAppRes{}

	err := json.Unmarshal(bresponse,
		&generateRandonAppRes)
	if err != nil {
		log.Println(err.Error())
	}
	return generateRandonAppRes

}

func GenerateKeysApp(
	clienteHTTP *http.Client,
	accessToken string,
	applicationId string) interfaces.GenerateKeysAppRes {

	bresponse := utils.PostGenerateKeysApp(
		clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.URI_GENERATE_KEYS_APP,
		accessToken,
		constants.POST,
		applicationId)

	generateKeysAppRes := interfaces.GenerateKeysAppRes{}

	err := json.Unmarshal(bresponse, &generateKeysAppRes)
	if err != nil {
		log.Println(err.Error())
	}
	return generateKeysAppRes

}

func AddSuscription(
	clienteHTTP *http.Client,
	accessToken string,
	applicationId string) interfaces.AddSuscriptionAppRes {

	bresponse := utils.PostAddSuscriptionApp(
		clienteHTTP,
		constants.HostWso2PrivatePortDocker,
		constants.URI_ADD_SUSCRIPTION_APP,
		accessToken,
		constants.POST,
		applicationId)

	addSuscriptionAppRes := interfaces.AddSuscriptionAppRes{}

	err := json.Unmarshal(bresponse, &addSuscriptionAppRes)
	if err != nil {
		log.Println(err.Error())
	}
	return addSuscriptionAppRes

}

func TestUserGetAccessToken(
	clienteHTTP *http.Client,
	keys interfaces.GenerateKeysAppRes) interfaces.AccessTokenPayloadRes {

	bresponse := utils.PostAccessToken(
		clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_ADMIN_URI_WSO2,
		GetUserBase64EncodedCredentials(keys),
		constants.POST)

	accessTokenPayloadRes := interfaces.AccessTokenPayloadRes{}

	err := json.Unmarshal(bresponse, &accessTokenPayloadRes)
	if err != nil {
		log.Println(err.Error())
	}
	return accessTokenPayloadRes

}

func GetUserBase64EncodedCredentials(keys interfaces.GenerateKeysAppRes) string {

	return b64.StdEncoding.EncodeToString(
		[]byte(
			keys.ConsumerKey + ":" +
				keys.ConsumerSecret))
}
