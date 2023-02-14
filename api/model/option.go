package model

type Option struct {
	ID          uint32 `gorm:"column:option_id;primary_key;auto_increment" json:"-"`
	VariableID  uint32 `gorm:"not null" json:"-"`
	OptionValue string `gorm:"size:255;not null" json:"OptionValue"`
}

func (Option) TableName() string {
	return "option"
}
