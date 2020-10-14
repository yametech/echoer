package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deployRequest struct {
	FlowId    string   `json:"flowId"`
	StepName  string   `json:"stepName"`
	AckStates []string `json:"ackStates"`
	UUID      string   `json:"uuid"`
	//deploy1 action args
	Project string `json:"project"`
	Version int64  `json:"version"`
}

func Deploy1(ctx *gin.Context) {
	var name = "deploy1"
	request := &deployRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		fmt.Printf("action (%s) request bind error (%s)\n", name, err)
		return
	}
	fmt.Printf("action (%s) recv (%v)\n", name, request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("deploy1", request.FlowId, request.StepName, request.AckStates[0], request.UUID, true)
}
