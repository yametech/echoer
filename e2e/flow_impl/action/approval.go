package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type approvalRequest struct {
	FlowId    string   `json:"flowId"`
	StepName  string   `json:"stepName"`
	AckStates []string `json:"ackStates"`
	UUID      string   `json:"uuid"`
	//approval action args
	WorkOrder string `json:"work_order"`
	Version   int64  `json:"version"`
}

func Approval(ctx *gin.Context) {
	var name = "approval"
	request := &approvalRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		fmt.Printf("action (%s) request bind error (%s)\n", name, err)
		return
	}
	fmt.Printf("action (%s) recv (%v)\n", name, request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("approval", request.FlowId, request.StepName, request.AckStates[0], request.UUID, true)
}
