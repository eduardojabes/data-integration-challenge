package company

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/eduardojabes/data-integration-challenge/entity"
)

type CompanyRepository interface {
	AddCompany(ctx context.Context, company entity.Companies) error
	ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error)
	GetCompany(ctx context.Context) ([]*entity.Companies, error)
	UpdateCompany(ctx context.Context, company entity.Companies) error
}

type CompanyCSVRepository struct{}

func NewCompanyCSVRepository() *CompanyCSVRepository {
	return &CompanyCSVRepository{}
}

func CreateCompanyEntityByCSV(ctx context.Context, fileData [][]string) []*entity.Companies {
	var companyData []*entity.Companies

	for i, line := range fileData {
		if i > 0 {
			lineRead := &entity.Companies{}

			for j, field := range line {
				switch j {
				case 0:
					lineRead.Name = strings.ToUpper(field)
				case 1:
					lineRead.Zip = field
				case 2:
					lineRead.Website = field
				}
			}

			companyData = append(companyData, lineRead)
		}
	}

	return companyData
}

func (ccCSV *CompanyCSVRepository) GetCompany(ctx context.Context) ([]*entity.Companies, error) {
	f, err := os.Open("/mnt/c/Golang/data-integration-chalenge/data-integration-challenge/data/q1_catalog.csv")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'

	data, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer f.Close()

	companyData := CreateCompanyEntityByCSV(ctx, data)

	return companyData, nil
}

func (ccCSV *CompanyCSVRepository) WriteCompany(ctx context.Context, company []entity.Companies) error {
	//To be implemented
	return nil
}

func (ccCSV *CompanyCSVRepository) AddCompany(ctx context.Context, company entity.Companies) error {
	//To be implemented
	return nil
}

func (ccCSV *CompanyCSVRepository) ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error) {
	//To be implemented
	return nil, nil
}

func (ccCSV *CompanyCSVRepository) UpdateCompany(ctx context.Context, company entity.Companies) error {
	//To be implemented
	return nil
}
