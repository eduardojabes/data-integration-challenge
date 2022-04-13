package company

import (
	"context"
	"fmt"

	"github.com/eduardojabes/data-integration-challenge/entity"
	csvrepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/csv-repository"
	"github.com/google/uuid"
)

type CompanyRepository interface {
	AddCompany(ctx context.Context, company entity.Company) error
	ReadCompany(ctx context.Context, name string) (*entity.Company, error)
	GetCompany(ctx context.Context) ([]*entity.Company, error)
}

type CompanyService struct {
	repository CompanyRepository
}

func (s *CompanyService) InitializeDataBase(ctx context.Context) error {
	companies, err := s.repository.GetCompany(ctx)
	if err != nil {
		err = fmt.Errorf("error while checking database: %w", err)
		return err
	}

	if len(companies) == 0 {
		oldrepository := s.repository

		AlterCompanyRepository(s, csvrepository.NewCompanyCSVRepository())
		companies, err := s.repository.GetCompany(ctx)
		if err != nil {
			err = fmt.Errorf("error while get companies list: %w", err)
			return err
		}

		AlterCompanyRepository(s, oldrepository)

		for _, company := range companies {
			company.ID = uuid.New()
			err = s.repository.AddCompany(ctx, *company)

			if err != nil {
				err = fmt.Errorf("error while writing companies list: %w", err)
				return err
			}
		}
	}
	return nil
}

func NewCompanyService(repository CompanyRepository) *CompanyService {
	return &CompanyService{repository: repository}
}

func AlterCompanyRepository(service *CompanyService, newRepository CompanyRepository) *CompanyService {
	service.repository = newRepository
	return service
}
