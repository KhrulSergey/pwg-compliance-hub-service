package model

import (
	"database/sql/driver"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type RoleAllowed string

const (
	admin   RoleAllowed = "admin"
	seller  RoleAllowed = "seller"
	rider   RoleAllowed = "rider"
	regular RoleAllowed = "regular"
)

func (st *RoleAllowed) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		*st = RoleAllowed(b)
	}
	return nil
}

func (st RoleAllowed) Value() (driver.Value, error) {
	return string(st), nil
}

type StatusAllowed string

const (
	suspended StatusAllowed = "suspended"
	active    StatusAllowed = "active"
	inactive  StatusAllowed = "inactive"
)

func (st *StatusAllowed) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		*st = StatusAllowed(b)
	}
	return nil
}

func (st StatusAllowed) Value() (driver.Value, error) {
	return string(st), nil
}

type User struct {
	ID              uint           `json:"id" gorm:"type:bigserial;primaryKey;autoIncrement"`
	Email           string         `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password        string         `json:"password" gorm:"type:varchar(255);unique;not null""`
	OTPCode         int            `json:"otp_code"`
	Role            RoleAllowed    `json:"role" gorm:"type:role_allowed;default:'regular'"`
	Status          StatusAllowed  `json:"status" gorm:"type:status_allowed;default:'active'"`
	IsEmailVerified bool           `json:"isEmailVerified" gorm:"default:false"`
	OTPExpireTime   time.Time      `json:"otp_expire_time"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPasswordHash(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
