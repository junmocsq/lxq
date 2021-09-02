package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {

}

// User godoc
// @Summary Show an account
// @Description get string by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} user.User
// @Failure 400 {object} user.User
// @Router /user/{id} [get]
func (c *UserController) User(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(http.StatusOK, id)
}
