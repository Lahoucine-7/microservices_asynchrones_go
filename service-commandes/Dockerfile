FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

COPY wait-for-it.sh .

RUN go build -o main ./cmd/main.go

EXPOSE 8081

CMD ["sh", "wait-for-it.sh", "postgres-commandes:5432", "--", "./main"]
