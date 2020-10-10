package api

import (
	"github.com/yametech/echoer/pkg/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
)

func (h *Handle) eventCreate(g *gin.Context) {
	postRaw, err := g.GetRawData()
	if err != nil {
		RequestParamsError(g, "post data param is wrong", err)
		return
	}
	r := &resource.Event{}
	if err := core.JSONRawToResource(postRaw, r); err != nil {
		RequestParamsError(g, "post data is wrong, can't not unmarshal", err)
	}
	newObj, err := h.Create(common.DefaultNamespace, common.EventCollection, r)
	if err != nil {
		InternalError(g, "store error", err)
		return
	}
	g.JSON(http.StatusOK, newObj)
}

func (h *Handle) eventList(g *gin.Context) {
	results, err := h.List(common.DefaultNamespace, common.EventCollection, "")
	if err != nil {
		InternalError(g, "list data error", err)
		return
	}
	g.JSON(http.StatusOK, results)
}
