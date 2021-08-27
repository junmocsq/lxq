package lxq

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func group(engine *gin.Engine) {
	v1 := engine.Group("/v1", func(c *gin.Context) {
		fmt.Println("/v1")
	})
	{
		v1.GET("/test1", func(c *gin.Context) {
			fmt.Println("/v1/test1")
			c.JSON(200, map[string]string{
				"junmo": "csq",
			})
		})
	}
}
