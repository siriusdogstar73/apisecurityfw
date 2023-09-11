package cryptod

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"io"

	"constants"

	"log"

	"cored"

	"interfaces"

	"encoding/json"

	"fmt"

	"strings"

	ecies "github.com/ecies/go"
	uuid "github.com/google/uuid"

	"math/big"

	"crypto/ecdsa"

	base58 "github.com/btcsuite/btcutil/base58"

	"crypto/elliptic"

	b64 "encoding/base64"

	"crypto/sha256"
)

var InitPublicServerKeyHex string
var InitPrivateServerKey *ecies.PrivateKey

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
		log.Println(err)
	}
	log.Println(constants.KEY_PAIR_GENERATED)

	/* Get Hex Public Key. The client can get it by simple REST call
	   because when the server restart it changes */
	InitPublicServerKey := PairKeys.PublicKey
	InitPrivateServerKey = PairKeys
	InitPublicServerKeyHex = InitPublicServerKey.Hex(true)

}

func CreateCryptoOnboarding() (uuid.UUID, string, string) {
	uuidWithHyphen := uuid.New()

	InitCrypto()
	InitPrivateServerKeyHex := InitPrivateServerKey.Hex()
	return uuidWithHyphen, InitPublicServerKeyHex, InitPrivateServerKeyHex
}

func DecryptEcies(textBytes []byte,
	sUuid string) string {

	jsonPrivateServerKeyHex := cored.SearchServerKeysPublicCore(sUuid)
	sPublicKeyRequest := interfaces.PublicKeyRequest{}

	err := json.Unmarshal([]byte(jsonPrivateServerKeyHex), &sPublicKeyRequest)
	if err != nil {
		log.Println(err)
	}
	searcherPrivateServerKey, _ := ecies.NewPrivateKeyFromHex(sPublicKeyRequest.NextServerPrivateKey)

	cored.CleanerPublic(sPublicKeyRequest.Uuid)

	plaintext, err := ecies.Decrypt(searcherPrivateServerKey, textBytes)
	if err != nil {
		fmt.Println(err)
	}
	return string(plaintext)
}

func DecryptEciesSignature(textBytes []byte,
	sUuid string,
	signature string) string {

	jsonPrivateServerKeyHex := cored.SearchServerKeysPublicCore(sUuid)
	sPublicKeyRequest := interfaces.PublicKeyRequest{}

	err := json.Unmarshal([]byte(jsonPrivateServerKeyHex), &sPublicKeyRequest)
	if err != nil {
		log.Println(err)
	}
	searcherPrivateServerKey, _ := ecies.NewPrivateKeyFromHex(sPublicKeyRequest.NextServerPrivateKey)

	cored.CleanerPublic(sPublicKeyRequest.Uuid)

	plaintext, err := ecies.Decrypt(searcherPrivateServerKey, textBytes)
	if err != nil {
		fmt.Println(err)
	}
	if Verify(plaintext, signature) {
		log.Println("Verified...")
	} else {
		log.Println("NOT Verified...")
	}
	return string(plaintext)
}

func DecryptEciesSignatureTest(
	textBytes []byte,
	sUuid string,
	signature string) string {

	jsonPrivateServerKeyHex := cored.SearchServerKeysPublicCore(sUuid)
	sPublicKeyRequest := interfaces.PublicKeyRequest{}

	err := json.Unmarshal([]byte(
		jsonPrivateServerKeyHex),
		&sPublicKeyRequest)
	if err != nil {
		log.Println(err)
	}
	searcherPrivateServerKey, _ := ecies.NewPrivateKeyFromHex(
		sPublicKeyRequest.NextServerPrivateKey)

	cored.CleanerPublic(sPublicKeyRequest.Uuid)

	plaintext, err := ecies.Decrypt(searcherPrivateServerKey, textBytes)
	if err != nil {
		fmt.Println(err)
	}
	if Verify(plaintext, signature) {
		log.Println("Verified...")
	} else {
		log.Println("NOT Verified...")
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.FORBIDDEN,
					constants.FORBIDDEN_CODE,
					string(textBytes))()))
	}

	if signature == constants.EMPTY {
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					constants.SIGNATURE_FAULT)()))
	}

	return string(plaintext)
}

func Verify(data []byte, signature string) bool {
	parts := strings.Split(signature, ".")
	if len(parts) < 3 {
		return false
	}
	r := new(big.Int).SetBytes(base58.Decode(parts[0]))
	s := new(big.Int).SetBytes(base58.Decode(parts[1]))

	publicKey := parts[2]
	byteArray := base58.Decode(publicKey)
	var pubKey ecdsa.PublicKey
	pubKey.Curve = elliptic.P256()
	pubKey.X, pubKey.Y = elliptic.UnmarshalCompressed(elliptic.P256(), byteArray)
	return ecdsa.Verify(&pubKey, data, r, s)
}

