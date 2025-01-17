package company

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardojabes/data-integration-challenge/entity"
	"github.com/google/uuid"
)

type MockCompanyRepository struct {
	AddCompanyMock                func(ctx context.Context, company entity.Companies) error
	ReadCompanyByNameMock         func(ctx context.Context, name string) (*entity.Companies, error)
	SearchCompanyByNameAndZipMock func(ctx context.Context, name string, zip string) (*entity.Companies, error)
	UpdateCompanyMock             func(ctx context.Context, company entity.Companies) error
	GetCompanyMock                func(ctx context.Context, key string) ([]*entity.Companies, error)
	DeleteCompanyMock             func(ctx context.Context, company entity.Companies) error
}

func (mcr *MockCompanyRepository) AddCompany(ctx context.Context, company entity.Companies) error {
	if mcr.AddCompanyMock != nil {
		return mcr.AddCompanyMock(ctx, company)
	}
	return errors.New("AddCompanyMock must be set")
}

func (mcr *MockCompanyRepository) ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error) {
	if mcr.ReadCompanyByNameMock != nil {
		return mcr.ReadCompanyByNameMock(ctx, name)
	}
	return nil, errors.New("ReadCompanyByNameMock must be set")
}

func (mcr *MockCompanyRepository) UpdateCompany(ctx context.Context, company entity.Companies) error {
	if mcr.UpdateCompanyMock != nil {
		return mcr.UpdateCompanyMock(ctx, company)
	}
	return errors.New("UpdateCompanyMock must be set")
}

func (mcr *MockCompanyRepository) GetCompany(ctx context.Context, key string) ([]*entity.Companies, error) {
	if mcr.GetCompanyMock != nil {
		return mcr.GetCompanyMock(ctx, key)
	}
	return nil, errors.New("GetCompanyMock must be set")
}

func (mcr *MockCompanyRepository) SearchCompanyByNameAndZip(ctx context.Context, name string, zip string) (*entity.Companies, error) {
	if mcr.SearchCompanyByNameAndZipMock != nil {
		return mcr.SearchCompanyByNameAndZipMock(ctx, name, zip)
	}
	return nil, errors.New("SearchCompanyByNameAndZipMock must be set")
}
func (mcr *MockCompanyRepository) DeleteCompany(ctx context.Context, company entity.Companies) error {
	if mcr.DeleteCompanyMock != nil {
		return mcr.DeleteCompanyMock(ctx, company)
	}
	return errors.New("DeleteCompanyMock must be set")
}

type MockCsvCompanyRepository struct {
	GetCompanyMock func(ctx context.Context, key string) ([]*entity.Companies, error)
}

func (mcsvr *MockCsvCompanyRepository) GetCompany(ctx context.Context, key string) ([]*entity.Companies, error) {
	if mcsvr.GetCompanyMock != nil {
		return mcsvr.GetCompanyMock(ctx, key)
	}
	return nil, errors.New("GetCompanyMock must be set")
}

func TestCheckNameValidity(t *testing.T) {
	t.Run("Valid name", func(t *testing.T) {
		name := "COMPANY"
		got, _ := CheckNameValidity(name)

		if !got {
			t.Errorf("Expected true, but get %v", got)
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		name := "CoMPANY."
		got, _ := CheckNameValidity(name)

		if got {
			t.Errorf("Expected false, but get %v", got)
		}
	})
}

func TestCheckZipValidity(t *testing.T) {
	t.Run("Valid zip", func(t *testing.T) {
		zip := "12345"
		got, _ := CheckZipValidity(zip)

		if !got {
			t.Errorf("Expected true, but get %v", got)
		}
	})

	t.Run("Invalid zip", func(t *testing.T) {
		zip := "123456"
		got, _ := CheckZipValidity(zip)

		if got {
			t.Errorf("Expected false, but get %v", got)
		}
	})
}

func TestCheckWebsiteValidity(t *testing.T) {
	t.Run("Valid website", func(t *testing.T) {
		website := "http://www.company.com"
		got := CheckWebsiteValidity(website)

		if !got {
			t.Errorf("Expected true, but get %v", got)
		}
	})

	t.Run("Invalid website", func(t *testing.T) {
		website := "http://www.om*pany.coam"
		got := CheckWebsiteValidity(website)

		if got {
			t.Errorf("Expected false, but get %v", got)
		}
	})
}

func TestInitializeDataBase(t *testing.T) {

	t.Run("Error getting data from database", func(t *testing.T) {
		want := errors.New("error")

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, want
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, nil
			},
		}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.InitializeDataBase(context.Background(), "")

		if errors.Is(nil, err) {
			t.Errorf("expected %v, but got %v", ERR_WHILE_GETTING_COMPANIES, err)
		}
	})

	t.Run("Error opening csv", func(t *testing.T) {
		want := errors.New("error")

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, want
			},
		}
		service := NewCompanyService(dbRepository, csvRepository)

		err := service.InitializeDataBase(context.Background(), "")

		if err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})

	t.Run("Sucessfull Init database", func(t *testing.T) {
		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dataread := []*entity.Companies{company}

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return dataread, nil
			},
		}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.InitializeDataBase(context.Background(), "test_CSV.csv")

		if &err == nil {
			t.Errorf("expected nil, but got %v", err)
		}
	})
}

