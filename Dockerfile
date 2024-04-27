FROM golang:alpine

WORKDIR /opt/app

COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download

COPY . .
RUN go install ./cmd/app/main
ENTRYPOINT ["/go/bin/main"]

EXPOSE 9090
EXPOSE 9091
