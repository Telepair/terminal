package e2e_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/server/api/tenant"
	"github.com/telepair/terminal/server/api/user"
	"github.com/telepair/terminal/server/middleware"

	_ "github.com/mattn/go-sqlite3" // 导入 SQLite3 驱动
)

// global variables for testing
var (
	ctx        context.Context
	client     *ent.Client
	testRouter *gin.Engine
	srv        *httptest.Server
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API E2E Suite")
}

var _ = BeforeSuite(func() {
	// Use test mode for Gin
	gin.SetMode(gin.TestMode)

	// Create a context with timeout for setup
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	DeferCleanup(cancel)

	// Setup in-memory SQLite for testing
	var err error
	client, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	Expect(err).NotTo(HaveOccurred())
	DeferCleanup(func() {
		Expect(client.Close()).To(Succeed())
	})

	// Run schema migration
	err = client.Schema.Create(ctx)
	Expect(err).NotTo(HaveOccurred())

	// Create test router
	testRouter = gin.New()
	testRouter.Use(gin.Recovery())

	// Setup routes similar to the actual application
	setupTestRoutes(testRouter, client)

	// Create test server
	srv = httptest.NewServer(testRouter)
	DeferCleanup(srv.Close)
})

// 在每个测试前重置数据库
var _ = BeforeEach(func() {
	// 清除数据库中的所有数据
	_, err := client.User.Delete().Exec(ctx)
	Expect(err).NotTo(HaveOccurred())

	_, err = client.Tenant.Delete().Exec(ctx)
	Expect(err).NotTo(HaveOccurred())
})

// setupTestRoutes sets up the routes for testing
func setupTestRoutes(r *gin.Engine, c *ent.Client) {
	// 测试服务器需要使用完整的 URL 路径，去掉 API 前缀
	// Public routes (no authentication required)
	public := r.Group("")
	{
		// Tenant registration
		public.POST("/tenants/register", tenant.Register(c))

		// User registration
		public.POST("/users/register", user.Register(c))

		// User login
		public.POST("/users/login", user.Login(c))
	}

	// System admin routes
	sysAdmin := r.Group("/admin")
	sysAdmin.Use(middleware.JWTAuth(c, true))
	{
		// Tenant management
		sysAdmin.GET("/tenants", tenant.List(c))
		sysAdmin.PUT("/tenants/:id/enable", tenant.Enable(c))
		sysAdmin.PUT("/tenants/:id/disable", tenant.Disable(c))
	}

	// Tenant admin routes
	tenantAdmin := r.Group("/tenants/:tenant_id/admin")
	tenantAdmin.Use(middleware.JWTAuth(c, false))
	tenantAdmin.Use(middleware.TenantAdminOnly(c))
	{
		// Tenant management
		tenantAdmin.PUT("", tenant.Update(c))
		tenantAdmin.PUT("/admin-email", tenant.UpdateAdminEmail(c))

		// User management
		tenantAdmin.GET("/users", user.List(c))
		tenantAdmin.POST("/users", user.Create(c))
		tenantAdmin.POST("/users/batch", user.BatchCreate(c))
		tenantAdmin.PUT("/users/:id/enable", user.Enable(c))
		tenantAdmin.PUT("/users/:id/disable", user.Disable(c))
		tenantAdmin.PUT("/users/:id/lock", user.Lock(c))
		tenantAdmin.PUT("/users/:id/unlock", user.Unlock(c))
		tenantAdmin.PUT("/users/:id/admin", user.SetAdmin(c))
		tenantAdmin.DELETE("/users/:id/admin", user.RemoveAdmin(c))
	}

	// User routes
	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.JWTAuth(c, false))
	{
		userRoutes.GET("/me", user.GetProfile(c))
		userRoutes.PUT("/me", user.UpdateProfile(c))
		userRoutes.PUT("/me/password", user.UpdatePassword(c))
		userRoutes.POST("/logout", user.Logout(c))
	}
}