func TestUpdateDataBase(t *testing.T) {

	t.Run("Error getting data from database", func(t *testing.T) {
		want := errors.New("error")

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, want
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, nil
			},
		}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateDataBaseFromCSV(context.Background(), "test_CSV.csv")

		if &err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})

	t.Run("Error while acessing csv", func(t *testing.T) {
		want := errors.New("error")

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dataread := []*entity.Companies{company}

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return dataread, nil
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, want
			},
		}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateDataBaseFromCSV(context.Background(), "test_CSV.csv")

		if err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})
	t.Run("Error while updating database", func(t *testing.T) {
		want := errors.New("error")

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dataread := []*entity.Companies{company}

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return dataread, nil
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return want
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return dataread, want
			},
		}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateDataBaseFromCSV(context.Background(), "test_CSV.csv")

		if err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})

	t.Run("Sucessfull Update database", func(t *testing.T) {
		want := errors.New("error")

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dataread := []*entity.Companies{company}

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return dataread, nil
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return dataread, want
			},
		}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateDataBaseFromCSV(context.Background(), "test_CSV.csv")

		if &err == nil {
			t.Errorf("expected nil, but got %v", err)
		}
	})
}

func TestAddCompany(t *testing.T) {
	t.Run("Adding Company", func(t *testing.T) {
		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12346",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			AddCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
			SearchCompanyByNameAndZipMock: func(ctx context.Context, name, zip string) (*entity.Companies, error) {
				return nil, nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{}
		service := NewCompanyService(dbRepository, csvRepository)

		err := service.AddCompany(context.Background(), company)

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}
	})

	t.Run("with_error", func(t *testing.T) {
		want := errors.New("error")

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			AddCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return want
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)
		err := service.AddCompany(context.Background(), company)

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}

func TestUpdateCompany(t *testing.T) {
	t.Run("Error getting data from database", func(t *testing.T) {
		want := errors.New("error")

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			ReadCompanyByNameMock: func(ctx context.Context, name string) (*entity.Companies, error) {
				return nil, want
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateCompany(context.Background(), company)

		if &err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})

	t.Run("Error not existscompany", func(t *testing.T) {

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			ReadCompanyByNameMock: func(ctx context.Context, name string) (*entity.Companies, error) {
				return nil, nil
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateCompany(context.Background(), company)

		if !errors.Is(err, ERR_COMPANY_NOT_EXISTS) {
			t.Errorf("expected %v, but got %v", ERR_COMPANY_NOT_EXISTS, err)
		}
	})

	t.Run("Error while updating database", func(t *testing.T) {
		want := errors.New("error")

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			ReadCompanyByNameMock: func(ctx context.Context, name string) (*entity.Companies, error) {
				return company, nil
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return want
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateCompany(context.Background(), company)

		if &err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})

	t.Run("Sucessful update on database", func(t *testing.T) {
		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			ReadCompanyByNameMock: func(ctx context.Context, name string) (*entity.Companies, error) {
				return company, nil
			},
			UpdateCompanyMock: func(ctx context.Context, company entity.Companies) error {
				return nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		err := service.UpdateCompany(context.Background(), company)

		if err != nil {
			t.Errorf("not expected an error, but got %v", err)
		}
	})
}

func TestGetCompanies(t *testing.T) {
	t.Run("Error getting data from database", func(t *testing.T) {
		want := errors.New("error")

		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return nil, want
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		_, err := service.GetCompanies()

		if &err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})
	t.Run("Getting data from database", func(t *testing.T) {

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		readData := []*entity.Companies{company}
		dbRepository := &MockCompanyRepository{
			GetCompanyMock: func(ctx context.Context, key string) ([]*entity.Companies, error) {
				return readData, nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		got, err := service.GetCompanies()

		if err != nil {
			t.Errorf("not expected an error, but got %v", err)
		}
		if !reflect.DeepEqual(*company, got[0]) {
			t.Errorf("expected %v, but got %v", *company, got[0])
		}
	})
}

func TestFindByNameAndZip(t *testing.T) {
	t.Run("Error getting data from database", func(t *testing.T) {
		want := errors.New("error")

		dbRepository := &MockCompanyRepository{
			SearchCompanyByNameAndZipMock: func(ctx context.Context, name, zip string) (*entity.Companies, error) {
				return nil, want
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		_, err := service.FindByNameAndZip("", "")

		if &err == nil {
			t.Errorf("expected an error, but got %v", err)
		}
	})
	t.Run("Getting data from database", func(t *testing.T) {

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "http://www.company.com",
		}

		dbRepository := &MockCompanyRepository{
			SearchCompanyByNameAndZipMock: func(ctx context.Context, name, zip string) (*entity.Companies, error) {
				return company, nil
			},
		}

		csvRepository := &MockCsvCompanyRepository{}

		service := NewCompanyService(dbRepository, csvRepository)

		got, err := service.FindByNameAndZip(company.Name, company.Zip)

		if err != nil {
			t.Errorf("not expected an error, but got %v", err)
		}
		if !reflect.DeepEqual(company, got) {
			t.Errorf("expected %v, but got %v", company, got)
		}
	})
}
