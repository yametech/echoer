package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type request struct {
	FlowId           string   `json:"flowId"`
	StepName         string   `json:"stepName"`
	AckStates        []string `json:"ackStates"`
	UUID             string   `json:"uuid"`
	Pipeline         string   `json:"pipeline"`
	PipelineResource string   `json:"pipelineResource"`
}

type response struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	Done     bool   `json:"done"`
}

var (
	currentReq   *request
	responseChan = make(chan struct{})
)

func resp(url string) {
	for {
		<-responseChan
		req := resty.New()
		resp := &response{
			FlowId:   currentReq.FlowId,
			StepName: currentReq.StepName,
			AckState: currentReq.AckStates[0],
			UUID:     currentReq.UUID,
			Done:     true,
		}
		body, _ := json.Marshal(resp)
		fmt.Printf("response to api-server %s\n", body)
		_, err := req.R().SetHeader("Accept", "application/json").SetBody(body).Post(url)
		fmt.Println(err)
	}
}

func main() {
	go resp("http://127.0.0.1:8080/step")

	route := gin.New()
	route.POST("/", func(ctx *gin.Context) {
		ci := &request{}
		if err := ctx.BindJSON(ci); err != nil {
			ctx.JSON(http.StatusBadRequest, "")
			return
		}
		fmt.Printf("recv (%v)\n", ci)
		currentReq = ci
		ctx.JSON(http.StatusOK, "")
		responseChan <- struct{}{}
	})

	route.Run(":18080")
}
