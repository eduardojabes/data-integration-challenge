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

type Handler struct {
	connector CompanyConnector.CompanyHandler
	route     Routes
}

func NewHandler() *Handler {
	return &Handler{
		connector: *CompanyConnector.NewCompanyHandler(),
	}
}

func (c *Handler) AddRoutesToConnector() {
	c.route = Routes{
		Route{
			"GetCompanies",
			"GET",
			"/v1/companies",
			c.connector.GetCompanies,
		},
		Route{
			"SearchCompany",
			"GET",
			"/v1/companies/search",
			c.connector.GetCompanyByNameAndZip,
		},
		Route{
			"CreateCompany",
			"POST",
			"/v1/companies",
			c.connector.CreateCompany,
		},
		Route{
			"MergeCompany",
			"POST",
			"/v1/companies/merge-all-companies",
			c.connector.MergeCompanies,
		},
	}
}

func (c *Handler) ImplementConnector(service companyService.CompanyService) {
	c.connector.Register(service)
	c.AddRoutesToConnector()
}
