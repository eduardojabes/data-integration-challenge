package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	company "github.com/eduardojabes/data-integration-challenge/internal/pkg/repository/company/postgreSQL"
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

	repository := company.NewPostgreCompanyRepository(conn)
	companyService := companyService.NewCompanyService(repository)

	httpConector := routes.NewConnector()
	httpConector.ImplementConnector(*companyService)

	path := "/mnt/c/Golang/data-integration-chalenge/data-integration-challenge/data/q1_catalog.csv"
	companyService.InitializeDataBase(ctx, path)

	path = "/mnt/c/Golang/data-integration-chalenge/data-integration-challenge/data/q2_clientData.csv"
	companyService.UpdateDataBase(ctx, path)

	//router := httpConector.NewRouter()

	log.Print("The server has started")
	//log.Fatal(http.ListenAndServe(":5000", router))

	<-ctx.Done()

	shutdownCtx := context.Background()

	err = conn.Close(shutdownCtx)
	if err != nil {
		panic(err)
	}
}
