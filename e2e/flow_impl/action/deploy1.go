package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deployRequest struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	//deploy1 action args
	Project string `json:"project"`
	Version int    `json:"version"`
}

func Deploy1(ctx *gin.Context) {
	request := &deployRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	fmt.Printf("deploy1 action recv (%v)\n", request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("deploy1", request.FlowId, request.StepName, request.AckState, request.UUID, true)
}
