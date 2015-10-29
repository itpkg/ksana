package base

import (
	"github.com/gin-gonic/gin"
)

func (p *BaseEngine) Mount(rt *gin.Engine) {

	rt.Static("/assets", "assets")

	//----------------------------------------
	users := rt.Group("/users")
	users.POST("/users/sign_in", func(c *gin.Context) {
	})
	users.POST("/users/sign_up", func(c *gin.Context) {
	})
	users.DELETE("/users/sign_out", func(c *gin.Context) {
	})
	users.POST("/users/confirm", func(c *gin.Context) {
	})
	users.POST("/users/unlock", func(c *gin.Context) {
	})
	users.POST("/users/forgot_password", func(c *gin.Context) {
	})
	users.POST("/users/reset_password", func(c *gin.Context) {
	})
	users.POST("/users/profile", func(c *gin.Context) {
	})
}
