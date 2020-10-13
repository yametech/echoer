package action

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ciRequest struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	//ci action args
	Project    string `json:"project"`
	Version    string `json:"version"`
	RetryCount string `json:"retry_count"`
}

func CI(ctx *gin.Context) {
	request := &ciRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	fmt.Printf("recv (%v)\n", request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("ci", request.FlowId, request.StepName, request.AckState, request.UUID, true)
}
