package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gejunqiang/npc-go"
)

// Environments Variables:
// - NPC_API_ENDPOINT
// - NPC_API_KEY
// - NPC_API_SECRET
// - NPC_API_REGION
// - NPC_API_CONFIG

const (
	EnvNpcApiEndpoint    = "NPC_API_ENDPOINT"
	EnvNpcApiKey         = "NPC_API_KEY"
	EnvNpcApiSecret      = "NPC_API_SECRET"
	EnvNpcApiRegion      = "NPC_API_REGION"
	EnvNpcApiConfig      = "NPC_API_CONFIG"

	NpcApiEndpoint       = "api_endpoint"
	NpcApiKey            = "api_key"
	NpcApiSecret         = "api_secret"
	NpcApiRegion         = "api_region"

	DefaultApiEndpoint   = "open.c.163.com"
	DefaultApiRegion     = "cn-east-1"

	EmptyString          = ""

)

var(
	endpoint, accessKey, secretKey, region string
	apiConfigFile string
	DefaultApiConfigFile = os.Getenv("HOME") + "/.npc/api.key"
)

var(
	method, uri, requestBody string
)

func variableFromEnv(variable string) string {
	return os.Getenv(variable)
}

func variableFromFile(filePath string) map[string]string {
	if filePath == EmptyString {
		return nil
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}
	variableMap := make(map[string]string, 0)
	err = json.Unmarshal(b, &variableMap)
	if err != nil {
		return nil
	}
	return variableMap
}

func loadVariable() {
	endpoint = variableFromEnv(EnvNpcApiEndpoint)
	accessKey = variableFromEnv(EnvNpcApiKey)
	secretKey = variableFromEnv(EnvNpcApiSecret)
	region = variableFromEnv(EnvNpcApiRegion)
	apiConfigFile = variableFromEnv(EnvNpcApiConfig)
	if apiConfigFile == EmptyString {
		apiConfigFile = DefaultApiConfigFile
	}
	variableMap := variableFromFile(apiConfigFile)
	if endpoint == EmptyString && variableMap != nil {
		endpoint = variableMap[NpcApiEndpoint]
	}
	if endpoint == EmptyString {
		endpoint = DefaultApiEndpoint
	}
	if accessKey == EmptyString && variableMap != nil {
		accessKey = variableMap[NpcApiKey]
	}
	if secretKey == EmptyString && variableMap != nil {
		secretKey = variableMap[NpcApiSecret]
	}
	if region == EmptyString && variableMap != nil {
		region = variableMap[NpcApiRegion]
	}
	if region == EmptyString {
		region = DefaultApiRegion
	}
}

func getArgs() {
	for index, arg := range os.Args {
		if index == 1 {
			method = arg
		}
		if index == 2 {
			uri = arg
		}
		if index == 3 {
			requestBody = arg
		}
	}
}

func getServiceAndParamsFromUri(uri string) (service string, params map[string]string){
	if uri == EmptyString {
		return
	}
	tmp := strings.Split(uri, "?")
	service = tmp[0]
	if len(tmp) == 2 {
		paramStr := tmp[1]
		params = make(map[string]string, 0)
		for _, param := range strings.Split(paramStr, "&") {
			singleParam := strings.Split(param, "=")
			if len(singleParam) == 2 {
				params[singleParam[0]] = singleParam[1]
			}
		}
	}
	return
}

func main() {
	loadVariable()
	getArgs()
	npcHttp := npc.NewNpc(endpoint, accessKey, secretKey, region)
	service, params := getServiceAndParamsFromUri(uri)
	var resp *http.Response
	var err error
	switch method {
	case http.MethodGet:
		resp, err = npcHttp.Get(service, params)
	case http.MethodPost:
		resp, err = npcHttp.Post(service, params, requestBody)
	default:
		log.Fatalf("HTTP method: %s is not supported now", method)
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(b))
}