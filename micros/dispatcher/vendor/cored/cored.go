package cored

import (
	"log"

	"constants"

	"context"

	"interfaces"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"

	"encoding/json"
	"fmt"
)

func SearchRegisterCore(sHash string) string {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_SEACHER_IP, 7879).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple Request Response.
	// Send request
	result, err := cli.RequestResponse(payload.NewString(sHash, "Metadata")).Block(context.Background())
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Println(err)
	}

	return result.DataUTF8()
}

func SearchServerKeysPublicCore(uuid string) string {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_PUBLIC_IP, 7880).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	result, err := cli.RequestResponse(payload.NewString(uuid, "Metadata")).Block(context.Background())
	if err != nil {
		log.Println(err)
	}

	return result.DataUTF8()
}

func CleanerPublic(uuid string) {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_CLEANER_IP, 7882).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}
	/*
		//Dont kill connection
		defer func() {
			_ = cli.Close()
		}()
	*/
	// Simple FireAndForget.
	cli.FireAndForget(payload.NewString(uuid, "Metadata"))
	if err != nil {
		log.Println(err)
	}
}

func SaveServerKeys(uuid string, publicKHex string, privateKHex string) string {
	var jsonStringToSave = `{"uuid": "` +
		uuid + `", ` +
		`"nextServerPrivateKey": "` +
		privateKHex + `", ` +
		`"nextServerPublicKey": "` +
		publicKHex +
		`"}`

	textPkBytes := []byte(jsonStringToSave)
	sPublicKeyRequest := interfaces.PublicKeyRequest{}
	err := json.Unmarshal(textPkBytes, &sPublicKeyRequest)
	if err != nil {
		fmt.Println(err)
	}

	return SaveServerKeysRegisterCore(jsonStringToSave, uuid)
}

func SaveServerKeysRegisterCore(jsonStringToSave string, uuid string) string {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_PUBLIC_IP, 7880).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	cli.FireAndForget(payload.NewString(jsonStringToSave, "Metadata"))
	if err != nil {
		log.Println(err)
	}

	return constants.OK
}

func SaveServerKeysSign(uuid string, publicKHex string, privateKHex string,
	privateKSign string, publicKSign string) string {
	var jsonStringToSave = `{"uuid": "` +
		uuid + `", ` +
		`"privateSignKey": "` +
		privateKSign + `", ` +
		`"publicSignKey": "` +
		publicKSign + `", ` +
		`"nextServerPrivateKey": "` +
		privateKHex + `", ` +
		`"nextServerPublicKey": "` +
		publicKHex +
		`"}`

	textPkBytes := []byte(jsonStringToSave)
	sPublicKeyRequest := interfaces.PublicKeyRequest{}
	err := json.Unmarshal(textPkBytes, &sPublicKeyRequest)
	if err != nil {
		fmt.Println(err)
	}

	return SaveServerKeysRegisterCore(jsonStringToSave, uuid)
}

func SaveLoginCore(jsonStringToSave string) {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_LOGIN_IP, 7881).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple FireAndForget.
	cli.FireAndForget(payload.NewString(jsonStringToSave, "Metadata"))
	if err != nil {
		log.Println(err)
	}

}

func SearchLoginCore(sHash string) string {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_LOGIN_IP, 7881).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple Request Response.
	// Send request
	result, err := cli.RequestResponse(payload.NewString(sHash, "Metadata")).Block(context.Background())
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Println(err)
	}

	return result.DataUTF8()
}

func ConnectRegisterCore(jsonStringToSave string) {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_REGISTER_IP, 7878).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple FireAndForget.
	cli.FireAndForget(payload.NewString(jsonStringToSave, "Metadata"))
	if err != nil {
		log.Println(err)
	}
}

func ClosureJsonString(
	deviceInfo string,
	email string,
	faceInfo string,
	locationInfo string,
	nextFrontPublicKey string,
	uuid string,
	sHash string,
	privateKSign string,
	privateKHex string,
	publicKHex string) func() string {

	return func() string {
		return `{"deviceInfo": "` +
			deviceInfo + `", ` +
			`"email": "` +
			email + `", ` +
			`"faceInfo": "` +
			faceInfo + `", ` +
			`"locationInfo": "` +
			locationInfo + `", ` +
			`"nextFrontPublicKey": "` +
			nextFrontPublicKey + `", ` +
			`"step": "02", ` +
			`"uuid": "` +
			uuid + `", ` +
			`"key": "` +
			sHash + `", ` +
			`"privateSignKey": "` +
			privateKSign + `", ` +
			`"nextServerPrivateKey": "` +
			privateKHex + `", ` +
			`"nextServerPublicKey": "` +
			publicKHex +
			`"}`
	}
}

