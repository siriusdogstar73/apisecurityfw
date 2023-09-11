package interfaces

type LoginKeyRequest struct {
	ClientId     string `json:"clientId"`
	DeviceInfo   string `json:"deviceInfo"`
	FaceInfo     string `json:"faceInfo"`
	LocationInfo string `json:"locationInfo"`
	Login        string `json:"login"`
}
