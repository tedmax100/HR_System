package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hrdemo/internal/handler/entity"
	"github.com/hrdemo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) login(c *gin.Context) {
	var loginReq entity.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var employee entity.Employee
	if err := s.db.Preload("Roles").Where("employee_id = ?", loginReq.EmployeeID).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 比較密碼
	if err := bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 獲取角色列表
	roles := make([]string, len(employee.Roles))
	for i, role := range employee.Roles {
		roles[i] = role.Name
	}

	// 生成 JWT token
	expiration := time.Duration(s.cfg.JWTExpiration) * time.Minute
	token, err := utils.GenerateToken(employee.EmployeeID, roles, s.cfg.JWTSecret, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 返回登錄響應
	c.JSON(http.StatusOK, entity.LoginResponse{
		Token:      token,
		EmployeeID: employee.EmployeeID,
		Name:       employee.FirstName + " " + employee.LastName,
		Roles:      roles,
	})
}

// 添加身份驗證中間件
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 解析 Bearer token
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]

		// 驗證 token
		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("employee_id", claims.EmployeeID)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// 獲取當前用戶信息
func (s *Server) getCurrentUser(c *gin.Context) {
	employeeID := c.GetString("employee_id")

	var employee entity.Employee
	if err := s.db.Preload("Roles").Where("employee_id = ?", employeeID).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee.ToResponse())
}
