package api

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/hrdemo/internal/config"
)

type RbacEnforcer struct {
	Enforcer *casbin.Enforcer
}

func NewRbacEnforcer(cfg *config.Config) (*RbacEnforcer, error) {
	adapter, err := gormadapter.NewAdapter("postgres", cfg.DataBaseURL)
	if err != nil {
		return nil, err
	}

	// load Casbin model and policy
	enforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// sync policy
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return &RbacEnforcer{
		Enforcer: enforcer,
	}, nil
}

func (r *RbacEnforcer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetHeader("user-role")
		obj := c.Request.URL.Path
		act := c.Request.Method

		allowed, err := r.Enforcer.Enforce(sub, obj, act)
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
	}
}
