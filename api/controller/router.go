package controller

import (
	//  "github.com/NEHA20-1992/form_generator/api/middleware"
	"net/http"

	"github.com/NEHA20-1992/form_generator/api/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitializeRouters(db *gorm.DB, router *mux.Router) {
	registerAuthenticationHandlers(db, router)
	registerUserHandlers(db, router)
	registerFormHandlers(db, router)
	registerFormDataHandlers(db, router)
}
func registerUserHandlers(db *gorm.DB, router *mux.Router) {
	router.Use(middleware.Cors)
	router.HandleFunc("/singUP", middleware.JSONResponder(GetUserHandlerInstance(db).Create)).Methods(http.MethodPost)

}
func registerAuthenticationHandlers(db *gorm.DB, router *mux.Router) {
	router.Use(middleware.Cors)

	router.HandleFunc("/auth/login",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).GenerateToken)).
		Methods(http.MethodPost)

	router.HandleFunc("/auth/getResetCode/{userEmail}",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).GetResetCode)).
		Methods(http.MethodGet)

	router.HandleFunc("/auth/resetPassword",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).ResetPassword)).
		Methods(http.MethodPost)

	router.HandleFunc("/auth/refresh",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).RefreshToken)).
		Methods(http.MethodPost)
	router.HandleFunc("/auth/user",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetAuthenticationHandlerInstance(db).GetCurrentUser))).
		Methods(http.MethodGet)

	router.HandleFunc("/auth/changePassword",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserHandlerInstance(db).ChangePassword))).
		Methods(http.MethodPost)
}

func registerFormHandlers(db *gorm.DB, router *mux.Router) {

	router.HandleFunc("/createForm", middleware.JSONResponder(middleware.Authenticate(getFormHandlerInstances(db).Create))).Methods(http.MethodPost)
	router.HandleFunc("/viewForm/{id}", middleware.JSONResponder(middleware.Authenticate(getFormHandlerInstances(db).Get))).Methods(http.MethodGet)
	router.HandleFunc("/viewAllForms", middleware.JSONResponder(middleware.Authenticate(getFormHandlerInstances(db).GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/updateForm/{id}", middleware.JSONResponder(middleware.Authenticate(getFormHandlerInstances(db).Update))).Methods(http.MethodPut)
	router.HandleFunc("/deleteForm/{id}", middleware.JSONResponder(middleware.Authenticate(getFormHandlerInstances(db).Delete))).Methods(http.MethodDelete)
}
func registerFormDataHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/insertFormData/{id}", middleware.JSONResponder(middleware.Authenticate(getFormDataHandlerInstances(db).Create))).Methods(http.MethodPost)
	router.HandleFunc("/viewFormData/{id}", middleware.JSONResponder(middleware.Authenticate(getFormDataHandlerInstances(db).GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/viewFormData/{id}/{dataID}", middleware.JSONResponder(middleware.Authenticate(getFormDataHandlerInstances(db).Get))).Methods(http.MethodGet)
	router.HandleFunc("/updateFormData/{id}/{dataID}", middleware.JSONResponder(middleware.Authenticate(getFormDataHandlerInstances(db).Update))).Methods(http.MethodPut)
	router.HandleFunc("/deleteFormData/{id}", middleware.JSONResponder(middleware.Authenticate(getFormDataHandlerInstances(db).Delete))).Methods(http.MethodDelete)
	router.HandleFunc("/downloadFormData/{id}", middleware.JSONResponder(middleware.Authenticate(getFormDataHandlerInstances(db).GetAllExcel))).Methods(http.MethodGet)
}
