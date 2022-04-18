package company

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	FindByNameAndZipMock func(name string, zip string) (*entity.Companies, error)
	FindByNameMock       func(name string) (*entity.Companies, error)
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

func (mcs *MockCompanyService) FindByNameAndZip(name string, zip string) (*entity.Companies, error) {
	if mcs.FindByNameAndZipMock != nil {
		return mcs.FindByNameAndZipMock(name, zip)
	}
	return nil, errors.New("FindByNameAndZipMock")
}

func (mcs *MockCompanyService) FindByName(name string) (*entity.Companies, error) {
	if mcs.FindByNameMock != nil {
		return mcs.FindByNameMock(name)
	}
	return nil, errors.New("FindByNameMock")
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

func CreatTestFile(data string) string {
	file, _ := ioutil.TempFile("./", "test_file_*.csv")

	buffer := &bytes.Buffer{}
	buffer.WriteString(data)

	file.Write(buffer.Bytes())

	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		panic(err)
	}

	file.Close()

	return file.Name()
}

func CreateHttpRequestAndResponse(fileName string) (*http.Request, *httptest.ResponseRecorder) {
	file, _ := os.Open(fileName)

	body := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(body)

	ioWriter, _ := mpWriter.CreateFormFile("csv", fileName)
	io.Copy(ioWriter, file)

	mpWriter.Close()
	request := httptest.NewRequest(http.MethodPost, "/v1/companies/merge-all-companies", bytes.NewReader(body.Bytes()))
	request.Header.Add("Content-Type", mpWriter.FormDataContentType())

	response := httptest.NewRecorder()

	return request, response
}

func TestMergeCompanies(t *testing.T) {
	t.Run("error in database", func(t *testing.T) {
		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return errors.New("error")
			},
		}
		data := "name;addresszip;website \n tola sales group;78229;http://repsources.com"
		fileName := CreatTestFile(data)

		request, response := CreateHttpRequestAndResponse(fileName)

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)
		if response.Result().StatusCode != http.StatusOK {
			t.Errorf(`got "%d", want error"`, response.Result().StatusCode)
		}
		os.Remove(fileName)
	})

	t.Run("Error in Formfile", func(t *testing.T) {

		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return errors.New("error")
			},
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/v1/companies/merge-all-companies", bytes.NewReader(body.Bytes()))
		request.Header.Add("Content-Type", writer.FormDataContentType())

		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)

		if response.Code != http.StatusInternalServerError {
			t.Errorf("got: %d, want: %d", response.Code, http.StatusInternalServerError)
		}

	})

	t.Run("Error in Lenght Data", func(t *testing.T) {
		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return nil
			},
		}

		data := ""
		fileName := CreatTestFile(data)

		request, response := CreateHttpRequestAndResponse(fileName)

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)

		if response.Code != http.StatusBadRequest {
			t.Errorf("got: %d, want: %d", response.Code, http.StatusBadRequest)
		}
		os.Remove(fileName)
	})
	t.Run("Correrct Update Data", func(t *testing.T) {
		companyService := &MockCompanyService{
			UpdateCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return nil
			},
		}

		data := "name;addresszip;website \n tola sales group;78229;http://repsources.com"
		fileName := CreatTestFile(data)
		request, response := CreateHttpRequestAndResponse(fileName)
		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.MergeCompanies(response, request)

		if response.Code == http.StatusBadRequest {
			t.Errorf("got: %d, want: %d", response.Code, http.StatusBadRequest)
		}
		os.Remove(fileName)
	})
}

func TestCreateCompany(t *testing.T) {
	t.Run("AddCompany", func(t *testing.T) {
		companyService := &MockCompanyService{
			AddCompanyMock: func(ctx context.Context, company *entity.Companies) error {
				return nil
			},
		}

		company := entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}
		companyJSON, _ := json.Marshal(company)

		request := httptest.NewRequest(http.MethodPost, "/v1/companies", bytes.NewBuffer(companyJSON))
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

		company := entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}
		companyJSON, _ := json.Marshal(company)

		request := httptest.NewRequest(http.MethodPost, "/v1/companies", bytes.NewBuffer(companyJSON))
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

		errorbytes := []byte("error")
		company := entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}
		companyJSON, _ := json.Marshal(company)
		companyJSON = append(companyJSON, errorbytes...)

		request := httptest.NewRequest(http.MethodPost, "/v1/companies", bytes.NewBuffer(companyJSON))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.CreateCompany(response, request)
		if response.Result().StatusCode == http.StatusCreated {
			t.Errorf(`got "%d", want error"`, response.Result().StatusCode)
		}
	})
}

