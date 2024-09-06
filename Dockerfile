FROM golang:1.22-alpine AS builder

LABEL maintainer="Nam HQ <namhq.1989@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

FROM alpine:latest

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
