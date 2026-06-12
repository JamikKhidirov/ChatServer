FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ChatServer .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/ChatServer .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/frontend ./frontend
RUN mkdir -p uploads

EXPOSE 8080
ENV GIN_MODE=release
HEALTHCHECK --interval=30s --timeout=3s CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./ChatServer"]