func ClosureKeyStringToSave(deviceInfo string,
	faceInfo string,
	locationInfo string) func() string {

	return func() string {
		return `{"deviceInfo": "` +
			deviceInfo + `", ` +
			`"faceInfo": "` +
			faceInfo + `", ` +
			`"locationInfo": "` +
			locationInfo +
			`"}`
	}
}

func Step02Objects(deviceInfo string,
	email string,
	faceInfo string,
	locationInfo string,
	nextFrontPublicKey string,
	uuid string,
	sHash string,
	privateKHex string,
	publicKHex string) (string, string) {

	var jsonStringFunc = ClosureJsonString(
		deviceInfo,
		email,
		faceInfo,
		locationInfo,
		nextFrontPublicKey,
		uuid,
		constants.EMPTY,
		constants.EMPTY,
		constants.EMPTY,
		publicKHex)

	var jsonStringToSaveFunc = ClosureJsonString(
		deviceInfo,
		email,
		faceInfo,
		locationInfo,
		nextFrontPublicKey,
		uuid,
		sHash,
		constants.EMPTY,
		privateKHex,
		publicKHex)

	return jsonStringFunc(), jsonStringToSaveFunc()
}

func GetJsonKeyStringToFindFunc(deviceInfo string,
	faceInfo string,
	locationInfo string) func() string {

	return func() string {
		return `{"deviceInfo": "` +
			deviceInfo + `", ` +
			`"faceInfo": "` +
			faceInfo + `", ` +
			`"locationInfo": "` +
			locationInfo +
			`"}`
	}
}
func GetJsonKeyLoginStringToHashFunc(deviceInfo string,
	faceInfo string,
	consumerKey string,
	locationInfo string) func() string {

	return func() string {
		return `{"deviceInfo": "` +
			deviceInfo + `", ` +
			`"faceInfo": "` +
			faceInfo + `", ` +
			`"clientId": "` +
			consumerKey + `", ` +
			`"locationInfo": "` +
			locationInfo +
			`"}`
	}
}

func GetJsonKeyLoginStringToReturnFunc(
	deviceInfo string,
	email string,
	faceInfo string,
	locationInfo string,
	jwt string,
	nextFrontPublicKey string,
	uuid string,
	sHash string,
	privateKSign string,
	privateKHex string,
	publicKHex string) func() string {

	return func() string {
		ret := `{"deviceInfo": "` +
			deviceInfo + `", ` +
			`"email": "` +
			email + `", ` +
			`"faceInfo": "` +
			faceInfo + `", ` +
			`"locationInfo": "` +
			locationInfo + `", ` +
			`"jwt": "` +
			jwt + `", ` +
			`"login": "` +
			sHash + `", ` +
			`"privateSignKey": "` +
			privateKSign + `", ` +
			`"nextFrontPublicKey": "` +
			nextFrontPublicKey + `", ` +
			`"uuid": "` +
			uuid + `", ` +
			`"nextServerPrivateKey": "` +
			privateKHex + `", ` +
			`"nextServerPublicKey": "` +
			publicKHex +
			`"}`

		return ret
	}

}

func GetJsonKeyLoginStringToSaveFunc(deviceInfo string,
	faceInfo string,
	consumerKey string,
	sHash string,
	locationInfo string) func() string {

	return func() string {
		return `{"deviceInfo": "` +
			deviceInfo + `", ` +
			`"faceInfo": "` +
			faceInfo + `", ` +
			`"clientId": "` +
			consumerKey + `", ` +
			`"login": "` +
			sHash + `", ` +
			`"locationInfo": "` +
			locationInfo +
			`"}`
	}
}

func GetResGeneralErrorKeyValuePayloadFunc(
	key string,
	message string,
	code string,
	request string) func() string {

	return func() string {
		return `{"key": "` +
			key + `", ` +
			`"value": {` +
			`"message": "` +
			message + `", ` +
			`"code": "` +
			code + `"}, ` +
			`"request": "` +
			request + `" ` +
			`}`
	}
}

func ConnectErrorCore(jsonStringToSave string) string {
	// Connect to rsocket server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("Hello", "World")).
		Transport(rsocket.TCPClient().SetHostAndPort(constants.DOCKER_GW_ERROR_IP, 7883).Build()).
		Start(context.Background())
	if err != nil {
		log.Println(err)
	}

	// Simple Request Response.
	result, err := cli.RequestResponse(
		payload.NewString(
			jsonStringToSave,
			"Metadata")).Block(context.Background())
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Println(err)
	}

	return result.DataUTF8()
}
