FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache ca-certificates git tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main .

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./main"] 