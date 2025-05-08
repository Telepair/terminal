// Package user implements APIs for user management
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/ent/user"
	"golang.org/x/crypto/bcrypt"
)

// Response represents the JSON response for a user
type Response struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name,omitempty"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Enabled       bool   `json:"enabled"`
	IsLocked      bool   `json:"is_locked"`
	IsSuperuser   bool   `json:"is_superuser,omitempty"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	GivenName string `json:"given_name"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

// BatchCreateUserRequest represents the request to batch create users
type BatchCreateUserRequest struct {
	Users []CreateUserRequest `json:"users" binding:"required,dive"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	GivenName string `json:"given_name"`
}

// UpdatePasswordRequest represents the request to update a user's password
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// LockUserRequest represents the request to lock a user
type LockUserRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// GetProfile returns the current user's profile
func GetProfile(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Get user
		u, err := client.User.Get(c.Request.Context(), uid)
		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
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

// UpdateProfile updates the current user's profile
func UpdateProfile(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Update user
		u, err := client.User.UpdateOneID(uid).
			SetGivenName(req.GivenName).
			Save(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败", "details": err.Error()})
			}
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

// UpdatePassword updates the current user's password
func UpdatePassword(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdatePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// Get user
		u, err := client.User.Get(c.Request.Context(), uid)
		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
			return
		}

		// Verify old password
		err = bcrypt.CompareHashAndPassword([]byte(u.PasswordBcrypt), []byte(req.OldPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "原密码错误"})
			return
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
			return
		}

		// Update password
		_, err = client.User.UpdateOneID(uid).
			SetPasswordBcrypt(string(hashedPassword)).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "密码已更新"})
	}
}

// Create creates a new user (tenant admin only)
func Create(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
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

		// Check if username already exists for this tenant
		exists, err := client.User.Query().
			Where(
				user.TenantID(tid),
				user.UsernameEQ(req.Username),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}

		// Check if email already exists for this tenant
		exists, err = client.User.Query().
			Where(
				user.TenantID(tid),
				user.EmailEQ(req.Email),
			).Exist(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
			return
		}

		// Create user
		u, err := client.User.Create().
			SetTenantID(tid).
			SetUsername(req.Username).
			SetEmail(req.Email).
			SetPasswordBcrypt(string(hashedPassword)).
			SetEmailVerified(false).
			SetGivenName(req.GivenName).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Response{
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

// BatchCreate batch creates users (tenant admin only)
func BatchCreate(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchCreateUserRequest
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

		// Create users in a transaction
		tx, err := client.Tx(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "开始事务失败"})
			return
		}

		responses := make([]Response, 0, len(req.Users))
		errors := make([]map[string]string, 0)

		for _, userReq := range req.Users {
			// Check if username already exists for this tenant
			exists, err := tx.User.Query().
				Where(
					user.TenantID(tid),
					user.UsernameEQ(userReq.Username),
				).Exist(c.Request.Context())

			if err != nil {
				errors = append(errors, map[string]string{
					"username": userReq.Username,
					"email":    userReq.Email,
					"error":    "检查用户名失败",
				})
				continue
			}
			if exists {
				errors = append(errors, map[string]string{
					"username": userReq.Username,
					"email":    userReq.Email,
					"error":    "用户名已存在",
				})
				continue
			}

			// Check if email already exists for this tenant
			exists, err = tx.User.Query().
				Where(
					user.TenantID(tid),
					user.EmailEQ(userReq.Email),
				).Exist(c.Request.Context())

			if err != nil {
				errors = append(errors, map[string]string{
					"username": userReq.Username,
					"email":    userReq.Email,
					"error":    "检查邮箱失败",
				})
				continue
			}
			if exists {
				errors = append(errors, map[string]string{
					"username": userReq.Username,
					"email":    userReq.Email,
					"error":    "邮箱已存在",
				})
				continue
			}

			// Hash password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
			if err != nil {
				errors = append(errors, map[string]string{
					"username": userReq.Username,
					"email":    userReq.Email,
					"error":    "密码处理失败",
				})
				continue
			}

			// Create user
			u, err := tx.User.Create().
				SetTenantID(tid).
				SetUsername(userReq.Username).
				SetEmail(userReq.Email).
				SetPasswordBcrypt(string(hashedPassword)).
				SetEmailVerified(false).
				SetGivenName(userReq.GivenName).
				Save(c.Request.Context())

			if err != nil {
				errors = append(errors, map[string]string{
					"username": userReq.Username,
					"email":    userReq.Email,
					"error":    "创建用户失败: " + err.Error(),
				})
				continue
			}

			responses = append(responses, Response{
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

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": responses,
			"errors":  errors,
		})
	}
}
