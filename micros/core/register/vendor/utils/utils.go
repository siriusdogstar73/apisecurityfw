package utils

import (
	"constants"
	"fmt"

	"log"
	"net/http"

	"encoding/json"

	"io/ioutil"

	"bytes"

	"strings"

	uuid "github.com/google/uuid"

	//http2curl "moul.io/http2curl"
	"context"

	"time"

	"interfaces"

	clientv3 "go.etcd.io/etcd/client/v3"
)

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
	body.ApiId = "c2f7ed30-4ab9-4f43-98dc-3ac45e8aeec6"

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
		`"process": "onboarding"}`

	onboarding := interfaces.OnboardingReqPrivate{}
	err := json.Unmarshal([]byte(jsonStringToSave), &onboarding)
	if err != nil {
		log.Println(err)
	}

	initEtcdClient("onboarding"+onboarding.LoginInfo.Key, jsonStringToSave)

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
	/*
		response, err := cli.Get(ctx, key)

		for _, ev := range response.Kvs {
			log.Printf("%s : %s\n", ev.Key, ev.Value)
		}
	*/

}

func initEtcdClientWirhControl(key string, value string) string {
	result := "OK"
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
	/*
		response, err := cli.Get(ctx, key)

		for _, ev := range response.Kvs {
			log.Printf("%s : %s\n", ev.Key, ev.Value)
		}
	*/
	if err == nil {
		result = "OK"
	} else {
		result = "KO"
	}
	return result

}

func SearchProcess(sHash string) string {
	jsonString := SearchEtcd(sHash)

	return jsonString
}

func SearchEtcd(key string) string {
	log.Println("key: ", key)

	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	response, err := cli.Get(ctx, key)

	sum := 1
	for sum < 100 {
		sum += sum
		response, err = cli.Get(ctx, key)
	}
	log.Println("response.Count: ", response.Count)
	if err != nil {
		log.Println("err: ", err.Error())
	}
	var res string
	for _, ev := range response.Kvs {
		log.Println("ev.Value: ", ev.Value)
		res = fmt.Sprintf("%s\n", ev.Value)
		log.Println("res: ", res)
	}

	return ""

}

func SaveServerKeys(input string) string {

	sPublicKeyRequest := interfaces.PublicKeyRequest{}
	err := json.Unmarshal([]byte(input), &sPublicKeyRequest)
	if err != nil {
		log.Println(err)
	}
	return initEtcdClientWirhControl("public"+sPublicKeyRequest.Uuid, input)
}

func ExistsOnboardingKey(input string) bool {

	flag := false

	onboardingRequest := interfaces.OnboardingRequest{}
	err := json.Unmarshal([]byte(input), &onboardingRequest)
	if err != nil {
		log.Println(err)
	}

	key := "onboarding" + onboardingRequest.Key

	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	resp, _ := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))

	if resp.Count == 0 {
		flag = false
	} else {
		flag = true
	}

	return flag
}
