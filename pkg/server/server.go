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

	// Auth required requests
	g1 := r.Group("/api", middleware.DBMiddleware(db), middleware.AuthMiddleware(db))
	g1.GET("/users", user.Get)
	g1.GET("/users/:id", user.GetByID)
	g1.PUT("/users", user.Update)
	g1.DELETE("/users", user.Delete)

	g2 := r.Group("/api", middleware.DBMiddleware(db))
	g2.GET("/logout", auth.Logout)
	g2.POST("/login", auth.Login)
	g2.POST("/users", user.Post)

	return r
}
