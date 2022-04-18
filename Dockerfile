FROM golang:latest

RUN go get github.com/georgysavva/scany/pgxscan
RUN go get ggithub.com/pashagolub/pgxmock
RUN go get github.com/gorilla/mux

RUN mkdir /

ADD . /

WORKDIR /
RUN go build cmd/main.go
ENTRYPOINT ./

EXPOSE 5000