package dto

import (
	"errors"

	"github.com/medivh13/koalatest/pkg/common/crypto"
	"github.com/medivh13/koalatest/pkg/common/env"
	"github.com/medivh13/koalatest/pkg/common/validator"
	util "github.com/medivh13/koalatest/pkg/utils"
)

type CustomersReqDTO struct {
	CustomerId   string `json:"customer_id"`
	CustomerName string `json:"customer_name" valid:"required" validname="name"`
	Email        string `json:"email" valid:"required,email" validname="email"`
	PhoneNumber  string `json:"phone_number" valid:"required,phoneno" validname="phone"`
	Dob          string `json:"dob" valid:"required" validname="dob"`
	Sex          string `json:"sex" valid:"required" validname="sex"`
	Salt         string `json:"salt"`
	Password     string `json:"password" valid:"required" validname="password"`
	CreatedDate  string `json:"created_date"`
	Signature    string `json:"signature" valid:"required" validname="signature"`
}

func (dto *CustomersReqDTO) Validate() error {
	v := validator.NewValidate(dto)
	v.SetCustomValidation(true, func() error {
		return dto.customValidation()
	})
	return v.Validate()
}
func (dto *CustomersReqDTO) customValidation() error {

	signature := crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), dto.CreatedDate)
	if signature != dto.Signature {
		if env.IsProduction() {
			return errors.New("invalid signature")
		}
		return errors.New("invalid signature" + " --> " + signature)
	}

	return nil
}
