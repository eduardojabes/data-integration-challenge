package company

import (
	"net/http"

	CompanyConnector "github.com/eduardojabes/data-integration-challenge/internal/pkg/connectors/company"
)

type Route struct {
	Name        string
	Method      string
	HandlerFunc http.HandlerFunc
}

var connector = CompanyConnector.NewCompanyConnector()

type Routes []Route

var hook = Routes{
	Route{
		"MergeCompany",
		"POST",
		connector.MergeCompanies,
	},
}
