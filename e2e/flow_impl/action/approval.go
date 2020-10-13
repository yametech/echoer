package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type approvalRequest struct {
	FlowId   string `json:"flowId"`
	StepName string `json:"stepName"`
	AckState string `json:"ackState"`
	UUID     string `json:"uuid"`
	//approval action args
	WorkOrder string `json:"work_order"`
	Version   int    `json:"version"`
}

func Approval(ctx *gin.Context) {
	request := &approvalRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}
	fmt.Printf("approval action recv (%v)\n", request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("approval", request.FlowId, request.StepName, request.AckState, request.UUID, true)
}
