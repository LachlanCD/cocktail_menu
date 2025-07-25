FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ /app/

RUN go build -o server /app/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

COPY --from=builder /app/internal ./internal/
COPY --from=builder /app/assets ./assets/

COPY --from=builder /app/data ./data

EXPOSE 4000

CMD ["./server"]

