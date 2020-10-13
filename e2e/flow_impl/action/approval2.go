package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type approval2Request struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	//approval2 action args
	Project string `json:"project"`
	Version int    `json:"version"`
}

func Approval2(ctx *gin.Context) {
	request := &approval2Request{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	fmt.Printf("approval2 action recv (%v)\n", request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("approval", request.FlowId, request.StepName, request.AckState, request.UUID, true)
}
