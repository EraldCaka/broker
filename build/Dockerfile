FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o bin/api cmd/main.go

EXPOSE 5555

CMD ./bin/api