package handler

func (s *Server) initRoutes() {
	s.engine.GET("/health", s.healthCheck)

	v1Group := s.engine.Group("/api/v1")

	authGroup := v1Group.Group("/auth")
	authGroup.POST("/login", s.login)
	authGroup.GET("/me", s.authMiddleware(), s.getCurrentUser)

	attendanceGroup := v1Group.Group("/attendance", s.authMiddleware())
	attendanceGroup.GET("/:employ_id", s.rbacEnforcer.Middleware())
	attendanceGroup.GET("/:employ_id/:date", s.rbacEnforcer.Middleware())
	attendanceGroup.POST("/:employ_id", s.recordAttendance)

	leaveGroup := v1Group.Group("/leave", s.authMiddleware(), s.rbacEnforcer.Middleware())
	leaveGroup.GET("/:employ_id")
	leaveGroup.POST("/:employ_id")
	leaveGroup.PUT("/:leave_id/approve")

	//s.engine.GET("/api/v1/health", s.healthCheck)
}
