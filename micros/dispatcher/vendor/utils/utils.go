package utils

import (
	"constants"
	"fmt"

	"log"
	"net/http"

	"encoding/json"

	"io/ioutil"

	b64 "encoding/base64"

	ecies "github.com/ecies/go"

	"interfaces"

	"cored"
	"cryptod"
	"errord"
	"strconv"
)

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

func GetOnboardingReq(
	r *http.Request,
	sUuid string) interfaces.ResGeneralPayload {

	brequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	resGeneralPayload := interfaces.ResGeneralPayload{}
	err = json.Unmarshal(brequest, &resGeneralPayload)
	if err != nil {
		fmt.Println(err)
	}

	if resGeneralPayload.Payload == constants.EMPTY {
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					constants.PAYLOAD_FAULT)()))
	}

	return resGeneralPayload

}

func GetOnboardingReqTest(r *http.Request) interfaces.ResGeneralPayload {

	brequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	resGeneralPayload := interfaces.ResGeneralPayload{}
	err = json.Unmarshal(brequest, &resGeneralPayload)
	if err != nil {
		fmt.Println(err)
	}

	return resGeneralPayload

}

func ReceiveReques(r *http.Request,
	sUuid string,
	signature string) string {

	var onboardingReq interfaces.ResGeneralPayload
	var jsonPayload string

	if sUuid == constants.EMPTY {
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					constants.UUID_FAULT_CODE,
					constants.UUID_FAULT,
					constants.BAD_REQUEST_CODE,
					constants.PAYLOAD_NA)()))
	} else {
		onboardingReq = GetOnboardingReq(r, sUuid)

		sOnboardingReq := cryptod.DesencrypGenericPayload(
			onboardingReq,
			sUuid,
			signature)

		if sOnboardingReq.Step == constants.Step02 {
			jsonPayload = Step02Onboarding(r, sUuid, sOnboardingReq)

		} else if sOnboardingReq.Step == constants.Step01 {
			jsonPayload = Step01Onboarding(r, sUuid, sOnboardingReq)
		}
	}

	return jsonPayload
}

func Step01Onboarding(r *http.Request,
	sUuid string,
	sOnboardingReq interfaces.OnboardingReq) string {

	uuid, publicKHex, privateKHex := cryptod.CreateCryptoOnboarding()
	privateKSign, publicKSign := cryptod.CreateCryptoSign()
	log.Println(
		cored.SaveServerKeysSign(
			uuid.String(),
			publicKHex,
			privateKHex,
			privateKSign,
			publicKSign))

	return getJsonPayload(
		sOnboardingReq.NextFrontPublicKey,
		cored.ClosureJsonString(
			sOnboardingReq.DeviceInfo,
			sOnboardingReq.Email,
			sOnboardingReq.FaceInfo,
			sOnboardingReq.LocationInfo,
			sOnboardingReq.NextFrontPublicKey,
			uuid.String(),
			constants.EMPTY,
			privateKSign,
			constants.EMPTY,
			publicKHex)())
}

func getJsonPayload(nextFrontPublicKey string, jsonString string) string {

	pkFront, err := ecies.NewPublicKeyFromHex(nextFrontPublicKey)
	if err != nil {
		log.Println(err)
	}
	ciphertext, err := ecies.Encrypt(pkFront, []byte(jsonString))
	if err != nil {
		log.Println(err)
	}

	return `{"payload": "` +
		b64.StdEncoding.EncodeToString([]byte(ciphertext)) +
		`"}`
}

func Step02Onboarding(r *http.Request,
	sUuid string,
	sOnboardingReq interfaces.
		OnboardingReq) string {

	uuid, publicKHex, privateKHex := cryptod.CreateCryptoOnboarding()

	var jsonKeyStringToSaveFunc = cored.ClosureKeyStringToSave(
		sOnboardingReq.
			DeviceInfo,
		sOnboardingReq.FaceInfo,
		sOnboardingReq.LocationInfo)

	validateJsonStringToAdd(cored.SearchRegisterCore(
		cryptod.
			GetHashFromString(jsonKeyStringToSaveFunc())), sUuid)

	var jsonStringFunc, jsonStringToSaveFunc = cored.
		Step02Objects(
			sOnboardingReq.DeviceInfo,
			sOnboardingReq.Email,
			sOnboardingReq.FaceInfo,
			sOnboardingReq.LocationInfo,
			sOnboardingReq.NextFrontPublicKey,
			uuid.String(),
			cryptod.GetHashFromString(jsonKeyStringToSaveFunc()),
			privateKHex,
			publicKHex)

	cored.ConnectRegisterCore(jsonStringToSaveFunc)

	return getJsonPayload(
		sOnboardingReq.NextFrontPublicKey,
		jsonStringFunc)
}

