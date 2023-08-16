package server

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/auth"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/user"
)

func New(db *gendb.Queries) *gin.Engine {
	r := gin.Default()

	r.GET("/users", middleware.DBMiddleware(db), user.Get)
	r.GET("/users/:id", middleware.DBMiddleware(db), user.GetByID)
	r.PUT("/users", middleware.DBMiddleware(db), user.Update)
	r.POST("/users", middleware.DBMiddleware(db), auth.GetByUsername)
	r.DELETE("/users/:id", middleware.DBMiddleware(db), user.Delete)

	return r
}
