package service

import (
	"errors"
	"time"

	"github.com/NEHA20-1992/form_generator/api/auth"
	"github.com/NEHA20-1992/form_generator/api/model"
	"gorm.io/gorm"
)

type FormService interface {
	Create(claim *auth.AuthenticatedClaim, form *model.Form) (string, error)
	GetAll(claim *auth.AuthenticatedClaim) ([]model.Form, error)
	Get(claim *auth.AuthenticatedClaim, id uint32) (*model.Form, error)
	Update(claim *auth.AuthenticatedClaim, id uint32, form *model.Form) (string, error)
	Delete(claim *auth.AuthenticatedClaim, id uint32) (string, error)
}

type FormServiceImpl struct {
	db *gorm.DB
}

func GetFormService(db *gorm.DB) FormService {
	return &FormServiceImpl{
		db: db,
	}
}
func (m *FormServiceImpl) Create(claim *auth.AuthenticatedClaim, form *model.Form) (result string, err error) {
	if form == nil {
		return
	}

	form.CreatedByID = claim.UserId
	form.CreatedAt = time.Now()
	tx := m.db.Begin()
	err = tx.Model(&form).Omit("updated_by_id", "updated_at").Create(&form).Error
	if err != nil {
		return
	}

	var variables []model.Variable = make([]model.Variable, len(form.Variables))
	for indx, vi := range form.Variables {
		var newRecord *model.Variable = &vi
		newRecord.FormID = form.ID
		if form.Variables[indx].VariableTypeID == 2 || form.Variables[indx].VariableTypeID == 3 || form.Variables[indx].VariableTypeID == 6 || form.Variables[indx].VariableTypeID == 7 {
			if form.Variables[indx].Options == nil {
				err = errors.New("please enter option value!")
				return
			}
		} else {
			if form.Variables[indx].Options != nil {
				err = errors.New("option value shoulde be null!")
				return
			}

		}
		// err = m.db.Model(&model.VariableType{}).Select("variable_type_name").Where("variable_type_id = ?", form.Variables[indx].VariableTypeID).Find(&newRecord.VariableType).Error
		variables[indx] = *newRecord
		err = tx.Model(&model.Variable{}).Create(newRecord).Error
		if err != nil {
			return
		}
		var option []model.Option = make([]model.Option, len(variables[indx].Options))
		for in, op := range variables[indx].Options {
			var newOption *model.Option = &op
			newOption.VariableID = vi.ID
			option[in] = *newOption
			err = tx.Model(&model.Option{}).Create(newOption).Error
			if err != nil {
				return
			}
		}

	}

	tx.Commit()

	if err == nil {
		result = "Created Successfully!"
	}

	return
}

func (m *FormServiceImpl) Get(claim *auth.AuthenticatedClaim, ID uint32) (result *model.Form, err error) {
	var existingRecord model.Form
	err = m.db.Model(&model.Form{}).
		Where("form_id = ?", ID).
		Find(&existingRecord).
		Error
	if err != nil {
		return
	}
	finalRecord, err := m.updatedRecord(m.db, claim, &existingRecord)
	if err != nil {

		return
	}

	result = finalRecord

	return
}

