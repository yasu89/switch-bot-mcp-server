FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /workspace
COPY . /workspace/
RUN go mod download
RUN go build -o switch-bot-mcp-server cmd/switch-bot-mcp-server/main.go

FROM alpine:3.22 AS runtime
COPY --from=builder /workspace/switch-bot-mcp-server /usr/bin/switch-bot-mcp-server

CMD ["switch-bot-mcp-server"]
