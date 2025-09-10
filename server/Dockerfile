FROM golang:1.25

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/server

RUN go build -o main .

CMD ["./main"]