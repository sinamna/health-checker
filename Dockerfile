
## Build
FROM m.docker-registry.ir/golang:1.18.5-alpine3.15

ENV GO111MODULE=auto

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /health_checker

EXPOSE 8080

ENTRYPOINT ["/health_checker"]