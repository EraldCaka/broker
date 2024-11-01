FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o bin/broker cmd/broker/main.go

EXPOSE 5555

CMD ./bin/api