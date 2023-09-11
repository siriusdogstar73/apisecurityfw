package utils

import (
	"constants"
	"fmt"

	"log"
	"net/http"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"encoding/json"

	b64 "encoding/base64"

	"bytes"

	"io/ioutil"

	"interfaces"

	ecies "github.com/ecies/go"
	//"moul.io/http2curl"

	"crypto/ecdsa"
	"crypto/elliptic"

	"math/big"

	base58 "github.com/btcsuite/btcutil/base58"
)

var InitPublicServerKeyHex string

var InitPublicFrontKeyHex string
var InitPrivateFrontKey *ecies.PrivateKey

func Health(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != constants.SLASH {
		http.Error(w, constants.NOT_FOUND_TXT, http.StatusNotFound)
		return
	}

	switch r.Method {
	case constants.GET:
		fmt.Fprintf(w, constants.HtmlStr)
	case constants.POST:
		if err := r.ParseForm(); err != nil {
			log.Fatal(constants.FATAL_ERROR_TXT)
			return
		}

	default:
		fmt.Fprintf(w, constants.SORRY_NOT_SOPPORTED)
	}
}

func Encrypt(plaintext []byte,
	key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func EncryptPublicKey(text []byte) []byte {

	ciphertext, err := Encrypt(text, constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err)
	}

	return ciphertext

}

func InitCrypto() {

	/* init crypto */
	/* Generate key pair */
	PairKeys, err := ecies.GenerateKey()
	if err != nil {
		panic(err)
	}
	//log.Println(constants.KEY_PAIR_GENERATED)

	/* Get Hex Public Key. The client can get it by simple REST call
	   because when the server restart it changes */
	InitPublicServerKey := PairKeys.PublicKey
	InitPublicServerKeyHex = InitPublicServerKey.Hex(true)
}

func InitFrontCrypto() {

	/* init crypto */
	/* Generate key pair */
	PairKeys, err := ecies.GenerateKey()
	if err != nil {
		panic(err)
	}
	//log.Println(constants.KEY_PAIR_GENERATED)

	/* Get Hex Public Key. The client can get it by simple REST call
	   because when the server restart it changes */
	InitPublicFrontKey := PairKeys.PublicKey
	InitPrivateFrontKey = PairKeys
	InitPublicFrontKeyHex = InitPublicFrontKey.Hex(true)

}

func PostOnboarding(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string) ResGeneralPayload {

	frontPk := []byte(InitPublicFrontKeyHex)

	pkServer, err := ecies.NewPublicKeyFromHex(clearServerJPk)
	if err != nil {
		log.Println(err)
	}
	input := new(interfaces.OnboardingReq)
	input.DeviceInfo = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"
	input.Email = "sirius@sirius.es"
	input.FaceInfo = "974d1cc1-0f03-4278-8c84-93a28127a53d"
	input.LocationInfo = "40.41131467490109, -3.6932489733328464"
	input.NextFrontPublicKey = string(frontPk)
	input.Step = step

	var jsonPayloadString = `{"deviceInfo": "` +
		input.DeviceInfo + `", ` +
		`"email": "` +
		input.Email + `", ` +
		`"faceInfo": "` +
		input.FaceInfo + `", ` +
		`"locationInfo": "` +
		input.LocationInfo + `", ` +
		`"step": "` +
		input.Step + `", ` +
		`"nextFrontPublicKey": "` +
		input.NextFrontPublicKey +
		`"}`

	ciphertext, err := ecies.Encrypt(pkServer, []byte(jsonPayloadString))
	if err != nil {
		log.Println(err)
	}

	body := new(interfaces.ResGetPublicKeyPayload)
	body.Payload = b64.StdEncoding.EncodeToString(ciphertext)

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST",
		constants.HostWso2Local+constants.ONBOARDING_UR_BY_API_QS+sUuid,
		bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(constants.AUTH, "Bearer "+clearServerJwt)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error in response: %v", err)
	}

	return GetOnboardingReq(response)
}

type ResGeneralPayload struct {
	Payload string `json:"payload"`
}

func GetOnboardingReq(r *http.Response) ResGeneralPayload {

	bresponse, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	responseString := string(bresponse)
	textBytes := []byte(responseString)

	resGeneralPayload := ResGeneralPayload{}
	err = json.Unmarshal(textBytes, &resGeneralPayload)
	if err != nil {
		fmt.Println(err)
	}

	return resGeneralPayload

}

func DecryptEcies(textBytes []byte) ([]byte, error) {

	plaintextByte, err := ecies.Decrypt(InitPrivateFrontKey, textBytes)
	if err != nil {
		return plaintextByte, err
	}
	return plaintextByte, err
}

