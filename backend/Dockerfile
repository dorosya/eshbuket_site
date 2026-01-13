FROM golang:1.25.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -v -o /app/server /app/main.go


FROM alpine
WORKDIR /small
COPY --from=0 /app/server ./binary
ENTRYPOINT [ "./binary" ] 