FROM golang:1.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

Run go mod tidy

COPY . .

RUN go build -o app

EXPOSE 3000

CMD ["./app"]
