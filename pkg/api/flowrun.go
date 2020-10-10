package api

import (
	"fmt"
	"github.com/yametech/echoer/pkg/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/service"
	pb "github.com/yametech/echoer/proto"
)

func (h *Handle) flowRunCreate(g *gin.Context) {
	postRaw, err := g.GetRawData()
	if err != nil {
		RequestParamsError(g, "post data is wrong", err)
	}
	r := &resource.FlowRun{}
	if err := core.JSONRawToResource(postRaw, r); err != nil {
		RequestParamsError(g, "post data is wrong, can't not unmarshal", err)
		return
	}
	newObj, isUpdate, err := h.Apply(common.DefaultNamespace, common.FlowRunCollection, r.GetName(), r)
	if err != nil {
		InternalError(g, "store error", err)
		return
	}
	var envType = pb.EventType_Added
	if isUpdate {
		envType = pb.EventType_Modified
	}
	if err := service.NewService(h.IStorage).RecordEvent(envType, newObj, "create flow run item object"); err != nil {
		InternalError(g, "record event error", err)
		return
	}
	g.JSON(http.StatusOK, newObj)
}

func (h *Handle) flowRunList(g *gin.Context) {
	metadataName := "metadata.name"
	query := g.Query(metadataName)
	if query != "" {
		query = fmt.Sprintf("%s=%s", metadataName, query)
	}
	results, err := h.List(common.DefaultNamespace, common.FlowRunCollection, query)
	if err != nil {
		InternalError(g, "list data error", err)
		return
	}
	g.JSON(http.StatusOK, results)
}

func (h *Handle) flowRunGet(g *gin.Context) {
	var result = &resource.FlowRun{}
	name := g.Param("name")
	if name == "" {
		RequestParamsError(g, "get data param is wrong", nil)
		return
	}
	err := h.Get(common.DefaultNamespace, common.FlowRunCollection, name, result)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}

func (h *Handle) flowRunDelete(g *gin.Context) {
	var result = &resource.Flow{}
	name := g.Param("name")
	uuid := g.Param("uuid")
	if name == "" {
		RequestParamsError(g, "delete data param is wrong", nil)
		return
	}
	err := h.Delete(common.DefaultNamespace, common.FlowRunCollection, name, uuid)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}
