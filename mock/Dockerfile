# ステージ1: ビルド環境
FROM golang:1.18 as builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係の管理 (go.sum がない場合はこのステップをスキップ)
COPY go.mod .
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# ステージ2: 実行環境
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたバイナリをステージ1からコピー
COPY --from=builder /app/server .

# ポート8080を公開
EXPOSE 8080

# サーバー実行
CMD ["./server"]
