package model

import (
	"time"
)

type Form struct {
	ID              uint32     `gorm:"column:form_id;primary_key;auto_increment" json:"form_id"`
	FormTitle       string     `gorm:"size:255;not null" json:"FormTitle"`
	FormDescription string     `gorm:"size:255;" json:"FormDescription"`
	Variables       []Variable `gorm:"-" json:"variables,omitempty"`
	CreatedByID     uint32     `gorm:"size:255;not null" json:"-,omitempty"`
	CreatedAt       time.Time  `gorm:";not null" json:"-, omitempty"`
	UpdatedByID     uint32     ` json:"-,omitempty"`
	UpdatedAt       time.Time  `json:"-,omitempty" `
}

func (Form) TableName() string {
	return "form"
}
