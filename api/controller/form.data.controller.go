package controller

import (
	"encoding/json"
	"errors"

	"net/http"
	"net/url"
	"strconv"

	"github.com/NEHA20-1992/form_generator/api/auth"
	"github.com/NEHA20-1992/form_generator/api/model"
	"github.com/NEHA20-1992/form_generator/api/response"
	"github.com/NEHA20-1992/form_generator/api/service"
	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

var errFormNotExsist = errors.New("Enter Valid Form ID")

type FormDataHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	GetAllExcel(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

type FormDataHandlerImpl struct {
	service     service.FormDataService
	formService service.FormService
}

func getFormDataHandlerInstances(db *gorm.DB) (Handler FormDataHandler) {

	return FormDataHandlerImpl{
		service:     service.GetFormDataService(db),
		formService: service.GetFormService(db)}
}
func (h FormDataHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idD string
	var vars = mux.Vars(req)
	idD = vars["id"]
	ID, _ := strconv.ParseUint(idD, 10, 32)
	formID := uint32(ID)
	var request []*model.FormData

	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	value, err := h.formService.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errFormNotExsist)
		return
	}
	result, err := h.service.Create(claim, formID, request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h FormDataHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idF string
	var vars = mux.Vars(req)
	idF = vars["id"]

	formFiletrData := req.FormValue("formFiletrs")
	formFiletrData, err = url.PathUnescape(formFiletrData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var formDataByte []byte = []byte(formFiletrData)
	var request = []model.FormFetchdata{}
	err = json.Unmarshal(formDataByte, &request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	paginationData := req.FormValue("pagination")
	paginationData, err = url.PathUnescape(paginationData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var pageDataByte []byte = []byte(paginationData)
	var request1 = model.Pagination{}
	err = json.Unmarshal(pageDataByte, &request1)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	ID, _ := strconv.ParseUint(idF, 10, 32)
	formID := uint32(ID)

	value, err := h.formService.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errFormNotExsist)
		return
	}
	result, err := h.service.GetAll(claim, formID, request, &request1)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h FormDataHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idF, idD string
	var vars = mux.Vars(req)
	idF = vars["id"]
	idD = vars["dataID"]
	ID, _ := strconv.ParseUint(idF, 10, 32)
	formID := uint32(ID)
	ID1, _ := strconv.ParseUint(idD, 10, 32)
	dataID := uint32(ID1)
	value, err := h.formService.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errFormNotExsist)
		return
	}
	result, err := h.service.Get(claim, formID, dataID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h FormDataHandlerImpl) Update(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idF, idD string
	var vars = mux.Vars(req)
	idF = vars["id"]
	idD = vars["dataID"]
	ID, _ := strconv.ParseUint(idF, 10, 32)
	formID := uint32(ID)
	ID1, _ := strconv.ParseUint(idD, 10, 32)
	dataID := uint32(ID1)
	var request []*model.FormData

	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	value, err := h.formService.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errFormNotExsist)
		return
	}
	result, err := h.service.Update(claim, formID, dataID, request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h FormDataHandlerImpl) Delete(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idF, idD string
	var vars = mux.Vars(req)
	idF = vars["id"]
	idD = vars["dataID"]
	ID, _ := strconv.ParseUint(idF, 10, 32)
	formID := uint32(ID)
	ID1, _ := strconv.ParseUint(idD, 10, 32)
	dataID := uint32(ID1)
	value, err := h.formService.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errFormNotExsist)
		return
	}
	result, err := h.service.Get(claim, formID, dataID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h FormDataHandlerImpl) GetAllExcel(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var idF string
	var vars = mux.Vars(req)
	idF = vars["id"]

	formFiletrData := req.FormValue("formFiletrs")
	formFiletrData, err = url.PathUnescape(formFiletrData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var formDataByte []byte = []byte(formFiletrData)
	var request = []model.FormFetchdata{}
	err = json.Unmarshal(formDataByte, &request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	paginationData := req.FormValue("pagination")
	paginationData, err = url.PathUnescape(paginationData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var pageDataByte []byte = []byte(paginationData)
	var request1 = model.Pagination{}
	err = json.Unmarshal(pageDataByte, &request1)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	ID, _ := strconv.ParseUint(idF, 10, 32)
	formID := uint32(ID)

	value, err := h.formService.Get(claim, uint32(ID))
	if err != nil {
		return
	}
	if value.ID != uint32(ID) {
		response.ERROR(w, http.StatusOK, errFormNotExsist)
		return
	}
	result, err := h.service.GetAllExcel(claim, formID, request, &request1)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSONDOWNLOAD(w, http.StatusOK, result)

}
