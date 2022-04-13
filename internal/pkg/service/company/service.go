package company

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/eduardojabes/data-integration-challenge/entity"
	csvrepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/csv-repository"
	"github.com/google/uuid"
)

type CompanyRepository interface {
	AddCompany(ctx context.Context, company entity.Companies) error
	ReadCompany(ctx context.Context, name string) (*entity.Companies, error)
	GetCompany(ctx context.Context) ([]*entity.Companies, error)
}

type CompanyService struct {
	repository CompanyRepository
}

var (
	ERR_WHILE_MATCHING_NAME = errors.New("Error While Matching Company Name")
	ERR_WHILE_MATCHING_ZIP  = errors.New("Error While Matching Company ZIP")
)

func CheckNameValidity(name string) (bool, error) {
	ok, err := regexp.MatchString("^[A-Z&' ]*$", name)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_MATCHING_NAME, err)
		return false, err
	}
	return ok, nil
}

func CheckZipValidity(zip string) (bool, error) {
	ok, err := regexp.MatchString("^[0-9]{5}$", zip)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_MATCHING_ZIP, err)
		return false, err
	}
	return ok, nil
}

func CheckWebsiteValidity(website string) bool {
	register := regexp.MustCompile(`(?mi)^((http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?)?$`)
	ok := register.MatchString(website)
	return ok
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

			CompanyNameIsValid, err := CheckNameValidity(company.Name)
			if err != nil {
				err = fmt.Errorf("%w", err)
				return err
			}
			CompanyZIPIsValid, err := CheckZipValidity(company.Zip)
			if err != nil {
				err = fmt.Errorf("%w", err)
				return err
			}

			if CompanyNameIsValid && CompanyZIPIsValid {
				company.ID = uuid.New()
				err = s.repository.AddCompany(ctx, *company)

				if err != nil {
					err = fmt.Errorf("error while writing companies list: %w", err)
					return err
				}
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