func validateJsonStringToAdd(
	jsonStringToAdd string,
	sUuid string) {

	onboardingPrivate := interfaces.OnboardingReqPrivate{}
	err := json.Unmarshal([]byte(jsonStringToAdd), &onboardingPrivate)
	if err != nil {
		log.Println("err: ", err)
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					constants.PAYLOAD_FAULT)()))
	}
}

func HandleTest(
	r *http.Request,
	sUuid string,
	signature string) string {

	var ciphertext = []byte(constants.EMPTY)
	resp := constants.EMPTY

	if sUuid == constants.EMPTY {
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					constants.UUID_FAULT_CODE,
					constants.UUID_FAULT,
					constants.BAD_REQUEST_CODE,
					constants.PAYLOAD_NA)()))
	} else {

		request := GetRequestHandleTest(
			r,
			sUuid,
			signature)

		pkFront, err := ecies.NewPublicKeyFromHex(
			request.NextFrontPublicKey)

		if GetSavedTest(
			r,
			sUuid,
			signature,
			request).DeviceInfo != constants.EMPTY {

			uuid, publicKHex := GetCryptoTest(r, sUuid, signature)
			privateKSign, _ := cryptod.CreateCryptoSign()

			ciphertext, err = ecies.Encrypt(
				pkFront,
				[]byte(
					cored.GetJsonKeyLoginStringToReturnFunc(
						GetSavedTest(
							r,
							sUuid,
							signature,
							request).DeviceInfo,
						constants.EMPTY,
						GetSavedTest(
							r,
							sUuid,
							signature,
							request).FaceInfo,
						GetSavedTest(
							r,
							sUuid,
							signature,
							request).LocationInfo,
						request.Jwt,
						request.NextFrontPublicKey,
						uuid,
						GetSavedTest(
							r,
							sUuid,
							signature,
							request).Login,
						privateKSign,
						constants.EMPTY,
						publicKHex)()))

			if err != nil {
				log.Println(err)
			}
		}

		resp = `{"payload": "` +
			b64.StdEncoding.EncodeToString([]byte(ciphertext)) +
			`"}`
	}

	return resp
}

func GetCryptoTest(r *http.Request,
	sUuid string,
	signature string) (string, string) {

	uuid, publicKHex, _ := cryptod.CreateCryptoOnboarding()

	return uuid.String(), publicKHex
}

func GetSavedTest(r *http.Request,
	sUuid string,
	signature string,
	request interfaces.ByLoginTest) interfaces.TestSaved {

	testSaved := interfaces.TestSaved{}

	search := cored.SearchLoginCore(
		request.Login)

	err := json.Unmarshal([]byte(search),
		&testSaved)
	if err != nil {
		log.Println("err: ", err)

		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.UNAUTHORIZED,
					constants.UNAUTHORIZED_CODE,
					constants.UNAUTHORIZED_MESSAGE)()))

	}
	return testSaved
}

func GetRequestHandleTest(
	r *http.Request,
	sUuid string,
	signature string) interfaces.ByLoginTest {

	sByLoginTest := interfaces.ByLoginTest{}
	req := GetOnboardingReqTest(r)
	desencr := cryptod.DesencrypGenericPayloadHandleTest(
		req,
		sUuid,
		signature)

	err := json.Unmarshal(
		[]byte(desencr),
		&sByLoginTest)

	if err != nil {
		log.Println(err)
	}
	return sByLoginTest
}
func ReceiveLogin(r *http.Request,
	sUuid string,
	signature string) string {

	ret := ""
	if sUuid == constants.EMPTY {
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					constants.UUID_FAULT_CODE,
					constants.UUID_FAULT,
					constants.BAD_REQUEST_CODE,
					constants.PAYLOAD_NA)()))
	} else {
		req := cryptod.
			DesencrypGenericPayload(
				GetOnboardingReq(r, sUuid),
				sUuid,
				signature)

		ret = Step01Login(r,
			sUuid,
			req)
	}

	return ret
}

func Step01Login(r *http.Request,
	sUuid string,
	sOnboardingReq interfaces.OnboardingReq) string {

	var ciphertext = []byte(constants.EMPTY)
	generateKeysJwtRes := searchFromOnboarding(
		sOnboardingReq,
		sUuid)

	if generateKeysJwtRes.ConsumerKey != constants.EMPTY {
		uuid, publicKHex, privateKHex := cryptod.CreateCryptoOnboarding()
		privateKSign, publicKSign := cryptod.CreateCryptoSign()
		log.Println(
			cored.SaveServerKeysSign(
				uuid.String(),
				publicKHex,
				privateKHex,
				privateKSign,
				publicKSign))

		ciphertext = manageLoginResponse(
			sOnboardingReq,
			uuid.String(),
			publicKHex,
			privateKHex,
			privateKSign,
			generateKeysJwtRes)
	}

	return `{"payload": "` +
		b64.StdEncoding.EncodeToString([]byte(ciphertext)) +
		`"}`
}

