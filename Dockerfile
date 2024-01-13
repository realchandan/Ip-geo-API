FROM golang:1.21.6 as builder

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

ENV GOPROXY https://proxy.golang.org,direct

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN GOOS=linux go build -o app .

FROM alpine

WORKDIR /app

COPY --from=builder /app/app .

COPY *.csv ./

CMD ["/app/app"]
