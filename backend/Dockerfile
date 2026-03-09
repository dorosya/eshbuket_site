FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

# COPY frontend ./frontend

RUN go build -v -o server ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server ./server
# COPY --from=builder /app/frontend ./frontend

EXPOSE 8080
ENTRYPOINT ["./server"]
