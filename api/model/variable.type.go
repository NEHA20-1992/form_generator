package model

type VariableType struct {
	ID               uint32 `gorm:"column:variable_type_id;primary_key;auto_increment" json:"variable_type_id"`
	VariableTypeName string `gorm:"size:255;not null" json:""`
}

func (VariableType) TableName() string {
	return "variable_type"
}
