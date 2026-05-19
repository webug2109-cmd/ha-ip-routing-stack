# Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Download dependencies (mocking go.sum generation)
RUN go mod download || true
COPY . .
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ip-lookup-api .

# Run Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ip-lookup-api .
EXPOSE 8080
CMD ["./ip-lookup-api"]
