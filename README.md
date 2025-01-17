
# Data-Integration-Api

The objective of this application is to expose a RESTful API to perform operation over companies data.

## Stack
- Go
- PostgreSQL
- Docker
- Goose

## Endpoints
After start up, the API will be avaible listening port 5000 for following endpoints.

| Name | Path | Method | Content-Type | Description |
| ------ | ------ | ------ | ------ | ------ |
| List all companies| /v1/companies | GET | application/json | Retrieve all companies stored in the database. |
| Search company by name and zip | /v1/companies/search?name={value}&zip={value} | GET | application/json | Provides companies informations based on query parameters values. Company name can be part of the company's name but zip needs to be the entire zip code of the company|
| Create company | /v1/companies | POST | application/json | Create a new company. [here](#post-v1companies)|
| Merge companies with CSV | /v1/companies/merge-all-companies | POST | multipart/form-data | Parses a valid CSV file and integrate its in the actual database. If the will be discarded if ir doesn't exist. The key of the file must be named "csv". See example [here](#post-v1companiesmerge)|

### GET /v1/companies

Response body:

    [{
        "ID": "5e6ab36fe5574a0006e920e7",
        "name":"TOLA SALES GROUP",
        "zip":"78229",
        "website":"http://repsources.com"
    }, ...]

### GET /v1/companies/search

Example: /v1/companies/search?name=TOLA&zip=78229

Response body:

    {
        "ID":"5e6ab36fe5574a0006e920e7",
        "name":"TOLA SALES GROUP",
        "zip":"78229",
        "website":"http://repsources.com"
    }

Example: /v1/companies/search?name=TOLA

Response body:

    {
        "ID":"5e6ab36fe5574a0006e920e7",
        "name":"TOLA SALES GROUP",
        "zip":"78229",
        "website":"http://repsources.com"
    }

### POST /v1/companies

Request body:

    {
        "name": "TOLA SALES GROUP",
        "zipCode": "78229"        
    }

### POST /v1/companies/merge

CSV format:
    
| Name | Address Zip | Website |
| ------ | ------ | ------ |
| TOLA SALES GROUP | 78229 | http://repsources.com |

## Setup

First, you need to have docker and docker-compose installed. The instructions can be found [here](https://docs.docker.com/install/)

You need to have Goose installed too for the SQL table migrations. The instructions can be found [here](https://github.com/pressly/goose#install)

## Container

To run the application execute:

```sh
$ docker-compose up -d
$ make
```
The first command will build the PostgreSQL database.
The second command will construct the table used in this application with the migrations configurations

On first time the application will load data in **q1_catalog.csv**

## Tests

To perform tests with go, run from project root:

```sh
go test ./...
```

When executing tests, on integrations tests the project will merge data with **q2_clientData**

All the queries expected for the server will be tested too

# Data integration challenge


Welcome to Data Integration challenge.

Yawoen company has hired you to implement a Data API for Data Integration team.

Data Integration team is focused on combining data from different heterogeneous sources and providing it to an unified view into entities.

## The challenge

It would be really good if you try to make the code using Go language :)
The other technologies you can feel free to choose.

### 1 - Load treated company data in a database

Read data from CSV file and load into the database to create an entity named **companies**.

This entity should contain the following fields: id, company name and zip code. 

- The loaded data should have the following treatment:
    - **Name:** upper case text
    - **zip:** a five digit text

support file: q1_catalog.csv


### 2 - An API to integrate data using a database

Yawoen now wants to get website data from another source and integrate it with the entity you've just created on the database. When the requirements are met, it's **mandatory** that the **data are merged**.

This new source data must meet the following requirements:

- Input file format: CSV
- Data treatment
    - **Name:** upper case text
    - **zip:** a five digit text
    - **website:** lower case text
- Parameters
    - Name: string
    - Zip: string 
    - Website: string

Build an API to **integrate** `website` data field into the entity records you've just created using **HTTP protocol**.

The `id` field is non existent on the data source, so you'll have to use the available fields to aggregate the new attribute **website** and store it. If the record doesn't exist, discard it.

support file: q2_clientData.csv


### Extra - Matching API to retrieve data

Now Yawoen wants to create an API to provide information getting companies information from the entity to a client. 
The parameters would be `name` and `zip` fields. To query on the database an **AND** logic operator must be used between the fields.

You will need to have a matching strategy because the client might only have a part of the company name. 
Example: "Yawoen" string from "Yawoen Business Solutions".

Output example: 
 ```
 {
 	"id": "abc-1de-123fg",
 	"name": "Yawoen Business Solutions",
 	"zip":"10023",
 	"website": "www.yawoen.com"
 }
 ```

## Notes


- Make sure other developers can easily run the application locally.
- Yawoen isn't picky about the programming language, the database and other tools that you might choose. Just take notice of the market before making your decision.
- Automated tests are mandatory.
- Document your API: fill out a **README.md** file with instructions on how to install and use it.


## Deliverable


- :heavy_check_mark: It would be REALLY nice if it was hosted in a git repo of your **own**. You can create a new empty project, create a branch and Pull Request it to the new master branch you have just created. Provide the PR URL for us so we can discuss the code :grin:. BUT if you'd rather, just compress this directory and send it back to us.
- :heavy_check_mark: Make sure Yawoen folks will have access to the source code.
- :heavy_check_mark: Fill the **Makefile** targets with the apropriated commands (**TODO** tags). That is for easy executing the deliverables (tests and execution). If you have other ideas besides a Makefile feel free to use and reference it on your documentation.
- :x: **Do not** start a Pull Request to this project.

Have fun!
