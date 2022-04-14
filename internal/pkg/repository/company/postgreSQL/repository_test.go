package company

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardojabes/data-integration-challenge/entity"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func TestAddCompany(t *testing.T) {
	mock, _ := pgxmock.NewConn()

	t.Run("Adding Company", func(t *testing.T) {
		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "www.company.com",
		}

		mock.ExpectExec("INSERT INTO companies_catalog_table").
			WithArgs(company.ID, company.Name, company.Zip, company.Website).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		repository := NewPostgreCompanyRepository(mock)
		err := repository.AddCompany(context.Background(), *company)

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}
	})

	t.Run("with_error", func(t *testing.T) {

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "www.company.com",
		}

		mock.ExpectQuery("SELECT (.+) FROM companies_catalog_table WHERE (.+)").
			WillReturnError(errors.New("error"))

		repository := NewPostgreCompanyRepository(mock)
		err := repository.AddCompany(context.Background(), *company)

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}

func TestReadCompanyByName(t *testing.T) {
	t.Run("no_rows", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()
		mock.ExpectQuery("SELECT (.+) FROM companies_catalog_table WHERE (.+)").
			WillReturnRows(mock.NewRows([]string{"cc_company_id", "cc_name", "cc_zip", "cc_website"}))

		repository := NewPostgreCompanyRepository(mock)

		user, err := repository.ReadCompanyByName(context.Background(), "company")

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}

		if user != nil {
			t.Errorf("got %v want nil", user)
		}
	})

	t.Run("with_company", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		company := &entity.Companies{
			ID:      uuid.New(),
			Name:    "Company",
			Zip:     "12345",
			Website: "www.company.com",
		}

		mock.ExpectQuery("SELECT (.+) FROM companies_catalog_table WHERE (.+)").
			WillReturnRows(mock.NewRows([]string{"cc_company_id", "cc_name", "cc_zip", "cc_website"}).
				AddRow(company.ID, company.Name, company.Zip, company.Website))

		repository := NewPostgreCompanyRepository(mock)

		got, err := repository.ReadCompanyByName(context.Background(), "company")

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}

		if !reflect.DeepEqual(company, got) {
			t.Errorf("got %v want %v", got, company)
		}
	})

	t.Run("with_error", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		mock.ExpectQuery("SELECT (.+) FROM companies_catalog_table WHERE (.+)").
			WillReturnError(errors.New("error"))

		repository := NewPostgreCompanyRepository(mock)

		_, err := repository.ReadCompanyByName(context.Background(), "company")

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}

func TestUpdateComapany(t *testing.T) {
	mock, _ := pgxmock.NewConn()

	company := &entity.Companies{
		ID:      uuid.New(),
		Name:    "Company",
		Zip:     "12345",
		Website: "www.company.com",
	}

	repository := NewPostgreCompanyRepository(mock)

	t.Run("Updating Company", func(t *testing.T) {

		mock.ExpectExec("UPDATE companies_catalog_table SET ").
			WithArgs(company.ID, company.Name, company.Zip, company.Website).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repository.UpdateCompany(context.Background(), *company)

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}
	})

	t.Run("with_error", func(t *testing.T) {
		mock, _ := pgxmock.NewConn()

		mock.ExpectQuery("UPDATE companies_catalog_table SET (.+) WHERE (.+)").
			WillReturnError(errors.New("error"))

		err := repository.UpdateCompany(context.Background(), *company)

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}

func TestGetCompany(t *testing.T) {
	mock, _ := pgxmock.NewConn()

	company := &entity.Companies{
		ID:      uuid.New(),
		Name:    "Company",
		Zip:     "12345",
		Website: "www.company.com",
	}

	repository := NewPostgreCompanyRepository(mock)

	t.Run("Getting Company", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM companies_catalog_table").
			WillReturnRows(mock.NewRows([]string{"cc_company_id", "cc_name", "cc_zip", "cc_website"}).
				AddRow(company.ID, company.Name, company.Zip, company.Website))

		got, err := repository.GetCompany(context.Background(), company.Name)

		if err != nil {
			t.Errorf("got %v error, it should be nil", err)
		}

		if !reflect.DeepEqual(company, got[0]) {
			t.Errorf("got %v want %v", got[0], company)
		}
	})

	t.Run("with_error", func(t *testing.T) {

		mock.ExpectQuery("SELECT (.+) FROM companies_catalog_table WHERE (.+)").
			WillReturnError(errors.New("error"))

		_, err := repository.GetCompany(context.Background(), company.Name)

		if err == nil {
			t.Errorf("got %v want nil", err)
		}
	})
}
