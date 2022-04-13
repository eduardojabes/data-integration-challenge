package company

import (
	"context"
	"fmt"

	"github.com/eduardojabes/data-integration-challenge/entity"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

type CompanyRepository interface {
	AddCompany(ctx context.Context, company entity.Companies) error
	ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error)
	GetCompany(ctx context.Context) ([]*entity.Companies, error)
	UpdateCompany(ctx context.Context, company entity.Companies) error
}

type CompanyModel struct {
	CompanyID      uuid.UUID `db:"cc_company_id"`
	ComapanyName   string    `db:"cc_name"`
	CompanyZIP     string    `db:"cc_zip"`
	CompanyWebSite string    `db:"cc_website"`
}

type PostgreCompanyRepository struct {
	conn connector
}

type connector interface {
	pgxscan.Querier
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func NewPostgreCompanyRepository(conn connector) *PostgreCompanyRepository {
	return &PostgreCompanyRepository{conn}
}

func (r *PostgreCompanyRepository) AddCompany(ctx context.Context, company entity.Companies) error {
	_, err := r.conn.Exec(ctx, `INSERT INTO companies_catalog_table(cc_company_id, cc_name, cc_zip) values($1, $2, $3)`, company.ID, company.Name, company.Zip)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreCompanyRepository) ReadCompanyByName(ctx context.Context, name string) (*entity.Companies, error) {
	var company []*CompanyModel
	err := pgxscan.Select(ctx, r.conn, &company, `SELECT * FROM companies_catalog_table WHERE cc_name = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(company) == 0 {
		return nil, nil
	}

	return &entity.Companies{
		ID:   company[0].CompanyID,
		Name: company[0].ComapanyName,
		Zip:  company[0].CompanyZIP,
	}, nil
}

func (r PostgreCompanyRepository) UpdateCompany(ctx context.Context, company entity.Companies) error {
	_, err := r.conn.Exec(ctx, `UPDATE companies_catalog_table SET cc_name = $2, cc_zip = $3,  cc_website = $4 WHERE cc_company_id = $1`, company.ID, company.Name, company.Zip, company.Website)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreCompanyRepository) GetCompany(ctx context.Context) ([]*entity.Companies, error) {
	var companyModel []*CompanyModel
	company := []*entity.Companies{}
	err := pgxscan.Select(ctx, r.conn, &companyModel, `SELECT * FROM companies_catalog_table`)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(companyModel) == 0 {
		return nil, nil
	}

	for index := range companyModel {
		company = append(company, &entity.Companies{
			ID:   companyModel[index].CompanyID,
			Name: companyModel[index].ComapanyName,
			Zip:  companyModel[index].CompanyZIP,
		})
	}
	return company, nil
}
