package interfaces

type PublicKeyRequest struct {
	NextServerPublicKey  string `json:"nextServerPublicKey"`
	Uuid                 string `json:"uuid"`
	NextServerPrivateKey string `json:"nextServerPrivateKey"`
}

type ResGeneralPayload struct {
	Payload string `json:"payload"`
}