func DesencryptedOneOnboarding(response string) interfaces.OnboardingSignFront {

	sDec, _ := b64.StdEncoding.DecodeString(response)

	textBytes := []byte(sDec)
	plaintextByte, err := DecryptEcies(textBytes)

	if err != nil {
		log.Println("Error in desencrypt")
	}

	sOnboardingReq := interfaces.OnboardingSignFront{}

	err = json.Unmarshal(plaintextByte, &sOnboardingReq)
	if err != nil {
		log.Println(err.Error())
	}

	return sOnboardingReq
}

func PostTest(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	login string) ResGeneralPayload {

	frontPk := []byte(InitPublicFrontKeyHex)

	pkServer, err := ecies.NewPublicKeyFromHex(clearServerJPk)
	if err != nil {
		log.Println(err)
	}
	input := new(interfaces.OnboardingByLoginFront)
	input.Jwt = clearServerJwt
	input.Uuid = sUuid
	input.Login = login
	input.NextFrontPublicKey = string(frontPk)

	var jsonPayloadString = `{"jwt": "` +
		input.Jwt + `", ` +
		`"uuid": "` +
		input.Uuid + `", ` +
		`"nextFrontPublicKey": "` +
		input.NextFrontPublicKey + `", ` +
		`"login": "` +
		input.Login +
		`"}`

	ciphertext, err := ecies.Encrypt(pkServer, []byte(jsonPayloadString))
	if err != nil {
		log.Println(err)
	}

	body := new(interfaces.ResGetPublicKeyPayload)
	body.Payload = b64.StdEncoding.EncodeToString(ciphertext)

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST",
		constants.HostWso2Local+constants.TEST_UR_BY_API_QS+sUuid,
		bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(constants.AUTH, "Bearer "+clearServerJwt)
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error in response: %v", err)
	}

	return GetOnboardingReq(response)
}

func PostTestSign(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	login string,
	privateSignKey string) ResGeneralPayload {

	frontPk := []byte(InitPublicFrontKeyHex)

	pkServer, err := ecies.NewPublicKeyFromHex(clearServerJPk)
	if err != nil {
		log.Println(err)
	}
	input := new(interfaces.OnboardingByLoginFront)
	input.Jwt = clearServerJwt
	input.Uuid = sUuid
	input.Login = login
	input.NextFrontPublicKey = string(frontPk)

	var jsonPayloadString = `{"jwt": "` +
		input.Jwt + `", ` +
		`"uuid": "` +
		input.Uuid + `", ` +
		`"nextFrontPublicKey": "` +
		input.NextFrontPublicKey + `", ` +
		`"login": "` +
		input.Login +
		`"}`

	ciphertext, err := ecies.Encrypt(pkServer, []byte(jsonPayloadString))
	if err != nil {
		log.Println(err)
	}

	body := new(interfaces.ResGetPublicKeyPayload)
	body.Payload = b64.StdEncoding.EncodeToString(ciphertext)

	signature, err := Sign([]byte(jsonPayloadString), privateSignKey)

	if err != nil {
		log.Println(err)
	}

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST",
		constants.HostWso2Local+constants.TEST_UR_BY_API_QS+sUuid,
		bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(constants.AUTH, "Bearer "+clearServerJwt)
	request.Header.Add("X-Sirius-Signature", signature)
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error in response: %v", err)
	}

	return GetOnboardingReq(response)
}

func PostLogin(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string) ResGeneralPayload {

	frontPk := []byte(InitPublicFrontKeyHex)

	pkServer, err := ecies.NewPublicKeyFromHex(clearServerJPk)
	if err != nil {
		log.Println(err)
	}
	input := new(interfaces.OnboardingReq)
	input.DeviceInfo = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"
	input.Email = "sirius@sirius.es"
	input.FaceInfo = "974d1cc1-0f03-4278-8c84-93a28127a53d"
	input.LocationInfo = "40.41131467490109, -3.6932489733328464"
	input.NextFrontPublicKey = string(frontPk)
	input.Step = step

	var jsonPayloadString = `{"deviceInfo": "` +
		input.DeviceInfo + `", ` +
		`"email": "` +
		input.Email + `", ` +
		`"faceInfo": "` +
		input.FaceInfo + `", ` +
		`"locationInfo": "` +
		input.LocationInfo + `", ` +
		`"step": "` +
		input.Step + `", ` +
		`"nextFrontPublicKey": "` +
		input.NextFrontPublicKey +
		`"}`

	ciphertext, err := ecies.Encrypt(pkServer, []byte(jsonPayloadString))
	if err != nil {
		log.Println(err)
	}

	body := new(interfaces.ResGetPublicKeyPayload)
	body.Payload = b64.StdEncoding.EncodeToString(ciphertext)

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST",
		constants.HostWso2Local+constants.LOGIN_UR_BY_API_QS+sUuid,
		bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(constants.AUTH, "Bearer "+clearServerJwt)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error in response: %v", err)
	}

	return GetOnboardingReq(response)
}

