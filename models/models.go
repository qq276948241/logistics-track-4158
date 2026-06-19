package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type Waybill struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TrackingNumber string         `gorm:"size:50;uniqueIndex;not null" json:"tracking_number"`
	Carrier        string         `gorm:"size:50;not null" json:"carrier"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Trackings      []Tracking     `json:"trackings,omitempty"`
}

type Tracking struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	WaybillID      uint           `gorm:"index;not null" json:"waybill_id"`
	Status         string         `gorm:"size:50;not null" json:"status"`
	Location       string         `gorm:"size:255;not null" json:"location"`
	Description    string         `gorm:"size:500" json:"description"`
	TrackingTime   time.Time      `gorm:"not null" json:"tracking_time"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type WaybillRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required"`
	Carrier        string `json:"carrier" binding:"required"`
}

type WaybillResponse struct {
	ID             uint       `json:"id"`
	TrackingNumber string     `json:"tracking_number"`
	Carrier        string     `json:"carrier"`
	CreatedAt      time.Time  `json:"created_at"`
	Trackings      []Tracking `json:"trackings"`
}
