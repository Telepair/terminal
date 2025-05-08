// Package middleware provides HTTP middleware functions for authentication and authorization
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/telepair/terminal/ent"
	"github.com/telepair/terminal/ent/user"
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID      string `json:"user_id"`
	Tenant      string `json:"tenant"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	IsSuperuser bool   `json:"is_superuser"`
	jwt.RegisteredClaims
}

// JWTAuth middleware for JWT authentication
func JWTAuth(client *ent.Client, requireSuperUser bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Authorization头"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization格式无效"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
			}

			// TODO: Replace with actual secret from config
			return []byte("supersecretkey-change-me-in-production"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌验证失败"})
			c.Abort()
			return
		}

		// Check if user exists and is not locked
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
			c.Abort()
			return
		}

		u, err := client.User.Query().
			Where(user.ID(userID), user.Enabled(true), user.IsLockedEQ(false)).
			Only(c.Request.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在或已被禁用"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
			c.Abort()
			return
		}

		// Check if superuser is required
		if requireSuperUser && !u.IsSuperuser {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("tenant", claims.Tenant)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("is_superuser", claims.IsSuperuser)

		c.Next()
	}
}

// TenantAdminOnly middleware to check if user is tenant admin
func TenantAdminOnly(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user info from context
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "租户信息缺失"})
			c.Abort()
			return
		}

		// Get tenant ID from URL param
		urlTenantID := c.Param("tenant_id")

		// Ensure user is operating on their own tenant
		if tenantID != urlTenantID {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问其他租户"})
			c.Abort()
			return
		}

		// Check if user is tenant admin by checking if their email matches tenant admin_email
		userEmail, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息缺失"})
			c.Abort()
			return
		}

		// Parse tenant ID
		tid, err := uuid.Parse(urlTenantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的租户ID"})
			c.Abort()
			return
		}

		// Get tenant
		tenant, err := client.Tenant.Get(c.Request.Context(), tid)
		if err != nil {
			if ent.IsNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "租户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
			}
			c.Abort()
			return
		}

		// Check if user email matches tenant admin email
		if tenant.AdminEmail != userEmail {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要租户管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