func PostLoginSign(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string,
	privateSignKey string) ResGeneralPayload {

	frontPk := []byte(InitPublicFrontKeyHex)

	pkServer, err := ecies.NewPublicKeyFromHex(clearServerJPk)
	if err != nil {
		log.Println(err)
	}
	input := new(interfaces.OnboardingReq)
	input.DeviceInfo = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"
	input.Email = "sirius@sirius.es"
	input.FaceInfo = "974d1cc1-0f03-4278-8c84-93a28127a53d"
	input.LocationInfo = "40.41131467490109, -3.6932489733328464"
	input.NextFrontPublicKey = string(frontPk)
	input.Step = step
	input.PrivateSignKey = privateSignKey

	var jsonPayloadString = `{"deviceInfo": "` +
		input.DeviceInfo + `", ` +
		`"email": "` +
		input.Email + `", ` +
		`"faceInfo": "` +
		input.FaceInfo + `", ` +
		`"locationInfo": "` +
		input.LocationInfo + `", ` +
		`"step": "` +
		input.Step + `", ` +
		`"privateSignKey": "` +
		input.PrivateSignKey + `", ` +
		`"nextFrontPublicKey": "` +
		input.NextFrontPublicKey +
		`"}`

	ciphertext, err := ecies.Encrypt(pkServer, []byte(jsonPayloadString))
	if err != nil {
		log.Println(err)
	}

	body := new(interfaces.ResGetPublicKeyPayload)
	body.Payload = b64.StdEncoding.EncodeToString(ciphertext)

	signature, err := Sign([]byte(jsonPayloadString), privateSignKey)

	if err != nil {
		log.Println(err)
	}

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST",
		constants.HostWso2Local+constants.LOGIN_UR_BY_API_QS+sUuid,
		bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(constants.AUTH, "Bearer "+clearServerJwt)
	request.Header.Add("X-Sirius-Signature", signature)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error in response: %v", err)
	}

	return GetOnboardingReq(response)
}

func Sign(data []byte,
	privateKey string) (string, error) {

	var pri ecdsa.PrivateKey
	pri.D = new(big.Int).SetBytes(base58.Decode(privateKey))
	pri.PublicKey.Curve = elliptic.P256()
	pri.PublicKey.X, pri.PublicKey.Y = pri.PublicKey.Curve.ScalarBaseMult(pri.D.Bytes())

	r, s, err := ecdsa.Sign(rand.Reader, &pri, data)
	if err != nil {
		return "", err
	}
	pubKeyStr := base58.Encode(elliptic.MarshalCompressed(elliptic.P256(), pri.PublicKey.X, pri.PublicKey.Y))
	sig2 := fmt.Sprintf("%s.%s.%s", base58.Encode(r.Bytes()), base58.Encode(s.Bytes()), pubKeyStr)
	return sig2, nil
}

func PostOnboardingSign(clienteHTTP *http.Client,
	clearServerJPk string,
	clearServerJwt string,
	sUuid string,
	step string,
	privateSignKey string) ResGeneralPayload {

	frontPk := []byte(InitPublicFrontKeyHex)

	pkServer, err := ecies.NewPublicKeyFromHex(clearServerJPk)
	if err != nil {
		log.Println(err)
	}
	input := new(interfaces.OnboardingReq)
	input.DeviceInfo = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"
	input.Email = "sirius@sirius.es"
	input.FaceInfo = "974d1cc1-0f03-4278-8c84-93a28127a53d"
	input.LocationInfo = "40.41131467490109, -3.6932489733328464"
	input.NextFrontPublicKey = string(frontPk)
	input.Step = step
	input.PrivateSignKey = privateSignKey

	var jsonPayloadString = `{"deviceInfo": "` +
		input.DeviceInfo + `", ` +
		`"email": "` +
		input.Email + `", ` +
		`"faceInfo": "` +
		input.FaceInfo + `", ` +
		`"locationInfo": "` +
		input.LocationInfo + `", ` +
		`"step": "` +
		input.Step + `", ` +
		`"privateSignKey": "` +
		input.PrivateSignKey + `", ` +
		`"nextFrontPublicKey": "` +
		input.NextFrontPublicKey +
		`"}`

	ciphertext, err := ecies.Encrypt(pkServer, []byte(jsonPayloadString))
	if err != nil {
		log.Println(err)
	}

	body := new(interfaces.ResGetPublicKeyPayload)
	body.Payload = b64.StdEncoding.EncodeToString(ciphertext)

	signature, err := Sign([]byte(jsonPayloadString), privateSignKey)

	if err != nil {
		log.Println(err)
	}

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST",
		constants.HostWso2Local+constants.ONBOARDING_UR_BY_API_QS+sUuid,
		bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(constants.AUTH, "Bearer "+clearServerJwt)
	request.Header.Add("X-Sirius-Signature", signature)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error in response: %v", err)
	}

	return GetOnboardingReq(response)
}
