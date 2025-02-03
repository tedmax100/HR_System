package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	EmployeeID string   `json:"employee_id"`
	Roles      []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(employeeID string, roles []string, secret string, expiration time.Duration) (string, error) {
	claims := Claims{
		EmployeeID: employeeID,
		Roles:      roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
