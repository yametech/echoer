package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type notifyRequest struct {
	FlowId    string   `json:"flowId"`
	StepName  string   `json:"stepName"`
	AckStates []string `json:"ackStates"`
	UUID      string   `json:"uuid"`
	//notify action args
	Project string `json:"project"`
	Version int64  `json:"version"`
}

func Notify(ctx *gin.Context) {
	var name = "notify"
	request := &notifyRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		fmt.Printf("action (%s) request bind error (%s)\n", name, err)
		return
	}
	fmt.Printf("action (%s) recv (%v)\n", name, request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("notify", request.FlowId, request.StepName, request.AckStates[0], request.UUID, true)
}
