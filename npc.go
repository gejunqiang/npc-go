package npc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const EmptyString = ""

type Npc struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
}

func NewNpc(endpoint string, accessKey string, secretKey string, region string) *Npc {
	if accessKey == EmptyString || secretKey == EmptyString {
		log.Fatal("api_key & api_secret can't be empty!")
	}
	return &Npc{Endpoint: endpoint, AccessKey: accessKey, SecretKey: secretKey, Region: region}
}

func (npc *Npc) getUrlQueryString(params map[string]string, orderedKeys []string) string {
	if params == nil{
		return EmptyString
	}
	v := url.Values{}
	switch orderedKeys {
	case nil:
		for key, value := range params {
			v.Add(key, value)
		}
	default:
		for _, key := range orderedKeys {
			if value, ok := params[key]; ok {
				v.Add(key,value)
			}
		}
	}
	return v.Encode()
}

func (npc *Npc) getTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func (npc *Npc) getUUID(length int) string{
	if length <= 0 {
		length = 10
	}
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (npc *Npc) getHashPayload(requestBody string) string {
	hash := sha256.New()
	hash.Write([]byte(requestBody))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

func (npc *Npc) getCanonicalizedQueryString(params map[string]string) string {
	fixedParams := make(map[string]string, 0)
	for key, value := range params {
		fixedParams[key] = value
	}
	fixedParams["AccessKey"] = npc.AccessKey
	fixedParams["Timestamp"] = npc.getTimestamp()
	fixedParams["SignatureVersion"] = "1.0"
	fixedParams["SignatureMethod"] = "HMAC-SHA256"
	fixedParams["SignatureNonce"] = npc.getUUID(10)
	fixedParams["Region"] = npc.Region
	orderedKeys := make([]string, 0)
	for key, _ := range fixedParams {
		orderedKeys = append(orderedKeys, key)
	}
	sort.Strings(orderedKeys)
	return npc.getUrlQueryString(fixedParams, orderedKeys)
}

func (npc *Npc) getString2Sign(method, endpoint, service, canonicalizedQueryString, hashPayload string) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, endpoint, service, canonicalizedQueryString, hashPayload)
}

func (npc *Npc) getSignature(string2Sign, secretKey string) string {
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(string2Sign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// example:
// service = /keypair
// params = {Action: ListKeyPair, Version: 2018-02-08}
func (npc *Npc) Get(service string, params map[string]string)(*http.Response, error) {
	method := http.MethodGet
	canonicalizedQueryString := npc.getCanonicalizedQueryString(params)
	string2sign := npc.getString2Sign(method, npc.Endpoint, service, canonicalizedQueryString, npc.getHashPayload(EmptyString))
	signature := npc.getSignature(string2sign, npc.SecretKey)
	signatureParams := map[string]string {
		"Signature": signature,
	}
	httpUrl := "https://" + npc.Endpoint + service + "?" + canonicalizedQueryString + "&" + npc.getUrlQueryString(signatureParams, nil)
	fmt.Println(httpUrl)
	client := &http.Client{}
	req, err := http.NewRequest(method, httpUrl, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// example:
// service = /keypair
// Action = {Action: UploadKeyPair, Version: 2018-02-08}
// body = jsonFormat string
func (npc *Npc) Post(service string, params map[string]string, body string)(*http.Response, error) {
	method := http.MethodPost
	canonicalizedQueryString := npc.getCanonicalizedQueryString(params)
	string2sign := npc.getString2Sign(method, npc.Endpoint, service, canonicalizedQueryString, npc.getHashPayload(body))
	signature := npc.getSignature(string2sign, npc.SecretKey)
	signatureParams := map[string]string {
		"Signature": signature,
	}
	httpUrl := "https://" + npc.Endpoint + service + "?" + canonicalizedQueryString + "&" + npc.getUrlQueryString(signatureParams, nil)
	fmt.Println(httpUrl)
	client := &http.Client{}
	req, err := http.NewRequest(method, httpUrl, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return client.Do(req)
}
