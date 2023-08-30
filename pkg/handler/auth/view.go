package auth

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

func LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", pongo2.Context{})
}
