FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /chaos-backend

FROM alpine:3.19.0 AS production

WORKDIR /

COPY --from=builder /chaos-backend /chaos-backend

CMD ["/chaos-backend"]