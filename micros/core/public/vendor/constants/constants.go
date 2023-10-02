package constants

var HtmlStr = ``

var KeySymetricPass = []byte("XpJTzHH70fOTOazpaRvLbwW7HmvgnZIx")

const (
	UrlLocal                  = `http://localhost:33000/publicKey`
	HostLocal                 = `http://localhost:33000`
	HostWso2Local             = `https://localhost:8243`
	HostWso2Docker            = `https://172.18.0.1:8243`
	HostWso2PrivatePortDocker = `https://172.18.0.1:9443`
	DOCKER_GW_IP              = `172.22.0.1`
	DOCKER_GW_REGISTER_IP     = `172.21.0.1`
	CONTENT_TYPE              = "Content-Type"
	AUTH                      = "Authorization"
	BASIC_GENERIC_CREDENTIALS = `Basic dHNNWEFDQ1g1MFFmZ1d3aHp2MHlMdGhDSGtjYTpPc3NYR21TZlFOSnhQaDcxbnJrSE5yalhJRFFh`
	BASIC_GENERIC_ADMIN       = `Basic YWRtaW46YWRtaW4=`
	BASIC_OPENID              = `Basic dVFKX09GRzRlTVpMRWpsbHZLYXFOdlg4ZmNZYTo3M1R6Q1lKdExVQlh4aW9nV25jRDI5bFBtTVVh`
	APPLICATION_JSON          = "application/json"
	SLASH                     = "/"
	ONBOARDING_URI            = "/onboarding"
	TOKEN_URI_INTERNAL        = "/token"
	TOKEN_URI_WSO2            = "/token?grant_type=password&username=admin&password=admin&scope=apim%3Aapi_key%2520apim%3Aapp_import_export%2520apim%3Aapp_manage%2520apim%3Astore_settings%2520apim%3Asub_alert_manage%2520apim%3Asub_manage%2520apim%3Asubscribe%2520openid"
	CLIENT_REGISTER_URI_WSO2  = `/client-registration/v0.17/register`
	URI_GENERATE_RANDOM_APP   = `/api/am/store/v1/applications`
	URI_GENERATE_KEYS_APP     = `/api/am/store/v1/applications/applicationId/generate-keys`
	URI_ADD_SUSCRIPTION_APP   = `/api/am/store/v1/subscriptions`
	TOKEN_ADMIN_URI_WSO2      = `/token?grant_type=password&username=admin&password=admin&scope=apim%3Aapi_key+apim%3Aapp_import_export+apim%3Aapp_manage+apim%3Astore_settings+apim%3Asub_alert_manage+apim%3Asub_manage+apim%3Asubscribe+openid`
	JSON_STRING_PAYLOAD       = `{"payload": "payload"}`
	DISPATCHER_PORT           = ":33000"
	REGISTER_PORT             = ":33001"
	GET                       = "GET"
	CIPHERTEXT_TOO_SHORT      = "ciphertext too short"
	NOT_FOUND_TXT             = "404 not found."
	POST                      = "POST"
	SORRY_NOT_SOPPORTED       = "Sorry, only GET and POST methods are supported."
	FATAL_ERROR_TXT           = "Fatal Error"
	KEY_PAIR_GENERATED        = "key pair has been generated"
)
