FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /chaos-backend

FROM alpine:3.19.0 AS production

WORKDIR /

RUN apk add chromium
COPY --from=builder /chaos-backend /chaos-backend

CMD ["/chaos-backend"]