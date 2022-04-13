#Challenge Makefile
.PHONY: migrate
migrate: 
	goose -dir deployment/migrations postgres "user=postgres password=postgres dbname=data-integration-challenge sslmode=disable" up

start:
#TODO: commands necessary to start the API

check:
#TODO: include command to test the code and show the results

#setup:
#if needed to setup the enviroment before starting it
