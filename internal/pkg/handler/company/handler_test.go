package company

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"os"

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
	return nil, errors.New("GetCompaniesMock")
}
func (mcs *MockCompanyService) AddCompany(ctx context.Context, company *entity.Companies) error {
	if mcs.AddCompanyMock != nil {
		return mcs.AddCompanyMock(ctx, company)
	}
	return errors.New("AddCompanyMock")
}

func (mcs *MockCompanyService) FindByNameAndZip(name string, zip string) ([]entity.Companies, error) {
	if mcs.FindByNameAndZipMock != nil {
		return mcs.FindByNameAndZipMock(name, zip)
	}
	return nil, errors.New("FindByNameAndZipMock")
}

func (mcs *MockCompanyService) UpdateCompany(ctx context.Context, company *entity.Companies) error {
	if mcs.UpdateCompanyMock != nil {
		return mcs.UpdateCompanyMock(ctx, company)
	}
	return errors.New("UpdateCompanyMock")
}

type Service struct {
	service CompanyService
}

func TestMergeCompanies(t *testing.T) {

	t.Run("functional CSV", func(t *testing.T) {
		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return nil
			},
		}

		file, _ := os.Open("test_CSV.csv")

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		ff, _ := writer.CreateFormFile("csv", "test_CSV.csv")
		io.Copy(ff, file)

		writer.Close()
		request := httptest.NewRequest(http.MethodPost, "/v1/companies/merge-all-companies", bytes.NewReader(body.Bytes()))
		request.Header.Add("Content-Type", writer.FormDataContentType())

		response := httptest.NewRecorder()
		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)
		if response.Result().StatusCode != http.StatusOK {
			t.Errorf(`got "%d", want "%d"`, response.Result().StatusCode, http.StatusOK)
		}
	})

	t.Run("error in database", func(t *testing.T) {
		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return errors.New("error")
			},
		}

		file, _ := os.Open("test_CSV.csv")

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		ff, _ := writer.CreateFormFile("csv", "test_CSV.csv")
		io.Copy(ff, file)

		writer.Close()
		request := httptest.NewRequest(http.MethodPost, "/v1/companies/merge-all-companies", bytes.NewReader(body.Bytes()))
		request.Header.Add("Content-Type", writer.FormDataContentType())

		response := httptest.NewRecorder()
		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)
		if response.Result().StatusCode == http.StatusOK {
			t.Errorf(`got "%d", want error"`, response.Result().StatusCode)
		}
	})

	t.Run("error in CSV", func(t *testing.T) {
		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return nil
			},
		}

		file, _ := os.Open("CSV.csv")

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		ff, _ := writer.CreateFormFile("csv", "test_CSV.csv")
		io.Copy(ff, file)

		writer.Close()
		request := httptest.NewRequest(http.MethodPost, "/v1/companies/merge-all-companies", bytes.NewReader(body.Bytes()))
		request.Header.Add("Content-Type", writer.FormDataContentType())

		response := httptest.NewRecorder()
		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)
		if response.Result().StatusCode == http.StatusOK {
			t.Errorf(`got "%d", want error"`, response.Result().StatusCode)
		}
	})

}

func TestCreateCompany(t *testing.T) {
	t.Run("AddCompany", func(t *testing.T) {
		companyService := &MockCompanyService{
			AddCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return nil
			},
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		company := Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}
		companyJSON, _ := json.Marshal(company)

		request := httptest.NewRequest(http.MethodPost, "/v1/companies", bytes.NewBuffer(companyJSON))
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.CreateCompany(response, request)
		if response.Result().StatusCode != http.StatusCreated {
			t.Errorf(`got "%d", but don't want an error"`, response.Result().StatusCode)
		}
	})

	t.Run("error with server", func(t *testing.T) {
		err := errors.New("error")

		companyService := &MockCompanyService{
			AddCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return err
			},
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		company := Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}
		companyJSON, _ := json.Marshal(company)

		request := httptest.NewRequest(http.MethodPost, "/v1/companies", bytes.NewBuffer(companyJSON))
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.CreateCompany(response, request)
		if response.Result().StatusCode == http.StatusCreated {
			t.Errorf(`got "%d", want error"`, response.Result().StatusCode)
		}
	})

	t.Run("error with json", func(t *testing.T) {
		err := errors.New("error")

		companyService := &MockCompanyService{
			AddCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return err
			},
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		errorbytes := []byte("error")
		company := Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}
		companyJSON, _ := json.Marshal(company)
		companyJSON = append(companyJSON, errorbytes...)

		request := httptest.NewRequest(http.MethodPost, "/v1/companies", bytes.NewBuffer(companyJSON))
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.CreateCompany(response, request)
		if response.Result().StatusCode == http.StatusCreated {
			t.Errorf(`got "%d", want error"`, response.Result().StatusCode)
		}
	})
}
