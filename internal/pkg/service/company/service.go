package company

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/eduardojabes/data-integration-challenge/entity"
	csvrepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/csv"
	"github.com/google/uuid"
)

type CompanyRepositoryBasicOperations interface {
	AddCompany(ctx context.Context, company entity.Companies) error
	ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error)
	UpdateCompany(ctx context.Context, company entity.Companies) error
}
type CompanyRepositoryImplementation interface {
	GetCompany(ctx context.Context, key string) ([]*entity.Companies, error)
}

type CompanyRepository interface {
	CompanyRepositoryBasicOperations
	CompanyRepositoryImplementation
}

type CompanyService struct {
	repository CompanyRepository
}

var (
	ERR_COMPANY_NOT_EXISTS      = errors.New("Erro: there is no company with this name")
	ERR_WHILE_MATCHING_NAME     = errors.New("Error while matching company Name")
	ERR_WHILE_MATCHING_ZIP      = errors.New("Error while matching company ZIP")
	ERR_WHILE_WRITING           = errors.New("Error while writing company")
	ERR_NOT_VALID_COMPANY       = errors.New("Error: There is invalid company camps")
	ERR_WHILE_GETTING_COMPANIES = errors.New("Error while getting companies from repositoru")
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

func (s *CompanyService) InitializeDataBase(ctx context.Context, key string) error {
	companies, err := s.repository.GetCompany(ctx, key)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
		return err
	}

	if len(companies) == 0 {
		oldrepository := s.repository

		AlterCompanyRepository(s, csvrepository.NewCompanyCSVRepository())
		companies, err := s.repository.GetCompany(ctx, key)
		if err != nil {
			err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
			return err
		}

		AlterCompanyRepository(s, oldrepository)

		for _, company := range companies {
			company.ID = uuid.New()

			CompanyNameIsValid, err := CheckNameValidity(company.Name)
			if err != nil {
				err = fmt.Errorf("%v: %w", ERR_WHILE_MATCHING_NAME, err)
				return err
			}
			CompanyZIPIsValid, err := CheckZipValidity(company.Zip)
			if err != nil {
				err = fmt.Errorf("%v: %w", ERR_WHILE_MATCHING_ZIP, err)
				return err
			}

			if CompanyNameIsValid && CompanyZIPIsValid {
				company.ID = uuid.New()
				err = s.repository.AddCompany(ctx, *company)

				if err != nil {
					err = fmt.Errorf("%v: %w", ERR_WHILE_WRITING, err)
					return err
				}
			}
		}
	}
	return nil
}

func (s *CompanyService) UpdateDataBase(ctx context.Context, key string) error {
	oldrepository := s.repository
	AlterCompanyRepository(s, csvrepository.NewCompanyCSVRepository())
	companies, err := s.repository.GetCompany(ctx, key)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
		return err
	}

	AlterCompanyRepository(s, oldrepository)

	for _, company := range companies {
		s.UpdateCompany(ctx, company)
	}

	return nil
}

func CheckAllValidity(company *entity.Companies) (bool, error) {
	CompanyNameIsValid, err := CheckNameValidity(company.Name)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_MATCHING_NAME, err)
		return false, err
	}
	CompanyZIPIsValid, err := CheckZipValidity(company.Zip)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_MATCHING_ZIP, err)
		return false, err
	}

	CompanyWebsiteIsValid := CheckWebsiteValidity(company.Website)

	if CompanyNameIsValid || CompanyZIPIsValid || CompanyWebsiteIsValid {
		return false, ERR_NOT_VALID_COMPANY
	}

	return true, nil
}

func (s *CompanyService) AddCompany(ctx context.Context, company *entity.Companies) error {
	ok, err := CheckAllValidity(company)

	if !ok {
		return err
	}

	company.ID = uuid.New()
	err = s.repository.AddCompany(ctx, *company)

	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_NOT_VALID_COMPANY, err)
		return err
	}

	return nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, company *entity.Companies) error {
	company.Name = strings.ToUpper(company.Name)

	readCompany, err := s.repository.ReadCompanyByName(ctx, company.Name)
	//fmt.Printf("ID: %v, Name %v \n", readCompany.ID, readCompany.Name)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
		return err
	}

	if readCompany == nil {
		return ERR_COMPANY_NOT_EXISTS
	}

	company.ID = readCompany.ID

	ok, err := CheckAllValidity(company)
	if !ok {
		return err
	}

	err = s.repository.UpdateCompany(ctx, *company)

	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_WRITING, err)
		return err
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
