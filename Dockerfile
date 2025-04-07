FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:latest AS runner

WORKDIR /root

COPY --from=builder /app/main .

# LOCAL DEV
#COPY .env .

CMD sh -c "sleep 5 && ./main"