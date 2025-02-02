package api

func (s *Server) initRoutes() {
	apiGroup := s.engine.Group("/api")
	// apiGroup.GET("/health", s.healthCheck)

	v1Group := apiGroup.Group("/v1")

	authGroup := v1Group.Group("/auth")
	authGroup.POST("/login")
	authGroup.GET("/authme")

	attendanceGroup := v1Group.Group("/attendance")
	attendanceGroup.GET("/:employ_id")
	attendanceGroup.GET("/:employ_id/:date")
	attendanceGroup.POST("/:employ_id")

	leaveGroup := v1Group.Group("/leave")
	leaveGroup.GET("/:employ_id")
	leaveGroup.POST("/:employ_id")
	leaveGroup.PUT("/:leave_id/approve")

	//s.engine.GET("/api/v1/health", s.healthCheck)
}
