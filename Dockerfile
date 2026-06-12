FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o iris-apis .

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Taipei
WORKDIR /app
COPY --from=builder /app/iris-apis .
EXPOSE 8800
ENTRYPOINT ["./iris-apis"]
