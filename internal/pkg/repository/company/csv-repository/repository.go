package company

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eduardojabes/data-integration-challenge/entity"
)

type CompanyRepository interface {
	AddCompany(ctx context.Context, company entity.Company) error
	ReadCompany(ctx context.Context, name string) (*entity.Company, error)
	GetCompany(ctx context.Context) ([]*entity.Company, error)
}

type CompanyCSVRepository struct{}

func NewCompanyCSVRepository() *CompanyCSVRepository {
	return &CompanyCSVRepository{}
}

func createCompanyEntityByCSV(ctx context.Context, fileData [][]string) []*entity.Company {
	var companyData []*entity.Company

	for i, line := range fileData {
		if i > 0 {
			lineRead := &entity.Company{}

			for j, field := range line {
				if j == 0 {
					lineRead.Name = strings.ToUpper(field)
				} else if j == 1 {
					lineRead.Zip = field
				}
			}

			if len(lineRead.Name) > 0 && len(lineRead.Zip) == 5 {
				companyData = append(companyData, lineRead)
				fmt.Printf("Company: %s, zip:%s\n", lineRead.Name, lineRead.Zip)
			} else {
				fmt.Printf("Warning! There is some data inconsistence while importing from CSV. Check Company: %s, zip:%s\n", lineRead.Name, lineRead.Zip)
			}
		}
	}
	return companyData
}

func (ccCSV *CompanyCSVRepository) GetCompany(ctx context.Context) ([]*entity.Company, error) {
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

func (ccCSV *CompanyCSVRepository) WriteCompany(ctx context.Context, company []entity.Company) error {
	return nil
}

func (ccCSV *CompanyCSVRepository) AddCompany(ctx context.Context, company entity.Company) error {
	return nil
}

func (ccCSV *CompanyCSVRepository) ReadCompany(ctx context.Context, name string) (*entity.Company, error) {
	return nil, nil
}