func TestGetCompanyByNameAndZip(t *testing.T) {
	t.Run("get company by name and zip", func(t *testing.T) {
		nameForSearch := "Company"
		zipForSearch := "12345"
		company := &entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}

		companyService := &MockCompanyService{
			FindByNameAndZipMock: func(name string, zip string) (*entity.Companies, error) {
				return company, nil
			},
		}

		queryURL := fmt.Sprintf("/v1/companies/search?name=%s&zip=%s", nameForSearch, zipForSearch)
		dummyIO := []byte{}

		request := httptest.NewRequest(http.MethodPost, queryURL, bytes.NewBuffer(dummyIO))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.GetCompanyByNameAndZip(response, request)

		var readCompany entity.Companies
		body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

		err := json.Unmarshal(body, &readCompany)
		if err != nil {
			t.Errorf(`got "%v", but expected none"`, err)
		}
		if !reflect.DeepEqual(company, &readCompany) {
			t.Errorf(`got "%s", want %s"`, readCompany, company)
		}
	})

	t.Run("error with server", func(t *testing.T) {
		nameForSearch := "Company"
		zipForSearch := "12345"
		company := &entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}

		companyService := &MockCompanyService{
			FindByNameAndZipMock: func(name string, zip string) (*entity.Companies, error) {
				return nil, errors.New("error")
			},
		}

		queryURL := fmt.Sprintf("/v1/companies/search?name=%s&zip=%s", nameForSearch, zipForSearch)
		dummyIO := []byte{}

		request := httptest.NewRequest(http.MethodPost, queryURL, bytes.NewBuffer(dummyIO))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.GetCompanyByNameAndZip(response, request)

		var readCompany entity.Companies
		body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

		json.Unmarshal(body, &readCompany)
		if http.StatusInternalServerError != response.Code {
			t.Errorf("expected error got %v, but expected %v", response.Code, http.StatusInternalServerError)
		}
		if reflect.DeepEqual(company, &readCompany) {
			t.Errorf(`got "%s", want empty company`, readCompany)
		}
	})
	t.Run("not exist company", func(t *testing.T) {
		nameForSearch := "Company"
		zipForSearch := "12345"
		company := &entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}

		companyService := &MockCompanyService{
			FindByNameAndZipMock: func(name string, zip string) (*entity.Companies, error) {
				return nil, nil
			},
		}

		queryURL := fmt.Sprintf("/v1/companies/search?name=%s&zip=%s", nameForSearch, zipForSearch)
		dummyIO := []byte{}

		request := httptest.NewRequest(http.MethodPost, queryURL, bytes.NewBuffer(dummyIO))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.GetCompanyByNameAndZip(response, request)

		var readCompany entity.Companies
		body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

		json.Unmarshal(body, &readCompany)
		if http.StatusInternalServerError == response.Code {
			t.Errorf("not expected error got %v", response.Code)
		}
		if reflect.DeepEqual(company, &readCompany) {
			t.Errorf(`got "%s", want empty company`, readCompany)
		}
	})
}

func TestGetCompanyByName(t *testing.T) {
	t.Run("get company by name", func(t *testing.T) {
		nameForSearch := "Company"
		company := &entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}

		companyService := &MockCompanyService{
			FindByNameMock: func(name string) (*entity.Companies, error) {
				return company, nil
			},
		}

		queryURL := fmt.Sprintf("/v1/companies/search?name=%s", nameForSearch)
		dummyIO := []byte{}

		request := httptest.NewRequest(http.MethodPost, queryURL, bytes.NewBuffer(dummyIO))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.GetCompanyByName(response, request)

		var readCompany entity.Companies
		body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

		err := json.Unmarshal(body, &readCompany)
		if err != nil {
			t.Errorf(`got "%v", but expected none"`, err)
		}
		if !reflect.DeepEqual(company, &readCompany) {
			t.Errorf(`got "%s", want %s"`, readCompany, company)
		}
	})

	t.Run("error with server", func(t *testing.T) {
		nameForSearch := "Company"
		company := &entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}

		companyService := &MockCompanyService{
			FindByNameMock: func(name string) (*entity.Companies, error) {
				return nil, errors.New("error")
			},
		}

		queryURL := fmt.Sprintf("/v1/companies/search?name=%s", nameForSearch)
		dummyIO := []byte{}

		request := httptest.NewRequest(http.MethodPost, queryURL, bytes.NewBuffer(dummyIO))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.GetCompanyByName(response, request)

		var readCompany entity.Companies
		body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

		if http.StatusInternalServerError != response.Code {
			t.Errorf("expected error got %v, but expected %v", response.Code, http.StatusInternalServerError)
		}
		json.Unmarshal(body, &readCompany)
		if reflect.DeepEqual(company, &readCompany) {
			t.Errorf(`got "%s", want empty company`, readCompany)
		}
	})
	t.Run("not exist company", func(t *testing.T) {
		nameForSearch := "Company"
		company := &entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://new_website.com"}

		companyService := &MockCompanyService{
			FindByNameMock: func(name string) (*entity.Companies, error) {
				return nil, nil
			},
		}

		queryURL := fmt.Sprintf("/v1/companies/search?name=%s", nameForSearch)
		dummyIO := []byte{}

		request := httptest.NewRequest(http.MethodPost, queryURL, bytes.NewBuffer(dummyIO))
		response := httptest.NewRecorder()

		companyHandler := NewCompanyHandler()
		companyHandler.Register(companyService)

		companyHandler.GetCompanyByName(response, request)

		var readCompany entity.Companies
		body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

		json.Unmarshal(body, &readCompany)

		if http.StatusInternalServerError == response.Code {
			t.Errorf("not expected error got %v", response.Code)
		}
		if reflect.DeepEqual(company, &readCompany) {
			t.Errorf(`got "%s", want empty company`, readCompany)
		}
	})
}
