FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o migrator cmd/migrator/main.go

EXPOSE 8080
EXPOSE 50051

CMD ["./migrator"]