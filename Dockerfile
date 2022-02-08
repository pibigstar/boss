# myboss
FROM golang:1.16

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

CMD ["go", "run", "main.go"]