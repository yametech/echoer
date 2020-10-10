package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
)

func (h *Handle) stepList(g *gin.Context) {
	metadataName := "metadata.name"
	query := g.Query(metadataName)
	if query != "" {
		query = fmt.Sprintf("%s=%s", metadataName, query)
	}
	results, err := h.List(common.DefaultNamespace, common.Step, query)
	if err != nil {
		InternalError(g, "list data error", err)
		return
	}
	g.JSON(http.StatusOK, results)
}

func (h *Handle) stepGet(g *gin.Context) {
	var result = &resource.Step{}
	uuid := g.Param("uuid")
	if uuid == "" {
		RequestParamsError(g, "get data param is wrong", nil)
		return
	}
	err := h.GetByUUID(common.DefaultNamespace, common.Step, uuid, result)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}

type ackStepState struct {
	common.Common
	Done bool `json:"done"`
}

func (h *Handle) ackStep(g *gin.Context) {
	ackStep := &ackStepState{}
	if err := g.ShouldBindJSON(ackStep); err != nil {
		RequestParamsError(g, "bind data param is wrong", err)
		return
	}
	step := &resource.Step{}
	err := h.GetByUUID(common.DefaultNamespace, common.Step, ackStep.UUID, step)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	step.Spec.Response.State = ackStep.AckState
	step.Spec.ActionRun.Done = ackStep.Done

	_, _, err = h.Apply(common.DefaultNamespace, common.Step, step.GetName(), step)
	if err != nil {
		InternalError(g, "apply data error", err)
		return
	}

	g.JSON(http.StatusOK, "")
}
