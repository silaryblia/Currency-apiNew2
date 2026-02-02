FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

##RUN go build -o app .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app /app/app

RUN apk --no-cache add tzdata ca-certificates

EXPOSE 50051

CMD ["./app"]