func (m FormServiceImpl) updatedRecord(tx *gorm.DB, claim *auth.AuthenticatedClaim, createdRecord *model.Form) (result *model.Form, err error) {
	err = tx.Model(&model.Variable{}).Where("form_id =?", createdRecord.ID).Find(&createdRecord.Variables).Error
	if err != nil {

		return
	}
	for indx, v := range createdRecord.Variables {
		err = tx.Model(&model.Option{}).Where("variable_id =?", v.ID).Find(&createdRecord.Variables[indx].Options).Error
		if err != nil {

			return
		}
	}
	result = createdRecord

	return
}
func (m FormServiceImpl) GetAll(claim *auth.AuthenticatedClaim) (result []model.Form, err error) {
	var List []model.Form
	err = m.db.Model(&model.Form{}).Find(&List).Error
	if err != nil {
		return
	}
	var resultList []model.Form = make([]model.Form, len(List))
	for indx, record := range List {
		var recValue *model.Form
		recValue, err = m.updatedRecord(m.db, claim, &record)
		if err != nil {

			return
		}
		resultList[indx] = *recValue
	}

	result = resultList
	return
}
func (m FormServiceImpl) Update(claim *auth.AuthenticatedClaim, id uint32, form *model.Form) (result string, err error) {
	if form == nil {
		return
	}
	var existingRecord *model.Form
	err = m.db.Model(&model.Form{}).Where("form_id = ?", id).Find(&existingRecord).Error
	if err != nil {
		return
	}
	if existingRecord == nil {
		err = errors.New("Data Not Exsist!")
		return
	}
	finalRecord1, err := m.updatedRecord(m.db, claim, existingRecord)
	if err != nil {
		m.db.Rollback()
		return
	}

	form.ID = existingRecord.ID

	for _, variable := range finalRecord1.Variables {
		err = m.db.Debug().
			Delete(&model.Option{}, &model.Option{VariableID: variable.ID}).
			Error
		if err != nil {
			return
		}

	}
	err = m.db.Debug().
		Delete(&model.Variable{}, &model.Variable{FormID: form.ID}).
		Error
	if err != nil {
		return
	}

	m.db.Debug().Model(&model.Form{}).Where("form_id = ?", form.ID).Take(&model.Form{}).UpdateColumns(
		map[string]interface{}{

			"form_title":       form.FormTitle,
			"form_description": form.FormDescription,
			"updated_by_id":    claim.UserId,
			"updated_at":       time.Now(),
		})

	var variables []model.Variable = make([]model.Variable, len(form.Variables))
	for indx, vi := range form.Variables {
		var newRecord *model.Variable = &vi
		newRecord.FormID = form.ID
		if form.Variables[indx].VariableTypeID == 3 || form.Variables[indx].VariableTypeID == 6 || form.Variables[indx].VariableTypeID == 7 {
			if form.Variables[indx].Options == nil {
				err = errors.New("please enter option value!")
				return
			}
		} else if form.Variables[indx].VariableTypeID != 2 {
			if form.Variables[indx].Options != nil {
				err = errors.New("option value should be null!")
				return
			}

		}
		// err = m.db.Model(&model.VariableType{}).Select("variable_type_name").Where("variable_type_id = ?", form.Variables[indx].VariableTypeID).Find(&newRecord.VariableType).Error

		variables[indx] = *newRecord
		err = m.db.Model(&model.Variable{}).Create(newRecord).Error
		if err != nil {
			return
		}
		var option []model.Option = make([]model.Option, len(variables[indx].Options))
		for in, op := range variables[indx].Options {
			var newOption *model.Option = &op
			newOption.VariableID = vi.ID
			option[in] = *newOption
			err = m.db.Model(&model.Option{}).Create(newOption).Error
			if err != nil {
				return
			}
		}

		variables = append(variables, *newRecord)
	}

	if err == nil {
		result = "Updated Successfully!"
	}

	return
}
func (m FormServiceImpl) Delete(claim *auth.AuthenticatedClaim, id uint32) (result string, err error) {
	var form *model.Form
	err = m.db.Model(&model.Form{}).Where("form_id =?", id).Find(&form).Error
	if err != nil {
		return
	}
	finalRecord, err := m.updatedRecord(m.db, claim, form)
	if err != nil {
		m.db.Rollback()
		return
	}
	if form == nil {
		result = "data not exsist!"
	} else {
		for _, variable := range finalRecord.Variables {
			err = m.db.Debug().
				Delete(&model.Option{}, &model.Option{VariableID: variable.ID}).
				Error
			if err != nil {
				return
			}

		}
		err = m.db.Debug().
			Delete(&model.Variable{}, &model.Variable{FormID: form.ID}).
			Error
		if err != nil {
			return
		}

		err = m.db.Delete(&model.Form{}, &model.Form{ID: form.ID}).Error
		if err != nil {
			return
		}
		result = "Deleted succesfully!"
	}
	return
}