func DecryptEciesSign(
	textBytes []byte,
	sUuid string,
	signature string,
	request string) (string, string) {

	jsonPrivateServerKeyHex := cored.SearchServerKeysPublicCore(sUuid)
	sPublicKeyRequest := interfaces.PublicKeyRequest{}

	err := json.Unmarshal([]byte(jsonPrivateServerKeyHex), &sPublicKeyRequest)
	if err != nil {
		log.Println(err)
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					request)()))
	}
	searcherPrivateServerKey, _ := ecies.NewPrivateKeyFromHex(
		sPublicKeyRequest.NextServerPrivateKey)

	cored.CleanerPublic(sPublicKeyRequest.Uuid)

	plaintext, err := ecies.Decrypt(searcherPrivateServerKey, textBytes)

	if err != nil {
		log.Println(err)
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.FORBIDDEN,
					constants.FORBIDDEN_CODE,
					request)()))
	}

	if Verify(plaintext, signature) {
		log.Println("Verified...")
	} else {
		log.Println("NOT Verified...")
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.FORBIDDEN,
					constants.FORBIDDEN_CODE,
					request)()))
	}
	return string(plaintext),
		sPublicKeyRequest.PrivateSignKey
}

func DesencrypGenericPayload(
	resGeneralPayload interfaces.ResGeneralPayload,
	sUuid string,
	signature string) interfaces.OnboardingReq {

	var plaintext, privateSignKey string
	sDec, _ := b64.StdEncoding.DecodeString(resGeneralPayload.Payload)

	//plaintext, err := Decrypt([]byte(sDec), constants.KeySymetricPass)

	textBytes := []byte(sDec)
	if signature != constants.EMPTY {
		plaintext, privateSignKey =
			DecryptEciesSign(
				textBytes,
				sUuid,
				signature,
				resGeneralPayload.Payload)

	} else {
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					constants.SIGNATURE_FAULT)()))
	}

	textBytes = []byte(plaintext)

	sOnboardingReq := interfaces.OnboardingReq{}
	err := json.Unmarshal(textBytes, &sOnboardingReq)
	if err != nil {
		log.Println(err)
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					resGeneralPayload.Payload)()))
	}

	if privateSignKey != sOnboardingReq.PrivateSignKey {
		log.Println("The privateSignKey is not like de etcd value...")
		log.Println(
			cored.ConnectErrorCore(
				cored.GetResGeneralErrorKeyValuePayloadFunc(
					sUuid,
					constants.BAD_REQUEST,
					constants.BAD_REQUEST_CODE,
					constants.PRIVATEKEY_FAULT)()))
	} else {
		log.Println("The privateSignKey is like de etcd value...")
	}

	return sOnboardingReq
}

func GetHashFromString(jsonKeyStringToSave string) string {
	data := []byte(jsonKeyStringToSave)
	hash := sha256.Sum256(data)
	sHash := fmt.Sprintf("%x\n", hash[:]) // 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	sHash = b64.StdEncoding.EncodeToString([]byte(sHash))
	return sHash
}

func DesencrypGenericPayloadHandleTest(resGeneralPayload interfaces.ResGeneralPayload,
	sUuid string,
	signature string) string {

	sDec, _ := b64.StdEncoding.DecodeString(resGeneralPayload.Payload)

	return DecryptEciesSignatureTest([]byte(sDec), sUuid, signature)
}

func CreateCryptoSign() (string, string) {

	newPrivateKey := NewPrivateKey()

	newPublicKey := GetPublicKey(newPrivateKey)
	return newPrivateKey, newPublicKey
}

func NewPrivateKey() string {
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return base58.Encode(pk.D.Bytes())
}

func GetPublicKey(privateKey string) string {
	var pri ecdsa.PrivateKey
	pri.D = new(big.Int).SetBytes(base58.Decode(privateKey))
	pri.PublicKey.Curve = elliptic.P256()
	pri.PublicKey.X, pri.PublicKey.Y = pri.PublicKey.Curve.ScalarBaseMult(pri.D.Bytes())
	return base58.Encode(elliptic.MarshalCompressed(elliptic.P256(), pri.PublicKey.X, pri.PublicKey.Y))
}

func Sign(data []byte, privateKey string) (string, error) {
	var pri ecdsa.PrivateKey
	pri.D = new(big.Int).SetBytes(base58.Decode(privateKey))
	pri.PublicKey.Curve = elliptic.P256()
	pri.PublicKey.X, pri.PublicKey.Y = pri.PublicKey.Curve.ScalarBaseMult(pri.D.Bytes())

	r, s, err := ecdsa.Sign(rand.Reader, &pri, data)
	if err != nil {
		return "", err
	}
	pubKeyStr := base58.Encode(elliptic.MarshalCompressed(elliptic.P256(), pri.PublicKey.X, pri.PublicKey.Y))
	sig2 := fmt.Sprintf("%s.%s.%s", base58.Encode(r.Bytes()), base58.Encode(s.Bytes()), pubKeyStr)
	return sig2, nil
}
