package services

import "github.com/medivh13/koalatest/pkg/dto"

type Services interface {
	Register(req *dto.CustomersReqDTO) error
}
