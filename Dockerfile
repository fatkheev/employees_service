FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN apk add --no-cache git make

COPY . .

RUN go build -o /employee-service

# Установка golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

EXPOSE 8080

CMD ["/app/start.sh"]
