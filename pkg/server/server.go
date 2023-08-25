package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/go-box/pongo2gin/v6"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/env"
	"github.com/TulgaCG/add-drop-classes-api/pkg/handler/auth"
	"github.com/TulgaCG/add-drop-classes-api/pkg/handler/index"
	"github.com/TulgaCG/add-drop-classes-api/pkg/handler/lecture"
	"github.com/TulgaCG/add-drop-classes-api/pkg/handler/user"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/role"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const (
	gracefulShutdownTimeoutDuration = 3 * time.Second
	readHeaderTimeout               = 1 * time.Second
)

type Server struct {
	opts    *Opts
	handler *gin.Engine
}

func (s Server) Run() error {
	sv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.opts.Port),
		Handler:           s.handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	if s.opts.Environment == env.LevelProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	go func() {
		s.opts.Log.Info("server is listening")
		if err := sv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.opts.Log.Error(fmt.Sprintf("server failed to listen and serve: %s", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.opts.Log.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeoutDuration)
	defer cancel()

	if err := sv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	<-ctx.Done()
	return nil
}

func New(opts ...Opt) *Server {
	options := &Opts{}
	for _, opt := range opts {
		opt(options)
	}

	r := gin.Default()
	r.HTMLRender = pongo2gin.Default()

	ac := &appctx.Context{
		DB:  options.DB,
		Log: options.Log,
	}
	{
		api := r.Group("/api", middleware.AppContext(ac))

		api.GET("/logout", auth.Logout)
		api.POST("/login", auth.Login)
		api.POST("/register", auth.Register)

		api.GET("/lectures/:id", middleware.Auth(options.DB), middleware.Authz(options.DB), lecture.GetFromUser)
		api.POST("/lectures", middleware.Auth(options.DB), middleware.Authz(options.DB), lecture.AddToUser)
		api.DELETE("/lectures/:id/:uid", middleware.Auth(options.DB), middleware.Authz(options.DB), lecture.RemoveFromUser)

		api.GET("/users", middleware.AuthzWithRoles(options.DB, types.RoleAdmin, types.RoleTeacher), user.List)
		api.GET("/users/:id", middleware.Auth(options.DB), middleware.Authz(options.DB), user.Get)
		api.PUT("/users", middleware.Auth(options.DB), middleware.Authz(options.DB), user.Update)

		api.DELETE("/roles/:uid/:rid", middleware.Auth(options.DB), middleware.AuthzWithRoles(options.DB, types.RoleAdmin), role.RemoveFromUserHandler)
		api.POST("/roles", middleware.Auth(options.DB), middleware.AuthzWithRoles(options.DB, types.RoleAdmin), role.AddToUserHandler)
	}

	{
		view := r.Group("/", middleware.AppContext(ac))

		view.GET("/", index.View)
		view.GET("/login", auth.LoginView)
		view.GET("/register", auth.LoginView)
	}

	return &Server{
		handler: r,
		opts:    options,
	}
}

type Opts struct {
	DB          *database.DB
	Log         *slog.Logger
	Environment env.Level
	Port        uint32
}

type Opt func(*Opts)

func WithDB(db *database.DB) Opt {
	return func(o *Opts) {
		o.DB = db
	}
}

func WithEnv(environment env.Level) Opt {
	return func(o *Opts) {
		o.Environment = environment
	}
}

func WithLogger(logger *slog.Logger) Opt {
	return func(o *Opts) {
		o.Log = logger
	}
}

func WithPort(port uint32) Opt {
	return func(o *Opts) {
		o.Port = port
	}
}
