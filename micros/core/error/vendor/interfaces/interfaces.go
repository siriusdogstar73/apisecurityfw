package interfaces

type ResGeneralErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ResGeneralErrorKeyValuePayload struct {
	Key     string                 `json:"key"`
	Value   ResGeneralErrorPayload `json:"value"`
	Request string                 `json:"request"`
}
