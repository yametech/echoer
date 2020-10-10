package api

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
)

type Handle struct{ *Server }

// watch
/*
	watch provide resource stream
	example:
      /watch?resource=action?version=1597920529&resource=workflow?version=1597920529
      res: action=>1597920529
           workflow=>1597920529
*/

func (h *Handle) watch(g *gin.Context) {
	objectChan := make(chan core.IObject, 32)
	closed := make(chan struct{})
	resources := g.QueryArray("resource")

	for _, res := range resources {
		go func(res string) {
			resList := strings.Split(res, "?")
			if len(resList) != 2 {
				return
			}
			versionList := strings.Split(resList[1], "=")
			if len(versionList) != 2 {
				return
			}
			version, err := strconv.ParseInt(versionList[1], 10, 64)
			if err != nil {
				return
			}
			resType := resList[0]
			coder := storage.GetResourceCoder(resType)
			if coder == nil {
				return
			}
			wc := storage.NewWatch(coder)
			h.Watch2(common.DefaultNamespace, resType, version, wc)
			for {
				select {
				case <-closed:
					return
				case err := <-wc.ErrorStop():
					fmt.Printf("[ERROR] watch type: (%s) version: (%d) error: (%s)\n", resType, version, err)
					close(objectChan)
					return
				case item, ok := <-wc.ResultChan():
					if !ok {
						return
					}
					objectChan <- item
				}
			}
		}(res)
	}

	streamEndEvent := "STREAM_END"

	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Writer.CloseNotify():
			closed <- struct{}{}
			close(closed)
			g.SSEvent("", streamEndEvent)
			return false
		case object, ok := <-objectChan:
			if !ok {
				g.SSEvent("", streamEndEvent)
				return false
			}
			g.SSEvent("", object)
		}
		return true
	},
	)

}
