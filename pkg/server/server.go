package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/auth"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/user"
)

func New(db *gendb.Queries, log *slog.Logger) *gin.Engine {
	r := gin.Default()

	g1 := r.Group("/api", middleware.LogMiddleware(log), middleware.DBMiddleware(db))
	g1.GET("/logout", auth.LogoutHandler)
	g1.POST("/login", auth.LoginHandler)
	g1.POST("/users", user.CreateHandler)

	g2 := r.Group("/api", middleware.LogMiddleware(log), middleware.DBMiddleware(db), middleware.AuthMiddleware(db))
	g2.GET("/users", user.ListHandler)
	g2.GET("/users/:id", user.GetHandler)
	g2.PUT("/users", user.UpdateHandler)

	return r
}
