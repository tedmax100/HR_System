package main

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	//"github.com/casbin/casbin/v2/persist"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 初始化 GORM Adapter
	adapter, err := gormadapter.NewAdapter("postgres", "user=your_user password=your_password dbname=your_db sslmode=disable")
	if err != nil {
		panic(err)
	}

	// 加載 Casbin 模型與策略
	enforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		panic(err)
	}

	// 同步策略到數據庫
	err = enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}

	// 初始化 Gin
	r := gin.Default()

	// 中間件檢查權限
	r.Use(func(c *gin.Context) {
		sub := c.GetHeader("user-role") // 獲取用戶角色
		obj := c.Request.URL.Path       // 獲取資源
		act := c.Request.Method         // 獲取行為

		// 檢查權限
		allowed, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	})

	// 定義 API 路由
	r.GET("/api/v1/employees", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Employee List"})
	})

	// 啟動服務
	r.Run(":8080")
}
