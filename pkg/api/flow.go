package api

import (
	"fmt"
	"net/http"

	"github.com/yametech/echoer/pkg/core"

	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/service"
	pb "github.com/yametech/echoer/proto"
)

func (h *Handle) flowCreate(g *gin.Context) {
	postRaw, err := g.GetRawData()
	if err != nil {
		RequestParamsError(g, "post data is wrong", err)
	}
	r := &resource.Flow{}
	if err := core.JSONRawToResource(postRaw, r); err != nil {
		RequestParamsError(g, "post data is wrong, can't not unmarshal", err)
		return
	}
	newObj, isUpdate, err := h.Apply(common.DefaultNamespace, common.FlowCollection, r.GetName(), r)
	if err != nil {
		InternalError(g, "store error", err)
		return
	}
	var envType = pb.EventType_Added
	if isUpdate {
		envType = pb.EventType_Modified
	}
	if err := service.NewService(h.IStorage).RecordEvent(envType, newObj, "create flow item object"); err != nil {
		InternalError(g, "record event error", err)
		return
	}
	g.JSON(http.StatusOK, newObj)
}

func (h *Handle) flowList(g *gin.Context) {
	metadataName := "metadata.name"
	query := g.Query(metadataName)
	if query != "" {
		query = fmt.Sprintf("%s=%s", metadataName, query)
	}
	results, err := h.List(common.DefaultNamespace, common.FlowCollection, query)
	if err != nil {
		InternalError(g, "list data error", err)
		return
	}
	g.JSON(http.StatusOK, results)
}

func (h *Handle) flowGet(g *gin.Context) {
	var result = &resource.Flow{}
	name := g.Param("name")
	if name == "" {
		RequestParamsError(g, "get data param is wrong", nil)
		return
	}
	err := h.Get(common.DefaultNamespace, common.FlowCollection, name, result)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}

func (h *Handle) flowDelete(g *gin.Context) {
	var result = &resource.Flow{}
	name := g.Param("name")
	uuid := g.Param("uuid")
	if name == "" {
		RequestParamsError(g, "delete data param is wrong", nil)
		return
	}
	err := h.Delete(common.DefaultNamespace, common.FlowCollection, name, uuid)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}
