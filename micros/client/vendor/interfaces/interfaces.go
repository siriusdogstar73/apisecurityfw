package interfaces

type ResGetPublicKeyWithUuid struct {
	PublicKey string `json:"publicKey"`
	Jwt       string `json:"jwt"`
	Uuid      string `json:"uuid"`
}

type ResGetPublicKeyWithUuidAndSign struct {
	PublicKey      string `json:"publicKey"`
	Jwt            string `json:"jwt"`
	Uuid           string `json:"uuid"`
	PrivateSignKey string `json:"privateSignKey"`
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

type ResGeneralBody struct {
	Body string `json:"body"`
}

type OnboardingReq struct {
	DeviceInfo         string `json:"deviceInfo"`
	FaceInfo           string `json:"faceInfo"`
	NextFrontPublicKey string `json:"nextFrontPublicKey"`
	Email              string `json:"email"`
	LocationInfo       string `json:"locationInfo"`
	Step               string `json:"step"`
	PrivateSignKey     string `json:"privateSignKey"`
}

type OnboardingFront struct {
	DeviceInfo          string `json:"deviceInfo"`
	FaceInfo            string `json:"faceInfo"`
	NextServerPublicKey string `json:"nextServerPublicKey"`
	Email               string `json:"email"`
	LocationInfo        string `json:"locationInfo"`
	Uuid                string `json:"uuid"`
	Jwt                 string `json:"jwt"`
	Login               string `json:"login"`
	Step                string `json:"step"`
}

type OnboardingSignFront struct {
	DeviceInfo          string `json:"deviceInfo"`
	FaceInfo            string `json:"faceInfo"`
	NextServerPublicKey string `json:"nextServerPublicKey"`
	Email               string `json:"email"`
	LocationInfo        string `json:"locationInfo"`
	Uuid                string `json:"uuid"`
	Jwt                 string `json:"jwt"`
	Login               string `json:"login"`
	Step                string `json:"step"`
	PrivateSignKey      string `json:"privateSignKey"`
}

type OnboardingByLoginFront struct {
	Uuid               string `json:"uuid"`
	Jwt                string `json:"jwt"`
	Login              string `json:"login"`
	NextFrontPublicKey string `json:"nextFrontPublicKey"`
}
