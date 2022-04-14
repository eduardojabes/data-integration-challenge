package company

import (
	"context"
	"encoding/csv"
	"errors"
	"os"
	"reflect"

	"github.com/eduardojabes/data-integration-challenge/entity"

	"net/http"
	"net/http/httptest"
	"testing"
)

type MockCompanyService struct {
	GetCompaniesMock     func() ([]entity.Companies, error)
	AddCompanyMock       func(ctx context.Context, company *entity.Companies) error
	FindByNameAndZipMock func(name string, zip string) ([]entity.Companies, error)
	UpdateCompanyMock    func(ctx context.Context, company *entity.Companies) error
}

func (mcs *MockCompanyService) GetCompanies() ([]entity.Companies, error) {
	if mcs.GetCompaniesMock != nil {
		return mcs.GetCompaniesMock()
	}
	return nil, errors.New("GetCodeFileMock must be set")
}
func (mcs *MockCompanyService) AddCompany(ctx context.Context, company *entity.Companies) error {
	if mcs.AddCompanyMock != nil {
		return mcs.AddCompanyMock(ctx, company)
	}
	return errors.New("GetCodeFileMock must be set")
}

func (mcs *MockCompanyService) FindByNameAndZip(name string, zip string) ([]entity.Companies, error) {
	if mcs.FindByNameAndZipMock != nil {
		return mcs.FindByNameAndZipMock(name, zip)
	}
	return nil, errors.New("GetCodeFileMock must be set")
}

func (mcs *MockCompanyService) UpdateCompany(ctx context.Context, company *entity.Companies) error {
	if mcs.UpdateCompanyMock != nil {
		return mcs.UpdateCompanyMock(ctx, company)
	}
	return errors.New("GetCodeFileMock must be set")
}

type Service struct {
	service CompanyService
}

func TestMergeCompanies(t *testing.T) {
	companyService := &MockCompanyService{}

	t.Run("functional CSV", func(t *testing.T) {
		f, _ := os.Open("test_CSV.csv")
		csvReader := csv.NewReader(f)
		data, _ := csvReader.ReadAll()

		request := httptest.NewRequest(http.MethodPost, "/v1/companies/merge-all-companies", f)
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)
		if !reflect.DeepEqual(response.Body.String(), data) {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}
	})
}
