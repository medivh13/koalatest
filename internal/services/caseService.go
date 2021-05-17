package services

import (
	"github.com/medivh13/koalatest/internal/repository"
	"github.com/medivh13/koalatest/pkg/dto"
	"github.com/medivh13/koalatest/pkg/dto/assembler"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Services {
	return &service{repo}
}

func (s *service) Register(req *dto.CustomersReqDTO) error {

	err := req.Validate()
	if err != nil {
		return err
	}

	err = s.repo.Register(assembler.ToSaveCustomer(req))
	if err != nil {
		return err
	}
	return nil
}
