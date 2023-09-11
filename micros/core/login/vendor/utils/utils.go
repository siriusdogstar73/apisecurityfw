package utils

import (
	"constants"
	"fmt"

	"log"
	"net/http"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"

	//b64 "encoding/base64"

	//"bytes"

	//"strings"

	ecies "github.com/ecies/go"
	uuid "github.com/google/uuid"

	//http2curl "moul.io/http2curl"
	"context"

	"time"

	"interfaces"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var InitPublicServerKeyHex string
var InitPrivateServerKey *ecies.PrivateKey

func Health(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != constants.SLASH {
		http.Error(w, constants.NOT_FOUND_TXT, http.StatusNotFound)
		return
	}

	switch r.Method {
	case constants.GET:
		fmt.Fprintf(w, constants.HtmlStr)
	case constants.POST:
		if err := r.ParseForm(); err != nil {
			log.Fatal(constants.FATAL_ERROR_TXT)
			return
		}

	default:
		fmt.Fprintf(w, constants.SORRY_NOT_SOPPORTED)
	}
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func EncryptPublicKey(text []byte) []byte {

	ciphertext, err := Encrypt(text, constants.KeySymetricPass)
	if err != nil {
		log.Fatal(err)
	}

	return ciphertext

}

func InitCrypto() {
	/* init crypto */
	/* Generate key pair */
	PairKeys, err := ecies.GenerateKey()
	if err != nil {
		panic(err)
	}
	log.Println(constants.KEY_PAIR_GENERATED)

	/* Get Hex Public Key. The client can get it by simple REST call
	   because when the server restart it changes */
	InitPublicServerKey := PairKeys.PublicKey
	InitPrivateServerKey = PairKeys
	InitPublicServerKeyHex = InitPublicServerKey.Hex(true)

}

func DecryptEcies(textBytes []byte) string {

	plaintext, err := ecies.Decrypt(InitPrivateServerKey, textBytes)
	if err != nil {
		fmt.Println(err)
	}
	return string(plaintext)
}

func CreateCryptoOnboarding() (uuid.UUID, string, string) {
	uuidWithHyphen := uuid.New()

	InitCrypto()
	InitPrivateServerKeyHex := InitPrivateServerKey.Hex()
	return uuidWithHyphen, InitPublicServerKeyHex, InitPrivateServerKeyHex
}

func PostAccessToken(clienteHTTP *http.Client, host string,
	uri string, authHeader string,
	method string) []byte {

	request, err := http.NewRequest("POST", host+constants.TOKEN_ADMIN_URI_WSO2, nil)
	if err != nil {
		log.Println("Error creando petición postApisName: %v", err)
	}

	/*
		NOT application/json!
		request.Header.Add("Content-Type", "application/json")
	*/

	request.Header.Add(constants.AUTH, "Basic "+authHeader)

	//TO DEBUG CURL
	/*
		command, _ := http2curl.GetCurlCommand(request)
		fmt.Println(command)
	*/

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

func initEtcdClient(key string, value string) {
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := cli.Put(ctx, key, value)

	if err != nil {
		switch err {
		case context.Canceled:
			log.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Println("ctx is attached with a deadline is exceeded: %v", err)

		default:
			log.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	/*
		response, err := cli.Get(ctx, key)

		for _, ev := range response.Kvs {
			log.Printf("%s : %s\n", ev.Key, ev.Value)
		}
	*/

}

func initEtcdClientWirhControl(key string, value string) string {
	result := "OK"
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := cli.Put(ctx, key, value)

	if err != nil {
		switch err {
		case context.Canceled:
			log.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Println("ctx is attached with a deadline is exceeded: %v", err)

		default:
			log.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	/*
		response, err := cli.Get(ctx, key)

		for _, ev := range response.Kvs {
			log.Printf("%s : %s\n", ev.Key, ev.Value)
		}
	*/
	if err == nil {
		result = "OK"
	} else {
		result = "KO"
	}
	return result

}

func searchEtcdClientLogin(input string) string {
	var k string
	var v string
	key := "login" + input
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	resp, _ := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))

	for _, ev := range resp.Kvs {
		k = fmt.Sprintf("%s", ev.Key)
		log.Println("encontrado: ", k)
		v = fmt.Sprintf("%s", ev.Value)
	}
	return v
}

func SaveLogin(input string) string {

	sLoginKeyRequest := interfaces.LoginKeyRequest{}
	err := json.Unmarshal([]byte(input), &sLoginKeyRequest)
	if err != nil {
		log.Println(err)
	}
	return initEtcdClientWirhControl("login"+sLoginKeyRequest.Login, input)
}

func SearchLoginCore(input string) string {

	return searchEtcdClientLogin(input)
}
