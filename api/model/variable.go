package model

type Variable struct {
	ID             uint32 `gorm:"column:variable_id;primary_key;auto_increment" json:"-"`
	FormID         uint32 `gorm:"not null" json:"-"`
	VariableTypeID uint32 `gorm:"size:255;not null" json:"VariableTypeID"`
	SerialNumber   uint32 `gorm:"not null" json:"SerialNumber"`
	// VariableType   string   `gorm:"size:255;not null" json:"VariableType"`
	VariableLabel string   `gorm:"size:255;not null" json:"VariableLabel"`
	VariableName  string   `gorm:"size:255;not null" json:"VariableName"`
	Options       []Option `gorm:"-" json:"Options,omitempty"`
}

func (Variable) TableName() string {
	return "variable"
}
