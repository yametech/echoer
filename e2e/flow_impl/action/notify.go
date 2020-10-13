package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type notifyRequest struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	//notify action args
	Project string `json:"project"`
	Version int    `json:"version"`
}

func Notify(ctx *gin.Context) {
	request := &notifyRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	fmt.Printf("notify action recv (%v)\n", request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("notify", request.FlowId, request.StepName, request.AckState, request.UUID, true)
}
