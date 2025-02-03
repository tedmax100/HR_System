package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrdemo/internal/handler/entity"
)

// 記錄打卡事件
func (s *Server) recordAttendance(c *gin.Context) {
	var req entity.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EmployeeID != c.GetString("employee_id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	now := time.Now()
	// 獲取今天的開始和結束時間（使用當地時區）
	loc, _ := time.LoadLocation("Asia/Taipei")
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	tomorrow := today.Add(24 * time.Hour)

	// 獲取今天的打卡記錄
	var events []entity.AttendanceEvent
	if err := s.db.Model(&entity.AttendanceEvent{}).
		Where("employee_id = ? AND event_time >= ? AND event_time < ?",
			req.EmployeeID, today, tomorrow).
		Order("event_time ASC").
		Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 檢查打卡次數
	if len(events) >= 2 {
		// 如果已經打卡兩次，返回最後一次的下班打卡記錄
		c.JSON(http.StatusOK, gin.H{
			"message": "You have already checked out for today",
			"event":   events[len(events)-1],
		})
		return
	}

	// 決定事件類型
	eventType := "CHECK_IN"
	if len(events) > 0 {
		eventType = "CHECK_OUT"
	}

	// 記錄打卡事件
	event := entity.AttendanceEvent{
		EmployeeID: c.GetString("employee_id"),
		EventType:  eventType,
		EventTime:  now,
		Location:   req.Location,
		DeviceInfo: req.DeviceInfo,
		Note:       req.Note,
	}

	if err := s.db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s recorded successfully", eventType),
		"event":   event,
	})
}
