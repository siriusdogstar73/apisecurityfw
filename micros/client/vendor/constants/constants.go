package constants

var HtmlStr = ``

var KeySymetricPass = []byte("XpJTzHH70fOTOazpaRvLbwW7HmvgnZIx")

const (
	Step01        = `01`
	Step02        = `02`
	UrlLocal      = `http://localhost:33000/publicKey`
	HostLocal     = `http://localhost:33000`
	HostWso2Local = `https://localhost:8243`
	//HostWso2Local             = `https://sirius-wso2.herokuapp.com`
	HostWso2Docker            = `https://172.18.0.1:8243`
	CONTENT_TYPE              = "Content-Type"
	AUTH                      = "Authorization"
	API_KEY_HEADER            = "apikey"
	BASIC_GENERIC_CREDENTIALS = "Basic dHNNWEFDQ1g1MFFmZ1d3aHp2MHlMdGhDSGtjYTpPc3NYR21TZlFOSnhQaDcxbnJrSE5yalhJRFFh"
	APPLICATION_JSON          = "application/json"
	SLASH                     = "/"
	ONBOARDING_URI            = "/onboarding"
	ONBOARDING_UR_BY_API      = "/sirius/dispatcher/v1/onboarding"
	ONBOARDING_UR_BY_API_QS   = "/sirius/dispatcher/v1/onboarding?uuid="
	LOGIN_UR_BY_API_QS        = "/sirius/dispatcher/v1/login?uuid="
	TEST_UR_BY_API_QS         = "/sirius/dispatcher/v1/test?uuid="
	TOKEN_URI_INTERNAL        = "/token"
	TOKEN_URI_CRYPTO          = "/sirius/login/v1/token"
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
	API_KEY_CRYPTO            = `eyJ4NXQiOiJOVGRtWmpNNFpEazNOalkwWXpjNU1tWm1PRGd3TVRFM01XWXdOREU1TVdSbFpEZzROemM0WkE9PSIsImtpZCI6ImdhdGV3YXlfY2VydGlmaWNhdGVfYWxpYXMiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZSI6bnVsbCwidGllciI6IjEwUGVyTWluIiwibmFtZSI6ImNyeXB0byIsImlkIjozNiwidXVpZCI6ImU4OTA0Mjg3LTgyNTAtNDY0NS1hMGRkLTE5NmRiYTU0MzgzNCJ9LCJpc3MiOiJodHRwczpcL1wvbG9jYWxob3N0Ojk0NDNcL29hdXRoMlwvdG9rZW4iLCJ0aWVySW5mbyI6eyJCcm9uemUiOnsidGllclF1b3RhVHlwZSI6InJlcXVlc3RDb3VudCIsImdyYXBoUUxNYXhDb21wbGV4aXR5IjowLCJncmFwaFFMTWF4RGVwdGgiOjAsInN0b3BPblF1b3RhUmVhY2giOnRydWUsInNwaWtlQXJyZXN0TGltaXQiOjAsInNwaWtlQXJyZXN0VW5pdCI6bnVsbH19LCJrZXl0eXBlIjoiUFJPRFVDVElPTiIsInBlcm1pdHRlZFJlZmVyZXIiOiIiLCJzdWJzY3JpYmVkQVBJcyI6W3sic3Vic2NyaWJlclRlbmFudERvbWFpbiI6ImNhcmJvbi5zdXBlciIsIm5hbWUiOiJhcGktc2lyaXVzLWxvZ2luIiwiY29udGV4dCI6Ilwvc2lyaXVzXC9sb2dpblwvdjEiLCJwdWJsaXNoZXIiOiJhZG1pbiIsInZlcnNpb24iOiJ2MSIsInN1YnNjcmlwdGlvblRpZXIiOiJCcm9uemUifV0sInBlcm1pdHRlZElQIjoiIiwiaWF0IjoxNjUzMjM2NTg0LCJqdGkiOiI4YjJhMzE3Yy04YmJmLTQwZDMtODlmYy01Y2IxZDQxYjE5YTMifQ==.jiol-dJEIG9pFTfcVbmg8lbWj5PhUJvFoFicpb2I7hIxgB2zdbA_lych2COsWYSvBxd1m1F_GmMKSUzVMM2CiHXgrkz-nQP6ICq1dALaOZZSIm5T-ViOvB3ZpIjLyUGJMzQyxZHWx1v2DqqLCXD3jDeE7fg7js6qA5KlToxuDDWRYV5S8VRtq3ezrmF3Gh5zeJWKUUbymEFooBOvwDosIA8B7B57qD9QbTK4dJbWH26DMeg1kmQtAvTtp2nkWWLIm1A9H4cearC3GCHsnq8sxsAUVcdOSiAs9gk1K2MPNVGxWIVu6S-V5vYuddQAyfVaJCDYW7VzaGIev92dAaIUsw==`
)
