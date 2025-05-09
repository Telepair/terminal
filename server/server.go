package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/telepair/terminal/pkg/version"
)

// StartServer initializes and starts the HTTP server.
func StartServer(addr string) error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/health", healthHandler)
		api.GET("/version", versionHandler)
	}

	return router.Run(addr)
}

// StartServerWithContext initializes and starts the HTTP server with graceful shutdown support.
func StartServerWithContext(ctx context.Context, addr string) error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/health", healthHandler)
		api.GET("/version", versionHandler)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    version.Version,
		"git_commit": version.GitCommit,
		"git_tag":    version.GitTag,
		"git_branch": version.GitBranch,
		"build_date": version.BuildDate,
	})
}
