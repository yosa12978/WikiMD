FROM golang:1.19-alpine3.17

WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o ./bin/wikimd ./cmd/WikiMD/main.go
EXPOSE 8089
CMD ["./bin/wikimd"]