package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Employee 員工實體
type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	EmployeeID   string     `json:"employee_id" gorm:"unique;not null;type:varchar(50)"`
	FirstName    string     `json:"first_name" gorm:"not null;type:varchar(50)"`
	LastName     string     `json:"last_name" gorm:"not null;type:varchar(50)"`
	Email        string     `json:"email" gorm:"unique;not null;type:varchar(100)"`
	PasswordHash string     `json:"-" gorm:"not null;type:varchar(255)"` // 不在JSON中顯示密碼
	Phone        string     `json:"phone" gorm:"type:varchar(20)"`
	Department   string     `json:"department" gorm:"type:varchar(50)"`
	Position     string     `json:"position" gorm:"type:varchar(50)"`
	HireDate     time.Time  `json:"hire_date" gorm:"not null"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:ACTIVE"`
	ManagerID    *int       `json:"manager_id" gorm:"default:null"`
	Manager      *Employee  `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Subordinates []Employee `json:"subordinates,omitempty" gorm:"foreignKey:ManagerID"`
	Roles        []Role     `json:"roles" gorm:"many2many:employee_roles"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Role 角色實體
type Role struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"unique;not null;type:varchar(50)"`
	Description string     `json:"description" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Employees   []Employee `json:"-" gorm:"many2many:employee_roles"`
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employees"
}

func (Role) TableName() string {
	return "roles"
}

// LoginRequest 登錄請求結構
type LoginRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// LoginResponse 登錄響應結構
type LoginResponse struct {
	Token        string   `json:"token"`
	EmployeeID   string   `json:"employee_id"`
	Name         string   `json:"name"`
	Roles        []string `json:"roles"`
	RefreshToken string   `json:"refresh_token,omitempty"`
}

// CreateEmployeeRequest 創建員工請求結構
type CreateEmployeeRequest struct {
	EmployeeID string   `json:"employee_id" binding:"required"`
	FirstName  string   `json:"first_name" binding:"required"`
	LastName   string   `json:"last_name" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Password   string   `json:"password" binding:"required,min=6"`
	Phone      string   `json:"phone"`
	Department string   `json:"department" binding:"required"`
	Position   string   `json:"position" binding:"required"`
	HireDate   string   `json:"hire_date" binding:"required"` // 前端傳入的日期字符串
	ManagerID  *int     `json:"manager_id"`
	RoleNames  []string `json:"role_names" binding:"required"`
}

// UpdateEmployeeRequest 更新員工請求結構
type UpdateEmployeeRequest struct {
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email" binding:"omitempty,email"`
	Phone      string   `json:"phone"`
	Department string   `json:"department"`
	Position   string   `json:"position"`
	Status     string   `json:"status"`
	ManagerID  *int     `json:"manager_id"`
	RoleNames  []string `json:"role_names"`
}

// EmployeeResponse 員工信息響應結構
type EmployeeResponse struct {
	ID         int      `json:"id"`
	EmployeeID string   `json:"employee_id"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Department string   `json:"department"`
	Position   string   `json:"position"`
	HireDate   string   `json:"hire_date"`
	Status     string   `json:"status"`
	ManagerID  *int     `json:"manager_id,omitempty"`
	Roles      []string `json:"roles"`
}

// ToResponse 將 Employee 轉換為響應結構
func (e *Employee) ToResponse() *EmployeeResponse {
	roleNames := make([]string, len(e.Roles))
	for i, role := range e.Roles {
		roleNames[i] = role.Name
	}

	return &EmployeeResponse{
		ID:         e.ID,
		EmployeeID: e.EmployeeID,
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		Email:      e.Email,
		Phone:      e.Phone,
		Department: e.Department,
		Position:   e.Position,
		HireDate:   e.HireDate.Format("2006-01-02"),
		Status:     e.Status,
		ManagerID:  e.ManagerID,
		Roles:      roleNames,
	}
}

// BeforeCreate GORM 鉤子：創建前的處理
func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if e.Status == "" {
		e.Status = "ACTIVE"
	}
	return nil
}

// ValidatePassword 驗證密碼
func (e *Employee) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(password))
	return err == nil
}

// SetPassword 設置密碼
func (e *Employee) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.PasswordHash = string(hashedPassword)
	return nil
}
