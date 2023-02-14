package model

import "time"

type FormData struct {
	ID            uint32    `gorm:"column:form_data_id" json:"-"`
	FormID        uint32    `gorm:"not null" json:"-"`
	DataID        uint32    `gorm:"not null" json:"-"`
	VariableID    uint32    `gorm:"not null" json:""`
	OptionID      string    `gorm:"not null" json:"OptionID,omitempty"`
	VariableValue string    `gorm:"size:255;not null"`
	CreatedByID   uint32    `gorm:"not null" json:"omitempty"`
	CreatedAt     time.Time `gorm:"not null" json:"omitempty"`
	UpdatedByID   uint32    `gorm:"" json:"omitempty"`
	UpdatedAt     time.Time `gorm:"" json:"omitempty"`
}

func (FormData) TableName() string {
	return "form_data"
}

type DisplayFormData struct {
	FormID      uint32
	FormTitle   string
	FormDiscrip string
	Data        []Data
}
type Data struct {
	DataID    uint32
	Variables []FormData
}

type FormFetchdata struct {
	MinValue uint32
	MaxValue uint32
	Value    string
}
type Pagination struct {
	PageNumber uint32
	Size       uint32
	Sort       string
}
