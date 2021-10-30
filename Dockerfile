FROM golang:1.17

WORKDIR /app
COPY src .

RUN go get -d -v ./...
RUN go install -v ./...
