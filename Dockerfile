FROM golang:1.17.3-alpine3.15 AS builder
ENV GOPROXY "https://goproxy.cn,direct"
RUN apk add --no-cache g++ git
WORKDIR /go/src/app
COPY go.mod go.sum /go/src/app/
RUN go mod download
COPY . /go/src/app/
RUN CGO_ENABLED=1 GO111MODULE=on GOOS=linux go build -o main main.go

FROM alpine:3.12.0
RUN adduser -D -h /app -u 1000 app
WORKDIR /app
COPY --from=builder /go/src/app/main ./main
EXPOSE 8080
USER 1000
CMD ["/app/main"]
