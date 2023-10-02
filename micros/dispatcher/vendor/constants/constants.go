package constants

var HtmlStr = ``

var KeySymetricPass = []byte("XpJTzHH70fOTOazpaRvLbwW7HmvgnZIx")

const (
	Step02                    = `02`
	Step01                    = `01`
	UrlLocal                  = `http://localhost:33000/publicKey`
	HostLocal                 = `http://localhost:33000`
	HostWso2Local             = `https://localhost:8243`
	HostWso2Docker            = `https://172.18.0.1:8243`
	DOCKER_GW_IP              = `172.19.0.1`
	DOCKER_GW_REGISTER_IP     = `172.21.0.1`
	DOCKER_GW_LOGIN_IP        = `172.23.0.1`
	DOCKER_GW_PUBLIC_IP       = `172.22.0.1`
	DOCKER_GW_SEACHER_IP      = `172.20.0.1`
	DOCKER_GW_CLEANER_IP      = `172.25.0.1`
	DOCKER_GW_ERROR_IP        = `172.24.0.1`
	CONTENT_TYPE              = "Content-Type"
	AUTH                      = "Authorization"
	BASIC_GENERIC_CREDENTIALS = "Basic bGtKZTNycmR3aDk5UkpHeEsyUFZldnJBRExFYTpfaGZPOF9RMW5BaUlBVUNPYm5vRmxmZmJSQUlh"
	APPLICATION_JSON          = "application/json"
	OK                        = "OK"
	SLASH                     = "/"
	ONBOARDING_URI            = "/onboarding"
	LOGIN_URI                 = "/login"
	TEST_URI                  = "/test"
	TOKEN_URI_INTERNAL        = "/token"
	TOKEN_URI_WSO2            = "/token?grant_type=password&username=admin&password=admin&scope=apim%3Aapi_key%2520apim%3Aapp_import_export%2520apim%3Aapp_manage%2520apim%3Astore_settings%2520apim%3Asub_alert_manage%2520apim%3Asub_manage%2520apim%3Asubscribe%2520openid"
	JSON_STRING_PAYLOAD       = `{"payload": "payload"}`
	DISPATCHER_PORT           = ":33000"
	GET                       = "GET"
	CIPHERTEXT_TOO_SHORT      = "ciphertext too short"
	NOT_FOUND_TXT             = "404 not found."
	POST                      = "POST"
	SORRY_NOT_SOPPORTED       = "Sorry, only GET and POST methods are supported."
	FATAL_ERROR_TXT           = "Fatal Error"
	KEY_PAIR_GENERATED        = "key pair has been generated"
	EMPTY                     = ""
	BAD_REQUEST               = "Bad Request"
	BAD_REQUEST_CODE          = "400"
	PAYLOAD_FAULT             = "Payload fault"
	FORBIDDEN                 = "Forbidden"
	FORBIDDEN_CODE            = "403"
	SIGNATURE_FAULT           = "Signature fault"
	PRIVATEKEY_FAULT          = "Privatekey fault"
	AUDIT                     = "audit"
	UUID_FAULT                = "UUID fault"
	UUID_FAULT_CODE           = "5a0ca082-d297-4db9-bff5-5f6731144334"
	PAYLOAD_NA                = "Payload"
	ONBOARDING_FAULT          = "Onboarding fault"
	UNAUTHORIZED              = "Unauthorized"
	UNAUTHORIZED_CODE         = "401"
	UNAUTHORIZED_MESSAGE      = "Login fault"
)
