FROM golang:1.24-alpine

WORKDIR /app

COPY ../go.mod ../go.sum ./

RUN go mod download

COPY . .
COPY ./web ./web

RUN go build -o order-service ./cmd/order-service && \
    go build -o producer ./cmd/producer && \
    go clean -cache -modcache

EXPOSE 8081

CMD ["/app/order-service"]