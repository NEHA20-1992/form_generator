package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NEHA20-1992/form_generator/api/auth"
	"github.com/NEHA20-1992/form_generator/api/model"
	"github.com/NEHA20-1992/form_generator/api/response"
	"github.com/NEHA20-1992/form_generator/api/service"
	"github.com/NEHA20-1992/form_generator/api/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	ChangePassword(w http.ResponseWriter, req *http.Request)
}

type UserHandlerImpl struct {
	service service.UserService
}

var ErrNotAdmin = errors.New("not an administrator")

func GetUserHandlerInstance(db *gorm.DB) (handler UserHandler) {
	return UserHandlerImpl{
		service: service.GetUserService(db)}
}

func (h UserHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
	var err error

	var request = model.User{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = validator.ValidateUser(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Create(&request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h UserHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	result, err := h.service.GetAll(claim)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h UserHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var emailAddress string
	var vars = mux.Vars(req)

	emailAddress = vars["emailAddress"]

	result, err := h.service.Get(claim, emailAddress)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	result.Password = ""

	response.JSON(w, http.StatusOK, result)

}

func (h UserHandlerImpl) Update(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var emailAddress string
	var vars = mux.Vars(req)

	emailAddress = vars["emailAddress"]

	var request = model.User{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//request.Email = emailAddress

	result, err := h.service.Update(claim, emailAddress, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h UserHandlerImpl) ChangePassword(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var request = model.ChangePasswordRequest{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if request.NewPassword != request.ConfirmPassword {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationNewAndConfirmPasswordMismatch)
		return
	}

	result, err := h.service.ChangePassword(claim, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
