package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/server/api/tenant"
	"github.com/telepair/terminal/server/api/user"
	"github.com/telepair/terminal/server/middleware"
)

// Server represents the API server
type Server struct {
	config *Config
	router *gin.Engine
	client *ent.Client
	srv    *http.Server
}

// NewServer creates a new server instance
func NewServer(config *Config, client *ent.Client) *Server {
	router := gin.Default()

	// Enable CORS if needed
	router.Use(middleware.CORS())

	return &Server{
		config: config,
		router: router,
		client: client,
	}
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
	// Public routes (no authentication required)
	public := s.router.Group("/api/v1")
	{
		// Tenant registration
		public.POST("/tenants/register", tenant.Register(s.client))

		// User registration
		public.POST("/users/register", user.Register(s.client))

		// User login
		public.POST("/users/login", user.Login(s.client))
	}

	// System admin routes
	sysAdmin := s.router.Group("/api/v1/admin")
	sysAdmin.Use(middleware.JWTAuth(s.client, true))
	{
		// Tenant management
		sysAdmin.GET("/tenants", tenant.List(s.client))
		sysAdmin.PUT("/tenants/:id/enable", tenant.Enable(s.client))
		sysAdmin.PUT("/tenants/:id/disable", tenant.Disable(s.client))
	}

	// Tenant admin routes
	tenantAdmin := s.router.Group("/api/v1/tenants/:tenant_id/admin")
	tenantAdmin.Use(middleware.JWTAuth(s.client, false))
	tenantAdmin.Use(middleware.TenantAdminOnly(s.client))
	{
		// Tenant management
		tenantAdmin.PUT("", tenant.Update(s.client))
		tenantAdmin.PUT("/admin-email", tenant.UpdateAdminEmail(s.client))

		// User management
		tenantAdmin.GET("/users", user.List(s.client))
		tenantAdmin.POST("/users", user.Create(s.client))
		tenantAdmin.POST("/users/batch", user.BatchCreate(s.client))
		tenantAdmin.PUT("/users/:id/enable", user.Enable(s.client))
		tenantAdmin.PUT("/users/:id/disable", user.Disable(s.client))
		tenantAdmin.PUT("/users/:id/lock", user.Lock(s.client))
		tenantAdmin.PUT("/users/:id/unlock", user.Unlock(s.client))
		tenantAdmin.PUT("/users/:id/admin", user.SetAdmin(s.client))
		tenantAdmin.DELETE("/users/:id/admin", user.RemoveAdmin(s.client))
	}

	// User routes
	userRoutes := s.router.Group("/api/v1/users")
	userRoutes.Use(middleware.JWTAuth(s.client, false))
	{
		userRoutes.GET("/me", user.GetProfile(s.client))
		userRoutes.PUT("/me", user.UpdateProfile(s.client))
		userRoutes.PUT("/me/password", user.UpdatePassword(s.client))
		userRoutes.POST("/logout", user.Logout(s.client))
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.srv = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	slog.Info("Server starting on " + addr)

	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
