package utils

import (
	"constants"
	"fmt"

	"log"
	"net/http"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"

	b64 "encoding/base64"

	"bytes"

	"strings"

	ecies "github.com/ecies/go"
	uuid "github.com/google/uuid"

	//http2curl "moul.io/http2curl"
	"context"

	"time"

	"interfaces"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var InitPublicServerKeyHex string
var InitPrivateServerKey *ecies.PrivateKey

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

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
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
	log.Println(constants.KEY_PAIR_GENERATED)

	/* Get Hex Public Key. The client can get it by simple REST call
	   because when the server restart it changes */
	InitPublicServerKey := PairKeys.PublicKey
	InitPrivateServerKey = PairKeys
	InitPublicServerKeyHex = InitPublicServerKey.Hex(true)

}

func GetOnboardingReq(r *http.Request) interfaces.ResGeneralPayload {
	brequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	requestString := string(brequest)
	textBytes := []byte(requestString)

	resGeneralPayload := interfaces.ResGeneralPayload{}
	err = json.Unmarshal(textBytes, &resGeneralPayload)
	if err != nil {
		fmt.Println(err)
	}

	return resGeneralPayload

}

func DecryptEcies(textBytes []byte) string {

	plaintext, err := ecies.Decrypt(InitPrivateServerKey, textBytes)
	if err != nil {
		fmt.Println(err)
	}
	return string(plaintext)
}

func DesencrypGenericPayload(resGeneralPayload interfaces.ResGeneralPayload) string {
	sDec, _ := b64.StdEncoding.DecodeString(resGeneralPayload.Payload)

	//plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)

	textBytes := []byte(sDec)
	plaintext := DecryptEcies(textBytes)

	return plaintext
}
func ReceiveReques(r *http.Request) string {
	onboardingReq := GetOnboardingReq(r)

	des := DesencrypGenericPayload(onboardingReq)

	textBytes := []byte(des)

	sOnboardingReq := interfaces.OnboardingReq{}
	err := json.Unmarshal(textBytes, &sOnboardingReq)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf(sOnboardingReq.DeviceInfo)
	uuid, publicKHex, privateKHex := CreateCryptoOnboarding()
	log.Println(privateKHex)

	var jsonString = `{"deviceInfo": "` +
		sOnboardingReq.DeviceInfo + `", ` +
		`"email": "` +
		sOnboardingReq.Email + `", ` +
		`"faceInfo": "` +
		sOnboardingReq.FaceInfo + `", ` +
		`"locationInfo": "` +
		sOnboardingReq.NextFrontPublicKey + `", ` +
		`"uuid": "` +
		uuid.String() + `", ` +
		`"nextServerPublicKey": "` +
		publicKHex +
		`"}`

	var jsonStringToSave = `{"deviceInfo": "` +
		sOnboardingReq.DeviceInfo + `", ` +
		`"email": "` +
		sOnboardingReq.Email + `", ` +
		`"faceInfo": "` +
		sOnboardingReq.FaceInfo + `", ` +
		`"locationInfo": "` +
		sOnboardingReq.NextFrontPublicKey + `", ` +
		`"uuid": "` +
		uuid.String() + `", ` +
		`"nextServerPrivateKey": "` +
		privateKHex + `", ` +
		`"nextServerPublicKey": "` +
		publicKHex +
		`"}`
	log.Println(jsonStringToSave)
	pkFront, err := ecies.NewPublicKeyFromHex(sOnboardingReq.NextFrontPublicKey)
	if err != nil {
		log.Println(err)
	}

	ciphertext, err := ecies.Encrypt(pkFront, []byte(jsonString))
	if err != nil {
		log.Println(err)
	}
	var jsonStringEncode = b64.StdEncoding.EncodeToString([]byte(ciphertext))
	var jsonPayload = `{"payload": "` +
		jsonStringEncode + `"}`

	return jsonPayload
}

func CreateCryptoOnboarding() (uuid.UUID, string, string) {
	uuidWithHyphen := uuid.New()

	InitCrypto()
	InitPrivateServerKeyHex := InitPrivateServerKey.Hex()
	return uuidWithHyphen, InitPublicServerKeyHex, InitPrivateServerKeyHex
}

func PostRegister(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	body := new(interfaces.ReqRegisterPayload)
	body.CallbackUrl = `www.google.lk`
	body.ClientName = `rest_api_devportal`
	body.Owner = `admin`
	body.GrantType = `client_credentials password refresh_token`
	body.SaasApp = true

	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST", host+uri, bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")

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

func PostAccessToken(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	request, err := http.NewRequest("POST", host+constants.TOKEN_ADMIN_URI_WSO2, nil)
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}

	/*
		NOT application/json!
		request.Header.Add("Content-Type", "application/json")
	*/

	request.Header.Add(constants.AUTH, "Basic "+authHeader)

	//TO DEBUG CURL
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/

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

func PostGenerateRandonApp(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	body := new(interfaces.GenerateRandonAppRequest)
	name, err := uuid.New().MarshalText()
	body.Name = string(name)
	body.ThrottlingPolicy = "Unlimited"
	body.Description = "Sirius custom application"
	body.TokenType = "JWT"

	req, err := json.Marshal(body)

	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST", host+uri, bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")

	request.Header.Add(constants.AUTH, "Bearer "+authHeader)

	//TO DEBUG CURL
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/

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

