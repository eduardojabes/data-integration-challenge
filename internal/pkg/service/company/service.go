package company

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/eduardojabes/data-integration-challenge/entity"
	"github.com/google/uuid"
)

type dbCompanyRepository interface {
	AddCompany(ctx context.Context, company entity.Companies) error
	ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error)
	SearchCompanyByNameAndZip(ctx context.Context, name string, zip string) (*entity.Companies, error)
	UpdateCompany(ctx context.Context, company entity.Companies) error
	DeleteCompany(ctx context.Context, company entity.Companies) error
}
type csvCompanyRepository interface {
	GetCompany(ctx context.Context, key string) ([]*entity.Companies, error)
}

type CompanyRepository interface {
	dbCompanyRepository
	csvCompanyRepository
}

type CompanyService struct {
	dbRepository  CompanyRepository
	csvRepository csvCompanyRepository
}

var (
	ERR_COMPANY_NOT_EXISTS      = errors.New("Erro: there is no company with this name")
	ERR_COMPANY_EXISTS          = errors.New("Erro: there is a company with this name")
	ERR_WHILE_MATCHING_NAME     = errors.New("Error while matching company Name")
	ERR_WHILE_MATCHING_ZIP      = errors.New("Error while matching company ZIP")
	ERR_WHILE_WRITING           = errors.New("Error while writing company")
	ERR_NOT_VALID_COMPANY       = errors.New("Error: There is invalid company camps")
	ERR_WHILE_GETTING_COMPANIES = errors.New("Error while getting companies from repository")
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
	companies, err := s.dbRepository.GetCompany(ctx, key)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
		return err
	}

	if len(companies) == 0 {

		companies, err := s.csvRepository.GetCompany(ctx, key)
		if err != nil {
			return fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)

		}

		for _, company := range companies {
			company.ID = uuid.New()

			ok, _ := CheckAllValidity(company)

			if ok {
				company.ID = uuid.New()
				err = s.dbRepository.AddCompany(ctx, *company)

				if err != nil {
					return fmt.Errorf("%v: %w", ERR_WHILE_WRITING, err)
				}
			}
		}
	}
	return nil
}

func (s *CompanyService) UpdateDataBaseFromCSV(ctx context.Context, key string) error {
	companies, err := s.csvRepository.GetCompany(ctx, key)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
		return err
	}

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

	if !CompanyNameIsValid || !CompanyZIPIsValid || !CompanyWebsiteIsValid {
		return false, ERR_NOT_VALID_COMPANY
	}

	return true, nil
}

func (s *CompanyService) AddCompany(ctx context.Context, company *entity.Companies) error {
	company.Name = strings.ToUpper(company.Name)

	//fmt.Printf("service.go ID: %s, name: %s, zip: %s, webmail:%s\n", company.ID, company.Name, company.Zip, company.Website)
	readCompany, err := s.dbRepository.SearchCompanyByNameAndZip(ctx, company.Name, company.Zip)
	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_GETTING_COMPANIES, err)
		return err
	}

	if readCompany != nil {
		return ERR_COMPANY_EXISTS
	}

	ok, err := CheckAllValidity(company)

	if !ok {
		return err
	}

	company.ID = uuid.New()
	err = s.dbRepository.AddCompany(ctx, *company)

	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_NOT_VALID_COMPANY, err)
		return err
	}

	return nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, company *entity.Companies) error {
	company.Name = strings.ToUpper(company.Name)

	readCompany, err := s.dbRepository.ReadCompanyByName(ctx, company.Name)
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

	err = s.dbRepository.UpdateCompany(ctx, *company)

	if err != nil {
		err = fmt.Errorf("%v: %w", ERR_WHILE_WRITING, err)
		return err
	}
	return nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, entity entity.Companies) error {
	err := s.dbRepository.DeleteCompany(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (s *CompanyService) GetCompanies() ([]entity.Companies, error) {
	companiesReferences, err := s.dbRepository.GetCompany(context.Background(), "")
	var companies []entity.Companies
	for _, values := range companiesReferences {
		companies = append(companies, *values)
	}
	return companies, err
}

func (s *CompanyService) FindByNameAndZip(name string, zip string) (*entity.Companies, error) {
	companies, err := s.dbRepository.SearchCompanyByNameAndZip(context.Background(), name, zip)
	return companies, err
}

func (s *CompanyService) FindByName(name string) (*entity.Companies, error) {
	companies, err := s.dbRepository.ReadCompanyByName(context.Background(), name)
	return companies, err
}

func NewCompanyService(dbRepository CompanyRepository, csvRepository csvCompanyRepository) *CompanyService {
	return &CompanyService{
		dbRepository:  dbRepository,
		csvRepository: csvRepository,
	}
}
