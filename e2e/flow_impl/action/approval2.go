package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type approval2Request struct {
	FlowId    string   `json:"flowId"`
	StepName  string   `json:"stepName"`
	AckStates []string `json:"ackStates"`
	UUID      string   `json:"uuid"`
	//approval2 action args
	Project string `json:"project"`
	Version int64  `json:"version"`
}

func Approval2(ctx *gin.Context) {
	var name = "approval2"
	request := &approval2Request{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		fmt.Printf("action (%s) request bind error (%s)\n", name, err)
		return
	}
	fmt.Printf("action (%s) recv (%v)\n", name, request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("approval2", request.FlowId, request.StepName, request.AckStates[0], request.UUID, true)
}
