package user

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
)

type User struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(100);not null" json:"username"`
	Password  *string   `gorm:"type:varchar(255);null" json:"-"`
	Fullname  *string   `gorm:"type:varchar(255);null" json:"fullname"`
	Role      Role      `gorm:"type:enum('admin', 'user');default:'user'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeCreate memastikan default Role dan timestamp
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Role == "" {
		u.Role = USER // Default role ke "user"
	}
	u.CreatedAt = time.Now() // Set waktu sekarang untuk created_at
	u.UpdatedAt = time.Now() // Set waktu sekarang untuk updated_at
	return nil
}

// BeforeUpdate untuk memperbarui updated_at setiap update
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
