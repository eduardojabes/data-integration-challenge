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
	AddCompany(ctx context.Context, company entity.Company) error
	ReadCompany(ctx context.Context, name string) (*entity.Company, error)
	GetCompany(ctx context.Context) ([]*entity.Company, error)
}

type CompanyModel struct {
	CompanyID    uuid.UUID `db:"cc_company_id"`
	ComapanyName string    `db:"cc_name"`
	CompanyZIP   string    `db:"cc_zip"`
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

func (r *PostgreCompanyRepository) AddCompany(ctx context.Context, company entity.Company) error {
	_, err := r.conn.Exec(ctx, `INSERT INTO companies_catalog_table(cc_company_id, cc_name, cc_zip) values($1, $2, $3)`, company.ID, company.Name, company.Zip)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreCompanyRepository) ReadCompany(ctx context.Context, name string) (*entity.Company, error) {
	var company []*CompanyModel
	err := pgxscan.Select(ctx, r.conn, &company, `SELECT * FROM companies_catalog_table WHERE cc_name = $1`, name)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(company) == 0 {
		return nil, nil
	}

	return &entity.Company{
		ID:   company[0].CompanyID,
		Name: company[0].ComapanyName,
		Zip:  company[0].CompanyZIP,
	}, nil
}

func (r *PostgreCompanyRepository) GetCompany(ctx context.Context) ([]*entity.Company, error) {
	var companyModel []*CompanyModel
	company := []*entity.Company{}
	err := pgxscan.Select(ctx, r.conn, &companyModel, `SELECT * FROM companies_catalog_table`)
	if err != nil {
		return nil, fmt.Errorf("error while executing query: %w", err)
	}

	if len(companyModel) == 0 {
		return nil, nil
	}

	for index := range companyModel {
		company = append(company, &entity.Company{
			ID:   companyModel[index].CompanyID,
			Name: companyModel[index].ComapanyName,
			Zip:  companyModel[index].CompanyZIP,
		})
	}
	return company, nil
}
