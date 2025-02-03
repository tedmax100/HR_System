package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hrdemo/internal/handler/entity"
)

func (s Server) healthCheck(c *gin.Context) {
	c.JSON(200, entity.Status{
		Status:    "ok",
		Version:   s.Version,
		Name:      s.Name,
		GoVersion: s.GoVersion,
		Build:     s.Build,
	})
}
