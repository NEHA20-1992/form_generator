package model

import "time"

type User struct {
	ID        uint32    `gorm:"column:user_id,primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"-"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null;" json:"password"`
	ResetCode string    `gorm:"-" json:"resetCode,omitempty"`
	CreatedAt time.Time `gorm:"" json:"-"`
	UpdatedAt time.Time `gorm:"" json:"-"`
}
type ChangePasswordRequest struct {
	Password        string `gorm:"size:255;null" json:"password,omitempty"`
	NewPassword     string `gorm:"size:255;null" json:"newPassword,omitempty"`
	ConfirmPassword string `gorm:"size:255;null" json:"confirmPassword,omitempty"`
}

type ForgetPasswordRequest struct {
	Email string `gorm:"size:255;null" json:"email,omitempty"`
}

func (User) TableName() string {
	return "user"
}
