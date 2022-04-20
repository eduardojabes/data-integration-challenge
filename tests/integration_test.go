package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/eduardojabes/data-integration-challenge/entity"
	csvRepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/csv"
	dbRepository "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/postgreSQL"
	companyService "github.com/eduardojabes/data-integration-challenge/internal/pkg/service/company"
	routes "github.com/eduardojabes/data-integration-challenge/module/features/routes/company"
	"github.com/jackc/pgx/v4"
)

var (
	port        = int(50051)
	DatabaseUrl = "postgres://postgres:postgres@localhost:5432/data-integration-challenge"
)

func RESTMergeCompaniesWithCSV() (*http.Response, error) {
	path := "q2_clientData.csv"
	file, _ := os.Open(path)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	ff, _ := writer.CreateFormFile("csv", path)
	io.Copy(ff, file)

	writer.Close()
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:5000/v1/companies/merge-all-companies", bytes.NewReader(body.Bytes()))
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)

	return response, err

}

func RESTSearchIntegrationTest(nameForSearch string, zipForSearch string) (*http.Response, entity.Companies, error) {
	dummyIO := []byte{}
	queryURL := fmt.Sprintf("http://localhost:5000/v1/companies/search?name=%s&zip=%s", nameForSearch, zipForSearch)
	fmt.Printf("%s\n", queryURL)
	request, _ := http.NewRequest(http.MethodGet, queryURL, bytes.NewBuffer(dummyIO))

	client := http.Client{}
	response, err := client.Do(request)

	var readCompany entity.Companies
	body, _ := ioutil.ReadAll(io.LimitReader(response.Body, 128*1024*8)) //128kb

	err = json.Unmarshal(body, &readCompany)
	return response, readCompany, err
}

func CreateCompanyIntegrationTest(company entity.Companies) (*http.Response, error) {
	companyJSON, _ := json.Marshal(company)

	request, _ := http.NewRequest(http.MethodPost, "http://localhost:5000/v1/companies", bytes.NewBuffer(companyJSON))

	client := http.Client{}
	response, err := client.Do(request)

	return response, err
}

func TestIntegration(t *testing.T) {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	// Set up connection to DB
	conn, err := pgx.Connect(context.Background(), DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	dbRepository := dbRepository.NewPostgreCompanyRepository(conn)
	csvRepository := csvRepository.NewCompanyCSVRepository()
	companyService := companyService.NewCompanyService(dbRepository, csvRepository)

	httpConector := routes.NewHandler()
	httpConector.ImplementConnector(companyService)

	path := "q1_catalog.csv"
	companyService.InitializeDataBase(ctx, path)

	router := httpConector.NewRouter()

	log.Print("The server has started")
	server := &http.Server{Addr: ":5000", Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			server.Close()
			log.Print("The server has been closed due an error")
		}
	}()

	t.Run("Merging Companies", func(t *testing.T) {
		response, err := RESTMergeCompaniesWithCSV()

		if response.StatusCode != http.StatusOK {
			t.Errorf("error importing websites: got: %d, want: %d", response.StatusCode, http.StatusOK)
		}
		if err != nil {
			t.Errorf("got: %d, want: nil", err)
		}
	})
	t.Run("REST Searching with reponse", func(t *testing.T) {
		response, readCompany, err := RESTSearchIntegrationTest("TOLA", "78229")

		if response.StatusCode != http.StatusOK {
			t.Errorf("error reading companies, got: %d, want: %d", response.StatusCode, http.StatusOK)
		}
		if err != nil {
			t.Errorf("got: %d, want: nil", err)
		}

		var emptyCompany entity.Companies

		if readCompany == emptyCompany {
			t.Errorf("The request need a response, but got %s", readCompany)
		}
	})

	t.Run("REST Searching with empty response", func(t *testing.T) {
		response, readCompany, err := RESTSearchIntegrationTest("TOLA", "")

		if response.StatusCode != http.StatusOK {
			t.Errorf("error reading companies, got: %d, want: %d", response.StatusCode, http.StatusOK)
		}
		if err != nil {
			t.Errorf("got: %d, want: nil", err)
		}

		var emptyCompany entity.Companies

		if readCompany != emptyCompany {
			t.Errorf("The request need a response, but got %s", readCompany)
		}

	})
	t.Run("Creating Company", func(t *testing.T) {
		company := entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://newwebsite.com"}
		response, err := CreateCompanyIntegrationTest(company)

		if response.StatusCode != http.StatusCreated {
			t.Errorf("error creating company: got: %d, want: %d", response.StatusCode, http.StatusOK)
		}
		if err != nil {
			t.Errorf("got: %d, want: nil", err)
		}
	})

	t.Run("Error while Creating Company by having the company", func(t *testing.T) {
		company := entity.Companies{Name: "New Company Test", Zip: "12345", Website: "http://newwebsite.com"}
		response, err := CreateCompanyIntegrationTest(company)

		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("error creating company: got: %d, want: %d", response.StatusCode, http.StatusBadRequest)
		}
		if err != nil {
			t.Errorf("got: %d, want: nil", err)
		}

		_, readCompany, _ := RESTSearchIntegrationTest("Company", "12345")
		companyService.DeleteCompany(ctx, readCompany)
	})

	t.Run("Closing server", func(t *testing.T) {

	})
	cancel()
	<-ctx.Done()

	log.Print("The server has been closed")
	shutdownCtx := context.Background()
	err = conn.Close(shutdownCtx)
	if err != nil {
		panic(err)
	}
}
