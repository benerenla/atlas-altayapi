FROM golang:1.24-alpine

WORKDIR /app

# Air yüklemek için gerekli
RUN go install github.com/air-verse/air@v1.61.7

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Air'i başlat
CMD ["air", "-c", ".air.toml"]