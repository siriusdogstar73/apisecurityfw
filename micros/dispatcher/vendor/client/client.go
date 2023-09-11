package client

import (
	"bytes"

	"encoding/json"

	"io/ioutil"
	"log"
	"net/http"

	"constants"
	"interfaces"
)

func GetJWT(clienteHTTP *http.Client) string {
	bresponse := PostGenericBodySimple(clienteHTTP,
		constants.HostWso2Docker,
		constants.TOKEN_URI_WSO2,
		constants.BASIC_GENERIC_CREDENTIALS,
		constants.POST)
	responseString := string(bresponse)

	textBytes := []byte(responseString)

	jwtResponse := interfaces.JwtResponse{}

	err := json.Unmarshal(textBytes, &jwtResponse)
	if err != nil {
		log.Println(err.Error())
	}
	return jwtResponse.Access_token

}

func PostGenericBodySimple(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	body := new(interface{})
	req, err := json.Marshal(body)
	if err != nil {
		log.Println("Error codificando usuario como JSON: %v", err)
	}
	request, err := http.NewRequest("POST", host+uri, bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}
	/*
		NOT application/json!
		request.Header.Add("Content-Type", "application/json")
	*/
	request.Header.Add(constants.AUTH, authHeader)

	response, err := clienteHTTP.Do(request)
	if err != nil {
		log.Println("Error recibiendo respuesta postApisName: %v", err)
	}
	bresponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response: %v", err)
	}

	return bresponse

}