func PostGenerateKeysApp(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string,
	applicationId string) []byte {

	body := new(interfaces.GenerateKeysAppRequest)
	body.KeyType = "PRODUCTION"
	body.KeyManager = "Resident Key Manager"

	body.GrantTypesToBeSupported = []string{"password", "client_credentials"}
	body.CallbackUrl = "http://sample.com/callback/url"

	body.Scopes = []string{"am_application_scope", "default"}
	body.ValidityTime = 3600
	req, err := json.Marshal(body)

	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	resultUri := strings.Replace(uri, "applicationId", applicationId, 1)
	request, err := http.NewRequest("POST", host+resultUri, bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")

	request.Header.Add(constants.AUTH, "Bearer "+authHeader)

	//TO DEBUG CURL
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/
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

func PostAddSuscriptionApp(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string,
	applicationId string) []byte {

	body := new(interfaces.AddSuscriptionAppRequest)
	body.ApplicationId = applicationId
	//REMEMBER: Id API apiDispatcher
	body.ApiId = "c68b681e-5e44-492b-be60-afa185107d8f"

	body.ThrottlingPolicy = "Bronze"
	body.RequestedThrottlingPolicy = "requestedThrottlingPolicy"

	req, err := json.Marshal(body)

	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}

	request, err := http.NewRequest("POST", host+uri, bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}

	request.Header.Add("Content-Type", "application/json")

	request.Header.Add(constants.AUTH, "Bearer "+authHeader)

	//TO DEBUG CURL
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/

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

func SaveForSecondCall(input string,
	application interfaces.GenerateRandonAppRes,
	keys interfaces.GenerateKeysAppRes,
	subscription interfaces.AddSuscriptionAppRes,
	jwt string) {

	var jsonStringToSave = `{"loginInfo": ` +
		input + `, ` +
		`"applicationId": "` +
		application.ApplicationId + `", ` +
		`"consumerKey": "` +
		keys.ConsumerKey + `", ` +
		`"consumerSecret": "` +
		keys.ConsumerSecret + `", ` +
		`"apiId": "` +
		subscription.ApiId + `", ` +
		`"subscriptionId": "` +
		subscription.SubscriptionId + `", ` +
		`"throttlingPolicy": "` +
		subscription.ThrottlingPolicy + `", ` +
		`"jwt": "` +
		jwt + `", ` +
		`"process": "onboarding"` + `, ` +
		`"step": 1}`

	onboarding := interfaces.OnboardingReqPrivate{}
	err := json.Unmarshal([]byte(jsonStringToSave), &onboarding)
	if err != nil {
		log.Println(err)
	}

	initEtcdClient(onboarding.LoginInfo.Key, jsonStringToSave)

}

func initEtcdClient(key string, value string) {
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := cli.Put(ctx, key, value)

	if err != nil {
		switch err {
		case context.Canceled:
			log.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Println("ctx is attached with a deadline is exceeded: %v", err)

		default:
			log.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}

	response, err := cli.Get(ctx, key)

	for _, ev := range response.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
	}

}

func SearchProcess(sHash string) string {
	jsonString := SearchEtcd(sHash)

	return jsonString
}

func SearchEtcd(key string) string {
	key = strings.ReplaceAll(key, "\n", "")

	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	var k string
	var v string
	key = "onboarding" + key

	//Search with prefix
	//REMEMBER
	/*
		resp, _ := cli.Get(ctx, "onboarding", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))

		log.Println("Key: ", key)


		for _, ev := range resp.Kvs {

			k = fmt.Sprintf("%s", ev.Key)
			log.Println(k)
		}
	*/
	resp, _ := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))

	for _, ev := range resp.Kvs {
		k = fmt.Sprintf("%s", ev.Key)
		log.Println("encontrado: ", k)
		v = fmt.Sprintf("%s", ev.Value)
	}
	return v
}

func RenewalJwt(clienteHTTP *http.Client, sOnboarding string) string {

	onboardingPrivate := interfaces.OnboardingReqPrivate{}
	err := json.Unmarshal([]byte(sOnboarding), &onboardingPrivate)
	if err != nil {
		log.Println("err: ", err)

	}
	keys := interfaces.GenerateKeysAppRes{}
	keys.ConsumerKey = onboardingPrivate.ConsumerKey
	keys.ConsumerSecret = onboardingPrivate.ConsumerSecret

	userAccessToken := TestUserGetAccessToken(clienteHTTP, keys)

	jsonCredentialsString := `{"consumerKey": "` +
		keys.ConsumerKey + `", ` +
		`"consumerSecret": "` +
		keys.ConsumerSecret + `", ` +
		`"jwt": "` +
		userAccessToken.Access_token + `"}`

	return jsonCredentialsString

}

func TestUserGetAccessToken(clienteHTTP *http.Client, keys interfaces.GenerateKeysAppRes) interfaces.AccessTokenPayloadRes {
	base64EncodedCredentials := GetUserBase64EncodedCredentials(keys)

	bresponse := PostAccessToken(clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_ADMIN_URI_WSO2,
		base64EncodedCredentials,
		constants.POST)
	responseString := string(bresponse)

	textBytes := []byte(responseString)

	accessTokenPayloadRes := interfaces.AccessTokenPayloadRes{}

	err := json.Unmarshal(textBytes, &accessTokenPayloadRes)
	if err != nil {
		log.Println(err.Error())
	}
	return accessTokenPayloadRes

}

func GetUserBase64EncodedCredentials(keys interfaces.GenerateKeysAppRes) string {
	clientId := keys.ConsumerKey
	clientSecret := keys.ConsumerSecret

	return b64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
}
