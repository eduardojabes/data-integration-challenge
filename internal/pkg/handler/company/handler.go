package company

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/eduardojabes/data-integration-challenge/entity"
	csvRepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/csv"
)

type CompanyService interface {
	AddCompany(ctx context.Context, company *entity.Companies) error
	GetCompanies() ([]entity.Companies, error)
	FindByNameAndZip(name string, zip string) (*entity.Companies, error)
	UpdateCompany(ctx context.Context, company *entity.Companies) error
}

type CompanyHandler struct {
	service CompanyService
}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{}
}

func (c *CompanyHandler) Register(service CompanyService) {
	c.service = service
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(status)

	if bytes.Compare(response, []byte("null")) != 0 {
		w.Write([]byte(response))
	} else {
		w.Write([]byte("{}"))
	}
}

//RespondError makes the error response with payload as json format
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

func (c *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := c.service.GetCompanies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RespondJSON(w, http.StatusOK, companies)
	return
}

//GetCompanyByNameAndZip GET /v1/companies?name={value}&zip={value} application/json
func (c *CompanyHandler) GetCompanyByNameAndZip(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	zip := r.URL.Query().Get("zip")

	companies, err := c.service.FindByNameAndZip(name, zip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RespondJSON(w, http.StatusOK, companies)

	return
}

//CreateCompany POST /v1/companies application/json
func (c *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company entity.Companies
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 128*1024*8)) //128kb

	if err != nil {
		log.Fatal("Failed to add company")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatal("Failed to add company")
	}

	if err := json.Unmarshal(body, &company); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddCompany unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	result := c.service.AddCompany(context.Background(), &company)

	if result != nil {
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	return
}

//MergeCompanies POST /v1/companies multipart/form-data
func (c *CompanyHandler) MergeCompanies(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	companyData := csvRepository.CreateCompanyEntityByCSV(ctx, data)
	for _, company := range companyData {
		err = c.service.UpdateCompany(ctx, company)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
