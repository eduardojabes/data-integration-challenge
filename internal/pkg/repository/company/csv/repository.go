package company

import (
	"context"
	"encoding/csv"
	"os"
	"strings"

	"github.com/eduardojabes/data-integration-challenge/entity"
)

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

func (ccCSV *CompanyCSVRepository) GetCompany(ctx context.Context, key string) ([]*entity.Companies, error) {
	f, err := os.Open(key)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'

	data, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}
	defer f.Close()

	companyData := CreateCompanyEntityByCSV(ctx, data)

	return companyData, nil
}
