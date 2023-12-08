# ビルドステージ
FROM golang:1.18-alpine as builder

WORKDIR /app
COPY ./ ./
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o server

# 実行ステージ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