func searchFromOnboarding(
	sOnboardingReq interfaces.
		OnboardingReq,
	sUuid string) interfaces.GenerateKeysJwtRes {

	generateKeysJwtRes := interfaces.GenerateKeysJwtRes{}
	err := json.Unmarshal(
		[]byte(
			cored.SearchRegisterCore(
				cryptod.GetHashFromString(
					cored.GetJsonKeyStringToFindFunc(
						sOnboardingReq.DeviceInfo,
						sOnboardingReq.FaceInfo,
						sOnboardingReq.LocationInfo)()))),
		&generateKeysJwtRes)
	if err != nil ||
		generateKeysJwtRes.ConsumerKey == constants.EMPTY {

		log.Println("err: ", err)
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.FORBIDDEN,
					constants.FORBIDDEN_CODE,
					constants.ONBOARDING_FAULT)()))

	}
	return generateKeysJwtRes
}

func saveLoginCore(sOnboardingReq interfaces.OnboardingReq,
	uuid string,
	publicKHex string,
	privateKHex string,
	generateKeysJwtRes interfaces.GenerateKeysJwtRes) {

	log.Println(cored.SaveServerKeys(uuid, publicKHex, privateKHex))
	cored.SaveLoginCore(cored.GetJsonKeyLoginStringToSaveFunc(
		sOnboardingReq.DeviceInfo,
		sOnboardingReq.FaceInfo,
		generateKeysJwtRes.ConsumerKey,
		cryptod.GetHashFromString(
			cored.GetJsonKeyLoginStringToHashFunc(
				sOnboardingReq.DeviceInfo,
				sOnboardingReq.FaceInfo,
				generateKeysJwtRes.ConsumerKey,
				sOnboardingReq.LocationInfo)()),
		sOnboardingReq.LocationInfo)())
}

func encrypReturnLogin(
	sOnboardingReq interfaces.OnboardingReq,
	uuid string,
	publicKHex string,
	privateKHex string,
	privateKSign string,
	generateKeysJwtRes interfaces.GenerateKeysJwtRes) []byte {

	pkFront, err := ecies.NewPublicKeyFromHex(
		sOnboardingReq.NextFrontPublicKey)
	ciphertext, err := ecies.Encrypt(
		pkFront,
		[]byte(cored.GetJsonKeyLoginStringToReturnFunc(
			sOnboardingReq.DeviceInfo,
			sOnboardingReq.Email,
			sOnboardingReq.FaceInfo,
			sOnboardingReq.LocationInfo,
			generateKeysJwtRes.Jwt,
			sOnboardingReq.NextFrontPublicKey,
			uuid,
			cryptod.GetHashFromString(
				cored.GetJsonKeyLoginStringToHashFunc(
					sOnboardingReq.DeviceInfo,
					sOnboardingReq.FaceInfo,
					generateKeysJwtRes.ConsumerKey,
					sOnboardingReq.LocationInfo)()),
			privateKSign,
			privateKHex,
			publicKHex)()))

	if err != nil {
		log.Println(err)
	}
	return ciphertext
}

func manageLoginResponse(
	sOnboardingReq interfaces.OnboardingReq,
	uuid string,
	publicKHex string,
	privateKHex string,
	privateKSign string,
	generateKeysJwtRes interfaces.GenerateKeysJwtRes) []byte {

	saveLoginCore(
		sOnboardingReq,
		uuid,
		publicKHex,
		privateKHex,
		generateKeysJwtRes)

	return encrypReturnLogin(
		sOnboardingReq,
		uuid,
		publicKHex,
		privateKHex,
		privateKSign,
		generateKeysJwtRes)

}

func FindError(sUuid string) string {
	return errord.GetErrorPayload(sUuid)
}

func MapErrorMessage(errMsg string) int {

	resGeneralErrorKeyValuePayload := interfaces.ResGeneralErrorKeyValuePayload{}

	err := json.Unmarshal([]byte(errMsg), &resGeneralErrorKeyValuePayload)
	if err != nil {
		log.Println(err)
	}

	code, errConv := strconv.Atoi(resGeneralErrorKeyValuePayload.Value.Code)

	if errConv != nil {
		log.Println("Error during conversion")
	}

	Audit(constants.AUDIT + resGeneralErrorKeyValuePayload.Key)
	return code
}

func Audit(input string) {
	errord.GetErrorPayload(input)
}
