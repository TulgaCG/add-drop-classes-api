package server

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gen/db"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/user"
)

func New(db *db.Queries) *gin.Engine {
	r := gin.Default()

	r.GET("/users", middleware.DBMiddleware(db), user.Get)
	r.GET("/users/id/:id", middleware.DBMiddleware(db), user.GetByID)
	r.GET("/users/username/:username", middleware.DBMiddleware(db), user.GetUserByUsername)
	r.PUT("/users", middleware.DBMiddleware(db), user.Update)
	r.POST("/users", middleware.DBMiddleware(db), user.Post)
	r.DELETE("/users/id/:id", middleware.DBMiddleware(db), user.Delete)

	return r
}
