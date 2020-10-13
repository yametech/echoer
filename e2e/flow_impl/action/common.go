package action

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Response struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	Done     bool   `json:"done"`
}

func RespToApiServer(action, flowId, stepName, ackState, uuid string, done bool) {
	req := resty.New()
	resp := &Response{
		FlowId:   flowId,
		StepName: stepName,
		AckState: strings.Split(ackState, ",")[0],
		UUID:     uuid,
		Done:     done,
	}
	body, _ := json.Marshal(resp)
	fmt.Printf("response to api-server %s\n", body)
	_, err := req.R().SetHeader("Accept", "application/json").SetBody(body).Post(ApiServerAddr)
	if err != nil {
		fmt.Println(err)
	}

}

var ApiServerAddr = "http://localhost:8080/step"

func randomAckState(states string) string {
	stateList := strings.Split(states, ",")
	return stateList[generateLimitedRandNum(len(stateList)-1)]
}

func generateLimitedRandNum(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}
