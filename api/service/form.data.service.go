package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/NEHA20-1992/form_generator/api/auth"
	"github.com/NEHA20-1992/form_generator/api/model"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var DataID uint32 = 0

type FormDataService interface {
	GetAllExcel(claim *auth.AuthenticatedClaim, formID uint32, formFetchData []model.FormFetchdata, page *model.Pagination) (*excelize.File, error)
	Create(claim *auth.AuthenticatedClaim, formID uint32, formData []*model.FormData) (string, error)
	GetAll(claim *auth.AuthenticatedClaim, formID uint32, formFetchData []model.FormFetchdata, page *model.Pagination) (*model.DisplayFormData, error)
	Get(claim *auth.AuthenticatedClaim, formID uint32, dataID uint32) (model.DisplayFormData, error)
	Update(claim *auth.AuthenticatedClaim, formID uint32, dataID uint32, formData []*model.FormData) (string, error)
	Delete(claim *auth.AuthenticatedClaim, formID uint32, dataID uint32) (string, error)
}

type FormDataServiceImpl struct {
	db          *gorm.DB
	formService FormService
}

func GetFormDataService(db *gorm.DB) FormDataService {
	return &FormDataServiceImpl{
		db:          db,
		formService: GetFormService(db),
	}
}
func (m *FormDataServiceImpl) Create(claim *auth.AuthenticatedClaim, formID uint32, formData []*model.FormData) (result string, err error) {
	var data []model.FormData
	if formData == nil {
		return
	}
	value, err := m.formService.Get(claim, formID)
	if err != nil {
		return
	}
	if len(value.Variables) != len(formData) {
		err = errors.New("form structure is incorrect")
		return
	}
	DataID += 1
	for indx, varList := range value.Variables {
		// formData[indx].SNo = varList.SerialNumber
		formData[indx].VariableID = varList.ID
		formData[indx].CreatedByID = claim.UserId
		formData[indx].CreatedAt = time.Now()
		formData[indx].FormID = value.ID
		formData[indx].DataID = DataID
		// formData[indx].VariableName = varList.VariableName
		if varList.VariableTypeID == 6 || varList.VariableTypeID == 3 {
			err = m.db.Model(&model.Option{}).Select("option_value").Where("variable_id = ? AND option_id =?", varList.ID, formData[indx].OptionID).Find(&formData[indx].VariableValue).Error
			if err != nil {
				return
			}
			data = append(data, *formData[indx])

		} else if varList.VariableTypeID == 7 {

			var str string = ""
			optionList := strings.Split(formData[indx].OptionID, ",")

			for _, op := range optionList {
				err = m.db.Model(&model.Option{}).Select("option_value").Where("variable_id = ? AND option_id=?", varList.ID, op).Find(&str).Error
				if err != nil {
					return
				}

				formData[indx].VariableValue = str
				formData[indx].OptionID = op
				data = append(data, *formData[indx])

			}
		} else {
			data = append(data, *formData[indx])

		}

	}

	err = m.db.Model(&model.FormData{}).Omit("updated_by_id", "updated_at").CreateInBatches(data, len(data)).Error
	if err != nil {
		return
	} else {
		result = "Data saved succesfully!"
	}

	return
}
func (m FormDataServiceImpl) GetAll(claim *auth.AuthenticatedClaim, formID uint32, formFetchData []model.FormFetchdata, page *model.Pagination) (result *model.DisplayFormData, err error) {

	value, err := m.formService.Get(claim, formID)
	if err != nil {
		return
	}
	var query string = ""
	for indx, Data := range value.Variables {
		var varTypeID, varID uint32
		varID = Data.ID
		varTypeID = Data.VariableTypeID

		if varTypeID == 2 && formFetchData[indx].Value != "" {
			query = query + " and  variable_id = " + strconv.Itoa(int(varID)) + " and variable_value between ( " + strconv.Itoa(int(formFetchData[indx].MinValue)) + " , " + strconv.Itoa(int(formFetchData[indx].MaxValue)) + ")"
		} else if formFetchData[indx].Value != "" {
			query = query + " and  variable_id = " + strconv.Itoa(int(varID)) + " and variable_value = '" + formFetchData[indx].Value + "'"
		}

	}
	query1 := "form_id = " + strconv.Itoa(int(formID)) + query
	var List []model.FormData
	var List1 model.DisplayFormData
	List1.FormTitle = value.FormTitle
	List1.FormDiscrip = value.FormDescription
	var Data model.Data
	var DataList []model.Data
	var dataID []uint32
	err = m.db.Model(&model.FormData{}).Select("DISTINCT data_id").Where(query1).Find(&dataID).Limit(int(page.Size)).Offset(int((page.Size) * (page.PageNumber - 1))).Error
	for _, dataId := range dataID {
		err = m.db.Model(&model.FormData{}).Select("*").Where("data_id=?", dataId).Find(&List).Order(page.Sort).Error
		if err != nil {
			return
		}

		Data.DataID = dataId
		Data.Variables = List
		List1.FormID = formID
		List1.FormTitle = value.FormTitle
		List1.FormDiscrip = value.FormDescription
		DataList = append(DataList, Data)

	}

	List1.FormID = formID
	List1.Data = DataList
	result = &List1
	return
}
func (m FormDataServiceImpl) Get(claim *auth.AuthenticatedClaim, formID uint32, dataID uint32) (result model.DisplayFormData, err error) {
	value, err := m.formService.Get(claim, formID)
	if err != nil {
		return
	}
	var List []model.FormData
	var List1 model.DisplayFormData
	var Data model.Data
	var DataList []model.Data
	err = m.db.Model(&model.FormData{}).Where("form_id =? and data_id =?", formID, dataID).Find(&List).Error
	if err != nil {
		return
	}

	Data.DataID = dataID
	Data.Variables = List
	List1.FormID = formID
	List1.FormTitle = value.FormTitle
	List1.FormDiscrip = value.FormDescription
	DataList = append(DataList, Data)
	List1.Data = DataList
	result = List1
	return
}
func (m FormDataServiceImpl) Delete(claim *auth.AuthenticatedClaim, formID uint32, dataID uint32) (result string, err error) {
	err = m.db.Debug().
		Delete(&model.FormData{}, &model.FormData{FormID: formID,
			DataID: dataID}).
		Error
	if err != nil {
		return
	}
	if err == nil {
		result = "Deleted Succesfullly"
	}
	return
}
func (m FormDataServiceImpl) Update(claim *auth.AuthenticatedClaim, formID uint32, dataID uint32, formData []*model.FormData) (result string, err error) {
	var result1, data []*model.FormData
	err = m.db.Model(&model.FormData{}).Where("form_id =? AND data_id =?", formID, dataID).Find(&result1).Error
	if err != nil {
		return
	}
	if result1 == nil {
		err = errors.New("DATA NOT EXSIST")
	}

	value, err := m.formService.Get(claim, formID)
	if err != nil {
		return
	}

	for i, formdata1 := range value.Variables {
		formData[i].ID = result1[i].ID
		if formdata1.VariableTypeID == 6 || formdata1.VariableTypeID == 3 {
			err = m.db.Model(&model.Option{}).Select("option_value").Where("variable_id = ? AND option_id =?", formdata1.ID, formData[i].OptionID).Find(&formData[i].VariableValue).Error
			data = append(data, formData[i])
		} else if formdata1.VariableTypeID == 7 {
			var str string = ""
			optionList := strings.Split(formData[i].OptionID, ",")

			for _, op := range optionList {
				err = m.db.Model(&model.Option{}).Select("option_value").Where("variable_id = ? AND option_id = ?", formdata1.ID, op).Find(&str).Error
				formData[i].VariableValue = str
				formData[i].OptionID = op
				data = append(data, formData[i])
			}
		} else {
			data = append(data, formData[i])
		}
	}

	for i, data1 := range data {
		var varTypeID uint32
		err = m.db.Model(&value.Variables).Select("variable_type_id").Where("variable_id =?", data1.VariableID).Find(&varTypeID).Error
		if varTypeID == 7 && data1.OptionID == result1[i].OptionID {
			err = m.db.Model(&model.FormData{}).Create(data1).Error
			if err != nil {
				return
			}
		} else {

			m.db.Debug().Model(&model.Form{}).Where("form_id =? AND data_id =?", formID, dataID).Take(&model.FormData{}).UpdateColumns(
				map[string]interface{}{
					"option_id":      data1.OptionID,
					"variable_value": data1.VariableValue,
					"updated_by_id":  claim.UserId,
					"updated_at":     time.Now(),
				})

		}
	}

	if err == nil {
		result = "Updated Successfully!"
	}

	return
}
func (m FormDataServiceImpl) GetAllExcel(claim *auth.AuthenticatedClaim, formID uint32, formFetchData []model.FormFetchdata, page *model.Pagination) (result *excelize.File, err error) {

	value, err := m.formService.Get(claim, formID)
	if err != nil {
		return
	}
	var query string = ""
	for indx, Data := range value.Variables {
		var varTypeID, varID uint32
		varID = Data.ID
		varTypeID = Data.VariableTypeID

		if varTypeID == 2 && formFetchData[indx].Value != "" {
			query = query + " and  variable_id = " + strconv.Itoa(int(varID)) + " and variable_value between ( " + strconv.Itoa(int(formFetchData[indx].MinValue)) + " , " + strconv.Itoa(int(formFetchData[indx].MaxValue)) + ")"
		} else if formFetchData[indx].Value != "" {
			query = query + " and  variable_id = " + strconv.Itoa(int(varID)) + " and variable_value = '" + formFetchData[indx].Value + "'"
		}

	}
	file := excelize.NewFile()
	file.NewConditionalStyle(&excelize.Style{

		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
			WrapText:   true,
		},
	})

	file.MergeCell("Sheet1", "A1", string('A'+len(value.Variables))+strconv.Itoa(1))
	file.MergeCell("Sheet1", "A2", string('A'+len(value.Variables))+strconv.Itoa(2))

	file.SetCellValue("Sheet1", "A1", value.FormTitle)
	file.SetCellValue("Sheet1", "A2", value.FormDescription)

	file.SetCellValue("Sheet1", "A3", "Data ID")

	for i, varList := range value.Variables {

		file.SetCellValue("Sheet1", string('B'+i)+strconv.Itoa(3), varList.VariableName)
	}
	query1 := "form_id = " + strconv.Itoa(int(formID)) + query

	var dataID []uint32
	err = m.db.Model(&model.FormData{}).Select("DISTINCT data_id").Where(query1).Find(&dataID).Limit(int(page.Size)).Offset(int((page.Size) * (page.PageNumber - 1))).Error
	for indx, dataId := range dataID {

		file.SetCellValue("Sheet1", string('A')+strconv.Itoa(indx+4), dataId)

		for i, listData := range value.Variables {
			if listData.VariableTypeID == 7 {

				var opValue []string
				err = m.db.Model(&model.FormData{}).Select("variable_value").Where("variable_id =? and data_id =?", listData.ID, dataId).Find(&opValue).Error
				file.SetCellValue("Sheet1", string('B'+i)+strconv.Itoa(indx+4), opValue)
			} else {

				var opValue string
				err = m.db.Model(&model.FormData{}).Select("variable_value").Where("variable_id =? and data_id = ?", listData.ID, dataId).Find(&opValue).Error
				file.SetCellValue("Sheet1", string('B'+i)+strconv.Itoa(indx+4), opValue)
			}
		}

	}

	result = file
	return
}
