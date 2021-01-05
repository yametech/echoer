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
	name := g.Param("name")
	if name == "" {
		RequestParamsError(g, "get data param is wrong", nil)
		return
	}
	err := h.Get(common.DefaultNamespace, common.Step, name, result)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}

func (h *Handle) stepDelete(g *gin.Context) {
	var result = &resource.Step{}
	name := g.Param("name")
	uuid := g.Param("uuid")
	if name == "" || uuid == "" {
		RequestParamsError(g, "delete data param is wrong", nil)
		return
	}
	err := h.Delete(common.DefaultNamespace, common.Step, name, uuid)
	if err != nil {
		InternalError(g, "delete data error or maybe not found", err)
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
		fmt.Printf("[INFO] client (%s) post request (%s) wrong\n", g.Request.Host, g.Request.Body)
		return
	}

	step := &resource.Step{}
	err := h.GetByUUID(common.DefaultNamespace, common.Step, ackStep.UUID, step)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		fmt.Printf("[INFO] flowrun (%s) step (%s) query error (%s)\n", ackStep.FlowId, ackStep.StepName, err)
		return
	}

	step.Spec.Response.State = ackStep.AckState
	step.Spec.ActionRun.Done = ackStep.Done
	step.Spec.Data = ackStep.Data
	step.Spec.GlobalVariables = ackStep.GlobalVariables

	_, _, err = h.Apply(common.DefaultNamespace, common.Step, step.GetName(), step)
	if err != nil {
		InternalError(g, "apply data error", err)
		fmt.Printf("[INFO] flowrun (%s) step (%s) apply error (%s)\n", ackStep.FlowId, ackStep.StepName, err)
		return
	}

	g.JSON(http.StatusOK, "")
}
