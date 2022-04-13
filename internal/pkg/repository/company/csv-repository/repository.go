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
	ReadCompany(ctx context.Context, name string) (*entity.Companies, error)
	GetCompany(ctx context.Context) ([]*entity.Companies, error)
}

type CompanyCSVRepository struct{}

func NewCompanyCSVRepository() *CompanyCSVRepository {
	return &CompanyCSVRepository{}
}

func createCompanyEntityByCSV(ctx context.Context, fileData [][]string) []*entity.Companies {
	var companyData []*entity.Companies

	for i, line := range fileData {
		if i > 0 {
			lineRead := &entity.Companies{}

			for j, field := range line {
				if j == 0 {
					lineRead.Name = strings.ToUpper(field)
				} else if j == 1 {
					lineRead.Zip = field
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

	companyData := createCompanyEntityByCSV(ctx, data)

	return companyData, nil
}

func (ccCSV *CompanyCSVRepository) WriteCompany(ctx context.Context, company []entity.Companies) error {
	return nil
}

func (ccCSV *CompanyCSVRepository) AddCompany(ctx context.Context, company entity.Companies) error {
	return nil
}

func (ccCSV *CompanyCSVRepository) ReadCompany(ctx context.Context, name string) (*entity.Companies, error) {
	return nil, nil
}
