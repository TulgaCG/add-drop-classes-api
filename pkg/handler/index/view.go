package index

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

func View(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", pongo2.Context{})
}
