package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role int

const (
	RoleIntern Role = iota + 1
	RoleManager
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleIntern:
		return "intern"
	case RoleManager:
		return "manager"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

func ParseRole(s string) Role {
	switch s {
	case "intern":
		return RoleIntern
	case "manager":
		return RoleManager
	case "admin":
		return RoleAdmin
	default:
		return RoleIntern
	}
}

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex" json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `gorm:"not null" json:"-"`
	Role     Role   `json:"role"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
