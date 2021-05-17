package dto

import (
	"errors"

	"github.com/medivh13/koalatest/pkg/common/crypto"
	"github.com/medivh13/koalatest/pkg/common/env"
	"github.com/medivh13/koalatest/pkg/common/validator"
	util "github.com/medivh13/koalatest/pkg/utils"
)

type GetTokensReqDTO struct {
	PhoneNumberOrEmail string `json:"phone_number_or_email" valid:"required" validname="phone or email"`
	Password           string `json:"password" valid:"required" validname="password"`
	Signature          string `json:"signature" valid:"required" validname="signature"`
}

type GetTokenResponseDTO struct {
	Token string `json:"token"`
}

func (dto *GetTokensReqDTO) Validate() error {
	v := validator.NewValidate(dto)
	v.SetCustomValidation(true, func() error {
		return dto.customValidation()
	})
	return v.Validate()
}
func (dto *GetTokensReqDTO) customValidation() error {

	signature := crypto.EncodeSHA256HMAC(util.GetPrivKeySignature(), dto.PhoneNumberOrEmail)
	if signature != dto.Signature {
		if env.IsProduction() {
			return errors.New("invalid signature")
		}
		return errors.New("invalid signature" + " --> " + signature)
	}

	return nil
}
