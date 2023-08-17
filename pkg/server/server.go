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

	r.GET("/users", middleware.DBMiddleware(db), middleware.AuthMiddleware(db), user.Get)
	r.GET("/users/:id", middleware.DBMiddleware(db), middleware.AuthMiddleware(db), user.GetByID)
	r.PUT("/users", middleware.DBMiddleware(db), middleware.AuthMiddleware(db), user.Update)
	r.POST("/users", middleware.DBMiddleware(db), user.Post)
	r.POST("/login", middleware.DBMiddleware(db), auth.Login)
	r.GET("/logout", middleware.DBMiddleware(db), auth.Logout)
	r.DELETE("/users/:id", middleware.DBMiddleware(db), middleware.AuthMiddleware(db), user.Delete)

	return r
}
