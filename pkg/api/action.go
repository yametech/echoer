package api

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/echoer/pkg/factory"
	"github.com/yametech/echoer/pkg/fsm/fss"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
)

func (h *Handle) actionCreate(g *gin.Context) {
	postRaw, err := g.GetRawData()
	if err != nil {
		RequestParamsError(g, "post data is wrong", err)
	}
	r := &createRawData{}
	if err := json.Unmarshal(postRaw, r); err != nil {
		RequestParamsError(g, "post data is wrong, can't not unmarshal", err)
		return
	}
	fsl := r.Data
	stmt, err := fss.NewActionFSLParser().Parse(fsl)
	if err != nil {
		InternalError(g, "fsl parse error", err)
		return
	}
	if err := factory.NewTranslation(factory.NewStoreImpl(h.IStorage)).ToAction(stmt); err != nil {
		InternalError(g, "fsl parse error", err)
		return
	}
	g.JSON(http.StatusOK, "")
}

func (h *Handle) actionList(g *gin.Context) {
	metadataName := "metadata.name"
	query := g.Query(metadataName)
	if query != "" {
		query = fmt.Sprintf("%s=%s", metadataName, query)
	}
	results, err := h.List(common.DefaultNamespace, common.ActionCollection, query)
	if err != nil {
		InternalError(g, "list data error", err)
		return
	}
	g.JSON(http.StatusOK, results)
}

func (h *Handle) actionGet(g *gin.Context) {
	var result = &resource.Action{}
	name := g.Param("name")
	if name == "" {
		RequestParamsError(g, "get data param is wrong", nil)
		return
	}
	err := h.Get(common.DefaultNamespace, common.ActionCollection, name, result)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}

func (h *Handle) actionDelete(g *gin.Context) {
	var result = &resource.Action{}
	name := g.Param("name")
	uuid := g.Param("uuid")
	if name == "" || uuid == "" {
		RequestParamsError(g, "delete data param is wrong", nil)
		return
	}
	err := h.Delete(common.DefaultNamespace, common.ActionCollection, name, uuid)
	if err != nil {
		InternalError(g, "get data error or maybe not found", err)
		return
	}
	g.JSON(http.StatusOK, result)
}
