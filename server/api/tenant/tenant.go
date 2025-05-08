// Package tenant implements APIs for tenant management
package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/ent/tenant"
)

// Response represents the JSON response for a tenant
type Response struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	GivenName           string   `json:"given_name,omitempty"`
	AdminEmail          string   `json:"admin_email"`
	AllowRegistration   bool     `json:"allow_registration"`
	AllowedEmailDomains []string `json:"allowed_email_domains,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}

// RegisterRequest represents the request to register a new tenant
type RegisterRequest struct {
	Name                string   `json:"name" binding:"required"`
	GivenName           string   `json:"given_name" binding:"required"`
	AdminEmail          string   `json:"admin_email" binding:"required,email"`
	AllowRegistration   bool     `json:"allow_registration"`
	AllowedEmailDomains []string `json:"allowed_email_domains"`
	Description         string   `json:"description"`
}

// UpdateRequest represents the request to update a tenant
type UpdateRequest struct {
	GivenName           string   `json:"given_name"`
	AllowRegistration   bool     `json:"allow_registration"`
	AllowedEmailDomains []string `json:"allowed_email_domains"`
	Description         string   `json:"description"`
}

// UpdateAdminEmailRequest represents the request to update a tenant's admin email
type UpdateAdminEmailRequest struct {
	AdminEmail string `json:"admin_email" binding:"required,email"`
}

// Register handles tenant registration
func Register(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		// Check if tenant name already exists
		exists, err := client.Tenant.Query().Where(tenant.NameEQ(req.Name)).Exist(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "租户名称已存在"})
			return
		}

		// Create tenant
		t, err := client.Tenant.Create().
			SetName(req.Name).
			SetGivenName(req.GivenName).
			SetAdminEmail(req.AdminEmail).
			SetAllowRegistration(req.AllowRegistration).
			SetAllowedEmailDomains(req.AllowedEmailDomains).
			SetDescription(req.Description).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建租户失败", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Response{
			ID:                  t.ID.String(),
			Name:                t.Name,
			GivenName:           t.GivenName,
			AdminEmail:          t.AdminEmail,
			AllowRegistration:   t.AllowRegistration,
			AllowedEmailDomains: t.AllowedEmailDomains,
			Description:         t.Description,
			Enabled:             t.Enabled,
			CreatedAt:           t.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:           t.UpdatedAt.Format(http.TimeFormat),
		})
	}
}

// List returns all tenants (for system admin)
func List(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenants, err := client.Tenant.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取租户列表失败"})
			return
		}

		response := make([]Response, len(tenants))
		for i, t := range tenants {
			response[i] = Response{
				ID:                  t.ID.String(),
				Name:                t.Name,
				GivenName:           t.GivenName,
				AdminEmail:          t.AdminEmail,
				AllowRegistration:   t.AllowRegistration,
				AllowedEmailDomains: t.AllowedEmailDomains,
				Description:         t.Description,
				Enabled:             t.Enabled,
				CreatedAt:           t.CreatedAt.Format(http.TimeFormat),
				UpdatedAt:           t.UpdatedAt.Format(http.TimeFormat),
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// Update handles tenant information update (by tenant admin)
func Update(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Update tenant
		t, err := client.Tenant.UpdateOneID(tid).
			SetGivenName(req.GivenName).
			SetAllowRegistration(req.AllowRegistration).
			SetAllowedEmailDomains(req.AllowedEmailDomains).
			SetDescription(req.Description).
			Save(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新租户失败", "details": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:                  t.ID.String(),
			Name:                t.Name,
			GivenName:           t.GivenName,
			AdminEmail:          t.AdminEmail,
			AllowRegistration:   t.AllowRegistration,
			AllowedEmailDomains: t.AllowedEmailDomains,
			Description:         t.Description,
			Enabled:             t.Enabled,
			CreatedAt:           t.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:           t.UpdatedAt.Format(http.TimeFormat),
		})
	}
}

// UpdateAdminEmail handles updating the tenant admin email
func UpdateAdminEmail(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateAdminEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Update admin email
		t, err := client.Tenant.UpdateOneID(tid).
			SetAdminEmail(req.AdminEmail).
			Save(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新管理员邮箱失败", "details": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:                  t.ID.String(),
			Name:                t.Name,
			GivenName:           t.GivenName,
			AdminEmail:          t.AdminEmail,
			AllowRegistration:   t.AllowRegistration,
			AllowedEmailDomains: t.AllowedEmailDomains,
			Description:         t.Description,
			Enabled:             t.Enabled,
			CreatedAt:           t.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:           t.UpdatedAt.Format(http.TimeFormat),
		})
	}
}

// Enable enables a tenant (system admin only)
func Enable(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Enable tenant
		t, err := client.Tenant.UpdateOneID(tid).
			SetEnabled(true).
			Save(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "启用租户失败"})
			}
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:                  t.ID.String(),
			Name:                t.Name,
			GivenName:           t.GivenName,
			AdminEmail:          t.AdminEmail,
			AllowRegistration:   t.AllowRegistration,
			AllowedEmailDomains: t.AllowedEmailDomains,
			Description:         t.Description,
			Enabled:             t.Enabled,
			CreatedAt:           t.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:           t.UpdatedAt.Format(http.TimeFormat),
		})
	}
}

// Disable disables a tenant (system admin only)
func Disable(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Disable tenant
		t, err := client.Tenant.UpdateOneID(tid).
			SetEnabled(false).
			Save(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "禁用租户失败"})
			}
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:                  t.ID.String(),
			Name:                t.Name,
			GivenName:           t.GivenName,
			AdminEmail:          t.AdminEmail,
			AllowRegistration:   t.AllowRegistration,
			AllowedEmailDomains: t.AllowedEmailDomains,
			Description:         t.Description,
			Enabled:             t.Enabled,
			CreatedAt:           t.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:           t.UpdatedAt.Format(http.TimeFormat),
		})
	}
}
