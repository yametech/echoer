package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/echoer/e2e/flow_impl/action"
)

func main() {
	route := gin.New()

	route.POST("/ci", action.CI)
	route.POST("/deploy1", action.Deploy1)
	route.POST("/approval", action.Approval)
	route.POST("/approval2", action.Approval2)
	route.POST("/notify", action.Notify)

	if err := route.Run(":18080"); err != nil {
		panic(err)
	}
}
