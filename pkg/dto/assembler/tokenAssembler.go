package assembler

import (
	models "github.com/medivh13/koalatest/internal/models"
	"github.com/medivh13/koalatest/pkg/dto"
)

func ToGetToken(m *models.GetTokens) *dto.GetTokenResponseDTO {
	return &dto.GetTokenResponseDTO{
		Token: m.Token,
	}
}

func ToTokens(datas []*models.GetTokens) []*dto.GetTokenResponseDTO {
	var ds []*dto.GetTokenResponseDTO
	for _, m := range datas {
		ds = append(ds, ToGetToken(m))
	}
	return ds
}
