package assembler

import (
	"time"

	models "github.com/medivh13/koalatest/internal/models"
	"github.com/medivh13/koalatest/pkg/dto"
)

func ToSaveCustomer(d *dto.CustomersReqDTO) *models.Customers {
	return &models.Customers{
		CustomerId:   d.CustomerId,
		CustomerName: d.CustomerName,
		Email:        d.Email,
		PhoneNumber:  d.PhoneNumber,
		Dob:          d.Dob,
		Sex:          d.Sex,
		Salt:         d.Salt,
		Password:     d.Password,
		CreatedDate:  time.Now(),
	}
}
