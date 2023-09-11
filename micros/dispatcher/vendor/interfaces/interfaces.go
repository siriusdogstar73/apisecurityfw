package interfaces

type ResGetPublicKey struct {
	PublicKey string `json:"publicKey"`
	Jwt       string `json:"jwt"`
}

type JwtResponse struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Scope         string `json:"scope"`
	Token_type    string `json:"token_type"`
	Expires_in    int    `json:"expires_in"`
}

type ResGetPublicKeyPayload struct {
	Payload string `json:"payload"`
}

type OnboardingReq struct {
	DeviceInfo         string `json:"deviceInfo"`
	FaceInfo           string `json:"faceInfo"`
	NextFrontPublicKey string `json:"nextFrontPublicKey"`
	Key                string `json:"key"`
	Email              string `json:"email"`
	LocationInfo       string `json:"locationInfo"`
	Step               string `json:"step"`
	PrivateSignKey     string `json:"privateSignKey"`
}

type OnboardingKeyReq struct {
	DeviceInfo         string `json:"deviceInfo"`
	FaceInfo           string `json:"faceInfo"`
	LocationInfo       string `json:"locationInfo"`
	Uuid               string `json:"uuid"`
	NextFrontPublicKey string `json:"nextFrontPublicKey"`
}

type ResGeneralPayload struct {
	Payload string `json:"payload"`
}

type OnboardingReqPrivate struct {
	LoginInfo      OnboardingRequest `json:"loginInfo"`
	ApplicationId  string            `json:"applicationId"`
	ConsumerKey    string            `json:"consumerKey"`
	ConsumerSecret string            `json:"consumerSecret"`
	ApiId          string            `json:"apiId"`
	SubscriptionId string            `json:"subscriptionId"`
	Jwt            string            `json:"jwt"`
	Process        string            `json:"process"`
	Step           uint              `json:"step"`
}

type OnboardingRequest struct {
	DeviceInfo           string `json:"deviceInfo"`
	FaceInfo             string `json:"faceInfo"`
	NextFrontPublicKey   string `json:"nextFrontPublicKey"`
	Email                string `json:"email"`
	LocationInfo         string `json:"locationInfo"`
	Uuid                 string `json:"uuid"`
	Key                  string `json:"key"`
	NextServerPrivateKey string `json:"nextServerPrivateKey"`
}

type TestSaved struct {
	DeviceInfo   string `json:"deviceInfo"`
	FaceInfo     string `json:"faceInfo"`
	ClientId     string `json:"clientId"`
	Login        string `json:"login"`
	LocationInfo string `json:"locationInfo"`
	//PrivateSignKey string `json:"privateSignKey"`
}

type PublicKeyRequest struct {
	NextServerPublicKey  string `json:"nextServerPublicKey"`
	Uuid                 string `json:"uuid"`
	NextServerPrivateKey string `json:"nextServerPrivateKey"`
	PrivateSignKey       string `json:"privateSignKey"`
	PublicSignKey        string `json:"publicSignKey"`
}

type GenerateKeysJwtRes struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	Jwt            string `json:"jwt"`
}

type ByLoginTest struct {
	Uuid               string `json:"uuid"`
	Jwt                string `json:"jwt"`
	Login              string `json:"login"`
	NextFrontPublicKey string `json:"nextFrontPublicKey"`
}

type ResGeneralErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ResGeneralErrorKeyValuePayload struct {
	Key   string                 `json:"key"`
	Value ResGeneralErrorPayload `json:"value"`
}
