FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

##RUN go build -o app .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/currency

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8081

ENTRYPOINT ["./app"]