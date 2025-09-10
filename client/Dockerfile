FROM golang:1.25

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/client

RUN go build -o main .

CMD ["./main"]