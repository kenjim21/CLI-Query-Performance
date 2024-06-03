FROM golang:1.22.3-alpine3.20

# copy and set up base application
WORKDIR /cli

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
ADD /cmd /cli/cmd

RUN go build
RUN go install

# copy query_params.csv file if desired for testing
COPY query_params.csv ./