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
	BASIC_GENERIC_CREDENTIALS = "Basic bGtKZTNycmR3aDk5UkpHeEsyUFZldnJBRExFYTpfaGZPOF9RMW5BaUlBVUNPYm5vRmxmZmJSQUlh"
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
	API_KEY_CRYPTO            = `eyJ4NXQiOiJOVGRtWmpNNFpEazNOalkwWXpjNU1tWm1PRGd3TVRFM01XWXdOREU1TVdSbFpEZzROemM0WkE9PSIsImtpZCI6ImdhdGV3YXlfY2VydGlmaWNhdGVfYWxpYXMiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJhZG1pbkBjYXJib24uc3VwZXIiLCJhcHBsaWNhdGlvbiI6eyJvd25lciI6ImFkbWluIiwidGllclF1b3RhVHlwZSI6bnVsbCwidGllciI6IjEwUGVyTWluIiwibmFtZSI6InNpcml1cyIsImlkIjoyLCJ1dWlkIjoiYWU2ZDA1ZTgtZGQyMS00YmI3LWJmMDEtNGJiNTBlZWE5NjY4In0sImlzcyI6Imh0dHBzOlwvXC9sb2NhbGhvc3Q6OTQ0M1wvb2F1dGgyXC90b2tlbiIsInRpZXJJbmZvIjp7IkJyb256ZSI6eyJ0aWVyUXVvdGFUeXBlIjoicmVxdWVzdENvdW50IiwiZ3JhcGhRTE1heENvbXBsZXhpdHkiOjAsImdyYXBoUUxNYXhEZXB0aCI6MCwic3RvcE9uUXVvdGFSZWFjaCI6dHJ1ZSwic3Bpa2VBcnJlc3RMaW1pdCI6MCwic3Bpa2VBcnJlc3RVbml0IjpudWxsfX0sImtleXR5cGUiOiJQUk9EVUNUSU9OIiwic3Vic2NyaWJlZEFQSXMiOlt7InN1YnNjcmliZXJUZW5hbnREb21haW4iOiJjYXJib24uc3VwZXIiLCJuYW1lIjoiYXBpLXNpcml1cy1sb2dpbiIsImNvbnRleHQiOiJcL3Npcml1c1wvbG9naW5cL3YxIiwicHVibGlzaGVyIjoiYWRtaW4iLCJ2ZXJzaW9uIjoidjEiLCJzdWJzY3JpcHRpb25UaWVyIjoiQnJvbnplIn1dLCJpYXQiOjE2OTU5MDI5MTMsImp0aSI6IjEyZmVmYzQwLWY2NGItNGQ5Ny05ZTQxLTIzNzcwOGZiMzljOCJ9.NCGw4gbWhynsmQLFwIOW8xME9Qqchesuoc3oYd8JM_4kz7Fi60F2qEMT6q3QoJzLUL3Ee3UQO2IcdAVRaEzEn7xtNOlhTVn-6MyfRv69XnDNnmy9pltZphcit6_CqHF9ssOKNaBTj3nwfmSv0vRN8EV-OPAC2ZFOA81Bn9JIqUhq1Yl8Tlnd8a1UEijHS6qLx_CQfXwXE4BRjLepGeFVXHziWOq-oGu73H6xDquuLGL9kyV25nQxnrC6dbR5LWSW9J1n1UTbed1u4H5LqxqU5umbaLikL2CiIt20URC5krGq1z1_5On71q4hOtL8pSoVlcMtXURFF6tgDTxDG7rt-A==`
)
