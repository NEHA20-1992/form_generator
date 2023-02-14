package controller

import (
	"encoding/json"
	"errors"

	"net/http"
	"strconv"

	"github.com/NEHA20-1992/form_generator/api/auth"
	"github.com/NEHA20-1992/form_generator/api/model"
	"github.com/NEHA20-1992/form_generator/api/response"
	"github.com/NEHA20-1992/form_generator/api/service"
	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

var errNotExsist = errors.New("DATA NOT EXSIST")

type FormHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

type FormHandlerImpl struct {
	service service.FormService
}

func getFormHandlerInstances(db *gorm.DB) (Handler FormHandler) {

	return FormHandlerImpl{
		service: service.GetFormService(db)}
}
func (h FormHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var request = model.Form{}

	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Create(claim, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h FormHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
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

func (h FormHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idD string
	var vars = mux.Vars(req)
	idD = vars["id"]
	ID, _ := strconv.ParseUint(idD, 10, 32)
	ID1 := uint32(ID)
	value, err := h.service.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errNotExsist)
		return
	}
	result, err := h.service.Get(claim, ID1)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h FormHandlerImpl) Update(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var vars = mux.Vars(req)
	idD := vars["id"]
	ID, _ := strconv.ParseUint(idD, 10, 32)

	var request = model.Form{}

	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	value, err := h.service.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errNotExsist)
		return
	}

	result, err := h.service.Update(claim, uint32(ID), &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h FormHandlerImpl) Delete(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var vars = mux.Vars(req)
	idD := vars["id"]
	ID, _ := strconv.ParseUint(idD, 10, 32)
	value, err := h.service.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errNotExsist)
		return
	}
	result, err := h.service.Delete(claim, uint32(ID))
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
