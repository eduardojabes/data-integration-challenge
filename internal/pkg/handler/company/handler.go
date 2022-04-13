package company

import (
	"bufio"
	"context"
	"encoding/csv"
	"log"
	"net/http"

	csvRepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/csv"
	service "github.com/eduardojabes/data-integration-challenge/internal/pkg/service/company"
	"github.com/google/uuid"
)

type Companies struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Zip     string    `json:"zipCode"`
	Website string    `json:"website"`
}

type CompanyConnector struct {
	service service.CompanyService
}

func NewCompanyConnector() *CompanyConnector {
	return &CompanyConnector{}
}

func (c *CompanyConnector) Register(service service.CompanyService) {
	c.service = service
}

func (c *CompanyConnector) MergeCompanies(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	file, _, err := r.FormFile("csv")
	if err != nil {
		log.Fatalln("Error MergeCompany", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	csvreader := csv.NewReader(bufio.NewReader(file))
	csvreader.Comma = ';'
	data, err := csvreader.ReadAll()

	if err != nil {
		log.Fatalln("Error MergeCompany", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(data) == 0 {
		log.Fatalln("Empty file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	companyData := csvRepository.CreateCompanyEntityByCSV(ctx, data)
	for _, company := range companyData {
		err = c.service.UpdateCompany(ctx, company)
		if err != nil {
			log.Fatalln("Error UpdatingCompany", err)
		}
	}

	w.WriteHeader(http.StatusOK)

	return
}
