package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

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

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// Set up connection to DB
	conn, err := pgx.Connect(context.Background(), DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	dbRepository := dbRepository.NewPostgreCompanyRepository(conn)
	csvRepository := csvRepository.NewCompanyCSVRepository()
	companyService := companyService.NewCompanyService(dbRepository, csvRepository)

	httpConector := routes.NewConnector()
	httpConector.ImplementConnector(*companyService)

	path := "./data/q1_catalog.csv"
	companyService.InitializeDataBase(ctx, path)

	path = "./data/q2_clientData.csv"
	companyService.UpdateDataBaseFromCSV(ctx, path)

	router := httpConector.NewRouter()

	log.Print("The server has started")
	server := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			server.Close()
			log.Print("The server has been closed due an error")
		}
	}()

	<-ctx.Done()

	log.Print("The server has been closed")
	shutdownCtx := context.Background()
	err = conn.Close(shutdownCtx)
	if err != nil {
		panic(err)
	}
}
