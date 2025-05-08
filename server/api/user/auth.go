package user

import (
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/ent/tenant"
	"github.com/telepair/terminal/ent/user"
	"github.com/telepair/terminal/server/middleware"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the request to register a new user
type RegisterRequest struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	GivenName string `json:"given_name"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	TenantName string `json:"tenant_name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	GivenName    string `json:"given_name,omitempty"`
	TenantID     string `json:"tenant_id"`
	TenantName   string `json:"tenant_name"`
	IsSuperuser  bool   `json:"is_superuser"`
}

// Register handles user registration
func Register(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		// Parse tenant ID
		tenantID, err := uuid.Parse(req.TenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			return
		}

		// Get tenant
		t, err := client.Tenant.Query().
			Where(tenant.ID(tenantID), tenant.Enabled(true)).
			Only(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在或已被禁用"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
			return
		}

		// Check if tenant allows registration
		if !t.AllowRegistration {
			c.JSON(http.StatusForbidden, gin.H{"error": "该租户不允许用户注册"})
			return
		}

		// Check email domain if allowed domains are specified
		if len(t.AllowedEmailDomains) > 0 {
			emailParts := strings.Split(req.Email, "@")
			if len(emailParts) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邮箱格式"})
				return
			}

			domain := emailParts[1]
			isAllowed := false
			for _, allowedDomain := range t.AllowedEmailDomains {
				if domain == allowedDomain {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "邮箱域名不在允许的列表中"})
				return
			}
		}

		// Check if username already exists for this tenant
		exists, err := client.User.Query().
			Where(
				user.TenantID(tenantID),
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
				user.TenantID(tenantID),
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
			SetTenantID(tenantID).
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

// Login handles user login
func Login(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
			return
		}

		// Get tenant by name
		t, err := client.Tenant.Query().
			Where(tenant.NameEQ(req.TenantName), tenant.Enabled(true)).
			Only(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在或已被禁用"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
			return
		}

		// Get user
		u, err := client.User.Query().
			Where(
				user.TenantID(t.ID),
				user.UsernameEQ(req.Username),
				user.Enabled(true),
				user.IsLockedEQ(false),
			).Only(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(u.PasswordBcrypt), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// Update last login
		_, err = u.Update().
			SetLastLoginAt(time.Now()).
			SetLastLoginIP(c.ClientIP()).
			Save(c.Request.Context())

		if err != nil {
			// Log error but continue (non-fatal)
			slog.Info("Failed to update last login", "error", err)
		}

		// Generate JWT token
		// TODO: Replace with actual secret from config
		tokenSecret := []byte("supersecretkey-change-me-in-production")

		claims := &middleware.JWTClaims{
			UserID:      u.ID.String(),
			Tenant:      t.ID.String(),
			Username:    u.Username,
			Email:       u.Email,
			IsSuperuser: u.IsSuperuser,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "terminal",
				Subject:   u.ID.String(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(tokenSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
			return
		}

		// Generate refresh token (simplified for now)
		refreshClaims := &middleware.JWTClaims{
			UserID: u.ID.String(),
			Tenant: t.ID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "terminal",
				Subject:   u.ID.String(),
			},
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		refreshTokenString, err := refreshToken.SignedString(tokenSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成刷新令牌失败"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token:        tokenString,
			RefreshToken: refreshTokenString,
			UserID:       u.ID.String(),
			Username:     u.Username,
			Email:        u.Email,
			GivenName:    u.GivenName,
			TenantID:     t.ID.String(),
			TenantName:   t.Name,
			IsSuperuser:  u.IsSuperuser,
		})
	}
}

// Logout handles user logout
func Logout(client *ent.Client) gin.HandlerFunc {
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

		// Update last logout time
		_, err = client.User.UpdateOneID(uid).
			SetLastLogoutAt(time.Now()).
			Save(c.Request.Context())

		if err != nil {
			// Log error but continue (non-fatal)
			slog.Info("Failed to update last logout", "error", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
	}
}
