package seeder

import (
	"time"

	"github.com/NEHA20-1992/form_generator/api/model"
	"gorm.io/gorm"
)

const (
	shouldDropTables = true
)

func dropTables(db *gorm.DB) (err error) {
	err = db.Migrator().DropTable(&model.User{}, &model.Form{}, &model.Variable{}, &model.VariableType{}, &model.FormData{}, &model.Option{})
	return
}

func Initialize(db *gorm.DB) {
	var err error
	if shouldDropTables {
		err = dropTables(db)
		if err != nil {
			panic(err)
		}
	}

	err = db.AutoMigrate(&model.User{}, &model.Form{}, &model.VariableType{}, &model.FormData{}, &model.Variable{}, &model.Option{})
	if err != nil {
		panic(err)
	}
	VariableType := []model.VariableType{

		{

			VariableTypeName: "text",
		},
		{
			VariableTypeName: "number"},
		{
			VariableTypeName: "radio"},
		{
			VariableTypeName: "email"},
		{
			VariableTypeName: "password"},
		{
			VariableTypeName: "select"},
		{
			VariableTypeName: "checkbox"},
		{
			VariableTypeName: "date"},
	}

	db.CreateInBatches(&VariableType, 100)

	user := model.User{
		Nickname: "Sophia",

		Email:    "nehamaltiyadav@gmail.com",
		Password: "hi123",

		CreatedAt: time.Now()}
	create(db, &user)
}

func create(db *gorm.DB, user *model.User) {
	passwordText := user.Password
	user.Password = "123"

	db.Model(&user).Omit("updated_by", "updated_at").Create(&user)
	user.Password = passwordText

	err := user.Hash()
	if err != nil {
		panic(err)
	}

	err = db.Model(&user).UpdateColumn("password", user.Password).Error
	if err != nil {
		panic(err)
	}

}
