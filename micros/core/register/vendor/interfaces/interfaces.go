package interfaces

type PublicKeyRequest struct {
	NextServerPublicKey  string `json:"nextServerPublicKey"`
	Uuid                 string `json:"uuid"`
	NextServerPrivateKey string `json:"nextServerPrivateKey"`
}

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
	Email              string `json:"email"`
	LocationInfo       string `json:"locationInfo"`
}

type ReqRegisterPayload struct {
	CallbackUrl string `json:"callbackUrl"`
	ClientName  string `json:"clientName"`
	Owner       string `json:"owner"`
	GrantType   string `json:"grantType"`
	SaasApp     bool   `json:"saasApp"`
}

type AccessTokenPayloadRes struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Scope         string `json:"scope"`
	Id_token      string `json:"id_token"`
	Token_type    string `json:"token_type"`
	Expires_in    int64  `json:"expires_in"`
}

type RegisterPayloadRes struct {
	ClientId          string `json:"clientId"`
	ClientName        string `json:"clientName"`
	CallBackURL       string `json:"callBackURL"`
	ClientSecret      string `json:"clientSecret"`
	IsSaasApplication bool   `json:"isSaasApplication"`
	JsonString        string `json:"jsonString"`
	JsonAppAttribute  string `json:"jsonAppAttribute"`
}

type GenerateKeysAppRes struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
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
	Step                 string `json:"step"`
}

type ResGeneralPayload struct {
	Payload string `json:"payload"`
}

type GenerateRandonAppRes struct {
	ApplicationId    string `json:"applicationId"`
	Name             string `json:"name"`
	ThrottlingPolicy string `json:"throttlingPolicy"`
	Description      string `json:"description"`
	TokenType        string `json:"tokenType"`
	Status           string `json:"status"`
}

type AddSuscriptionAppRes struct {
	SubscriptionId   string `json:"subscriptionId"`
	ApplicationId    string `json:"applicationId"`
	ApiId            string `json:"apiId"`
	ThrottlingPolicy string `json:"throttlingPolicy"`
}

type AddSuscriptionAppRequest struct {
	ApplicationId             string `json:"applicationId"`
	ApiId                     string `json:"apiId"`
	ThrottlingPolicy          string `json:"throttlingPolicy"`
	RequestedThrottlingPolicy string `json:"requestedThrottlingPolicy"`
}

type GenerateRandonAppRequest struct {
	Name             string `json:"name"`
	ThrottlingPolicy string `json:"throttlingPolicy"`
	Description      string `json:"description"`
	TokenType        string `json:"tokenType"`
}

type GenerateKeysAppRequest struct {
	KeyType                 string   `json:"keyType"`
	KeyManager              string   `json:"keyManager"`
	GrantTypesToBeSupported []string `json:"grantTypesToBeSupported"`
	CallbackUrl             string   `json:"callbackUrl"`
	Scopes                  []string `json:"scopes"`
	ValidityTime            uint     `json:"validityTime"`
}
