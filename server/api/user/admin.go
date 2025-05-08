// Package user implements APIs for user management with admin operations
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/ent/user"
)

// List returns all users for a tenant (tenant admin only)
func List(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get users for tenant
		users, err := client.User.Query().
			Where(user.TenantID(tid)).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
			return
		}

		response := make([]Response, len(users))
		for i, u := range users {
			response[i] = Response{
				ID:            u.ID.String(),
				Username:      u.Username,
				Email:         u.Email,
				GivenName:     u.GivenName,
				EmailVerified: u.EmailVerified,
				CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
				UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
				Enabled:       u.Enabled,
				IsLocked:      u.IsLocked,
				IsSuperuser:   u.IsSuperuser,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// Enable enables a user (tenant admin only)
func Enable(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get user ID from URL param
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Check if user belongs to tenant
		exists, err := client.User.Query().
			Where(
				user.ID(uid),
				user.TenantID(tid),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或不属于该租户"})
			return
		}

		// Enable user
		u, err := client.User.UpdateOneID(uid).
			SetEnabled(true).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "启用用户失败"})
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:            u.ID.String(),
			Username:      u.Username,
			Email:         u.Email,
			GivenName:     u.GivenName,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
			Enabled:       u.Enabled,
			IsLocked:      u.IsLocked,
		})
	}
}

// Disable disables a user (tenant admin only)
func Disable(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get user ID from URL param
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Check if user belongs to tenant
		exists, err := client.User.Query().
			Where(
				user.ID(uid),
				user.TenantID(tid),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或不属于该租户"})
			return
		}

		// Disable user
		u, err := client.User.UpdateOneID(uid).
			SetEnabled(false).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "禁用用户失败"})
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:            u.ID.String(),
			Username:      u.Username,
			Email:         u.Email,
			GivenName:     u.GivenName,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
			Enabled:       u.Enabled,
			IsLocked:      u.IsLocked,
		})
	}
}

// Lock locks a user (tenant admin only)
func Lock(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LockUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get user ID from URL param
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Check if user belongs to tenant
		exists, err := client.User.Query().
			Where(
				user.ID(uid),
				user.TenantID(tid),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或不属于该租户"})
			return
		}

		// Lock user
		u, err := client.User.UpdateOneID(uid).
			SetIsLocked(true).
			SetLockReason(req.Reason).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "锁定用户失败"})
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:            u.ID.String(),
			Username:      u.Username,
			Email:         u.Email,
			GivenName:     u.GivenName,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
			Enabled:       u.Enabled,
			IsLocked:      u.IsLocked,
		})
	}
}

// Unlock unlocks a user (tenant admin only)
func Unlock(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get user ID from URL param
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Check if user belongs to tenant
		exists, err := client.User.Query().
			Where(
				user.ID(uid),
				user.TenantID(tid),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或不属于该租户"})
			return
		}

		// Unlock user
		u, err := client.User.UpdateOneID(uid).
			SetIsLocked(false).
			ClearLockReason().
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解锁用户失败"})
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:            u.ID.String(),
			Username:      u.Username,
			Email:         u.Email,
			GivenName:     u.GivenName,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
			Enabled:       u.Enabled,
			IsLocked:      u.IsLocked,
		})
	}
}

// SetAdmin sets a user as admin (tenant admin only)
func SetAdmin(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get user ID from URL param
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Check if user belongs to tenant
		exists, err := client.User.Query().
			Where(
				user.ID(uid),
				user.TenantID(tid),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或不属于该租户"})
			return
		}

		// Set user as admin
		u, err := client.User.UpdateOneID(uid).
			SetIsSuperuser(true).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "设置管理员失败"})
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:            u.ID.String(),
			Username:      u.Username,
			Email:         u.Email,
			GivenName:     u.GivenName,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
			Enabled:       u.Enabled,
			IsLocked:      u.IsLocked,
			IsSuperuser:   u.IsSuperuser,
		})
	}
}

// RemoveAdmin removes admin status from a user (tenant admin only)
func RemoveAdmin(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from URL param
		tenantID := c.Param("tenant_id")
		tid, err := uuid.Parse(tenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get user ID from URL param
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Check if user belongs to tenant
		exists, err := client.User.Query().
			Where(
				user.ID(uid),
				user.TenantID(tid),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在或不属于该租户"})
			return
		}

		// Remove admin status
		u, err := client.User.UpdateOneID(uid).
			SetIsSuperuser(false).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "移除管理员失败"})
			return
		}

		c.JSON(http.StatusOK, Response{
			ID:            u.ID.String(),
			Username:      u.Username,
			Email:         u.Email,
			GivenName:     u.GivenName,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:     u.UpdatedAt.Format(http.TimeFormat),
			Enabled:       u.Enabled,
			IsLocked:      u.IsLocked,
			IsSuperuser:   u.IsSuperuser,
		})
	}
}
