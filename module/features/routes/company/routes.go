package company

import (
	"net/http"

	CompanyConnector "github.com/eduardojabes/data-integration-challenge/internal/pkg/handler/company"
	companyService "github.com/eduardojabes/data-integration-challenge/internal/pkg/service/company"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//var connector = CompanyConnector.NewCompanyConnector()

type Routes []Route

type Connector struct {
	connector CompanyConnector.CompanyConnector
	route     Routes
}

func NewConnector() *Connector {
	return &Connector{
		connector: *CompanyConnector.NewCompanyConnector(),
	}
}

func (c *Connector) AddRoutesToConnector() {
	c.route = Routes{
		Route{
			"MergeCompany",
			"POST",
			"/v1/companies/merge-all-companies",
			c.connector.MergeCompanies,
		},
	}
}

func (c *Connector) ImplementConnector(service companyService.CompanyService) {
	c.connector.Register(service)
	c.AddRoutesToConnector()
}
