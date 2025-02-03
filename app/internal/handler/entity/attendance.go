package entity

import "time"

type AttendanceEvent struct {
	ID         int64     `json:"id"`
	EmployeeID string    `json:"employee_id"`
	EventType  string    `json:"event_type"`
	EventTime  time.Time `json:"event_time"`
	Location   string    `json:"location"`
	DeviceInfo string    `json:"device_info"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

type DailyAttendance struct {
	EmployeeID      string     `json:"employee_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	EventDate       time.Time  `json:"event_date"`
	CheckInTime     *time.Time `json:"check_in_time"`
	CheckOutTime    *time.Time `json:"check_out_time"`
	DailyCheckCount int        `json:"daily_check_count"`
	Status          string     `json:"status"`
}

// 打卡請求結構
type CheckRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Location   string `json:"location"`
	DeviceInfo string `json:"device_info"`
	Note       string `json:"note"`
}
