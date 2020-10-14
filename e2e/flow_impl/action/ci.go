package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ciRequest struct {
	FlowId    string   `json:"flowId"`
	StepName  string   `json:"stepName"`
	AckStates []string `json:"ackStates"`
	UUID      string   `json:"uuid"`
	//ci action args
	Project    string `json:"project"`
	Version    string `json:"version"`
	RetryCount int64  `json:"retry_count"`
}

func CI(ctx *gin.Context) {
	var name = "ci"
	request := &ciRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		fmt.Printf("action (%s) request bind error (%s)\n", name, err)
		return
	}
	fmt.Printf("action (%s) recv (%v)\n", name, request)
	ctx.JSON(http.StatusOK, "")

	RespToApiServer("ci", request.FlowId, request.StepName, request.AckStates[0], request.UUID, true)
}
