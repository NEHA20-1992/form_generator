package api

import (
	"github.com/NEHA20-1992/form_generator/api/controller"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Initialize(DB *gorm.DB, router *mux.Router) {
	controller.InitializeRouters(DB, router)
}
